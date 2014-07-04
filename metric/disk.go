package metric

import (
	"encoding/json"
	"github.com/shirou/gopsutil"
)

type DiskIO struct {
	Partitions             []string
	PrevDiskIOCountersStat map[string]gopsutil.DiskIOCountersStat
}

func NewDiskIO(args []string) (DiskIO, error) {
	return DiskIO{Partitions: args}, nil
}

func (m DiskIO) Emit(metric *Metric) error {
	v, err := gopsutil.DiskIOCounters()
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

	return nil
}
