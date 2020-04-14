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
	"github.com/go-chi/chi"
	"github.com/mosn/binding"
	dubbologger "github.com/mosn/registry/dubbo/common/logger"
	"mosn.io/mosn/pkg/log"
	"mosn.io/pkg/utils"
	"net/http"
)

// init the http api for application when application bootstrap
// for sub/pub
func Init( /*port string, dubboLogPath string*/ ) {
	r := chi.NewRouter()
	r.Post("/sub", subscribe)
	r.Post("/unsub", unsubscribe)
	r.Post("/pub", publish)
	r.Post("/unpub", unpublish)

	_ = dubbologger.InitLog("./dubbogo.log")

	// FIXME make port configurable
	utils.GoWithRecover(func() {
		if err := http.ListenAndServe(":22222", r); err != nil {
			log.DefaultLogger.Infof("auto write config when updated")
		}
	}, nil)
}

// bind the struct content from http.Request body/uri
func bind(r *http.Request, data interface{}) error {
	b := binding.Default(r.Method, r.Header.Get("Content-Type"))
	return b.Bind(r, data)
}
