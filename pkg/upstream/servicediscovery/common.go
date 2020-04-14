package servicediscovery

import (
	"encoding/json"
	"github.com/valyala/fasttemplate"
	"net/http"
)

var (
	dubboPathTpl    = fasttemplate.New("dubbo://{{ip}}:{{port}}/{{interface}}.{{service_name}}", "{{", "}}")
	registryPathTpl = fasttemplate.New("registry://{{addr}}", "{{", "}}")
)

const (
	succ = iota
	fail
)

func response(w http.ResponseWriter, respBody interface{}) {
	bodyBytes, err := json.Marshal(respBody)
	if err != nil {
		_, _ = w.Write([]byte("response marshal failed, err: " + err.Error()))
	}

	_, _ = w.Write(bodyBytes)
}
