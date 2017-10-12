// Copyright the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
////////////////////////////////////////////////////////////////////////////////
//

package models

type Service struct {
	ID              string                  `json:"id" validate:"nonzero"`
	Name            string                  `json:"name" validate:"nonzero"`
	Description     string                  `json:"description"`
	Bindable        bool                    `json:"bindable"`
	Tags            []string                `json:"tags,omitempty"`
	PlanUpdatable   bool                    `json:"plan_updateable"`
	Plans           []ServicePlan           `json:"plans"`
	Requires        []RequiredPermission    `json:"requires,omitempty"`
	Metadata        *ServiceMetadata        `json:"metadata,omitempty"`
	DashboardClient *ServiceDashboardClient `json:"dashboard_client,omitempty"`
}

type ServiceDashboardClient struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	RedirectURI string `json:"redirect_uri"`
}

type ServicePlan struct {
	ID                string               `json:"id" validate:"nonzero"`
	Name              string               `json:"name" validate:"nonzero"`
	Description       string               `json:"description"`
	Free              *bool                `json:"free,omitempty"`
	Metadata          *ServicePlanMetadata `json:"metadata,omitempty"`
	ServiceProperties map[string]string    `json:"service_properties"`
}

type ServicePlanCandidate struct {
	Guid              string `json:"guid"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	DisplayName       string `json:"display_name"`
	ServiceProperties string `json:"service_properties"`
}

type ServicePlanMetadata struct {
	DisplayName string            `json:"displayName,omitempty"`
	Bullets     []string          `json:"bullets,omitempty"`
	Costs       []ServicePlanCost `json:"costs,omitempty"`
}

type ServicePlanCost struct {
	Amount map[string]float64 `json:"amount"`
	Unit   string             `json:"unit"`
}

type ServiceMetadata struct {
	DisplayName         string `json:"displayName,omitempty"`
	ImageUrl            string `json:"imageUrl,omitempty"`
	LongDescription     string `json:"longDescription,omitempty"`
	ProviderDisplayName string `json:"providerDisplayName,omitempty"`
	DocumentationUrl    string `json:"documentationUrl,omitempty"`
	SupportUrl          string `json:"supportUrl,omitempty"`
}

type RequiredPermission string

const (
	PermissionRouteForwarding = RequiredPermission("route_forwarding")
	PermissionSyslogDrain     = RequiredPermission("syslog_drain")
)
