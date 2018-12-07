// Copyright 2018 the Service Broker Project Authors.
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

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pivotal-cf/brokerapi"
	"github.com/spf13/viper"
)

const (
	// ClientsBrokerApiVersion is the minimum supported version of the client.
	// Note: This may need to be changed in the future as we use newer versions
	// of the OSB API, but should be kept near the lower end of the systems we
	// expect to be compatible with to ensure any reverse-compatibility measures
	// put in place work.
	ClientsBrokerApiVersion = "2.13"
)

// NewClientFromEnv creates a new client from the client configuration properties.
func NewClientFromEnv() (*Client, error) {
	user := viper.GetString("api.user")
	pass := viper.GetString("api.password")
	port := viper.GetInt("api.port")

	viper.SetDefault("api.hostname", "localhost")
	host := viper.GetString("api.hostname")
	return New(user, pass, host, port)
}

// New creates a new OSB Client connected to the given resource.
func New(username, password, hostname string, port int) (*Client, error) {
	base := fmt.Sprintf("http://%s:%s@%s:%d/v2/", username, password, hostname, port)
	baseUrl, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	return &Client{BaseUrl: baseUrl}, nil
}

type Client struct {
	BaseUrl *url.URL
}

// Catalog fetches the service catalog
func (client *Client) Catalog() *BrokerResponse {
	return client.makeRequest(http.MethodGet, "catalog", nil)
}

// Provision creates a new service with the given instanceId, of type serviceId,
// from the plan planId, with additional details provisioningDetails
func (client *Client) Provision(instanceId, serviceId, planId string, provisioningDetails json.RawMessage) *BrokerResponse {
	url := fmt.Sprintf("service_instances/%s?accepts_incomplete=true", instanceId)

	return client.makeRequest(http.MethodPut, url, brokerapi.ProvisionDetails{
		ServiceID:     serviceId,
		PlanID:        planId,
		RawParameters: provisioningDetails,
	})
}

// Deprovision destroys a service instance of type instanceId
func (client *Client) Deprovision(instanceId, serviceId, planId string) *BrokerResponse {
	url := fmt.Sprintf("service_instances/%s?accepts_incomplete=true&service_id=%s&plan_id=%s", instanceId, serviceId, planId)

	return client.makeRequest(http.MethodDelete, url, nil)
}

// Bind creates an account identified by bindingId and gives it access to instanceId
func (client *Client) Bind(instanceId, bindingId, serviceId, planId string, parameters json.RawMessage) *BrokerResponse {
	url := fmt.Sprintf("service_instances/%s/service_bindings/%s", instanceId, bindingId)

	return client.makeRequest(http.MethodPut, url, brokerapi.BindDetails{
		ServiceID:     serviceId,
		PlanID:        planId,
		RawParameters: parameters,
	})
}

// Unbind destroys an account identified by bindingId
func (client *Client) Unbind(instanceId, bindingId, serviceId, planId string) *BrokerResponse {
	url := fmt.Sprintf("service_instances/%s/service_bindings/%s?service_id=%s&plan_id=%s", instanceId, bindingId, serviceId, planId)

	return client.makeRequest(http.MethodDelete, url, nil)
}

// Update sends a patch request to change the plan
func (client *Client) Update(instanceId, serviceId, planId string, parameters json.RawMessage) *BrokerResponse {
	url := fmt.Sprintf("service_instances/%s", instanceId)

	return client.makeRequest(http.MethodPatch, url, brokerapi.UpdateDetails{
		ServiceID:     serviceId,
		PlanID:        planId,
		RawParameters: parameters,
	})
}

// LastOperation queries the status of a long-running job on the server
func (client *Client) LastOperation(instanceId string) *BrokerResponse {
	url := fmt.Sprintf("service_instances/%s/last_operation", instanceId)

	return client.makeRequest(http.MethodGet, url, nil)
}

func (client *Client) makeRequest(method, path string, body interface{}) *BrokerResponse {
	br := BrokerResponse{}

	req, err := client.newRequest(method, path, body)
	br.UpdateRequest(req)
	br.UpdateError(err)
	if br.InError() {
		return &br
	}

	resp, err := http.DefaultClient.Do(req)

	br.UpdateResponse(resp)
	br.UpdateError(err)

	return &br
}

func (client *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	url, err := client.BaseUrl.Parse(path)
	if err != nil {
		return nil, err
	}

	var buffer io.ReadWriter
	if body != nil {
		buffer = new(bytes.Buffer)
		enc := json.NewEncoder(buffer)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequest(method, url.String(), buffer)
	if err != nil {
		return nil, err
	}

	request.Header.Set("X-Broker-Api-Version", ClientsBrokerApiVersion)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	return request, nil
}
