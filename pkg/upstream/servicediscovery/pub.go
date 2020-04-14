package servicediscovery

import (
	"fmt"
	dubboregistry "github.com/mosn/registry/dubbo"
	dubbocommon "github.com/mosn/registry/dubbo/common"
	dubboconsts "github.com/mosn/registry/dubbo/common/constant"
	zkreg "github.com/mosn/registry/dubbo/zookeeper"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// publish a service to registry
func publish(w http.ResponseWriter, r *http.Request) {
	var req pubReq

	err := bind(r, &req)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: err.Error()})
		return
	}

	var registryPath = registryPathTpl.ExecuteString(map[string]interface{}{
		"addr": req.Registry.Addr,
	})
	registryURL, err := dubbocommon.NewURL(registryPath,
		dubbocommon.WithParams(url.Values{
			dubboconsts.GROUP_KEY:            []string{req.Service.Group},
			dubboconsts.ROLE_KEY:             []string{fmt.Sprint(dubbocommon.PROVIDER)},
			dubboconsts.REGISTRY_KEY:         []string{req.Registry.Type},
			dubboconsts.REGISTRY_TIMEOUT_KEY: []string{"5s"},
		}),
		dubbocommon.WithUsername(req.Registry.UserName),
		dubbocommon.WithPassword(req.Registry.Password),
		dubbocommon.WithLocation(req.Registry.Addr),
	)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: err.Error()})
		return
	}

	reg, err := zkreg.NewZkRegistry(&registryURL)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "publish fail, err: " + err.Error()})
		return
	}

	var (
		mosnIP, mosnPort = "127.0.0.1", "20000" // TODO, need to read from mosn config
		dubboPath        = dubboPathTpl.ExecuteString(map[string]interface{}{
			"ip":           mosnIP,
			"port":         mosnPort,
			"interface":    req.Service.Interface,
			"service_name": req.Service.Name,
		})
	)
	dubboURL, err := dubbocommon.NewURL(dubboPath,
		dubbocommon.WithParamsValue("timestamp", fmt.Sprint(time.Now().Unix())),
		dubbocommon.WithMethods(req.Service.Methods))
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "publish fail, err: "+ err.Error()})
		return
	}

	err = reg.Register(dubboURL)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "publish fail, err: " + err.Error()})
		return
	}

	response(w, resp{Errno: succ, ErrMsg: "publish success"})
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
		registryPath        = fmt.Sprintf("registry://%v", req.RegistryAddr)
	)

	registryURL, err := dubbocommon.NewURL(registryPath,
		dubbocommon.WithParamsValue(dubboconsts.ROLE_KEY, strconv.Itoa(dubbocommon.PROVIDER)),
	)
	if err != nil {
		_, _ = w.Write([]byte("unpublish fail" + err.Error()))
		return
	}

	reg := dubboregistry.BaseRegistry{
		URL: &registryURL,
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
