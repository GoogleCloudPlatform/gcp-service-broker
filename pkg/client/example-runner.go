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

func RunExamples(client *Client) error {
	return RunExamplesForService(client, "")
}

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

func pollUntilFinished(client *Client, instanceId string) error {
	timeout := time.After(15 * time.Minute)
	tick := time.Tick(15 * time.Second)

	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		case <-timeout:
			return errors.New("Timeout while waiting for result")
		case <-tick:
			log.Println("Polling for async job")
			resp := client.LastOperation(instanceId)
			if resp.InError() {
				return resp.Error
			}

			if resp.StatusCode != 200 {
				log.Printf("Bad status code %d, needed 200", resp.StatusCode)
				continue
			}

			var responseBody map[string]string
			err := json.Unmarshal(resp.ResponseBody, &responseBody)
			if err != nil {
				return err
			}

			state := responseBody["state"]
			eq := state == string(brokerapi.Succeeded)
			log.Printf("Last operation for %q was %q\n", instanceId, state)

			if eq {
				return nil
			}
		}
	}
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

	testid := rand.Uint32()

	return &exampleExecutor{
		Name:       fmt.Sprintf("%s/%s", service.Name, example.Name),
		ServiceId:  service.CatalogEntry().ID,
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
		return errors.New(fmt.Sprintf("Unexpected response code %d", resp.StatusCode))
	}
}

func (ee *exampleExecutor) pollUntilFinished() error {
	return pollUntilFinished(ee.client, ee.InstanceId)
}

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
		return errors.New(fmt.Sprintf("Unexpected response code %d", resp.StatusCode))
	}
}

func (ee *exampleExecutor) Unbind() error {
	log.Printf("Unbinding %s\n", ee.Name)
	resp := ee.client.Unbind(ee.InstanceId, ee.BindingId, ee.ServiceId, ee.PlanId)

	log.Println(resp.String())
	if resp.InError() {
		return resp.Error
	}

	if resp.StatusCode == 200 {
		return nil
	}

	return errors.New(fmt.Sprintf("Unexpected response code %d", resp.StatusCode))
}

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

	return nil, errors.New(fmt.Sprintf("Unexpected response code %d", resp.StatusCode))
}

func (ee *exampleExecutor) LogTestInfo() {
	log.Printf("Running Example: %s\n", ee.Name)

	ips := fmt.Sprintf("--instanceid %q --planid %q --serviceid %q", ee.InstanceId, ee.PlanId, ee.ServiceId)
	log.Printf("gcp-service-broker client provision %s --params %q\n", ips, ee.ProvisionParams)
	log.Printf("gcp-service-broker client bind %s --bindingid %q --params %q\n", ips, ee.BindingId, ee.BindParams)
	log.Printf("gcp-service-broker client unbind %s --bindingid %q\n", ips, ee.BindingId)
	log.Printf("gcp-service-broker client deprovision %s\n", ips)
}
