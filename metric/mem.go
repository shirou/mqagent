package metric

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/mem"

	"github.com/shirou/mqagent/transport"
)

type Memory struct {
	LastMemoryStat *mem.VirtualMemoryStat
	Data           map[string]string
}

type SwapMemory struct {
	LastSwapStat *mem.SwapMemoryStat
	Data         map[string]string
}

func NewMemory() (Memory, error) {
	return Memory{
		Data: make(map[string]string),
	}, nil
}
func NewSwapMemory() (SwapMemory, error) {
	return SwapMemory{
		Data: make(map[string]string),
	}, nil
}

const (
	MEMFLOATUNIT  = 1000000.0
	MEMFLOATCOUNT = 3
)

func (m Memory) Update() error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	m.LastMemoryStat = v

	m.Data["time"] = strconv.FormatInt(time.Now().Unix(), 10)

	m.Data["total"] = strconv.FormatFloat(float64(v.Total)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	m.Data["available"] = strconv.FormatFloat(float64(v.Available)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	m.Data["used"] = strconv.FormatFloat(float64(v.Used)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	m.Data["free"] = strconv.FormatFloat(float64(v.Free)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	m.Data["active"] = strconv.FormatFloat(float64(v.Active)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	if v.Cached > 0 {
		m.Data["cached"] = strconv.FormatFloat(float64(v.Cached)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	}
	if v.Wired > 0 {
		m.Data["wired"] = strconv.FormatFloat(float64(v.Wired)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	}

	return nil
}

func (m SwapMemory) Update() error {
	v, err := mem.SwapMemory()
	if err != nil {
		return err
	}

	m.LastSwapStat = v

	m.Data["time"] = strconv.FormatInt(time.Now().Unix(), 10)

	m.Data["total"] = strconv.FormatFloat(float64(v.Total)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	m.Data["used"] = strconv.FormatFloat(float64(v.Used)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	m.Data["free"] = strconv.FormatFloat(float64(v.Free)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	m.Data["sin"] = strconv.FormatFloat(float64(v.Sin)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)
	m.Data["sout"] = strconv.FormatFloat(float64(v.Sout)/MEMFLOATUNIT, 'f', MEMFLOATCOUNT, 64)

	return nil
}

func (m Memory) Emit(to chan transport.Message) error {
	j, err := json.Marshal(m.Data)
	if err != nil {
		return err
	}
	to <- transport.Message{
		Type:        "metric",
		Destination: "memory",
		Payload:     j,
	}
	return nil
}
func (m SwapMemory) Emit(to chan transport.Message) error {
	j, err := json.Marshal(m.Data)
	if err != nil {
		return err
	}
	to <- transport.Message{
		Type:        "metric",
		Destination: "swap",
		Payload:     j,
	}
	return nil
}
