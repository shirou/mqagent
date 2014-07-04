package metric

import (
	"encoding/json"
	"github.com/shirou/gopsutil"
	"strconv"
)

type LoadAvg struct{}

func NewLoadAvg() (LoadAvg, error) {
	return LoadAvg{}, nil
}

func (m LoadAvg) Get() (map[string]string, error) {
	ret := NewRetJson()
	v, err := gopsutil.LoadAvg()
	if err != nil {
		return ret, err
	}

	ret["load1"] = strconv.FormatFloat(v.Load1, 'f', 2, 64)
	ret["load5"] = strconv.FormatFloat(v.Load5, 'f', 2, 64)
	ret["load15"] = strconv.FormatFloat(v.Load15, 'f', 2, 64)

	return ret, nil
}
func (m LoadAvg) Emit(metric *Metric) error {
	j, err := m.Get()
	if err != nil {
		return err
	}
	topic := metric.TopicRootMetric + metric.ActionId
	payload, err := json.Marshal(j)
	if err != nil {
		return err
	}

	metric.Client.Send(topic, payload, 0)
	return nil
}
