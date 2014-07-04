package transport

type Message struct {
	Destination string
	Payload     string
	Qos         int
}
