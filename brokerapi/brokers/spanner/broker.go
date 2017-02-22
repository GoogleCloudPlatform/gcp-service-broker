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

package spanner

import (
	googlespanner "cloud.google.com/go/spanner/admin/instance/apiv1"
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"fmt"
	"gcp-service-broker/brokerapi/brokers/broker_base"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/brokerapi/brokers/name_generator"
	"gcp-service-broker/db_service"
	"net/http"
)

type SpannerBroker struct {
	Client         *http.Client
	ProjectId      string
	Logger         lager.Logger
	AccountManager models.AccountManager

	broker_base.BrokerBase
}

type InstanceInformation struct {
	InstanceId string `json:"instance_id"`
}

// Creates a new Spanner Instance identified by the name provided in details.RawParameters.name and
// an optional region (defaults to us-central1)
func (s *SpannerBroker) Provision(instanceId string, details models.ProvisionDetails, plan models.PlanDetails) (models.ServiceInstanceDetails, error) {
	var err error
	var params map[string]string

	if len(details.RawParameters) == 0 {
		params = map[string]string{}
	} else if err = json.Unmarshal(details.RawParameters, &params); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
	}

	// Ensure there is a name for this instance
	if _, ok := params["name"]; !ok {
		params["name"] = name_generator.Basic.InstanceNameWithSeparator("-")
	}

	// get plan parameters
	var planDetails map[string]string
	if err = json.Unmarshal([]byte(plan.Features), &planDetails); err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error unmarshalling plan features: %s", err)
	}

	// set up client

	// set up instance config

	// create instance


	// save off instance information
	ii := InstanceInformation{
		InstanceId: params["name"],
	}

	otherDetails, err := json.Marshal(ii)
	if err != nil {
		return models.ServiceInstanceDetails{}, fmt.Errorf("Error marshalling other details: %s", err)
	}

	i := models.ServiceInstanceDetails{
		Name:         params["name"],
		Url:          "",
		Location:     "",
		OtherDetails: string(otherDetails),
	}

	return i, nil
}

// gets the last operation for this instance and polls the status of it
func (s *SpannerBroker) PollInstance(instanceId string) (bool, error) {
	return false, nil
}

// deletes the instance associated with the given instanceID string
func (s *SpannerBroker) Deprovision(instanceID string, details models.DeprovisionDetails) error {
	var err error
	// set up client

	instance := models.ServiceInstanceDetails{}
	if err = db_service.DbConnection.Where("ID = ?", instanceID).First(&instance).Error; err != nil {
		return models.ErrInstanceDoesNotExist
	}

	// delete instance

	return nil
}



// Indicates that Spanner uses asynchronous provisioning
func (s *SpannerBroker) Async() bool {
	return true
}

type SpannerDynamicPlan struct {
	Guid        string `json:"guid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NumNodes    string `json:"num_nodes"`
	DisplayName string `json:"display_name"`
	ServiceId   string `json:"service"`
}

func MapPlan(details map[string]string) map[string]string {

	features := map[string]string{
		"num_nodes":    details["num_nodes"],
	}
	return features
}
