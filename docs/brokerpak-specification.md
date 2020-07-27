# Brokerpak V1 specification

This document will explain how a brokerpak is structured, and the schema it follows.

A brokerpak is comprised of a versioned Terraform binary and providers for one
or more platform, a manifest, one or more service definitions, and source code.
Here are the contents of an example brokerpak:

```
MODE        SIZE      NAME
drwx------  0         
drwxr-xr-x  0         bin
drwxr-xr-x  0         bin/linux
drwxr-xr-x  0         bin/linux/386
-rwxrwxr-x  75669952  bin/linux/386/terraform
-rwxr-xr-x  35176256  bin/linux/386/terraform-provider-google_v1.19.0_x4
drwxr-xr-x  0         bin/linux/amd64
-rwxrwxr-x  89397536  bin/linux/amd64/terraform
-rwxr-xr-x  41450144  bin/linux/amd64/terraform-provider-google_v1.19.0_x4
drwx------  0         definitions
-rw-------  5810      definitions/service-0.yml
-rw-------  467       manifest.yml
drwxr-xr-x  0         src
-rw-r--r--  8082786   src/terraform-provider-google.zip
-rw-r--r--  14418716  src/terraform.zip
```

You can create, inspect, validate, document and test brokerpaks using the `pak` sub-command.
Run the command `gcp-service-broker pak help` for more information about creating a pak.

## Manifest

Each brokerpak has a manifest file named `manifest.yml`. This file determines
which architectures your brokerpak will work on, which plugins it will use,
and which services it will provide.

### Schema

#### Manifest YAML file

| Field | Type | Description |
| --- | --- | --- |
| packversion* | int | The version of the schema the manifest adheres to. This MUST be set to `1` to be compatible with the brokerpak specification v1. |
| version* | string | The version of this brokerpak. It's RECOMMENDED you follow [semantic versioning](https://semver.org/) for your brokerpaks. |
| name* | string | The name of this brokerpak. It's RECOMMENDED that this be lower-case and include only alphanumeric characters, dashes, and underscores. |
| metadata | object | A free-form field for key/value pairs of additional information about this brokerpak. This could include the authors, creation date, source code repository, etc. |
| platforms* | array of platform | The platforms this brokerpak will be executed on. |
| terraform_binaries* | array of Terraform resource | The list of Terraform providers and Terraform that'll be bundled with the brokerpak. |
| service_definitions* | array of string | Each entry points to a file relative to the manifest that defines a service as part of the brokerpak. |
| parameters | array of parameter | These values are set as environment variables when Terraform is executed. |

#### Platform object

The platform OS and architecture follow Go's naming scheme.

| Field | Type | Description |
| --- | --- | --- |
| os* | string | The operating system of the platform. |
| arch* | string | The architecture of the platform. |

#### Terraform resource object

This structure holds information about a specific Terraform version or Resource.

| Field | Type | Description |
| --- | --- | --- |
| name* | string | The name of this resource. e.g. `terraform-provider-google-beta`. |
| version* | string | The version of the resource e.g. 1.19.0. |
| source* | string | The URL to a zip of the source code for the resource. |
| url_template | string | (optional) A custom URL template to get the release of the given tool. Available parameters are ${name}, ${version}, ${os}, and ${arch}. If unspecified the default Hashicorp Terraform download server is used. |

#### Parameter object

This structure holds information about an environment variable that the user can set on the Terraform instance.
These variables are first resolved from the configuration of the brokerpak then against a global set of values.

| Field | Type | Description |
| --- | --- | --- |
| name* | string | The environment variable that will be injected e.g. `PROJECT_ID`. |
| description* | string | A human readable description of what the variable represents. |

### Example

```yaml
packversion: 1
name: my-custom-services
version: 1.0.0
metadata:
  author: someone@my-company.com
platforms:
- os: linux
  arch: "386"
- os: linux
  arch: amd64
terraform_binaries:
- name: terraform
  version: 0.11.9
  source: https://github.com/hashicorp/terraform/archive/v0.11.9.zip
- name: terraform-provider-google
  version: 1.19.0
  source: https://github.com/terraform-providers/terraform-provider-google/archive/v1.19.0.zip
service_definitions:
- custom-cloud-storage.yml
- custom-redis.yml
- service-mesh.yml
parameters:
- name: TF_VAR_redis_version
  description: Set this to override the Redis version globally via injected Terraform variable.
```

## Services

### Schemas

#### Service YAML flie

