package servicediscovery

import (
	"encoding/json"
	"fmt"
	dubbocommon "github.com/mosn/registry/dubbo/common"
	zkreg "github.com/mosn/registry/dubbo/zookeeper"
	"github.com/valyala/fasttemplate"
	"math/rand"
	"net/http"
	"sync"
)

var (
	dubboPathTpl    = fasttemplate.New("dubbo://{{ip}}:{{port}}/{{interface}}.{{service_name}}", "{{", "}}")
	registryPathTpl = fasttemplate.New("registry://{{addr}}", "{{", "}}")
	dubboRouterConfigName = "dubbo" // keep the same with the router config name in mosn_config.json
)

var (
	mosnIP, mosnPort = "127.0.0.1", fmt.Sprint(rand.Int63n(30000)+1) // TODO, need to read from mosn config
)

const (
	succ = iota
	fail
)

// /com.test.cch.UserService --> zk client
var registryClientCache = sync.Map{}

func getRegistry(registryCacheKey string, role int, registryURL dubbocommon.URL) (*zkreg.ZkRegistry, error) {
	// do not cache provider registry, or it may collide with the consumer registry
	if role == dubbocommon.PROVIDER {
		return zkreg.NewZkRegistry(&registryURL)
	}

	regInterface, ok := registryClientCache.Load(registryCacheKey)

	var (
		reg *zkreg.ZkRegistry
		err error
	)

	if !ok {
		// init registry
		reg, err = zkreg.NewZkRegistry(&registryURL)
		// store registry object to global cache
		if err == nil {
			registryClientCache.Store(registryCacheKey, reg)
		}
	} else {
		reg = regInterface.(*zkreg.ZkRegistry)
	}

	return reg, err
}

func response(w http.ResponseWriter, respBody interface{}) {
	bodyBytes, err := json.Marshal(respBody)
	if err != nil {
		_, _ = w.Write([]byte("response marshal failed, err: " + err.Error()))
	}

	_, _ = w.Write(bodyBytes)
}


