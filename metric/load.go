package metric

import (
	"github.com/shirou/gopsutil"
)

type LoadAvg struct {
	conf MetricConf
}

func NewLoadAvg(mconf MetricConf) (*LoadAvg, error) {
	return &LoadAvg{conf: mconf}, nil
}

func (m *LoadAvg) Get() (string, error) {
	ret, err := gopsutil.LoadAvg()
	if err != nil {
		return "", err
	}
	return ret.String(), nil
}
func (m *LoadAvg) GetConf() MetricConf {
	return m.conf
}
