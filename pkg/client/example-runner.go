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
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/pivotal-cf/brokerapi"
)

// RunExamples runs all examples against the service broker pointed to by client.
func RunExamples(client *Client) error {
	return RunExamplesForService(client, "")
}

// RunExamplesForService runs all the exmaples for a given service name against
// the service broker pointed to by client. All examples get run if serviceName
// is blank.
func RunExamplesForService(client *Client, serviceName string) error {
	rand.Seed(time.Now().UTC().UnixNano())

	services := broker.GetAllServices()

	for _, service := range services {

		if serviceName != "" && serviceName != service.Name {
			continue
		}

		for _, example := range service.Examples {
			err := RunExample(client, example, service)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

// RunExample runs a single example against the given service on the broker
// pointed to by client.
func RunExample(client *Client, example broker.ServiceExample, service *broker.BrokerService) error {
	executor, err := newExampleExecutor(client, example, service)
	if err != nil {
		return err
	}

	executor.LogTestInfo()

	// Cleanup the test if it fails partway through
	defer func() {
		log.Println("Cleaning up the environment")
		executor.Unbind()
		executor.Deprovision()
	}()

	if err := executor.Provision(); err != nil {
		return err
	}

	bindResponse, err := executor.Bind()
	if err != nil {
		return err
	}

	if err := executor.Unbind(); err != nil {
		return err
	}

	if err := executor.Deprovision(); err != nil {
		return err
	}

	// Check that the binding response has the same fields as expected
	var binding brokerapi.Binding
	err = json.Unmarshal(bindResponse, &binding)
	if err != nil {
		return err
	}

	credentialsEntry := binding.Credentials.(map[string]interface{})

	allContained := true
	for _, v := range service.BindOutputVariables {
		_, ok := credentialsEntry[v.FieldName]
		if !ok {
			allContained = false
			log.Printf("Error: credentials were missing property: %q", v.FieldName)
		}
	}

	if !allContained {
		return errors.New("Not all properties were found in the bound credentials")
	}

	return nil
}

func retry(timeout, period time.Duration, function func() (tryAgain bool, err error)) error {
	to := time.After(timeout)
	tick := time.Tick(period)

	if tryAgain, err := function(); !tryAgain {
		return err
	}

	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		case <-to:
			return errors.New("Timeout while waiting for result")
		case <-tick:
			tryAgain, err := function()

			if !tryAgain {
				return err
			}
		}
	}
}

func pollUntilFinished(client *Client, instanceId string) error {
	return retry(15*time.Minute, 15*time.Second, func() (bool, error) {
		log.Println("Polling for async job")

		resp := client.LastOperation(instanceId)
		if resp.InError() {
			return false, resp.Error
		}

		if resp.StatusCode != 200 {
			log.Printf("Bad status code %d, needed 200", resp.StatusCode)
			return true, nil
		}

		var responseBody map[string]string
		err := json.Unmarshal(resp.ResponseBody, &responseBody)
		if err != nil {
			return false, err
		}

		state := responseBody["state"]
		eq := state == string(brokerapi.Succeeded)
		log.Printf("Last operation for %q was %q\n", instanceId, state)

		return !eq, nil

	})
}

func newExampleExecutor(client *Client, example broker.ServiceExample, service *broker.BrokerService) (*exampleExecutor, error) {
	provisionParams, err := json.Marshal(example.ProvisionParams)
	if err != nil {
		return nil, err
	}

	bindParams, err := json.Marshal(example.BindParams)
	if err != nil {
		return nil, err
	}

	catalog, err := service.CatalogEntry()
	if err != nil {
		return nil, err
	}

	testid := rand.Uint32()

	return &exampleExecutor{
		Name:       fmt.Sprintf("%s/%s", service.Name, example.Name),
		ServiceId:  catalog.ID,
		PlanId:     example.PlanId,
		InstanceId: fmt.Sprintf("ex%d", testid),
		BindingId:  fmt.Sprintf("ex%d", testid),

		ProvisionParams: provisionParams,
		BindParams:      bindParams,

		client: client,
	}, nil
}

type exampleExecutor struct {
	Name string

	ServiceId  string
	PlanId     string
	InstanceId string
	BindingId  string

	ProvisionParams json.RawMessage
	BindParams      json.RawMessage

	client *Client
}

// Provision attempts to create a service instance from the example.
// Multiple calls to provision will attempt to create a resource with the same
// ServiceId and details.
// If the response is an async result, Provision will attempt to wait until
// the Provision is complete.
func (ee *exampleExecutor) Provision() error {
	log.Printf("Provisioning %s\n", ee.Name)

	resp := ee.client.Provision(ee.InstanceId, ee.ServiceId, ee.PlanId, ee.ProvisionParams)

	log.Println(resp.String())
	if resp.InError() {
		return resp.Error
	}

	switch resp.StatusCode {
	case 201:
		return nil
	case 202:
		return ee.pollUntilFinished()
	default:
		return fmt.Errorf("Unexpected response code %d", resp.StatusCode)
	}
}

func (ee *exampleExecutor) pollUntilFinished() error {
	return pollUntilFinished(ee.client, ee.InstanceId)
}

// Deprovision destroys the instance created by a call to Provision.
func (ee *exampleExecutor) Deprovision() error {
	log.Printf("Deprovisioning %s\n", ee.Name)
	resp := ee.client.Deprovision(ee.InstanceId, ee.ServiceId, ee.PlanId)

	log.Println(resp.String())
	if resp.InError() {
		return resp.Error
	}

	switch resp.StatusCode {
	case 200:
		return nil
	case 202:
		return ee.pollUntilFinished()
	default:
		return fmt.Errorf("Unexpected response code %d", resp.StatusCode)
	}
}

// Unbind unbinds the exact binding created by a call to Bind.
func (ee *exampleExecutor) Unbind() error {
	// XXX(josephlewis42) Due to some unknown reason, binding Postgres and MySQL
	// don't wait for all operations to finish before returning even though it
	// looks like they do so we can get 500 errors back the first few times we try
	// to unbind. Issue #222 was opened to address this. In the meantime this
	// is a hack to get around it that will still fail if the 500 errors truly
	// occur because of a real, unrecoverable, server error.
	return retry(5*time.Minute, 15*time.Second, func() (bool, error) {
		log.Printf("Unbinding %s\n", ee.Name)
		resp := ee.client.Unbind(ee.InstanceId, ee.BindingId, ee.ServiceId, ee.PlanId)

		log.Println(resp.String())
		if resp.InError() {
			return false, resp.Error
		}

		if resp.StatusCode == 200 {
			return false, nil
		}

		if resp.StatusCode == 500 {
			return true, nil
		}

		return false, fmt.Errorf("Unexpected response code %d", resp.StatusCode)
	})
}

// Bind executes the bind portion of the create, this can only be called
// once successfully as subsequent binds will attempt to create bindings with
// the same ID.
func (ee *exampleExecutor) Bind() (json.RawMessage, error) {
	log.Printf("Binding %s\n", ee.Name)
	resp := ee.client.Bind(ee.InstanceId, ee.BindingId, ee.ServiceId, ee.PlanId, ee.BindParams)

	log.Println(resp.String())
	if resp.InError() {
		return nil, resp.Error
	}

	if resp.StatusCode == 201 {
		return resp.ResponseBody, nil
	}

	return nil, fmt.Errorf("Unexpected response code %d", resp.StatusCode)
}

// LogTestInfo writes information about the running example and a manual backout
// strategy if the test dies part of the way through.
func (ee *exampleExecutor) LogTestInfo() {
	log.Printf("Running Example: %s\n", ee.Name)

	ips := fmt.Sprintf("--instanceid %q --planid %q --serviceid %q", ee.InstanceId, ee.PlanId, ee.ServiceId)
	log.Printf("gcp-service-broker client provision %s --params %q\n", ips, ee.ProvisionParams)
	log.Printf("gcp-service-broker client bind %s --bindingid %q --params %q\n", ips, ee.BindingId, ee.BindParams)
	log.Printf("gcp-service-broker client unbind %s --bindingid %q\n", ips, ee.BindingId)
	log.Printf("gcp-service-broker client deprovision %s\n", ips)
}
