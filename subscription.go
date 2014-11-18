package main

import (
	"encoding/json"
	"fmt"
	"strings"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	log "github.com/Sirupsen/logrus"

	"github.com/shirou/mqagent/metric"
	"github.com/shirou/mqagent/transport"
)

const (
	MAX_METRICS_NUM int = 100
)

type Subscription struct {
	Name       string
	ToMainChan chan transport.Message // channel to the main routine
	Conf       *ClientConfig
	Topic      string
	Metrics    []string
}

const (
	RoleRoot   = "roles"
	RoleSubQos = 2
)

func (sub *Subscription) Subscribe(conf *ClientConfig, subch, ach chan transport.Message) error {
	topic := strings.Join([]string{conf.TopicRoot, RoleRoot, sub.Name}, "/")

	if receipt, err := conf.Transport.Subscribe(sub.messageHandler,
		topic, RoleSubQos); err != nil {
		return err
	} else {
		<-receipt
	}

	log.Infof("Role Subscribed: %s", topic)

	sub.ToMainChan = ach
	sub.Conf = conf
	sub.Topic = topic
	return nil
}

func (sub *Subscription) messageHandler(client *MQTT.MqttClient, msg MQTT.Message) {
	actions, err := sub.ParseRoleConfig(msg.Payload())
	if err != nil {
		log.Error(err)
		return
	}

	for _, actionJson := range actions {
		if actionJson.Type == "" {
			continue
		}
		switch actionJson.Type {
		case "metric":
			a, err := metric.NewMetric(actionJson.Byte(),
				sub.Conf.HostId, sub.ToMainChan)
			if err != nil {
				log.Error(err)
				continue
			}
			// if same actionid over all subscribes, error and skip
			MetricMapLock.Lock()
			_, exists := MetricMap[actionJson.Id]
			if exists {
				log.Errorf("%s is duplicated at %s", actionJson.Id, sub.Name)
			} else {
				log.Infof("subscribing %s", actionJson.Id)
				MetricMap[actionJson.Id] = a
			}
			MetricMapLock.Unlock()
		default:
			err := fmt.Errorf("No such a action type: %s", actionJson.Type)
			log.Error(err)
		}
	}

	if err != nil {
		log.Error(err)
		return
	}

	for id, m := range MetricMap {
		if m.Status == metric.StatusStarted {
			log.Infof("already started: %s", id)
			continue
		}
		log.Infof("Starting: %s", id)
		go m.Start()
	}

}
func (sub *Subscription) ParseRoleConfig(buf []byte) ([]ActionJson, error) {
	var actions []ActionJson
	err := json.Unmarshal([]byte(buf), &actions)
	if err != nil {
		return actions, err
	}

	return actions, nil
}
