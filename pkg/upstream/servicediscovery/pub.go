package servicediscovery

import (
	"fmt"
	dubbocommon "github.com/mosn/registry/dubbo/common"
	dubboconsts "github.com/mosn/registry/dubbo/common/constant"
	"net/http"
	"net/url"
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
		response(w, resp{Errno: fail, ErrMsg: "publish fail, err: " + err.Error()})
		return
	}

	// find registry from cache
	registryCacheKey := req.Service.Interface + "." + req.Service.Name
	reg ,err := getRegistry(registryCacheKey,dubbocommon.PROVIDER, registryURL)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "publish fail, err: " + err.Error()})
		return
	}

	var (
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

	// register service provider
	err = reg.Register(dubboURL)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "publish fail, err: " + err.Error()})
		return
	}

	response(w, resp{Errno: succ, ErrMsg: "publish success"})
}

// unpublish user service from registry
// FIXME, not supported
func unpublish(w http.ResponseWriter, r *http.Request) {
	var req unpubReq
	err := bind(r, &req)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "unpublish fail, err: " + err.Error()})
		return
	}

	response(w, resp{Errno: succ, ErrMsg: "unpublish success"})
}
