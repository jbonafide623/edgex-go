package container

import (
    "github.com/edgexfoundry/go-mod-bootstrap/di"
    "github.com/edgexfoundry/go-mod-messaging/messaging"
)

// MessagingClientName contains the name of the messaging client instance in the DIC.
var MessagingClientName = di.TypeInstanceToName((*messaging.MessageClient)(nil))

// MessagingClientFrom helper function queries the DIC and returns the messaging client.
func MessagingClientFrom(get di.Get) messaging.MessageClient {
    return get(MessagingClientName).(messaging.MessageClient)
}
