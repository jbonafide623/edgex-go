package device

import (
    "context"
    "fmt"
    "github.com/edgexfoundry/go-mod-core-contracts/clients"
    "github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
    "github.com/edgexfoundry/go-mod-core-contracts/models"
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

    ctx = context.WithValue(ctx, clients.ContentType, clients.ContentTypeJSON)

    uuid, err := uuid.GenerateUUID()
    if err != nil {
        fmt.Errorf("error generating uuid: %s", err.Error())
    }


    ctx = context.WithValue(ctx, clients.CorrelationHeader, uuid)

    n := models.NewDeviceOnboarded{Name:deviceName}
    b,_ := n.MarshalJSON()

    return EventMessagePublisher{
        client,
        types.NewMessageEnvelope(b, ctx),
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
    err = m.client.Publish(m.message, m.topic)
    if err != nil {
        m.lc.Error(fmt.Sprintf("error publishing MQTT message: %s", err.Error()))
        return
    }
    m.lc.Debug("Successfully published message")
    err = m.client.Disconnect()
    if err != nil {
        m.lc.Error("Error disconnecting " + err.Error())
    }
}
