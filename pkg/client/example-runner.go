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
	log.Printf("Running Example: %s/%s\n", service.Name, example.Name)
	catalogEntry := service.CatalogEntry()

	serviceId := catalogEntry.ID
	planId := example.PlanId

	testid := rand.Uint32()
	instanceId := fmt.Sprintf("ex%d", testid)
	bindingId := fmt.Sprintf("ex%d", testid)

	provisioningDetails, err := json.Marshal(example.ProvisionParams)
	if err != nil {
		return err
	}

	bindParams, err := json.Marshal(example.BindParams)
	if err != nil {
		return err
	}

	log.Printf("gcp-service-broker client provision --instanceid %q --planid %q --serviceid %q --params %q\n", instanceId, planId, serviceId, provisioningDetails)
	log.Printf("gcp-service-broker client bind --instanceid %q --planid %q --serviceid %q --bindingid %q --params %q\n", instanceId, planId, serviceId, bindingId, bindParams)
	log.Printf("gcp-service-broker client unbind --instanceid %q --planid %q --serviceid %q --bindingid %q \n", instanceId, planId, serviceId, bindingId)
	log.Printf("gcp-service-broker client deprovision --instanceid %q --planid %q --serviceid %q\n", instanceId, planId, serviceId)

	log.Println("Provisioning")
	resp := client.Provision(instanceId, serviceId, planId, provisioningDetails)
	if err := handleProvisionResponse(client, instanceId, resp); err != nil {
		return err
	}

	log.Println("Binding")
	resp = client.Bind(instanceId, bindingId, serviceId, planId, bindParams)
	bindErr := handleBindResponse(resp, service)

	if bindErr == nil {
		log.Println("Unbinding")
		resp := client.Unbind(instanceId, bindingId, serviceId, planId)
		bindErr = handleUnbindResponse(resp)
	}

	log.Println("Deprovisioning")
	resp = client.Deprovision(instanceId, serviceId, planId)
	if err := handleDeprovisionResponse(client, instanceId, resp); err != nil {
		return err
	}

	return bindErr
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

func handleProvisionResponse(client *Client, instanceId string, resp *BrokerResponse) error {
	log.Println(resp.String())
	if resp.InError() {
		return resp.Error
	}

	if resp.StatusCode == 201 {
		return nil
	}

	if resp.StatusCode == 202 {
		return pollUntilFinished(client, instanceId)
	}

	return errors.New(fmt.Sprintf("Unexpected response code %d", resp.StatusCode))
}

func handleDeprovisionResponse(client *Client, instanceId string, resp *BrokerResponse) error {
	log.Println(resp.String())
	if resp.InError() {
		return resp.Error
	}

	if resp.StatusCode == 200 {
		return nil
	}

	if resp.StatusCode == 202 {
		return pollUntilFinished(client, instanceId)
	}

	return errors.New(fmt.Sprintf("Unexpected response code %d", resp.StatusCode))
}

func handleBindResponse(resp *BrokerResponse, service *broker.BrokerService) error {
	log.Println(resp.String())
	if resp.InError() {
		return resp.Error
	}

	if resp.StatusCode != 201 {
		return errors.New(fmt.Sprintf("Unexpected response code %d", resp.StatusCode))
	}

	// Check all the response variables
	var binding brokerapi.Binding
	err := json.Unmarshal(resp.ResponseBody, &binding)
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

func handleUnbindResponse(resp *BrokerResponse) error {
	log.Println(resp.String())
	if resp.InError() {
		return resp.Error
	}

	if resp.StatusCode == 201 {
		return nil
	}

	return errors.New(fmt.Sprintf("Unexpected response code %d", resp.StatusCode))
}
