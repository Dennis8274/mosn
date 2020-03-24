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

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/urfave/cli"
	"mosn.io/mosn/pkg/admin/store"
	"mosn.io/mosn/pkg/configmanager"
	"mosn.io/mosn/pkg/featuregate"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/metrics"
	"mosn.io/mosn/pkg/mosn"
	"mosn.io/mosn/pkg/types"
)

var (
	cmdStart = cli.Command{
		Name:  "start",
		Usage: "start mosn proxy",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "config, c",
				Usage:  "Load configuration from `FILE`",
				EnvVar: "MOSN_CONFIG",
				Value:  "configs/mosn_config.json",
			}, cli.StringFlag{
				Name:   "service-cluster, s",
				Usage:  "sidecar service cluster",
				EnvVar: "SERVICE_CLUSTER",
			}, cli.StringFlag{
				Name:   "service-node, n",
				Usage:  "sidecar service node",
				EnvVar: "SERVICE_NODE",
			}, cli.StringSliceFlag{
				Name:   "service-meta, sm",
				Usage:  "sidecar service metadata",
				EnvVar: "SERVICE_META",
			}, cli.StringFlag{
				Name:   "feature-gates, f",
				Usage:  "config feature gates",
				EnvVar: "FEATURE_GATES",
			},
		},
		Action: func(c *cli.Context) error {
			// 各种配置的路径
			// 这么多变量其实不是很优雅
			// 配置繁琐了以后，应该是 structOpt 之类的更好用？
			var (
				configPath     = c.String("config") // 这里的 string 和上面是对应的，感觉还是挺容易拼错的
				serviceCluster = c.String("service-cluster")
				serviceNode    = c.String("service-node")
				serviceMeta    = c.StringSlice("service-meta")
			)

			// 这个 config 文件和 v2ray 的类似
			// 因此也和 v2ray 有一样的配置复杂的问题
			// 可以改成问答式的配置流程
			conf := configmanager.Load(configPath)

			// set feature gates
			// feature gate 是动态开关闭功能的一个特性
			// feature1=true,feature2=false
			// 设计的比较奇怪
			err := featuregate.Set(c.String("feature-gates"))
			if err != nil {
				log.StartLogger.Infof("[mosn] [start] parse feature-gates flag fail : %+v", err)
				os.Exit(1)
			}

			// start pprof
			// 配置了 pprof 的话会开启
			if conf.Debug.StartDebug {
				port := 9090 //default use 9090
				if conf.Debug.Port != 0 {
					port = conf.Debug.Port
				}
				addr := fmt.Sprintf("0.0.0.0:%d", port)
				s := &http.Server{Addr: addr, Handler: nil}
				store.AddService(s, "pprof", nil, nil)
			}

			// set mosn metrics flush
			metrics.FlushMosnMetrics = true
			// set version and go version
			// 这个 version 感觉也应该写在一个专门的文件而不是代码里
			metrics.SetVersion(Version)
			metrics.SetGoVersion(runtime.Version())

			// TODO, what is xds
			// https://zhaohuabing.com/post/2018-09-25-istio-traffic-management-impl-intro/#xds
			types.InitXdsFlags(serviceCluster, serviceNode, serviceMeta)

			// 实际的启动函数
			mosn.Start(conf)
			return nil
		},
	}

	cmdStop = cli.Command{
		Name:  "stop",
		Usage: "stop mosn proxy",
		Action: func(c *cli.Context) error {
			return nil
		},
	}

	cmdReload = cli.Command{
		Name:  "reload",
		Usage: "reconfiguration",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
)
