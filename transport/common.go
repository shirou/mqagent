package transport

type Message struct {
	Type        string // ex: subscribe, metric
	Destination string // ex: topic
	Payload     []byte // ex: json
	Qos         int    // ex: 0
}
