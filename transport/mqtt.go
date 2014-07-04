package transport

import (
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

type MQTTTransport struct {
	Opts   *MQTT.ClientOptions
	Client *MQTT.MqttClient
}

func NewMQTTTransport() *MQTTTransport {
	return &MQTTTransport{}
}

// connect to MQTT broker
func (m *MQTTTransport) Connect(brokerUri string, clientId string,
	user string, password string) (*MQTT.MqttClient, error) {

	m.Opts = MQTT.NewClientOptions()

	m.Opts.SetBroker(brokerUri)
	m.Opts.SetClientId(clientId)
	m.Opts.SetTraceLevel(MQTT.Critical)
	m.Opts.SetUsername(user)
	m.Opts.SetPassword(password)

	m.Client = MQTT.NewClient(m.Opts)
	_, err := m.Client.Start()
	if err != nil {
		return nil, err
	}
	return m.Client, nil
}

func (m *MQTTTransport) Send(topic string, payload []byte, qos int) error {
	mqttmsg := MQTT.NewMessage(payload)
	// FIXME: validate qos number
	mqttmsg.SetQoS(MQTT.QoS(qos))
	mqttmsg.SetRetainedFlag(true) // always true

	//	receipt := m.Client.PublishMessage(msg.Destination, mqttmsg)
	receipt := m.Client.PublishMessage(topic, mqttmsg)
	<-receipt

	return nil
}
