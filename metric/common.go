package metric

type MetricConf struct {
	MetricType string
	Interval   int // if 0, get onlye once
	QoS        int
	Topic      string
}

type Metric interface {
	Get() (string, error) // return JSON
	GetConf() MetricConf
}