| Field | Type | Description |
| --- | --- | --- |
| version* | int |  The version of the schema the service definition adheres to. This MUST be set to `1` to be compatible with the brokerpak specification v1. |
| name* | string | A CLI-friendly name of the service. MUST only contain alphanumeric characters, periods, and hyphens (no spaces). MUST be unique across all service objects returned in this response. MUST be a non-empty string. |
| id* | string | A UUID used to correlate this service in future requests to the Service Broker. This MUST be globally unique such that Platforms (and their users) MUST be able to assume that seeing the same value (no matter what Service Broker uses it) will always refer to this service. |
| description* | string | A short description of the service. MUST be a non-empty string. |
| tags | array of strings | Tags provide a flexible mechanism to expose a classification, attribute, or base technology of a service, enabling equivalent services to be swapped out without changes to dependent logic in applications, buildpacks, or other services. E.g. mysql, relational, redis, key-value, caching, messaging, amqp. |
| display_name* | string | The name of the service to be displayed in graphical clients. |
| image_url* | string | The URL to an image or a data URL containing an image. |
| documentation_url* | string | Link to documentation page for the service. |
| support_url* | string | Link to support page for the service. |
| plans* | array of plan objects | A list of plans for this service, schema is defined below. MUST contain at least one plan. |
| provision* | action object | Contains configuration for the provision operation, schema is defined below. |
| bind* | action object | Contains configuration for the bind operation, schema is defined below. |
| examples* | example object | Contains examples for the service, used in documentation and testing.  MUST contain at least one example. |

#### Plan object

A service plan in a human-friendly format that can be converted into an OSB compatible plan.

| Field | Type | Description |
| --- | --- | --- |
| name* | string | The CLI-friendly name of the plan. MUST only contain alphanumeric characters, periods, and hyphens (no spaces). MUST be unique within the service. MUST be a non-empty string. |
| id* | string | A GUID for this plan in UUID format. This MUST be globally unique such that Platforms (and their users) MUST be able to assume that seeing the same value (no matter what Service Broker uses it) will always refer to this plan. |
| description* | string | A short description of the plan. MUST be a non-empty string. |
| display_name* | string | The name of the plan to be displayed in graphical clients. |
| bullets | array of string | Features of this plan, to be displayed in a bulleted-list. |
| free | boolean | When false, Service Instances of this plan have a cost. The default is false. |
| properties* | map of string:string | Default values for the provision and bind calls. |

#### Action object

The Action object contains a Terraform template to execute as part of a
provision or bind action, and the inputs and outputs to that template.

| Field | Type | Description |
| --- | --- | --- |
| plan_inputs | array of variable | Defines constraints and settings for the variables plans provide in their properties map. |
| user_inputs | array of variable | Defines constraints and settings for the variables users provide as part of their request. |
| computed_inputs | array of computed variable | Defines default values or overrides that are executed before the template is run. |
| template | string | The complete HCL of the Terraform template to execute. |
| outputs | array of variable | Defines constraints and settings for the outputs of the Terraform template. This MUST match the Terraform outputs and the constraints WILL be used as part of integration testing. |

#### Variable object

The variable object describes a particular input or output variable. The
structure is turned into a JSONSchema to validate the inputs or outputs.
Outputs are _only_ validated on integration tests.

| Field | Type | Description |
| --- | --- | --- |
| required | boolean | Should the user request fail if this variable isn't provided? |
| field_name* | string | The name of the JSON field this variable serializes/deserializes to. |
| type* | string | The JSON type of the field. This MUST be a valid JSONSchema type excepting `null`. |
| details* | string | Provides explanation about the purpose of the variable. |
| default | any | The default value for this field. If `null`, the field MUST be marked as required. If a string, it will be executed as a HIL expression and cast to the appropriate type described in the `type` field. See the "Expression language reference" section for more information about what's available. |
| enum | map of any:string | Valid values for the field and their human-readable descriptions suitable for displaying in a drop-down list. |
| constraints | map of string:any | Holds additional JSONSchema validation for the field. The following keys are supported: `examples`, `const`, `multipleOf`, `minimum`, `maximum`, `exclusiveMaximum`, `exclusiveMinimum`, `maxLength`, `minLength`, `pattern`, `maxItems`, `minItems`, `maxProperties`, `minProperties`, and `propertyNames`. |


#### Computed Variable Object

Computed variables allow you to evaluate arbitrary HIL expressions against
variables or metadata about the provision or bind call.

| Field | Type | Description |
| --- | --- | --- |
| name* | string | The name of the variable. |
| default* | any | The value to set the variable to. If it's a string, it will be evaluated by the expression engine and cast to the provided type afterwards. See the "Expression language reference" section for more information about what's available. |
| overwrite | boolean | If a variable already exists with the same name, should this one replace it? |
| type | string | The JSON type of the field it will be cast to if evaluated as an expression. If defined, this MUST be a valid JSONSchema type excepting `null`. |

### Example

```yaml
version: 1
name: example-service
id: 00000000-0000-0000-0000-000000000000
description: a longer service description
display_name: Example Service
image_url: https://example.com/icon.jpg
documentation_url: https://example.com
support_url: https://example.com/support.html
tags: [gcp, example, service]
plans:
- name: example-email-plan
  id: 00000000-0000-0000-0000-000000000001
  description: Builds emails for example.com.
  display_name: example.com email builder
  bullets:
  - information point 1
  - information point 2
  - some caveat here
  properties:
    domain: example.com
provision:
  plan_inputs:
  - required: true
    field_name: domain
    type: string
    details: The domain name
  user_inputs:
  - required: true
    field_name: username
    type: string
    details: The username to create
  computed_inputs: []
  template: |-
    variable domain {type = "string"}
    variable username {type = "string"}

    output email {value = "${var.username}@${var.domain}"}

  outputs:
  - required: true
    field_name: email
    type: string
    details: The combined email address
bind:
  plan_inputs: []
  user_inputs: []
  computed_inputs:
  - name: address
    default: ${instance.details["email"]}
    overwrite: true
  template: |-
    resource "random_string" "password" {
      length = 16
      special = true
      override_special = "/@\" "
    }

    output uri {value = "smtp://${var.address}:${random_string.password.result}@smtp.mycompany.com"}

  outputs:
  - required: true
    field_name: uri
    type: string
    details: The uri to use to connect to this service
examples:
- name: Example
  description: Examples are used for documenting your service AND as integration tests.
  plan_id: 00000000-0000-0000-0000-000000000001
  provision_params:
    username: my-account
  bind_params: {}

```

