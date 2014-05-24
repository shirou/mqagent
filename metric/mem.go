package metric

import (
	"github.com/shirou/gopsutil"
)

type MetricConf struct {
	MetricType string
	Interval   int // if 0, get onlye once
	QoS        int
	Topic     string
}

type Metric interface {
	Get() (string, error) // return JSON
	GetConf() MetricConf
}

type Memory struct {
	conf MetricConf
}
type SwapMemory struct {
	conf MetricConf
}
type HostInfo struct {
	conf MetricConf
}
type LoadAvg struct {
	conf MetricConf
}

func NewMemory(mconf MetricConf) (*Memory, error) {
	return &Memory{conf: mconf}, nil
}

func (m *Memory) Get() (string, error) {
	ret, err := gopsutil.VirtualMemory()
	if err != nil {
		return "", err
	}
	return ret.String(), nil
}
func (m *Memory) GetConf() MetricConf {
	return m.conf
}


func NewSwapMemory(mconf MetricConf) (*SwapMemory, error) {
	return &SwapMemory{conf: mconf}, nil
}

func (m *SwapMemory) Get() (string, error) {
	ret, err := gopsutil.SwapMemory()
	if err != nil {
		return "", err
	}
	return ret.String(), nil
}
func (m *SwapMemory) GetConf() MetricConf {
	return m.conf
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
