// Copyright (C) 2015-Present Pivotal Software, Inc. All rights reserved.

// This program and the accompanying materials are made available under
// the terms of the under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package brokerapi

type EmptyResponse struct{}

type ErrorResponse struct {
	Error       string `json:"error,omitempty"`
	Description string `json:"description"`
}

type CatalogResponse struct {
	Services []Service `json:"services"`
}

type ProvisioningResponse struct {
	DashboardURL  string `json:"dashboard_url,omitempty"`
	OperationData string `json:"operation,omitempty"`
}

type UpdateResponse struct {
	OperationData string `json:"operation,omitempty"`
}

type DeprovisionResponse struct {
	OperationData string `json:"operation,omitempty"`
}

type LastOperationResponse struct {
	State       LastOperationState `json:"state"`
	Description string             `json:"description,omitempty"`
}

type ExperimentalVolumeMountBindingResponse struct {
	Credentials     interface{}               `json:"credentials"`
	SyslogDrainURL  string                    `json:"syslog_drain_url,omitempty"`
	RouteServiceURL string                    `json:"route_service_url,omitempty"`
	VolumeMounts    []ExperimentalVolumeMount `json:"volume_mounts,omitempty"`
}

type ExperimentalVolumeMount struct {
	ContainerPath string                         `json:"container_path"`
	Mode          string                         `json:"mode"`
	Private       ExperimentalVolumeMountPrivate `json:"private"`
}

type ExperimentalVolumeMountPrivate struct {
	Driver  string `json:"driver"`
	GroupID string `json:"group_id"`
	Config  string `json:"config"`
}
