package servicediscovery

import (
	"fmt"
	dubboregistry "github.com/mosn/registry/dubbo"
	"github.com/mosn/registry/dubbo/remoting"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	clusterAdapter "mosn.io/mosn/pkg/upstream/cluster"
)

// listener listens for registry subscription data change
type listener struct{}
var li = &listener{}

func (l *listener) Notify(event *dubboregistry.ServiceEvent) {
	var (
		err         error
		clusterName = event.Service.Service() // FIXME
		addr        = event.Service.Ip + ":" + event.Service.Port
	)

	fmt.Println("$$$$$", event, event.Service, event.Action)

	/*
	FIXME, if there is need to load full provider list, then need to finish this logic, and overwrite the hosts in cluster manager
	children, err := reg.(*zkreg.ZkRegistry).ZkClient().GetChildren(event.Service.Path + "/providers")
	if err != nil {}
	*/
	switch event.Action {
	case remoting.EventTypeAdd:
		err = clusterAdapter.GetClusterMngAdapterInstance().TriggerHostAppend(clusterName, []v2.Host{
			{
				HostConfig: v2.HostConfig{
					Address: addr,
				},
			},
		})
		if err != nil {
			// call the cluster manager to add the host
			err = clusterAdapter.GetClusterMngAdapterInstance().TriggerClusterAndHostsAddOrUpdate(
				v2.Cluster{
					Name : clusterName,
					ClusterType: v2.SIMPLE_CLUSTER,
					LbType: v2.LB_RANDOM,
					MaxRequestPerConn: 1024,
					ConnBufferLimitBytes: 32768,
				},
				[]v2.Host{
					{
						HostConfig: v2.HostConfig{
							Address: addr,
						},
					},
				},
			)
		}
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
