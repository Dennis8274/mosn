package servicediscovery

import (
	"fmt"
	dubboregistry "github.com/mosn/registry/dubbo"
	dubbocommon "github.com/mosn/registry/dubbo/common"
	dubboconsts "github.com/mosn/registry/dubbo/common/constant"
	zkreg "github.com/mosn/registry/dubbo/zookeeper"
	"net/http"
	"strconv"
)

// publish a service to registry
func publish(w http.ResponseWriter, r *http.Request) {
	var req pubReq
	err := bind(r, &req)
	if err != nil {
		_, _ = w.Write([]byte("publish fail" + err.Error()))
		return
	}

	var (
		dubboPathInRegistry = fmt.Sprintf("dubbo://127.0.0.1:20000/%v.%v", req.Interface, req.Method)
		registryPath = fmt.Sprintf("registry://%v", req.RegistryAddr)
	)

	registryURL, err := dubbocommon.NewURL(registryPath,
		dubbocommon.WithParamsValue(dubboconsts.ROLE_KEY, strconv.Itoa(dubbocommon.PROVIDER)),
	)
	if err != nil {
		_, _ = w.Write([]byte("publish fail" + err.Error()))
		return
	}

	fmt.Println(registryURL)
	reg, err := zkreg.NewZkRegistry(&registryURL)

	if err != nil {
		_, _ = w.Write([]byte("publish fail" + err.Error()))
		return
	}
	//reg.InitBaseRegistry(&registryURL, zookeeper.)

	url, err := dubbocommon.NewURL(dubboPathInRegistry,
		//dubbocommon.WithParamsValue(dubboconsts.CLUSTER_KEY, "cluster"), // need to read from user config
		dubbocommon.WithParamsValue("serviceid", req.Interface),
		dubbocommon.WithMethods([]string{req.Method}))
	if err != nil {
		_, _ = w.Write([]byte("publish fail" + err.Error()))
		return
	}

	fmt.Println("#####", url)
	err = reg.Register(url)
	if err != nil {
		_, _ = w.Write([]byte("publish fail" + err.Error()))
		return
	}

	_, _ = w.Write([]byte("publish succeed"))
}

// unpublish user service from registry
func unpublish(w http.ResponseWriter, r *http.Request) {
	var req unpubReq
	err := bind(r, &req)
	if err != nil {
		_, _ = w.Write([]byte("unpublish fail"))
		return
	}

	var (
		dubboPathInRegistry = fmt.Sprintf("/dubbo/%v.%v/providers", req.Interface, req.Method)
		registryPath = fmt.Sprintf("registry://%v", req.RegistryAddr)
	)

	registryURL, err := dubbocommon.NewURL(registryPath,
		dubbocommon.WithParamsValue(dubboconsts.ROLE_KEY, strconv.Itoa(dubbocommon.PROVIDER)),
	)
	if err != nil {
		_, _ = w.Write([]byte("unpublish fail" + err.Error()))
		return
	}

	reg := dubboregistry.BaseRegistry{
		URL : &registryURL,
	}

	url, err := dubbocommon.NewURL(dubboPathInRegistry,
		dubbocommon.WithParamsValue(dubboconsts.CLUSTER_KEY, "cluster"), // need to read from user config
		dubbocommon.WithParamsValue("serviceid", req.Interface),
		dubbocommon.WithMethods([]string{req.Method}))
	if err != nil {
		_, _ = w.Write([]byte("unpublish fail" + err.Error()))
		return
	}

	err = reg.Register(url)
	if err != nil {
		_, _ = w.Write([]byte("publish fail" + err.Error()))
		return
	}

	_, _ = w.Write([]byte("unpublish succeed"))
}
