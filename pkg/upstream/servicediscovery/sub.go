package servicediscovery

import (
	"fmt"
	dubbocommon "github.com/mosn/registry/dubbo/common"
	dubboconsts "github.com/mosn/registry/dubbo/common/constant"
	zkreg "github.com/mosn/registry/dubbo/zookeeper"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	routerAdapter "mosn.io/mosn/pkg/router"
	clusterAdapter "mosn.io/mosn/pkg/upstream/cluster"
	"net/http"
	"net/url"
	"time"
)

// subscribe a service from registry
func subscribe(w http.ResponseWriter, r *http.Request) {
	var req subReq
	err := bind(r, &req)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "subscribe fail, err: " + err.Error()})
		return
	}

	var registryPath = registryPathTpl.ExecuteString(map[string]interface{}{
		"addr": req.Registry.Addr,
	})
	registryURL, err := dubbocommon.NewURL(registryPath,
		dubbocommon.WithParams(url.Values{
			dubboconsts.GROUP_KEY:            []string{req.Service.Group},
			dubboconsts.ROLE_KEY:             []string{fmt.Sprint(dubbocommon.CONSUMER)},
			dubboconsts.REGISTRY_KEY:         []string{req.Registry.Type},
			dubboconsts.REGISTRY_TIMEOUT_KEY: []string{"5s"},
		}),
	)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "subscribe fail, err: " + err.Error()})
		return
	}

	servicePath := req.Service.Interface + "." + req.Service.Name // com.mosn.test.UserService
	reg, err := getRegistry(servicePath, registryURL)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "subscribe fail, err: " + err.Error()})
		return
	}

	var (
		dubboPath = dubboPathTpl.ExecuteString(map[string]interface{}{
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
		response(w, resp{Errno: fail, ErrMsg: "subscribe fail, err: " + err.Error()})
		return
	}

	// register consumer to registry
	err = reg.Register(dubboURL)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "subscribe fail, err: " + err.Error()})
		return
	}

	// listen to provider change events
	go reg.Subscribe(&dubboURL, listener{})

	err = routerAdapter.GetRoutersMangerInstance().AddRoute(dubboRouterConfigName, "*", &v2.Router{
		RouterConfig: v2.RouterConfig{
			Match: v2.RouterMatch{
				Headers: []v2.HeaderMatcher{
					{
						Name:  "service", // use the xprotocol header field "service"
						Value:  servicePath,
					},
				},
			},
			Route: v2.RouteAction{
				RouterActionConfig: v2.RouterActionConfig{
					ClusterName:  servicePath,
				},
			},
		},
	})
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "subscribe fail, err: " + err.Error()})
		return
	}

	response(w, resp{Errno: succ, ErrMsg: "subscribe success"})
}

// unsubscribe a service from registry
func unsubscribe(w http.ResponseWriter, r *http.Request) {
	var req unsubReq
	err := bind(r, &req)
	if err != nil {
		response(w, resp{Errno: fail, ErrMsg: "unsubscribe fail, err: " + err.Error()})
		return
	}

	// remove all registry cache
	registryClientCache.Range(func(k, v interface{}) bool {
		// remove cluster -> host info
		err := clusterAdapter.GetClusterMngAdapterInstance().RemovePrimaryCluster(k.(string))
		if err != nil {
			log.DefaultLogger.Errorf("remove cluster failed, err: %v", err.Error())
		}

		// destroy registry, stop subscribe
		r := v.(*zkreg.ZkRegistry)
		r.Destroy()
		return true
	})

	response(w, resp{Errno: succ, ErrMsg: "unsubscribe success"})
}

/*
func (c *ReferenceConfig) getUrlMap() url.Values {
	urlMap := url.Values{}
	//first set user params
	for k, v := range c.Params {
		urlMap.Set(k, v)
	}
	urlMap.Set(constant.INTERFACE_KEY, c.InterfaceName)
	urlMap.Set(constant.TIMESTAMP_KEY, strconv.FormatInt(time.Now().Unix(), 10))
	urlMap.Set(constant.CLUSTER_KEY, c.Cluster)
	urlMap.Set(constant.LOADBALANCE_KEY, c.Loadbalance)
	urlMap.Set(constant.RETRIES_KEY, c.Retries)
	urlMap.Set(constant.GROUP_KEY, c.Group)
	urlMap.Set(constant.VERSION_KEY, c.Version)
	urlMap.Set(constant.GENERIC_KEY, strconv.FormatBool(c.Generic))
	urlMap.Set(constant.ROLE_KEY, strconv.Itoa(common.CONSUMER))

	urlMap.Set(constant.RELEASE_KEY, "dubbo-golang-"+constant.Version)
	urlMap.Set(constant.SIDE_KEY, (common.RoleType(common.CONSUMER)).Role())

	if len(c.RequestTimeout) != 0 {
		urlMap.Set(constant.TIMEOUT_KEY, c.RequestTimeout)
	}
	//getty invoke async or sync
	urlMap.Set(constant.ASYNC_KEY, strconv.FormatBool(c.Async))
	urlMap.Set(constant.STICKY_KEY, strconv.FormatBool(c.Sticky))

	//application info
	urlMap.Set(constant.APPLICATION_KEY, consumerConfig.ApplicationConfig.Name)
	urlMap.Set(constant.ORGANIZATION_KEY, consumerConfig.ApplicationConfig.Organization)
	urlMap.Set(constant.NAME_KEY, consumerConfig.ApplicationConfig.Name)
	urlMap.Set(constant.MODULE_KEY, consumerConfig.ApplicationConfig.Module)
	urlMap.Set(constant.APP_VERSION_KEY, consumerConfig.ApplicationConfig.Version)
	urlMap.Set(constant.OWNER_KEY, consumerConfig.ApplicationConfig.Owner)
	urlMap.Set(constant.ENVIRONMENT_KEY, consumerConfig.ApplicationConfig.Environment)

	//filter
	var defaultReferenceFilter = constant.DEFAULT_REFERENCE_FILTERS
	if c.Generic {
		defaultReferenceFilter = constant.GENERIC_REFERENCE_FILTERS + "," + defaultReferenceFilter
	}
	urlMap.Set(constant.REFERENCE_FILTER_KEY, mergeValue(consumerConfig.Filter, c.Filter, defaultReferenceFilter))

	for _, v := range c.Methods {
		urlMap.Set("methods."+v.Name+"."+constant.LOADBALANCE_KEY, v.Loadbalance)
		urlMap.Set("methods."+v.Name+"."+constant.RETRIES_KEY, v.Retries)
		urlMap.Set("methods."+v.Name+"."+constant.STICKY_KEY, strconv.FormatBool(v.Sticky))
		if len(v.RequestTimeout) != 0 {
			urlMap.Set("methods."+v.Name+"."+constant.TIMEOUT_KEY, v.RequestTimeout)
		}
	}

	return urlMap
}
*/
