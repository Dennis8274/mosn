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

package configmanager

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"

	v2 "mosn.io/mosn/pkg/config/v2"
)

var (
	configPath     string
	configLock     sync.Mutex
	config         v2.MOSNConfig
	configLoadFunc ConfigLoadFunc = DefaultConfigLoad
)

// protetced configPath, read only
func GetConfigPath() string {
	return configPath
}

// ConfigLoadFunc parse a input(usually file path) into a mosn config
type ConfigLoadFunc func(path string) *v2.MOSNConfig

// RegisterConfigLoadFunc can replace a new config load function instead of default
func RegisterConfigLoadFunc(f ConfigLoadFunc) {
	configLoadFunc = f
}

// 默认的 config 加载方法
// 可以通过上面 RegisterConfigLoadFunc 来替换
// 然而目前却没有任何替换方法的使用
func DefaultConfigLoad(path string) *v2.MOSNConfig {
	log.Println("load config from : ", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("[config] [default load] load config failed, ", err)
	}
	cfg := &v2.MOSNConfig{}
	// translate to lower case
	err = json.Unmarshal(content, cfg)
	if err != nil {
		// 如果失败，fatalln 里直接 os.Exit 了
		// 所以这里不返回 err 也是可以的
		log.Fatalln("[config] [default load] json unmarshal config failed, ", err)
	}
	return cfg

}

// Load config file and parse
// 这里难道就没有加载错误的可能性？
func Load(path string) *v2.MOSNConfig {
	configPath, _ = filepath.Abs(path)
	if cfg := configLoadFunc(path); cfg != nil {
		config = *cfg
	}
	return &config
}
