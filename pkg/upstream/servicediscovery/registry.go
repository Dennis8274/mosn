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
	"context"
)

type registerType int

const (
	typeConsumer registerType = 0
	typeProvider              = 1
)

func init() {}

// Subscribe for Dependent services from service registry
func Subscribe(ctx context.Context, serviceName string) {}

func doUnpublish(ctx context.Context, t registerType) error {
	switch t {
	case typeConsumer:
		pubConsumer()
	case typeProvider:
		pubProvider()
	default:
		// TODO log error
	}
	return nil
}

func doPublish(ctx context.Context, t registerType) error {
	switch t {
	case typeConsumer:
		pubConsumer()
	case typeProvider:
		pubProvider()
	default:
		// TODO log error
	}
	return nil
}

func pubConsumer() {
}

func pubProvider() {}

// ChangeNotify accepts the registry notification
// call the cluster manager to do the job
func ChangeNotify(ctx context.Context) {}
