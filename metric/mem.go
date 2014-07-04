package metric

import (
	"encoding/json"
	"github.com/shirou/gopsutil"
	"strconv"
)

type Memory struct {
	LastMemoryStat *gopsutil.VirtualMemoryStat
}

type SwapMemory struct {
	LastSwapStat *gopsutil.SwapMemoryStat
}

func NewMemory() (Memory, error) {
	return Memory{}, nil
}
func NewSwapMemory() (SwapMemory, error) {
	return SwapMemory{}, nil
}

const (
	MEMFLOATUNIT  = 1000000.0
	MEMFLOATCOUNT = 3
)

func (m Memory) Get() (map[string]string, error) {
	ret := NewRetJson()
	v, err := gopsutil.VirtualMemory()
	if err != nil {
		return ret, err
	}

	m.LastMemoryStat = v

	ret["total"] = strconv.FormatFloat(float64(v.Total)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	ret["available"] = strconv.FormatFloat(float64(v.Available)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	ret["used"] = strconv.FormatFloat(float64(v.Used)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	ret["free"] = strconv.FormatFloat(float64(v.Free)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	ret["active"] = strconv.FormatFloat(float64(v.Active)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	if v.Cached > 0 {
		ret["cached"] = strconv.FormatFloat(float64(v.Cached)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	}
	if v.Wired > 0 {
		ret["wired"] = strconv.FormatFloat(float64(v.Wired)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	}

	return ret, nil
}

func (m SwapMemory) Get() (map[string]string, error) {
	ret := NewRetJson()
	v, err := gopsutil.SwapMemory()
	if err != nil {
		return ret, err
	}
	m.LastSwapStat = v

	ret["total"] = strconv.FormatFloat(float64(v.Total)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	ret["used"] = strconv.FormatFloat(float64(v.Used)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	ret["free"] = strconv.FormatFloat(float64(v.Free)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	ret["sin"] = strconv.FormatFloat(float64(v.Sin)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	ret["sout"] = strconv.FormatFloat(float64(v.Sout)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)

	return ret, nil
}

func (m Memory) Emit(metric *Metric) error {
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
func (m SwapMemory) Emit(metric *Metric) error {
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
