package device

import (
    "fmt"
    "github.com/edgexfoundry/edgex-go/internal/core/metadata/interfaces"
    "github.com/edgexfoundry/go-mod-messaging/messaging"
    "github.com/edgexfoundry/go-mod-messaging/pkg/types"
    "reflect"
)

type Observerable interface {
    registerObserver(o Observer)
    unregisterObserver(o Observer)
    notify(m types.MessageEnvelope)
}

type Observer interface {
    handleMessage(m types.MessageEnvelope)
}

type BlacklistDeviceObserverable struct {
    observers     []Observer
    messageClient messaging.MessageClient
    dbClient      interfaces.DBClient
    message       types.MessageEnvelope
}

func NewBlacklistDeviceObservable(mc messaging.MessageClient) *BlacklistDeviceObserverable {
    b := &BlacklistDeviceObserverable{
        observers:     make([]Observer, 0),
        messageClient: mc,
    }
    b.listen()
    return b
}

func (b *BlacklistDeviceObserverable) registerObserver(o Observer) {
    b.observers = append(b.observers, o)
}

func (b *BlacklistDeviceObserverable) unregisterObserver(o Observer) {
    var match int
    for index, i := range b.observers {
        if reflect.DeepEqual(i, o) {
            match = index
            return
        }
    }

    b.observers = append(b.observers[:match], b.observers[match+1:]...)
}

func (b *BlacklistDeviceObserverable) notify(m types.MessageEnvelope) {
    fmt.Printf("Unregistering device %s", m.Payload)
    for _, i := range b.observers {
        i.handleMessage(m)
    }
}

func (b *BlacklistDeviceObserverable) listen() {
    go func() {
        // Tie in to go-mod-messaging (from MQTT Server)
        fmt.Println("Connecting as subscriber")
        if err := b.messageClient.Connect(); err != nil {
            fmt.Errorf("Error: %s", err.Error())
        }

        msg := make(chan types.MessageEnvelope)
        topics := []types.TopicChannel {
            {
                Topic: "blacklist_device",
                Messages: msg,
            },
        }

        messageErrors := make(chan error)
        err := b.messageClient.Subscribe(topics, messageErrors)
        if err != nil {
            fmt.Errorf("Failed to subscribe for event messages: " + err.Error())
            return
        }

        for {
            fmt.Printf("Loop")
            select {
            case e := <- messageErrors:
                fmt.Errorf("Message Error: %s", e.Error())
            case msgEnv := <-msg:
                b.notify(msgEnv)
            }
        }

        // Add message to envelope
    }()
}

type UnregisterDeviceObserver struct {
    Observerable Observerable
    dbClient interfaces.DBClient
}

func NewUnregisterDeviceObserver(o Observerable, dbClient interfaces.DBClient) *UnregisterDeviceObserver {
    u := UnregisterDeviceObserver{
        Observerable: o,
        dbClient: dbClient,
    }
    o.registerObserver(&u)
    return &u
}

func (n *UnregisterDeviceObserver) handleMessage(m types.MessageEnvelope) {
    //payload := string(m.Payload)
    //payload = strings.Replace(payload, "\\", "", -1)
    //m.Payload = []byte(payload)
    //b, _ := base64.StdEncoding.DecodeString(string(m.Payload))
    //var dp DeviceMessagePayload
    //_ = json.Unmarshal(m.Payload, &dp)

    fmt.Printf("Unregister Handler: %s\n", m.Payload)
    d, err := n.dbClient.GetDeviceByName(string(m.Payload))
    if err != nil {
        fmt.Errorf("Error: %s", err.Error())
    }

    err = n.dbClient.DeleteDeviceById(d.Id)
    if err != nil {
        fmt.Errorf("Error: %s", err.Error())
    }
}