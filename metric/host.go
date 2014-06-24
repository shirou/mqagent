package metric

import (
	"github.com/shirou/gopsutil"
)

type HostInfo struct {
	conf MetricConf
}

func NewHostInfo(mconf MetricConf) (*HostInfo, error) {
	return &HostInfo{conf: mconf}, nil
}

func (m *HostInfo) Get() (string, error) {
	ret, err := gopsutil.HostInfo()
	if err != nil {
		return "", err
	}
	return ret.String(), nil
}
func (m *HostInfo) GetConf() MetricConf {
	return m.conf
}
