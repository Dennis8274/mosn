/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package servicediscovery

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-chi/chi"
	dubboregistry "github.com/mosn/registry/dubbo"
	dubbocommon "github.com/mosn/registry/dubbo/common"
	dubboconsts "github.com/mosn/registry/dubbo/common/constant"
	"github.com/mosn/registry/dubbo/remoting"
	zkreg "github.com/mosn/registry/dubbo/zookeeper"
	"mosn.io/mosn/pkg/log"
	"mosn.io/pkg/utils"
	"net/http"
	"strconv"
)

// init the http api for application when application bootstrap
// for sub/pub
func Init() {
	r := chi.NewRouter()
	r.Post("/sub", subscribe)
	r.Post("/unsub", unsubscribe)

	r.Post("/pub", publish)
	r.Post("/unpub", unpublish)

	// FIXME make port configurable
	utils.GoWithRecover(func(){
		if err := http.ListenAndServe(":22222", r); err != nil {
			log.DefaultLogger.Infof("auto write config when updated")
		}
	}, nil)
}

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
		registryPath = fmt.Sprintf("registry://%v", req.RegistryAddr)
	)

	registryURL, err := dubbocommon.NewURL(registryPath,
		dubbocommon.WithParamsValue(dubboconsts.ROLE_KEY, strconv.Itoa(dubbocommon.PROVIDER)),
	)
	if err != nil {
		_, _ = w.Write([]byte("sub fail" + err.Error()))
		return
	}

	fmt.Println("#####",registryURL)
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

	var notifyListener = listener{}
	go reg.Subscribe(&url, notifyListener)

	_, _ = w.Write([]byte("unsubscribe succeed"))
}

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


// user will always publish a provider
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

// bind the struct content from http.Request body/uri
func bind(r *http.Request, data interface{}) error {
	b := binding.Default(r.Method, r.Header.Get("Content-Type"))
	return b.Bind(r, data)
}
