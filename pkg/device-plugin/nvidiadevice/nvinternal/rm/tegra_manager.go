/*
 * SPDX-License-Identifier: Apache-2.0
 *
 * The HAMi Contributors require contributions made to
 * this file be licensed under the Apache-2.0 license or a
 * compatible open source license.
 */

/*
 * Licensed to NVIDIA CORPORATION under one or more contributor
 * license agreements. See the NOTICE file distributed with
 * this work for additional information regarding copyright
 * ownership. NVIDIA CORPORATION licenses this file to you under
 * the Apache License, Version 2.0 (the "License"); you may
 * not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * Modifications Copyright The HAMi Authors. See
 * GitHub history for details.
 */

package rm

import (
	"fmt"

	"github.com/Project-HAMi/HAMi/pkg/device/nvidia"
)

type tegraResourceManager struct {
	resourceManager
}

var _ ResourceManager = (*tegraResourceManager)(nil)

// NewTegraResourceManagers returns a set of ResourceManagers for tegra resources
func NewTegraResourceManagers(config *nvidia.DeviceConfig) ([]ResourceManager, error) {
	deviceMap, err := buildTegraDeviceMap(config)
	if err != nil {
		return nil, fmt.Errorf("error building Tegra device map: %v", err)
	}

	deviceMap, err = updateDeviceMapWithReplicas(config, deviceMap)
	if err != nil {
		return nil, fmt.Errorf("error updating device map with replicas from config.sharing.timeSlicing.resources: %v", err)
	}

	var rms []ResourceManager
	for resourceName, devices := range deviceMap {
		if len(devices) == 0 {
			continue
		}
		r := &tegraResourceManager{
			resourceManager: resourceManager{
				config:   config,
				resource: resourceName,
				devices:  devices,
			},
		}
		if len(devices) != 0 {
			rms = append(rms, r)
		}
	}

	return rms, nil
}

// GetPreferredAllocation returns a standard allocation for the Tegra resource manager.
func (r *tegraResourceManager) GetPreferredAllocation(available, required []string, size int) ([]string, error) {
	return r.distributedAlloc(available, required, size)
}

// GetDevicePaths returns an empty slice for the tegraResourceManager
func (r *tegraResourceManager) GetDevicePaths(ids []string) []string {
	return nil
}

// CheckHealth is disabled for the tegraResourceManager
func (r *tegraResourceManager) CheckHealth(stop <-chan any, unhealthy chan<- *Device, disableNVML <-chan bool, ackDisableHealthChecks chan<- bool) error {
	return nil
}
