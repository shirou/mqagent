package main

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/shirou/mqagent/metric"
	"github.com/shirou/mqagent/transport"
)

var usage = `
Usage here
`

var MetricMap map[string]*metric.Metric
var MetricMapLock *sync.Mutex

func init() {
	log.SetFormatter(&LTSVFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)

	MetricMap = make(map[string]*metric.Metric)
	MetricMapLock = &sync.Mutex{}
}

func loadConfig(confPath string) *ClientConfig {
	conf, err := NewConfig(confPath)
	if err != nil {
		log.Fatal(err)
	}
	user := conf.ServerAuth.Username
	password := conf.ServerAuth.Password

	conf.Transport = transport.NewMQTTTransport()

	log.Debug(user)
	log.Debug(conf.BrokerUri)

	_, err = conf.Transport.Connect(conf.BrokerUri,
		conf.HostId, user, password)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

func startAgent(c *cli.Context) {
	if c.Bool("d") {
		log.SetLevel(log.DebugLevel)
	}

	confPath := c.String("c")
	conf := loadConfig(confPath)

	actionChan := make(chan transport.Message)
	subChan := make(chan transport.Message)

	for _, sub := range conf.Subscriptions {
		err := sub.Subscribe(conf, subChan, actionChan)
		if err != nil {
			log.Error(err)
			log.Info("Skip subscribe: " + sub.Name)
		}
	}

	mainDispatcher(conf, subChan, actionChan)
}

func mainDispatcher(conf *ClientConfig, subChan, actionChan chan transport.Message) {
	for {
		select {
		case sub, ok := <-subChan:
			if !ok {
				log.Error("closed")
			}
			log.Info(sub, "hogehgoe")
		case action, ok := <-actionChan:
			if !ok {
				log.Error("closed")
			}
			topic := strings.Join([]string{conf.TopicRoot,
				action.Type,
				action.Destination,
			}, "/")
			log.WithFields(log.Fields{
				"topic": topic,
			}).Debug("send")
			conf.Transport.Send(topic, action.Payload, action.Qos)
		}
	}
}

func registerSubscribe(c *cli.Context) {
	if c.Bool("d") {
		log.SetLevel(log.DebugLevel)
	}
	confPath := c.String("c")
	conf := loadConfig(confPath)

	subName := c.String("s")
	jsonFile := c.String("f")

	content, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	topic := strings.Join([]string{conf.TopicRoot, RoleRoot, subName}, "/")
	log.Infof("Register to %s (file: %s)", topic, jsonFile)

	err = conf.Transport.Send(topic, content, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "mqagent"
	app.Usage = usage
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start agent",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "d", Usage: "verbose"},
				cli.StringFlag{
					Name:   "c",
					Value:  "/etc/mqagent/client.conf",
					Usage:  "client config path",
					EnvVar: "MQAGENT_CONF_PATH",
				},
			},
			Action: startAgent,
		},
		{
			Name:  "register",
			Usage: "register subscribe",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "d", Usage: "verbose"},
				cli.StringFlag{
					Name:   "c",
					Value:  "/etc/mqagent/client.conf",
					Usage:  "client config path",
					EnvVar: "MQAGENT_CONF_PATH",
				},
				cli.StringFlag{
					Name:  "s",
					Usage: "subscribe name",
				},
				cli.StringFlag{
					Name:  "f",
					Value: "/etc/mqagent/default.json",
					Usage: "subscribed content file",
				},
			},
			Action: registerSubscribe,
		},
	}
	app.Run(os.Args)
}
