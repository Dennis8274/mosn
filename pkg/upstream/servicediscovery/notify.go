package servicediscovery

import (
	dubboregistry "github.com/mosn/registry/dubbo"
	"github.com/mosn/registry/dubbo/remoting"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	clusterAdapter "mosn.io/mosn/pkg/upstream/cluster"
)

// listener listens for registry subscription data change
type listener struct{}

func (l listener) Notify(event *dubboregistry.ServiceEvent) {
	var (
		err         error
		clusterName = event.Service.Service() // FIXME
		addr        = event.Service.Ip + ":" + event.Service.Port
	)

	/*
	FIXME, if there is need to load full provider list, then need to finish this logic, and overwrite the hosts in cluster manager
	children, err := reg.(*zkreg.ZkRegistry).ZkClient().GetChildren(event.Service.Path + "/providers")
	if err != nil {}
	*/
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
		log.DefaultLogger.Warnf("not supported")
	}

	if err != nil {
		log.DefaultLogger.Errorf("process zk event fail, err: %v", err.Error())
	}
}
