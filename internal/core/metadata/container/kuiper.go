package container

import (
    "github.com/edgexfoundry/edgex-go/internal/core/metadata/operators/device"
    "github.com/edgexfoundry/go-mod-bootstrap/di"
)

var KuiperClientName = di.TypeInstanceToName((*device.KuiperClient)(nil))

func KuiperClientFrom(get di.Get) device.KuiperClient {
    return get(KuiperClientName).(device.KuiperClient)
}