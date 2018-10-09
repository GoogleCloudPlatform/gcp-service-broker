# Creating Custom Services (UNRELEASED)

This documentation is for an **UNRELEASED** upcoming feature for the GCP Service Broker and should not be considered complete.

## Variable resolution

The variables fed into your Terraform services file are resolved in the following order:

* Variables defined in your `computed_variables` JSON list.
* Variables defined by the selected service plan in its `service_properties` map.
* User defined variables (in `provision_input_variables` or `bind_input_variables`)
* **TODO** Operator default variables loaded from the environment.
* Default variables (in `provision_input_variables` or `bind_input_variables`).

Note that the order the variables are combined in code is slightly different.

* **TODO** Operator default variables loaded from the environment.
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

* `time.nano() -> string`
  * This function returns the current UNIX time in nanoseconds as a decimal string.
* `str.truncate(count, string) -> string`
  * Trims the given string to be at most `count` characters long.
  * If the string is already shorter, nothing is changed.
* `counter.next() -> int`
  * Provides a counter that increments once per call within the same call context.
  * The counter is reset on restart of the application.
* `rand.base64(count) -> string`
  * Generates `count` bytes of cryptographically secure randomness and converts it to [URL Encoded Base64](https://tools.ietf.org/html/rfc4648).
  * The randomness makes it suitable for using as passwords.
