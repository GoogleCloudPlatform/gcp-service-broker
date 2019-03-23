# Installation Customization

This file documents the various environment variables you can set to change the functionality of the service broker.
If you are using the PCF Tile deployment, then you can manage all of these options through the operator forms.
If you are running your own, then you can set them in the application manifest of a PCF deployment, or in your pod configuration for Kubernetes.


## Root Service Account

Please paste in the contents of the json keyfile (un-encoded) for your service account with owner credentials.

You can configure the following environment variables:

| Environment Variable | Type | Description |
|----------------------|------|-------------|
| <tt>ROOT_SERVICE_ACCOUNT_JSON</tt> <b>*</b> | text | <p>Root Service Account JSON. </p>|


\* = Required


## Database Properties

Connection details for the service broker's database. It must be MySQL compatible.

You can configure the following environment variables:

| Environment Variable | Type | Description |
|----------------------|------|-------------|
| <tt>DB_HOST</tt> <b>*</b> | string | <p>Database host. </p>|
| <tt>DB_USERNAME</tt> | string | <p>Database username. </p>|
| <tt>DB_PASSWORD</tt> | secret | <p>Database password. </p>|
| <tt>DB_PORT</tt> <b>*</b> | string | <p>Database port.  Default: <code>3306</code></p>|
| <tt>DB_NAME</tt> <b>*</b> | string | <p>Database name.  Default: <code>servicebroker</code></p>|
| <tt>CA_CERT</tt> | text | <p>Server CA cert. </p>|
| <tt>CLIENT_CERT</tt> | text | <p>Client cert. </p>|
| <tt>CLIENT_KEY</tt> | text | <p>Client key. </p>|


\* = Required


## Service Configuration

Configuration for built-in and Brokerpak services.

You can configure the following environment variables:

| Environment Variable | Type | Description |
|----------------------|------|-------------|
| <tt>GSB_BROKERPAK_CONFIG</tt> <b>*</b> | text | <p>Global Brokerpak Configuration. A JSON map of configuration key/value pairs for all brokerpaks. If a variable isn't found in the specific brokerpak's configuration it's looked up here. Default: <code>{}</code></p>|
| <tt>GSB_SERVICE_CONFIG</tt> <b>*</b> | text | <p>Service Configuration Options. See the configuration.md file or /configuration endpoint on the service for how to configure services using this field. Default: <code>{}</code></p>|


\* = Required


## Feature Flags

Service broker feature flags.

You can configure the following environment variables:

| Environment Variable | Type | Description |
|----------------------|------|-------------|
| <tt>GSB_COMPATIBILITY_ENABLE_BUILTIN_BROKERPAKS</tt> <b>*</b> | boolean | <p>enable-builtin-brokerpaks. Load brokerpaks that are built-in to the software. Default: <code>true</code></p>|
| <tt>GSB_COMPATIBILITY_ENABLE_BUILTIN_SERVICES</tt> <b>*</b> | boolean | <p>enable-builtin-services. Enable services that are built in to the broker i.e. not brokerpaks. Default: <code>true</code></p>|
| <tt>GSB_COMPATIBILITY_ENABLE_CATALOG_SCHEMAS</tt> <b>*</b> | boolean | <p>enable-catalog-schemas. Enable generating JSONSchema for the service catalog. Default: <code>false</code></p>|
| <tt>GSB_COMPATIBILITY_ENABLE_CF_SHARING</tt> <b>*</b> | boolean | <p>enable-cf-sharing. Set all services to have the Sharable flag so they can be shared across spaces in PCF. Default: <code>false</code></p>|
| <tt>GSB_COMPATIBILITY_ENABLE_EOL_SERVICES</tt> <b>*</b> | boolean | <p>enable-eol-services. Enable broker services that are end of life. Default: <code>false</code></p>|
| <tt>GSB_COMPATIBILITY_ENABLE_GCP_BETA_SERVICES</tt> <b>*</b> | boolean | <p>enable-gcp-beta-services. Enable services that are in GCP Beta. These have no SLA or support policy. Default: <code>true</code></p>|
| <tt>GSB_COMPATIBILITY_ENABLE_GCP_DEPRECATED_SERVICES</tt> <b>*</b> | boolean | <p>enable-gcp-deprecated-services. Enable services that use deprecated GCP components. Default: <code>false</code></p>|
| <tt>GSB_COMPATIBILITY_ENABLE_PREVIEW_SERVICES</tt> <b>*</b> | boolean | <p>enable-preview-services. Enable services that are new to the broker this release. Default: <code>true</code></p>|
| <tt>GSB_COMPATIBILITY_ENABLE_TERRAFORM_SERVICES</tt> <b>*</b> | boolean | <p>enable-terraform-services. Enable services that use the experimental, unstable, Terraform back-end. Default: <code>false</code></p>|
| <tt>GSB_COMPATIBILITY_ENABLE_UNMAINTAINED_SERVICES</tt> <b>*</b> | boolean | <p>enable-unmaintained-services. Enable broker services that are unmaintained. Default: <code>false</code></p>|


\* = Required



## Install Brokerpaks

You can install one or more brokerpaks using the <tt>GSB_BROKERPAK_SOURCES</tt>
environment variable.

The value should be a JSON array containing zero or more brokerpak configuration
objects with the following properties:


| Property | Type | Description |
|----------|------|-------------|
| <tt>uri</tt> <b>*</b> | string | <p>Brokerpak URI.  The URI to load. Supported protocols are http, https, gs, and git. Cloud Storage (gs) URIs follow the gs://<bucket>/<path> convention and will be read using the service broker service account. You can validate the checksum of any file on download by appending a checksum query parameter to the URI in the format type:value. Valid checksum types are MD5, SHA1, SHA256 and SHA512. e.g. gs://foo/bar.brokerpak?checksum=md5:3063a2c62e82ef8614eee6745a7b6b59</p>|
| <tt>service_prefix</tt> | string | <p>Service Prefix. A prefix to prepend to every service name. This will be exact, so you may want to include a trailing dash.</p>|
| <tt>excluded_services</tt> | text | <p>Excluded Services. A list of UUIDs of services to exclude, one per line.</p>|
| <tt>config</tt> <b>*</b> | text | <p>Brokerpak Configuration. A JSON map of configuration key/value pairs for the brokerpak. If a variable isn't found here, it's looked up in the global config. Default: <code>{}</code></p>|
| <tt>notes</tt> | text | <p>Notes. A place for your notes, not used by the broker.</p>|


\* = Required


### Example

Here is an example that loads three brokerpaks.

	[
		{
			"notes":"GA services for all users.",
			"uri":"https://link/to/artifact.brokerpak?checksum=md5:3063a2c62e82ef8614eee6745a7b6b59",
			"excluded_services":"00000000-0000-0000-0000-000000000000",
			"config":{}
		},
		{
			"notes":"Beta services for all users.",
			"uri":"gs://link/to/beta.brokerpak",
			"service_prefix":"beta-",
			"config":{}
		},
		{
			"notes":"Services for the marketing department. They use their own GCP Project.",
			"uri":"https://link/to/marketing.brokerpak",
			"service_prefix":"marketing-",
			"config":{"PROJECT_ID":"my-marketing-project"}
		},
	]

---------------------------------------

_Note: **Do not edit this file**, it was auto-generated by running <code>gcp-service-broker generate customization</code>. If you find an error, change the source code in <tt>customization-md.go</tt> or file a bug._
