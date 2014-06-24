package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/shirou/mqagent/metric"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	log "github.com/cihub/seelog"
)

func pub(client MQTT.MqttClient, qos int, topic string, payload string) {
}

func NewMetrics(metrics []metric.MetricConf) []metric.Metric {
	ret := make([]metric.Metric, 0, len(metrics))

	for _, mconf := range metrics {
		var mm metric.Metric
		var err error
		switch mconf.MetricType {
		case "memory":
			mm, err = metric.NewMemory(mconf)
			if err != nil {
				continue
			}
		case "swap":
			mm, err = metric.NewSwapMemory(mconf)
			if err != nil {
				continue
			}
		case "host":
			mm, err = metric.NewHostInfo(mconf)
			if err != nil {
				continue
			}
		case "load":
			mm, err = metric.NewLoadAvg(mconf)
			if err != nil {
				continue
			}
		}
		ret = append(ret, mm)
	}

	return ret
}

func main() {
	testConfig := `
<seelog>
  <outputs formatid="develop">
    <console />
    <rollingfile type="size" filename="./log/manyrolls.log" maxsize="100" maxrolls="4" />
  </outputs>
  <formats>
    <format id="develop" format="%Date %Time [%LEV] %Msg%n"/>
  </formats>
</seelog>
`
	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)

	broker := flag.String("broker", "", "The broker URI. ex: tcp://10.10.1.1:1883")
	password := flag.String("password", "", "The password (optional)")
	user := flag.String("user", "", "The User (optional)")
	clientid := flag.String("id", "", "The ClientID (optional)")
	//	cleansess := flag.Bool("clean", false, "Set Clean Session (default false)")
	store := flag.String("store", ":memory:", "The Store Directory (default use memory store)")
	flag.Parse()
	if *broker == "" {
		fmt.Println("Invalid setting for -broker")
		return
	}
	log.Infof("broker:    %s", *broker)
	log.Infof("clientid:  %s", *clientid)
	log.Infof("broker:    %s", *broker)
	log.Infof("clientid:  %s", *clientid)
	log.Infof("user:      %s", *user)
	log.Infof("password:  %s", *password)
	log.Infof("store:     %s", *store)
	opts := MQTT.NewClientOptions()
	opts.SetBroker(*broker)
	opts.SetTraceLevel(MQTT.Off)
	opts.SetClientId(*clientid)
	opts.SetUsername(*user)
	opts.SetPassword(*password)
	if *store != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(*store))
	}

	client := MQTT.NewClient(opts)
	_, err := client.Start()
	if err != nil {
		fmt.Println(err)
	}

	var metrics []metric.MetricConf

	metrics = append(metrics, metric.MetricConf{
		MetricType: "memory",
		Interval:   10,
		Topic:      "/mqagent/metrics/memory",
	})
	metrics = append(metrics, metric.MetricConf{
		MetricType: "swap",
		Interval:   10,
		Topic:      "/mqagent/metrics/swap",
	})
	metrics = append(metrics, metric.MetricConf{
		MetricType: "host",
		Topic:      "/mqagent/metrics/host",
	})
	metrics = append(metrics, metric.MetricConf{
		MetricType: "load",
		Interval:   10,
		Topic:      "/mqagent/metrics/load",
	})

	mm := NewMetrics(metrics)

	for _, m := range mm {
		go func(client MQTT.MqttClient, metric metric.Metric) {
			for {
				payload, err := metric.Get()
				if err != nil {
					continue
				}
				conf := metric.GetConf()

				msg := MQTT.NewMessage([]byte(payload))
				msg.SetQoS(MQTT.QoS(conf.QoS))
				msg.SetRetainedFlag(true) // always true

				receipt := client.PublishMessage(conf.Topic, msg)
				<-receipt
				if conf.Interval == 0 {
					return
				} else {
					time.Sleep(time.Duration(conf.Interval) * time.Second)
				}
			}
		}(*client, m)
	}
	for {
		time.Sleep(time.Second)
	}
}