## Expression language reference

The broker uses the [HIL expression language](https://github.com/hashicorp/hil) with a limited set of built-in functions.

### Functions

The following string interpolation functions are available for use:

* `assert(condition_bool, message_string) -> bool`
  * If the condition is false, then an error will be raised to the user containing `message_string`.
  * Avoid using this function. Instead, try to make it so your users can't get into a bad state to begin with. See the "design guidelines" section.  
  * In the words of [PEP-20](https://www.python.org/dev/peps/pep-0020/):
    * If the implementation is hard to explain, it's a bad idea.
    * If the implementation is easy to explain, it may be a good idea.
* `time.nano() -> string`
  * This function returns the current time as a Unix time, the number of nanoseconds elapsed since January 1, 1970 UTC, as a decimal string.
  * The result is undefined if the Unix time in nanoseconds cannot be represented by an int64 (a date before the year 1678 or after 2262).
* `regexp.matches(regex_string, string) -> bool`
  * Checks if the string matches the given regex.
* `str.truncate(count, string) -> string`
  * Trims the given string to be at most `count` characters long.
  * If the string is already shorter, nothing is changed.
* `counter.next() -> int`
  * Provides a counter that increments once per call within the same call context.
  * The counter is reset on restart of the application.
* `rand.base64(count) -> string`
  * Generates `count` bytes of cryptographically secure randomness and converts it to [URL Encoded Base64](https://tools.ietf.org/html/rfc4648).
  * The randomness makes it suitable for using as passwords.
* `json.marshal(type) -> string`
  * Returns a JSON marshaled string of the given type.
* `map.flatten(keyValueSeparator, tupleSeparator, map)`
  * Converts a map into a string with each key/value pair separated by `keyValueSeparator` and each entry separated by `tupleSeparator`.
  * The output is deterministic.
  * Example: if `labels = {"key1":"val1", "key2":"val2"}` then `map.flatten(":", ";", labels)` produces `key1:val1;key2:val2`.

### Variables

The broker makes additional variables available to be used during provision and bind calls.

#### Resolution

The variables fed into your Terraform services file are resolved in the following order:

* Variables defined in your `computed_variables` JSON list.
* Variables defined by the selected service plan in its `service_properties` map.
* Variables overridden by the plan (in `provision_overrides` or `bind_overrides`).
* User defined variables (in `provision_input_variables` or `bind_input_variables`).
* Operator default variables loaded from the environment.
* Default variables (in `provision_input_variables` or `bind_input_variables`).

Note that the order the variables are combined in code is slightly different.

* Operator default variables loaded from the environment.
* User defined variables (in `provision_input_variables` or `bind_input_variables`)
* **If the variables are not defined yet** default variables (in `provision_input_variables` or `bind_input_variables`).
* Variables defined by the selected service plan in its `service_properties` map.
* Variables defined in your `computed_variables` JSON list.

Moving default variables to be loaded third allow their computed values to make more sense.
This is because they can resolve variables to the user's values first.

#### Provision

* `request.service_id` - _string_ The GUID of the requested service.
* `request.plan_id` - _string_ The ID of the requested plan. Plan IDs are unique within an instance.
* `request.instance_id` - _string_ The ID of the requested instance. Instance IDs are unique within a service.
* `request.default_labels` - _map[string]string_ A map of labels that should be applied to the created infrastructure for billing/accounting/tracking purposes.

#### Bind

* `request.binding_id` - _string_ The ID of the new binding.
* `request.instance_id` - _string_ The ID of the existing instance to bind to.
* `request.service_id` - _string_ The GUID of the service this binding is for.
* `request.plan_id` - _string_ The ID of plan the instance was created with.
* `request.plan_properties` - _map[string]string_ A map of properties set in the service's plan.
* `request.app_guid` - _string_ The ID of the application this binding is for.
* `instance.name` - _string_ The name of the instance.
* `instance.details` - _map[string]any_ Output variables of the instance as specified by ProvisionOutputVariables.

## File format

The brokerpak itself is a zip file with the extension `.brokerpak`.
In the root is the `manifest.yml` file, which will specify the version of the pak.

There are three directories in the pak's root:

* `src/` an unstructured directory that holds source code for the bundled binaries, this is for complying with 3rd party licenses.
* `bin/` contains binaries under `bin/{os}/{arch}` sub-directories for each supported platform.
* `definitions/` contain the service definition YAML files.
