package metric

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil"
)

type Memory struct {
	conf MetricConf
}
type MemoryStat struct {
	Stat     gopsutil.VirtualMemoryStat `json:"stat"`
	ClientID string                     `json:"clientid"`
}

type SwapMemory struct {
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
	mem := MemoryStat{
		Stat:     *ret,
		ClientID: m.conf.Topic,
	}
	s, _ := json.Marshal(mem)
	fmt.Println(string(s))

	return string(s), nil
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
