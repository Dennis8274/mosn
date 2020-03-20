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

package featuregate

import (
	"time"
)

/*
1. 这个文件里的内容都是和 defaultFeatureGate 相关的，完全可以把这些方法挂在同一个 struct 类型下
要不看着这些方法非常奇怪
2. defaultFeatureGate 完全可以使用 sync.Map，不需要自己锁来锁去
*/
var defaultFeatureGate *FeatureGate = NewFeatureGate()

func Enabled(key Feature) bool {
	return defaultFeatureGate.Enabled(key)
}

func Subscribe(key Feature, timeout time.Duration) (bool, error) {
	return defaultFeatureGate.Subscribe(key, timeout)
}

func Set(value string) error {
	return defaultFeatureGate.Set(value)
}

func SetFromMap(m map[string]bool) error {
	return defaultFeatureGate.SetFromMap(m)
}

func AddFeatureSpec(key Feature, spec FeatureSpec) error {
	return defaultFeatureGate.AddFeatureSpec(key, spec)
}

func SetFeatureState(key Feature, enable bool) error {
	return defaultFeatureGate.SetFeatureState(key, enable)
}

func StartInit() {
	defaultFeatureGate.StartInit()
}

func WaitInitFinsh() error {
	return defaultFeatureGate.WaitInitFinsh()
}
