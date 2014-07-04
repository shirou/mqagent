package metric

import (
	"fmt"
	"time"

	"encoding/json"
	"errors"
	"github.com/shirou/mqagent/transport"
	"strconv"
)

type Metric struct {
	ActionId string   `json:"actionid"`
	Action   string   `json:"action"`
	Args     []string `json:"args"`
	Type     string   `json:"type"`
	Interval int      `json:"interval"`
	Ticker   *time.Ticker

	TopicRootMetric string

	HostId      string
	Client      *transport.MQTTTransport
	MetricValue MetricValue
}

type MetricValue interface {
	Emit(*Metric) error
}

func NewMetric(jsonBuf string,
	hostId string,
	client *transport.MQTTTransport,
	topicRoot string) (*Metric, error) {
	ret := &Metric{
		HostId:          hostId,
		Client:          client,
		TopicRootMetric: topicRoot + "/metrics/" + hostId + "/",
	}

	b := []byte(jsonBuf)
	err := json.Unmarshal(b, ret)
	if err != nil {
		return nil, err
	}

	// TODO: change to map
	var m MetricValue
	switch ret.ActionId {
	case "metrics.memory":
		m, err = NewMemory()
	case "metrics.swap":
		m, err = NewSwapMemory()
	case "metrics.load":
		m, err = NewLoadAvg()
	case "metrics.diskio":
		m, err = NewDiskIO(ret.Args)
	default:
		err := errors.New("No such actionid: " + ret.ActionId)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	ret.MetricValue = m

	return ret, nil
}

func (m *Metric) Start(ch chan string) error {
	m.Ticker = time.NewTicker(time.Duration(m.Interval) * time.Second)

	tickChan := m.Ticker.C

	for {
		select {
		case <-tickChan:
			err := m.MetricValue.Emit(m)
			if err != nil {
				// FIXME: logging
				continue
			}
		case <-ch:
			fmt.Println("From main routine")
		}
	}
	return nil
}

func (m *Metric) Stop() error {
	m.Ticker.Stop()
	return nil
}

func (m *Metric) SetInterval(interval int) error {
	m.Interval = interval

	return nil
}

func NewRetJson() map[string]string {
	ret := make(map[string]string, 100)

	utime := time.Now().Unix()
	ret["time"] = strconv.FormatInt(utime, 10)

	return ret
}
