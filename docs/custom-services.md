# Creating Custom Services (UNRELEASED)

This documentation is for an **UNRELEASED** upcoming feature for the GCP Service Broker and should not be considered complete.

## Variable resolution

The variables fed into your Terraform services file are resolved in the following order:

* Variables defined in your `computed_variables` JSON list.
* (Only for Provision) Variables defined by the selected service plan in its `service_properties` map.
* User defined variables (in `provision_input_variables` or `bind_input_variables`)
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

## Expression language reference

The broker uses the [HIL expression language](https://github.com/hashicorp/hil) with a limited set of built-in functions.

## Function reference

The following string interpolation functions are available for use:

* `assert(condition_bool, message_string) -> bool`
  * If the condition is false, then an error will be raised to the user containing `message_string`.
  * Avoid using this function. Instead, try to make it so your users can't get into a bad state to begin with. See the "design guidelines" section.  
  * In the words of [PEP-20](https://www.python.org/dev/peps/pep-0020/):
    * If the implementation is hard to explain, it's a bad idea.
    * If the implementation is easy to explain, it may be a good idea.
* `time.nano() -> string`
  * This function returns the current UNIX time in nanoseconds as a decimal string.
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

## Variable reference

The broker makes additional variables available to be used during provision and bind calls.

### Provision

* `request.service_id` - _string_ The GUID of the requested service.
* `request.plan_id` - _string_ The ID of the requested plan. Plan IDs are unique within an instance.
* `request.instance_id` - _string_ The ID of the requested instance. Instance IDs are unique within a service.
* `request.default_labels` - _map[string]string_ A map of labels that should be applied to the created infrastructure for billing/accounting/tracking purposes.

### Bind

* `request.binding_id` - _string_ The ID of the new binding.
* `request.instance_id` - _string_ The ID of the existing instance to bind to.
* `request.service_id` - _string_ The GUID of the service this binding is for.
* `request.plan_id` - _string_ The ID of plan the instance was created with.
* `request.app_guid` - _string_ The ID of the application this binding is for.
* `instance.name` - _string_ The name of the instance.
* `instance.details` - _map[string]any_ Output variables of the instance as specified by ProvisionOutputVariables.

## Design guidelines

When you're creating a new service for the broker you're designing for three separate sets of people:

* The users, developers who will use your service to provision things they work with day-to-day.
* The operators, the people who are responsible for approving services and plans for developers to use.
* Yourself, the person who has to maintain the service, strike the right balance of power between the operators and users, and make sure and make sure the new plans/services work as intended.

The following sections contain guidelines to help you out.

### Deciding what to include

Services don't need to map one-to-one with cloud products, and probably shouldn't.
Instead, services should be focused around particular workflows, allowing you to get a single, useful, task done.
Service plans allow you to scale that up or down.

For example, Google CloudSQL contains options for high availability, PostgreSQL and MySQL servers, managing on-prem servers, and read-only replication architectures.
These features all exist for different audiences, and a generic service trying to fit all the use-cases won't give a good experience to the users, operators, or maintainers.

If you find yourself wishing you could selectively enable or disable variables based on flags, it's a sign you should break down your code into another service.
For example, a Cloud Storage bucket can be configured to have a retention policy, a public-facing URL, and/or push file-change updates to a Pub/Sub queue.
It would be a good idea to break those features into multiple distinct services:

* One for hosting a static website with settings for URL, index/error pages, and CNAME.
* Another the other for general storage that has retention policies.
* A third that also provisions a Pub/Sub queue and acts as a staging area for data.

Breaking things down like this makes it easier to figure out what variables you need to expose, what risks they entail and what kind of plans you'll want:

* The static site plans could be simple, maybe containing different domain names and regions.
* The archive bucket plans could be for different retention policies and object durability.
* The staging bucket plans could include options for setting up alerting and the queue at the same time as the bucket is created.

Each cloud service you expose will have a plethora of tunable parameters to choose from.
Ideally, you should expose enough to be useful to developers and secure, but few enough that your service has a well defined use-case.
You can always add more parameters later, but you can never get rid of one.

### Deciding where to include things

Each parameter can either be set by the operator when they define plans (or in your plans that the operators enable for users) or by the user.

In general, properties which have monetary cost or affect the security of the platform should be put in the plan and properties affecting the behavior of the resource should be defined by the user.

In our static site bucket example the operator would create plans for different domain names (security) and bucket locations/durabilities (pricing) and the developer would get to set the parameters for the default index/error pages and maybe hostname. A full CNAME would be calculated from the hostname and domain name combination. It isn't clear who would get control over the Pub/Sub endpoint. On one hand, the developers might need it to update a search engine index but on the other the operator might to conduct ongoing security audits.

### Deciding on sensible defaults

The GCP Service Broker operates under the model that the users are benign but fallible.
Sensible defaults are secure and work well in the average use-case.
This highly depends on your target audience.

For example, a Pub/Sub instance with one-to-many semantics might default to a read-only role, assuming the default consumer is just going to be a worker node whereas a Pub/Sub instance with many-to-many semantics might default to a read/write role even if some consumers want to be read-only.

### Deciding on what your default plans will be

If you've gotten to this point, you should have a clear understanding of what your service is trying to accomplish, who the users are, and what variables are configurable in your plans.
It can be tempting to include every permutation of the variables for plans.
However, less is more.
Operators need to look at each plan, decide if it fits a distinct use-case, budget and security model then make it available to individual teams.
A few plans that hit key use-cases are much easier to grok.

Let's go back to our archival storage use-case. Instead of creating plans for every availability tier and zone, we'd create plans for these criteria:

* Companies hosting archives in the US
* Companies that do not want their data in the US
* Teams that need buckets they control for non-prod environments

We'd end up with something the following:

    (US | EU | Asia) x (high availability + legally mandated retention policy | standard availability + no retention policy) = 6 plans
