package main

import (
	"errors"

	"encoding/json"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/shirou/mqagent/metric"
)

const (
	MAX_METRICS_NUM int = 100
)

type Subscription struct {
	Name           string
	ToMainChan     chan string // channel to the main routine
	FromMainChan   chan string // channel from the main routine
	FromActionChan chan string

	Conf    *ClientConfig
	Metrics map[string]*metric.Metric
	Topic   string
}

const (
	TopicRoot     = "/mqvision"
	RoleSubPrefix = TopicRoot + "/role"
	RoleSubQos    = MQTT.QOS_TWO
)

func (sub *Subscription) Subscribe(conf *ClientConfig) (chan string, error) {
	topic := RoleSubPrefix + "/" + sub.Name

	if receipt, err := conf.Transport.Client.StartSubscription(sub.messageHandler,
		topic, RoleSubQos); err != nil {
		return nil, err
	} else {
		<-receipt
	}

	log.Printf("Role Subscribed: %s", topic)

	sub.ToMainChan = make(chan string)
	sub.FromActionChan = make(chan string)
	sub.Conf = conf
	sub.Topic = topic
	return sub.ToMainChan, nil
}

func (sub *Subscription) messageHandler(msg MQTT.Message) {

	for actionId, metric := range sub.Metrics {
		metric.Stop()
		log.Printf("%s Stopped", actionId)
	}

	err := sub.ParseRoleConfig(msg.Payload())
	if err != nil {
		log.Error(err)
		return
	}

	for actionId, m := range sub.Metrics {
		log.Printf("Starting: %s", actionId)
		go m.Start(sub.FromActionChan)
	}
}
func (sub *Subscription) ParseRoleConfig(buf []byte) error {
	actions := make([]ActionJson, 1)
	err := json.Unmarshal([]byte(buf), &actions)
	if err != nil {
		return err
	}

	sub.Metrics = make(map[string]*metric.Metric, MAX_METRICS_NUM)

	for _, actionJson := range actions {
		if actionJson.Type == "" {
			continue
		}

		switch actionJson.Type {
		case "metric":
			a, err := metric.NewMetric(actionJson.String(),
				sub.Conf.HostId, sub.Conf.Transport, TopicRoot)
			if err != nil {
				log.Error(err)
			} else {
				// if same actionid, override it. TODO: is it safe?
				sub.Metrics[actionJson.ActionId] = a
			}
		default:
			err := errors.New("No such a action type:" + actionJson.Type)
			log.Error(err)
		}
	}

	return nil
}
