package metric

import (
	"github.com/shirou/gopsutil/disk"

	"github.com/shirou/mqagent/transport"
)

type DiskIO struct {
	Partitions             []string
	PrevDiskIOCountersStat map[string]disk.DiskIOCountersStat
	Data                   map[string]map[string]string
}

func NewDiskIO(args []string) (DiskIO, error) {
	return DiskIO{Partitions: args}, nil
}

func (m DiskIO) Update() error {
	v, err := disk.DiskIOCounters()
	if err != nil {
		return err
	}

	m.PrevDiskIOCountersStat = v
	keys := make([]string, 0, len(v))

	// get target partitions
	if len(m.Partitions) == 0 { // if not specified, get from all
		for key := range v {
			keys = append(keys, key)
		}
		m.Partitions = keys
	} else {
		keys = m.Partitions
	}

	return nil
}

func (m DiskIO) Emit(ch chan transport.Message) error {
	/*
		for _, part := range keys {
			stat, ok := v[part]
			if !ok {
				// FIXME: Logging
				// log.Printf("disk partiton %s is not exists. skip it", part)
				continue
			}
			// ex: topicroot/metrics.disk.sda
			topic := metric.TopicRootMetric + metric.ActionId + "." + part
			payload, err := json.Marshal(stat)
			if err != nil {
				return err
			}

			metric.Client.Send(topic, payload, 0)
		}
	*/

	return nil
}
