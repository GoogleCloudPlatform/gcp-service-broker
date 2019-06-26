# Adding a built-in service to the service broker

## Background

The GCP service broker is a tool that developers use to provision access Google Cloud resources. The service broker currently supports a list of built-in services listed [here](use.md). 

This doc will cover how to add and test a built-in service to the GCP service broker.

## Overview

For the purposes of this tutorial, let’s say that the example service we want to add is [Google Memorystore for Redis](https://cloud.google.com/memorystore/). An example pull request can be found [here](https://github.com/GoogleCloudPlatform/gcp-service-broker/pull/477/files).

As you can see from the Redis example, you will generally be updating 3 files:

1. The service definition: `/pkg/providers/builtin/${SERVICE_NAME}/definition.go`
2. The service broker back-end: `/pkg/providers/builtin/${SERVICE_NAME}/broker.go`
3. The built-in broker registry: `/pkg/providers/builtin/registry.go`

The `definition.go` file is used to define the variables needed to provision a service (service ID, plan choices, etc). 
The `broker.go` file inherits [BrokerBase](../pkg/providers/builtin/base/broker_base.go) and parses the variables given in the Service Definition to make direct API calls to provision or bind services.
The `registry.go` file lists all built-in service providers. You will need to edit this file to include the new service.

#### Useful Resources
1. Google's **Go Reference** for the service. 

   To add a built-in service, you will essentially write a service template based on direct Go API calls. Therefore, it's important to be familiar with the Go API.
   Pay particular attention to the `CreateInstance` function and any required arguments. In the `Provision` function in `broker.go`, you will call `CreateInstance` to provision your new service. 

2. Any **Pricing Tiers**.

   Some resources will have different plan tiers or pricing tiers. 
   The built-in service may want to offer multiple plans to map to the different pricing tiers. For example, Redis has two [pricing tiers](https://cloud.google.com/memorystore/pricing), so the built-in Redis service offers two plans (BASIC and STANDARD HIGH AVAILABILITY) that map to the two pricing tiers.

3. Google **Documentation** for your new service.

   This will help you fill out fields in the `ServiceDefinition`. The Redis docs can be found [here](https://cloud.google.com/memorystore/docs/redis).
    
## Code Walkthrough
To get an idea of how to create the service template for your new service, we'll walk through through the changes made to `definition.go`, `broker.go`, and `registry.go` when implementing Redis as a new built-in service. At the end of the code walkthrough, you will also be able to test your new service!

### Service Definition

[Redis example: `../pkg/providers/builtin/redis/definition.go`](../pkg/providers/builtin/redis/definition.go)

The `ServiceDefinition` holds the necessary details to describe an OSB service and provision it. It is defined in [this struct](../pkg/broker/service_definition.go). Some of the basic fields are described below:

**Service Definition Basics**
* `ID`: Required: A unique Version 4 UUID. You can generate a UUID using the `uuidgen` CLI tool. If that's not installed, you can use this [Online UUID Generator](https://www.uuidgenerator.net/).
* `Name`: Required: A unique lowercase, hyphen-separated name for the service.
* `Description`: A description of the service can be found in the Google documentation.
* `DisplayName`: A human-friendly display name for the service.
* `ImageUrl`: The url of the `.svg` image of the service icon. This can often be found by inspecting the HTML element in the Google Cloud product page.
* `DocumentationUrl`: The url of the Google Cloud documentation.
* `SupportUrl`: The url of the Google Cloud support page. This can often be found through the Google Cloud documentation.
* `Tags`: Tags to help find the service.
* `Bindable`: Boolean value. Represents whether or not the service supports the bind call.
* `PlanUpdateable`: Boolean value. This should always be false, since the GCP service broker does not support updating plans.
* `DefaultRoleWhitelist`: String array of whitelisted roles. More about GCP roles can be found [here](https://cloud.google.com/iam/docs/understanding-roles), and the Redis roles are listed [here](https://cloud.google.com/iam/docs/understanding-roles#memorystore-redis-roles). 

**Provision Input Variables**

The `ProvisionInputVariables` are a string array of the `broker.BrokerVariable` struct. The struct definition can be found [here](../pkg/broker/variables.go), in [`../pkg/broker/variables.go`](../pkg/broker/variables.go).

Things to note:
* `Default` field specifies a default value in the case a value was not given.
* `Constraints` field specifies different constraints on the input, such as length or regex constraints.

The `ProvisionInputVariables` field is parsed in `broker.go` and used to make a direct API request. To construct the `ProvisionInputVariables` field, use the Go API to find the required input fields, and have a `broker.BrokerVariable` object for each required input field.
We can also add optional broker variables using our discretion.

From the Redis API, the [`CreateInstance` function](https://godoc.org/cloud.google.com/go/redis/apiv1beta1#CloudRedisClient.CreateInstance) 
takes as input a [`req *redispb.CreateInstanceRequest`](https://godoc.org/google.golang.org/genproto/googleapis/cloud/redis/v1beta1#CreateInstanceRequest).

The [`CreateInstanceRequest`](https://godoc.org/google.golang.org/genproto/googleapis/cloud/redis/v1beta1#CreateInstanceRequest) contains three required fields:

```Go
type CreateInstanceRequest struct {
    // Required. The resource name of the instance location using the form:
    //     `projects/{project_id}/locations/{location_id}`
    // where `location_id` refers to a GCP region.
    Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
    // Required. The logical name of the Redis instance in the customer project
    InstanceId string `protobuf:"bytes,2,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
    // Required. A Redis [Instance] resource
    Instance             *Instance `protobuf:"bytes,3,opt,name=instance,proto3" json:"instance,omitempty"`
}
```

The [`Instance field`](https://godoc.org/google.golang.org/genproto/googleapis/cloud/redis/v1beta1#Instance) also contains a few required fields:

```Go
type Instance struct {
    // Required. Unique name of the resource in this scope including project and
    // location using the form:
    //     `projects/{project_id}/locations/{location_id}/instances/{instance_id}`
    Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
    
    // Required. The service tier of the instance.
    Tier Instance_Tier `protobuf:"varint,17,opt,name=tier,proto3,enum=google.cloud.redis.v1.Instance_Tier" json:"tier,omitempty"`
    // Required. Redis memory size in GiB.
    MemorySizeGb int32 `protobuf:"varint,18,opt,name=memory_size_gb,json=memorySizeGb,proto3" json:"memory_size_gb,omitempty"`
}
```

Let's see how the Redis `ProvisionInputVariables` map to the required fields:

* `location_id`: Maps to the `region` field
* `MemorySizeGb`: Maps to the `memory_size_gb` field
* `InstanceId`: Maps to the `instance_id` field
* `Parent` and `Name`: Since we are given the `project_id`, `location_id`, and `instance_id`, we can specify these fields in `broker.go`.
* `Tier`: This is specified in the Plan Variables field, explained below.

**Plan and Plan Variables** 

Since Redis has two different [pricing tiers](https://cloud.google.com/memorystore/pricing), I decided to offer two service plans that map to the two pricing tiers. 
That way, operators and programmers can choose the plan that works best for their purposes.

**Note**: If you reference the Redis pricing page, you can see that each Redis instance is billed at a difference price depending on the capacity tier. There are 5 different capacity tiers. I could have offered 10 different plans to cover all
the pricing options (a plan for each service tier and capacity tier), but for simplicity's sake, I decided two offer two different service plans and specify the capacity in a separate `ProvisionInputVariable`.

There are two steps to adding service plans:

1. Instantiate the `Plans` field using an array of `broker.ServicePlan` structs:

    The `broker.ServicePlan` struct has a `brokerapi.ServicePlan` and a `ServiceProperties` field:
    
    ```Go
    type ServicePlan struct {
      brokerapi.ServicePlan
      ServiceProperties map[string]string `json:"service_properties"`
    }
    ```
    
    The `brokerapi.ServicePlan` struct contains information about the plan. The `ID` field here should also be a unique UUID, separate from the `ServiceDefinition` UUID.
    
    ```Go
    type ServicePlan struct {
      ID              string               `json:"id"`
      Name            string               `json:"name"`
      Description     string               `json:"description"`
      Free            *bool                `json:"free,omitempty"`
      Bindable        *bool                `json:"bindable,omitempty"`
      Metadata        *ServicePlanMetadata `json:"metadata,omitempty"`
      Schemas         *ServiceSchemas      `json:"schemas,omitempty"`
      MaintenanceInfo *MaintenanceInfo     `json:"maintenance_info,omitempty"`
    }
    ```
    
    We can use the `ServiceProperties` field to include information about the plan needed to make the Go API request. 
    
    For the Redis service, we need to include information about the `Tier` field, so the `ServiceProperties` map to the service tiers each plan offers.

2. Instantiate the `PlanVariables` field using an array of `broker.BrokerVariable` structs.

    This is so the `ServiceProperties` field can later be parsed by `broker.go`. Here is the `PlanVariables` field in the Redis example:
    
    ```Go
    PlanVariables: []broker.BrokerVariable{
          {
            FieldName: "service_tier",
            Type:      broker.JsonTypeString,
            Details:   "Either BASIC or STANDARD_HA. See: https://cloud.google.com/memorystore/pricing for more information.",
            Default:   "basic",
            Required:  true,
          },
        },
    ```

### Service broker back-end
[Redis example: `../pkg/providers/builtin/redis/broker.go`](../pkg/providers/builtin/redis/broker.go)

The service broker back-end inherits `base.BrokerBase`. We can extract `ProvisionInputVariables` and `PlanInputVariables` using the `provisionContext.GetString("${field_name}")` helper function. We can then make direct API calls to provision and deprovision the GCP service.

### Broker registry
[Go to: `../pkg/providers/builtin/registry.go`](../pkg/providers/builtin/registry.go)

Add the folder holding your `definition.go` and `service.go` to the `imports`, and your service definition to the `RegisterBuiltinBrokers` function.

## Testing your service

In the `Examples` field in your service definition, you can add a service definition to test out. 

```Go
Examples: []broker.ServiceExample{
	{
			Name:            "Basic Redis Configuration",
			Description:     "Create a Redis instance with basic service tier.",
			PlanId:          "dd1923b6-ac26-4697-83d6-b3a0c05c2c94",
			ProvisionParams: map[string]interface{}{},
			BindParams: map[string]interface{}{
			  "role": "redis.viewer",
			},
	},
},
```

From your command line, inside the gcp-service-broker folder, run:

1. `go build`
2. To serve: `./gcp-service-broker --config ../minimal.yml serve`

   You can view the auto-generated docs from your `ServiceDefinition` at http://localhost:8000/docs.
3. In a separate terminal, to run examples: `./gcp-service-broker --config config.yml client run-examples --service-name ${SERVICE_NAME}`

   The `${SERVICE_NAME}` should be the `Name` field specified in the `ServiceDefinition`.
