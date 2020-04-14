package servicediscovery

import (
	"fmt"
	dubboregistry "github.com/mosn/registry/dubbo"
	"github.com/mosn/registry/dubbo/remoting"
	v2 "mosn.io/mosn/pkg/config/v2"
	clusterAdapter "mosn.io/mosn/pkg/upstream/cluster"
)

// listener listens for registry subscription data change
type listener struct{}

func (l listener) Notify(event *dubboregistry.ServiceEvent) {
	fmt.Println("###", event.Action, event.Service, event.String())

	var (
		err         error
		clusterName = event.Service.Service() // FIXME
		addr        = event.Service.Ip + ":" + event.Service.Port
	)

	switch event.Action {
	case remoting.EventTypeAdd:
		// call the cluster manager to add the host
		err = clusterAdapter.GetClusterMngAdapterInstance().TriggerHostAppend(clusterName, []v2.Host{
			{
				HostConfig: v2.HostConfig{
					Address: addr,
				},
			},
		})
	case remoting.EventTypeDel:
		// call the cluster manager to del the host
		err = clusterAdapter.GetClusterMngAdapterInstance().TriggerHostDel(clusterName, []string{addr})
	case remoting.EventTypeUpdate:
		fallthrough
	default:
		fmt.Println("not supported")
	}

	if err != nil {
		// TODO, log error
	}
}
