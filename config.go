package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/shirou/mqagent/transport"
)

const (
	HostIdFileName = "hostid.txt"
)

type ServerAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ClientConfig struct {
	BrokerUri     string     `json:"server"`
	ServerAuth    ServerAuth `json:"serverauth"`
	ConfRoot      string     `json:"confroot"`
	SubList       []string   `json:"subscriptions"`
	Subscriptions []Subscription
	Transport     *transport.MQTTTransport
	HostId        string // HostID
}

func NewConfig(confPath string) (*ClientConfig, error) {
	content, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	var conf ClientConfig
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}

	conf.Subscriptions = make([]Subscription, 0, len(conf.SubList))

	for _, sub := range conf.SubList {
		conf.Subscriptions = append(conf.Subscriptions, Subscription{Name: sub})
	}

	err = conf.setup()
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// setup these.
// - create directories under the ConfRoot
// - create hostid if not created
func (conf *ClientConfig) setup() error {
	root := conf.ConfRoot

	if err := os.MkdirAll(root, 0777); err != nil {
		log.Println(root)
		return err
	}

	hostIdPath := filepath.Join(root, HostIdFileName)
	_, err := os.Stat(hostIdPath)
	if os.IsNotExist(err) {
		h, err := createHostId()
		if err != nil {
			return err
		}
		conf.HostId = h
		// Write to the file
		err = ioutil.WriteFile(hostIdPath, []byte(h), 0777)
		if err != nil {
			return err
		}
	} else {
		// hostid file is alread created. read from it
		h, err := ioutil.ReadFile(hostIdPath)
		if err != nil {
			return err
		}
		conf.HostId = string(h)
	}

	return nil
}

// createHostId creates host id from
//   md5(hostname + all of ipaddresses + random(10000))
// This will stored mqagent directory and used as hostid forever.
func createHostId() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	addrs, err := net.InterfaceAddrs()
	var ipaddress []string
	if err == nil {
		ipaddress = make([]string, 0, len(addrs))
		for _, a := range addrs {
			ipaddress = append(ipaddress, a.String())
		}
	} else {
		ipaddress = make([]string, 0)
	}

	h := md5.New()
	io.WriteString(h, hostname)
	io.WriteString(h, strings.Join(ipaddress, ""))
	io.WriteString(h, string(rand.Intn(10000)))

	hex := fmt.Sprintf("%x", h.Sum(nil))
	return hex, nil
}
