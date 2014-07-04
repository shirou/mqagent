package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/shirou/mqagent/transport"

	"github.com/codegangsta/cli"
	"io/ioutil"
)

var log = logrus.New()
var usage = `
Usage here
`

func initFunc() {
	log.Formatter = new(logrus.TextFormatter)
}

func loadConfig(confPath string) *ClientConfig {
	cliconfig, err := NewConfig(confPath)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	user := cliconfig.ServerAuth.Username
	password := cliconfig.ServerAuth.Password

	cliconfig.Transport = transport.NewMQTTTransport()

	_, err = cliconfig.Transport.Connect(cliconfig.BrokerUri,
		cliconfig.HostId, user, password)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	return cliconfig
}

func startAgent(c *cli.Context) {
	confPath := c.String("c")
	cliconfig := loadConfig(confPath)
	toMainChan := make(chan string)

	for _, sub := range cliconfig.Subscriptions {
		_, err := sub.Subscribe(cliconfig)
		if err != nil {
			log.Println(err)
			log.Println("Skip subscribe: " + sub.Name)
		}
		//		go sub.Start(toMainChan)
	}

	for {
		fromSub, ok := <-toMainChan
		if !ok {
			log.Println("closed")
		}
		log.Println(fromSub)
	}
}

func registerSubscribe(c *cli.Context) {
	confPath := c.String("c")
	cliconfig := loadConfig(confPath)

	subName := c.String("s")
	jsonFile := c.String("f")

	content, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	topic := RoleSubPrefix + "/" + subName
	log.Printf("Register to %s (file: %s)", topic, jsonFile)

	err = cliconfig.Transport.Send(topic, content, 2)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}

func main() {
	initFunc()

	app := cli.NewApp()
	app.Name = "mqagent"
	app.Usage = usage
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start agent",
			Flags: []cli.Flag{
				cli.StringFlag{"c", "/etc/mqagent/client.conf", "client config path"},
			},
			Action: startAgent,
		},
		{
			Name:  "register",
			Usage: "register subscribe",
			Flags: []cli.Flag{
				cli.StringFlag{"c", "/etc/mqagent/client.conf", "client config path"},
				cli.StringFlag{"s", "", "subscribe name"},
				cli.StringFlag{"f", "/etc/mqagent/default.json", "subscribed content file"},
			},
			Action: registerSubscribe,
		},
	}
	app.Run(os.Args)
}
