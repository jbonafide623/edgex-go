package device

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-messaging/messaging"
	"github.com/edgexfoundry/go-mod-messaging/pkg/types"
	"github.com/hashicorp/go-uuid"
)

type Command interface {
	execute()
}

type MessagePublisher struct {
	commands []Command
}

func NewMessagePublisher(commands []Command) MessagePublisher {
	return MessagePublisher{commands }
}

func (d MessagePublisher) Execute() {
	for _,c := range d.commands {
		c.execute()
	}
}

type EventMessagePublisher struct {
	client  messaging.MessageClient
	message types.MessageEnvelope
	topic   string
	lc      logger.LoggingClient
}

func NewEventMessagePublisher(
	client messaging.MessageClient,
	deviceName string,
	topic string,
	ctx context.Context,
	lc logger.LoggingClient) EventMessagePublisher {

	//payload := newDeviceMessagePayload(deviceName).toString()
	ctx = context.WithValue(ctx, clients.ContentType, clients.ContentTypeJSON)

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		fmt.Errorf("error generating uuid: %s", err.Error())
	}
	ctx = context.WithValue(ctx, clients.CorrelationHeader, uuid)


	return EventMessagePublisher{
		client,
		types.NewMessageEnvelope([]byte(deviceName), ctx),
		topic,
		lc,
	}
}

func (m EventMessagePublisher) execute() {
	m.lc.Debug("Connecting to MQTT Broker")
	err := m.client.Connect()

	if err != nil {
		m.lc.Error(fmt.Sprintf("MQTT broker connection error occurred: %s", err.Error()))
		return
	}

	m.lc.Debug(fmt.Sprintf("Publishing message: %s to topic: %s", m.message, m.topic))
	//payload := string(m.message.Payload)
	//payload = strings.Replace(payload, "\"", "\\\"", -1)
	//m.message.Payload = []byte(base64.StdEncoding.EncodeToString(m.message.Payload))
	err = m.client.Publish(m.message, m.topic)
	if err != nil {
		m.lc.Error(fmt.Sprintf("error publishing MQTT message: %s", err.Error()))
		return
	}
	m.lc.Debug("Successfully published message")
}

type DeviceMessagePayload struct {
	Name string `json:"name"`
}

func newDeviceMessagePayload(name string)  DeviceMessagePayload {
	return DeviceMessagePayload{
		Name: name,
	}
}

func (n DeviceMessagePayload) toString() []byte {
	b, err := json.Marshal(n)
	if err != nil {
		fmt.Errorf("error marshalling: %s", err.Error())
	}

	return b
}
