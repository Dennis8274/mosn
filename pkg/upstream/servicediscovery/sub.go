package servicediscovery

import (
	"fmt"
	dubbocommon "github.com/mosn/registry/dubbo/common"
	dubboconsts "github.com/mosn/registry/dubbo/common/constant"
	zkreg "github.com/mosn/registry/dubbo/zookeeper"
	"net/http"
	"strconv"
)

// unsubscribe a service from registry
func unsubscribe(w http.ResponseWriter, r *http.Request) {
	// TODO fetch this cluster config from registry
	// TODO convert it to v2.Cluster
	/*err := clusterAdapter.GetClusterMngAdapterInstance().AddOrUpdatePrimaryCluster(v2.Cluster{})
	if err != nil {
		_, _ = w.Write([]byte("json failed"))
	}
	*/
	// todo
	// 1. build url, then
	// 2. call dubbo-go service subscribe
	// 3. register current service as dependent service consumers
	var req subReq
	err := bind(r, &req)
	if err != nil {
		_, _ = w.Write([]byte("subscribe fail"))
		return
	}

	_, _ = w.Write([]byte("subscribe succeed"))
}

// subscribe a service from registry
func subscribe(w http.ResponseWriter, r *http.Request) {
	var req unsubReq
	err := bind(r, &req)
	if err != nil {
		_, _ = w.Write([]byte("subscribe fail"))
		return
	}

	// 在收到 notify 之后，调用 TriggerClusterAddOrUpdate 等等接口
	// err := clusterAdapter.GetClusterMngAdapterInstance().AddOrUpdatePrimaryCluster(v2.Cluster{})
	/*
		err = clusterAdapter.GetClusterMngAdapterInstance().TriggerClusterDel(req.Services...)
		if err != nil {
			_, _ = w.Write([]byte("json failed"))
		}
	*/
	var (
		dubboPathInRegistry = fmt.Sprintf("dubbo://127.0.0.1:20000/%v.%v", req.Interface, req.Method)
		registryPath        = fmt.Sprintf("registry://%v", req.RegistryAddr)
	)

	registryURL, err := dubbocommon.NewURL(registryPath,
		dubbocommon.WithParamsValue(dubboconsts.ROLE_KEY, strconv.Itoa(dubbocommon.PROVIDER)),
	)
	if err != nil {
		_, _ = w.Write([]byte("sub fail" + err.Error()))
		return
	}

	fmt.Println("#####", registryURL)
	reg, err := zkreg.NewZkRegistry(&registryURL)

	if err != nil {
		_, _ = w.Write([]byte("sub fail" + err.Error()))
		return
	}

	url, err := dubbocommon.NewURL(dubboPathInRegistry,
		//dubbocommon.WithParamsValue(dubboconsts.CLUSTER_KEY, "cluster"), // need to read from user config
		dubbocommon.WithParamsValue("serviceid", req.Interface),
		dubbocommon.WithMethods([]string{req.Method}))
	if err != nil {
		_, _ = w.Write([]byte("sub fail" + err.Error()))
		return
	}

	/*
		err = reg.Register(url)
		if err != nil {
			_, _ = w.Write([]byte("sub fail" + err.Error()))
			return
		}
	*/

	go reg.Subscribe(&url, listener{})

	_, _ = w.Write([]byte("unsubscribe succeed"))
}
