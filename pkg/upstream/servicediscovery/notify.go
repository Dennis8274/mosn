package servicediscovery

import (
	"fmt"
	dubboregistry "github.com/mosn/registry/dubbo"
	"github.com/mosn/registry/dubbo/remoting"
)

// listener listens for registry subscription data change
type listener struct {}

func (l listener) Notify(event *dubboregistry.ServiceEvent) {
	fmt.Println(event.Action, event.Service, event.String())

	switch event.Action {
	case remoting.EventTypeAdd:
		// call the cluster manager to add the host
	case remoting.EventTypeDel:
		// call the cluster manager to del the host
	case remoting.EventTypeUpdate:
		fmt.Println("not supported")
	default:
		fmt.Println("not supported")
	}
}

