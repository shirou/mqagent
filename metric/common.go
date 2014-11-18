package metric

import (
	"encoding/json"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/shirou/mqagent/transport"
)

type Updater interface {
	Update() error
	Emit(chan transport.Message) error
}

const (
	StatusUnknown = iota
	StatusStopped
	StatusStarted
)

type Metric struct {
	Id       string   `json:"id"`
	Action   string   `json:"action"`
	Args     []string `json:"args"`
	Type     string   `json:"type"`
	Interval int      `json:"interval"`

	LastUpdated time.Time
	MetricValue Updater

	ToChan chan transport.Message
	Ticker *time.Ticker
	HostId string
	Status int
}

func NewMetric(jsonBuf []byte, hostId string, toChan chan transport.Message) (*Metric, error) {
	ret := &Metric{
		HostId: hostId,
		ToChan: toChan,
	}

	err := json.Unmarshal(jsonBuf, ret)
	if err != nil {
		return nil, err
	}

	var m Updater
	switch ret.Action {
	case "metrics.memory":
		m, err = NewMemory()
	case "metrics.swap":
		m, err = NewSwapMemory()
		/*
			case "metrics.load":
						m, err = NewLoadAvg()
		*/
	case "metrics.diskio":
		m, err = NewDiskIO(ret.Args)
	default:
		err := errors.New("No such action: " + ret.Action)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	ret.MetricValue = m

	return ret, nil
}

func (m *Metric) Start() {
	m.Ticker = time.NewTicker(time.Duration(m.Interval) * time.Second)

	m.Status = StatusStarted
	tickChan := m.Ticker.C

	// emit first time
	m.MetricValue.Update()
	m.MetricValue.Emit(m.ToChan)

	for {
		select {
		case <-tickChan:
			err := m.MetricValue.Update()
			if err != nil {
				log.Error(err)
				continue
			}
			err = m.MetricValue.Emit(m.ToChan)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
}

func (m *Metric) Stop() error {
	m.Ticker.Stop()
	m.Status = StatusStopped
	return nil
}

func (m *Metric) SetInterval(interval int) error {
	m.Interval = interval

	return nil
}
