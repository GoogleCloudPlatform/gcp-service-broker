# Installation Customization

This file documents the various environment variables you can set to change the functionality of the service broker.
If you are using the PCF Tile deployment, then you can manage all of these options through the operator forms.
If you are running your own, then you can set them in the application manifest of a PCF deployment, or in your pod configuration for Kubernetes.


## Root Service Account

Please paste in the contents of the json keyfile (un-encoded) for your service account with owner credentials.

You can configure the following environment variables:

<b><tt>ROOT_SERVICE_ACCOUNT_JSON</tt></b> - <i>text</i> - Root Service Account JSON





<ul>
  <li><b>Required</b></li>
</ul>



## Database Properties

Connection details for the backing database for the service broker.

You can configure the following environment variables:

<b><tt>DB_HOST</tt></b> - <i>string</i> - Database host





<ul>
  <li><b>Required</b></li>
</ul>

<b><tt>DB_USERNAME</tt></b> - <i>string</i> - Database username





<ul>
  <li><i>Optional</i></li>
</ul>

<b><tt>DB_PASSWORD</tt></b> - <i>secret</i> - Database password





<ul>
  <li><i>Optional</i></li>
</ul>

<b><tt>DB_PORT</tt></b> - <i>string</i> - Database port (defaults to 3306)





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>3306</code></li>
</ul>

<b><tt>DB_NAME</tt></b> - <i>string</i> - Database name





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>servicebroker</code></li>
</ul>

<b><tt>CA_CERT</tt></b> - <i>text</i> - Server CA cert





<ul>
  <li><i>Optional</i></li>
</ul>

<b><tt>CLIENT_CERT</tt></b> - <i>text</i> - Client cert





<ul>
  <li><i>Optional</i></li>
</ul>

<b><tt>CLIENT_KEY</tt></b> - <i>text</i> - Client key





<ul>
  <li><i>Optional</i></li>
</ul>



## Enable Services

Enable or disable services.

You can configure the following environment variables:

<b><tt>GSB_SERVICE_GOOGLE_BIGQUERY_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Google BigQuery instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_BIGTABLE_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Google Bigtable instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_CLOUDSQL_MYSQL_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Google CloudSQL MySQL instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_CLOUDSQL_POSTGRES_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Google CloudSQL PostgreSQL instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_DATASTORE_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Google Cloud Datastore instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_DIALOGFLOW_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Google Cloud Dialogflow instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_ML_APIS_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Google Machine Learning APIs instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_PUBSUB_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Google PubSub instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_SPANNER_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Google Spanner instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_STACKDRIVER_DEBUGGER_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Stackdriver Debugger instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_STACKDRIVER_MONITORING_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Stackdriver Monitoring instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_STACKDRIVER_PROFILER_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Stackdriver Profiler instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_STACKDRIVER_TRACE_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Stackdriver Trace instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_STORAGE_ENABLED</tt></b> - <i>boolean</i> - Let the broker create and bind Google Cloud Storage instances.





<ul>
  <li><b>Required</b></li>
  <li>Default: <code>true</code></li>
</ul>



## Role Whitelisting

Enable or disable role whitelisting.

You can configure the following environment variables:

<b><tt>GSB_SERVICE_GOOGLE_BIGQUERY_WHITELIST</tt></b> - <i>string</i> - Role whitelist for Google BigQuery instances.

A comma delimited list of roles (minus the role/ prefix) that can be used when creating bound users for this service.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>bigquery.dataViewer,bigquery.dataEditor,bigquery.dataOwner,bigquery.user,bigquery.jobUser</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_BIGTABLE_WHITELIST</tt></b> - <i>string</i> - Role whitelist for Google Bigtable instances.

A comma delimited list of roles (minus the role/ prefix) that can be used when creating bound users for this service.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>bigtable.user,bigtable.reader,bigtable.viewer</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_CLOUDSQL_MYSQL_WHITELIST</tt></b> - <i>string</i> - Role whitelist for Google CloudSQL MySQL instances.

A comma delimited list of roles (minus the role/ prefix) that can be used when creating bound users for this service.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>cloudsql.editor,cloudsql.viewer,cloudsql.client</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_CLOUDSQL_POSTGRES_WHITELIST</tt></b> - <i>string</i> - Role whitelist for Google CloudSQL PostgreSQL instances.

A comma delimited list of roles (minus the role/ prefix) that can be used when creating bound users for this service.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>cloudsql.editor,cloudsql.viewer,cloudsql.client</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_ML_APIS_WHITELIST</tt></b> - <i>string</i> - Role whitelist for Google Machine Learning APIs instances.

A comma delimited list of roles (minus the role/ prefix) that can be used when creating bound users for this service.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>ml.developer,ml.viewer,ml.modelOwner,ml.modelUser,ml.jobOwner,ml.operationOwner</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_PUBSUB_WHITELIST</tt></b> - <i>string</i> - Role whitelist for Google PubSub instances.

A comma delimited list of roles (minus the role/ prefix) that can be used when creating bound users for this service.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>pubsub.publisher,pubsub.subscriber,pubsub.viewer,pubsub.editor</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_SPANNER_WHITELIST</tt></b> - <i>string</i> - Role whitelist for Google Spanner instances.

A comma delimited list of roles (minus the role/ prefix) that can be used when creating bound users for this service.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>spanner.databaseAdmin,spanner.databaseReader,spanner.databaseUser,spanner.viewer</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_STORAGE_WHITELIST</tt></b> - <i>string</i> - Role whitelist for Google Cloud Storage instances.

A comma delimited list of roles (minus the role/ prefix) that can be used when creating bound users for this service.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>storage.objectCreator,storage.objectViewer,storage.objectAdmin</code></li>
</ul>



## Compatibility

Legacy Compatibility Options

You can configure the following environment variables:

<b><tt>GSB_COMPATIBILITY_ENABLE_INPUT_VALIDATION</tt></b> - <i>boolean</i> - enable-input-validation

Enables validating user input variables against JSON Schema definitions.



<ul>
  <li><i>Optional</i></li>
  <li>Default: <code>true</code></li>
</ul>

<b><tt>GSB_COMPATIBILITY_THREE_TO_FOUR_LEGACY_PLANS</tt></b> - <i>boolean</i> - three-to-four.legacy-plans

Enable compatibility with the GCP Service Broker v3.x. Before version 4.0, each installation generated its own plan UUIDs, after 4.0 they have been standardized. This option installs a compatibility layer which checks if a service is using the correct plan GUID. If the service does not use the correct GUID, the request will fail with a message about how to upgrade.



<ul>
  <li><i>Optional</i></li>
  <li>Default: <code>false</code></li>
</ul>



## Default Overrides

Override the default values your users get when provisioning.

You can configure the following environment variables:

<b><tt>GSB_SERVICE_GOOGLE_BIGQUERY_PROVISION_DEFAULTS</tt></b> - <i>text</i> - Provision default override Google BigQuery instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined provision property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_BIGQUERY_BIND_DEFAULTS</tt></b> - <i>text</i> - Bind default override Google BigQuery instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined bind property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_BIGTABLE_PROVISION_DEFAULTS</tt></b> - <i>text</i> - Provision default override Google Bigtable instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined provision property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_BIGTABLE_BIND_DEFAULTS</tt></b> - <i>text</i> - Bind default override Google Bigtable instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined bind property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_CLOUDSQL_MYSQL_PROVISION_DEFAULTS</tt></b> - <i>text</i> - Provision default override Google CloudSQL MySQL instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined provision property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_CLOUDSQL_MYSQL_BIND_DEFAULTS</tt></b> - <i>text</i> - Bind default override Google CloudSQL MySQL instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined bind property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_CLOUDSQL_POSTGRES_PROVISION_DEFAULTS</tt></b> - <i>text</i> - Provision default override Google CloudSQL PostgreSQL instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined provision property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_CLOUDSQL_POSTGRES_BIND_DEFAULTS</tt></b> - <i>text</i> - Bind default override Google CloudSQL PostgreSQL instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined bind property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_ML_APIS_PROVISION_DEFAULTS</tt></b> - <i>text</i> - Provision default override Google Machine Learning APIs instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined provision property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_ML_APIS_BIND_DEFAULTS</tt></b> - <i>text</i> - Bind default override Google Machine Learning APIs instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined bind property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_PUBSUB_PROVISION_DEFAULTS</tt></b> - <i>text</i> - Provision default override Google PubSub instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined provision property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_PUBSUB_BIND_DEFAULTS</tt></b> - <i>text</i> - Bind default override Google PubSub instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined bind property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_SPANNER_PROVISION_DEFAULTS</tt></b> - <i>text</i> - Provision default override Google Spanner instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined provision property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_SPANNER_BIND_DEFAULTS</tt></b> - <i>text</i> - Bind default override Google Spanner instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined bind property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_STORAGE_PROVISION_DEFAULTS</tt></b> - <i>text</i> - Provision default override Google Cloud Storage instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined provision property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>

<b><tt>GSB_SERVICE_GOOGLE_STORAGE_BIND_DEFAULTS</tt></b> - <i>text</i> - Bind default override Google Cloud Storage instances.

A JSON object with key/value pairs. Keys MUST be the name of a user-defined bind property and values are the alternative default.



<ul>
  <li><b>Required</b></li>
  <li>Default: <code>{}</code></li>
</ul>




## Custom Plans

You can specify custom plans for the following services.
The plans MUST be an array of flat JSON objects stored in their associated environment variable e.g. <code>[{...}, {...},...]</code>.
Each plan MUST have a unique UUID, if you modify the plan the UUID should stay the same to ensure previously provisioned services continue to work.
If you are using the PCF tile, it will generate the UUIDs for you.
DO NOT delete plans, instead you should change their labels to mark them as deprecated.

### Google Bigtable Custom Plans

Generate custom plans for Google Bigtable.
To specify a custom plan manually, create the plan as JSON in a JSON array and store it in the environment variable: <tt>BIGTABLE_CUSTOM_PLANS</tt>.

For example:
<code>
[{"id":"00000000-0000-0000-0000-000000000000", "name": "custom-plan-1", "display_name": setme, "description": setme, "service": setme, "storage_type": setme, "num_nodes": setme},...]
</code>

<table>
<tr>
  <th>JSON Property</th>
  <th>Type</th>
  <th>Label</th>
  <th>Details</th>
</tr>
<tr>
  <td><tt>id</tt></td>
  <td><i>string</i></td>
  <td>Plan UUID</td>
  <td>
    The UUID of the custom plan, use the <tt>uuidgen</tt> CLI command or [uuidgenerator.net](https://www.uuidgenerator.net/) to create one.
    <ul><li><b>Required</b></li></ul>
  </td>
</tr>
<tr>
  <td><tt>name</tt></td>
  <td><i>string</i></td>
  <td>Plan CLI Name</td>
  <td>
    The name of the custom plan used to provision it, must be lower-case, start with a letter a-z and contain only letters, numbers and dashes (-).
    <ul><li><b>Required</b></li></ul>
  </td>
</tr>


<tr>
  <td><tt>display_name</tt></td>
  <td><i>string</i></td>
  <td>Display Name</td>
  <td>
  Name of the plan to be displayed to users.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>description</tt></td>
  <td><i>string</i></td>
  <td>Plan description</td>
  <td>
  The description of the plan shown to users.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>service</tt></td>
  <td><i>dropdown_select</i></td>
  <td>Service</td>
  <td>
  The service this plan is associated with.


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>b8e19880-ac58-42ef-b033-f7cd9c94d1fe</code></li>
  <li>This option _is not_ user configurable. It must be set to the default.</li>
  <li>Valid Values:
  <ul>
    <li><tt>b8e19880-ac58-42ef-b033-f7cd9c94d1fe</tt> - Google Bigtable</li>
  </ul>
  </li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>storage_type</tt></td>
  <td><i>dropdown_select</i></td>
  <td>Storage Type</td>
  <td>
  Either HDD or SSD. See: https://cloud.google.com/bigtable/pricing for more information.


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>SSD</code></li>
  <li>Valid Values:
  <ul>
    <li><tt>SSD</tt> - SSD - Solid-state Drive</li><li><tt>HDD</tt> - HDD - Hard Disk Drive</li>
  </ul>
  </li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>num_nodes</tt></td>
  <td><i>string</i></td>
  <td>Num Nodes</td>
  <td>
  Number of nodes, between 3 and 30. See: https://cloud.google.com/bigtable/pricing for more information.


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>3</code></li>
</ul>


  </td>
</tr>

</table>

### Google CloudSQL MySQL Custom Plans

Generate custom plans for Google CloudSQL MySQL.
To specify a custom plan manually, create the plan as JSON in a JSON array and store it in the environment variable: <tt>CLOUDSQL_MYSQL_CUSTOM_PLANS</tt>.

For example:
<code>
[{"id":"00000000-0000-0000-0000-000000000000", "name": "custom-plan-1", "display_name": setme, "description": setme, "service": setme, "tier": setme, "pricing_plan": setme, "max_disk_size": setme},...]
</code>

<table>
<tr>
  <th>JSON Property</th>
  <th>Type</th>
  <th>Label</th>
  <th>Details</th>
</tr>
<tr>
  <td><tt>id</tt></td>
  <td><i>string</i></td>
  <td>Plan UUID</td>
  <td>
    The UUID of the custom plan, use the <tt>uuidgen</tt> CLI command or [uuidgenerator.net](https://www.uuidgenerator.net/) to create one.
    <ul><li><b>Required</b></li></ul>
  </td>
</tr>
<tr>
  <td><tt>name</tt></td>
  <td><i>string</i></td>
  <td>Plan CLI Name</td>
  <td>
    The name of the custom plan used to provision it, must be lower-case, start with a letter a-z and contain only letters, numbers and dashes (-).
    <ul><li><b>Required</b></li></ul>
  </td>
</tr>


<tr>
  <td><tt>display_name</tt></td>
  <td><i>string</i></td>
  <td>Display Name</td>
  <td>
  Name of the plan to be displayed to users.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>description</tt></td>
  <td><i>string</i></td>
  <td>Plan description</td>
  <td>
  The description of the plan shown to users.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>service</tt></td>
  <td><i>dropdown_select</i></td>
  <td>Service</td>
  <td>
  The service this plan is associated with.


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>4bc59b9a-8520-409f-85da-1c7552315863</code></li>
  <li>This option _is not_ user configurable. It must be set to the default.</li>
  <li>Valid Values:
  <ul>
    <li><tt>4bc59b9a-8520-409f-85da-1c7552315863</tt> - Google CloudSQL MySQL</li>
  </ul>
  </li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>tier</tt></td>
  <td><i>string</i></td>
  <td>Tier</td>
  <td>
  Case-sensitive tier/machine type name (see https://cloud.google.com/sql/pricing for more information).


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>pricing_plan</tt></td>
  <td><i>dropdown_select</i></td>
  <td>Pricing Plan</td>
  <td>
  Select a pricing plan (only for 1st generation instances).


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>PER_USE</code></li>
  <li>Valid Values:
  <ul>
    <li><tt>PER_USE</tt> - Per-Use</li><li><tt>PACKAGE</tt> - Package</li>
  </ul>
  </li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>max_disk_size</tt></td>
  <td><i>string</i></td>
  <td>Max Disk Size</td>
  <td>
  Maximum disk size in GB (applicable only to Second Generation instances, 10 minimum/default).


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>10</code></li>
</ul>


  </td>
</tr>

</table>

### Google CloudSQL PostgreSQL Custom Plans

Generate custom plans for Google CloudSQL PostgreSQL.
To specify a custom plan manually, create the plan as JSON in a JSON array and store it in the environment variable: <tt>CLOUDSQL_POSTGRES_CUSTOM_PLANS</tt>.

For example:
<code>
[{"id":"00000000-0000-0000-0000-000000000000", "name": "custom-plan-1", "display_name": setme, "description": setme, "service": setme, "tier": setme, "pricing_plan": setme, "max_disk_size": setme},...]
</code>

<table>
<tr>
  <th>JSON Property</th>
  <th>Type</th>
  <th>Label</th>
  <th>Details</th>
</tr>
<tr>
  <td><tt>id</tt></td>
  <td><i>string</i></td>
  <td>Plan UUID</td>
  <td>
    The UUID of the custom plan, use the <tt>uuidgen</tt> CLI command or [uuidgenerator.net](https://www.uuidgenerator.net/) to create one.
    <ul><li><b>Required</b></li></ul>
  </td>
</tr>
<tr>
  <td><tt>name</tt></td>
  <td><i>string</i></td>
  <td>Plan CLI Name</td>
  <td>
    The name of the custom plan used to provision it, must be lower-case, start with a letter a-z and contain only letters, numbers and dashes (-).
    <ul><li><b>Required</b></li></ul>
  </td>
</tr>


<tr>
  <td><tt>display_name</tt></td>
  <td><i>string</i></td>
  <td>Display Name</td>
  <td>
  Name of the plan to be displayed to users.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>description</tt></td>
  <td><i>string</i></td>
  <td>Plan description</td>
  <td>
  The description of the plan shown to users.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>service</tt></td>
  <td><i>dropdown_select</i></td>
  <td>Service</td>
  <td>
  The service this plan is associated with.


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>cbad6d78-a73c-432d-b8ff-b219a17a803a</code></li>
  <li>This option _is not_ user configurable. It must be set to the default.</li>
  <li>Valid Values:
  <ul>
    <li><tt>cbad6d78-a73c-432d-b8ff-b219a17a803a</tt> - Google CloudSQL PostgreSQL</li>
  </ul>
  </li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>tier</tt></td>
  <td><i>string</i></td>
  <td>Tier</td>
  <td>
  A string of the form db-custom-[CPUS]-[MEMORY_MBS], where memory is at least 3840.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>pricing_plan</tt></td>
  <td><i>dropdown_select</i></td>
  <td>Pricing Plan</td>
  <td>
  The pricing plan.


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>PER_USE</code></li>
  <li>This option _is not_ user configurable. It must be set to the default.</li>
  <li>Valid Values:
  <ul>
    <li><tt>PER_USE</tt> - Per-Use</li>
  </ul>
  </li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>max_disk_size</tt></td>
  <td><i>string</i></td>
  <td>Max Disk Size</td>
  <td>
  Maximum disk size in GB, 10 is the minimum.


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>10</code></li>
</ul>


  </td>
</tr>

</table>

### Google Spanner Custom Plans

Generate custom plans for Google Spanner.
To specify a custom plan manually, create the plan as JSON in a JSON array and store it in the environment variable: <tt>SPANNER_CUSTOM_PLANS</tt>.

For example:
<code>
[{"id":"00000000-0000-0000-0000-000000000000", "name": "custom-plan-1", "display_name": setme, "description": setme, "service": setme, "num_nodes": setme},...]
</code>

<table>
<tr>
  <th>JSON Property</th>
  <th>Type</th>
  <th>Label</th>
  <th>Details</th>
</tr>
<tr>
  <td><tt>id</tt></td>
  <td><i>string</i></td>
  <td>Plan UUID</td>
  <td>
    The UUID of the custom plan, use the <tt>uuidgen</tt> CLI command or [uuidgenerator.net](https://www.uuidgenerator.net/) to create one.
    <ul><li><b>Required</b></li></ul>
  </td>
</tr>
<tr>
  <td><tt>name</tt></td>
  <td><i>string</i></td>
  <td>Plan CLI Name</td>
  <td>
    The name of the custom plan used to provision it, must be lower-case, start with a letter a-z and contain only letters, numbers and dashes (-).
    <ul><li><b>Required</b></li></ul>
  </td>
</tr>


<tr>
  <td><tt>display_name</tt></td>
  <td><i>string</i></td>
  <td>Display Name</td>
  <td>
  Name of the plan to be displayed to users.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>description</tt></td>
  <td><i>string</i></td>
  <td>Plan description</td>
  <td>
  The description of the plan shown to users.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>service</tt></td>
  <td><i>dropdown_select</i></td>
  <td>Service</td>
  <td>
  The service this plan is associated with.


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>51b3e27e-d323-49ce-8c5f-1211e6409e82</code></li>
  <li>This option _is not_ user configurable. It must be set to the default.</li>
  <li>Valid Values:
  <ul>
    <li><tt>51b3e27e-d323-49ce-8c5f-1211e6409e82</tt> - Google Spanner</li>
  </ul>
  </li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>num_nodes</tt></td>
  <td><i>string</i></td>
  <td>Num Nodes</td>
  <td>
  Number of nodes, a minimum of 3 nodes is recommended for production environments. See: https://cloud.google.com/spanner/pricing for more information.


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>1</code></li>
</ul>


  </td>
</tr>

</table>

### Google Cloud Storage Custom Plans

Generate custom plans for Google Cloud Storage.
To specify a custom plan manually, create the plan as JSON in a JSON array and store it in the environment variable: <tt>STORAGE_CUSTOM_PLANS</tt>.

For example:
<code>
[{"id":"00000000-0000-0000-0000-000000000000", "name": "custom-plan-1", "display_name": setme, "description": setme, "service": setme, "storage_class": setme},...]
</code>

<table>
<tr>
  <th>JSON Property</th>
  <th>Type</th>
  <th>Label</th>
  <th>Details</th>
</tr>
<tr>
  <td><tt>id</tt></td>
  <td><i>string</i></td>
  <td>Plan UUID</td>
  <td>
    The UUID of the custom plan, use the <tt>uuidgen</tt> CLI command or [uuidgenerator.net](https://www.uuidgenerator.net/) to create one.
    <ul><li><b>Required</b></li></ul>
  </td>
</tr>
<tr>
  <td><tt>name</tt></td>
  <td><i>string</i></td>
  <td>Plan CLI Name</td>
  <td>
    The name of the custom plan used to provision it, must be lower-case, start with a letter a-z and contain only letters, numbers and dashes (-).
    <ul><li><b>Required</b></li></ul>
  </td>
</tr>


<tr>
  <td><tt>display_name</tt></td>
  <td><i>string</i></td>
  <td>Display Name</td>
  <td>
  Name of the plan to be displayed to users.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>description</tt></td>
  <td><i>string</i></td>
  <td>Plan description</td>
  <td>
  The description of the plan shown to users.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>service</tt></td>
  <td><i>dropdown_select</i></td>
  <td>Service</td>
  <td>
  The service this plan is associated with.


<ul>
  <li><b>Required</b></li>
  <li>Default: <code>b9e4332e-b42b-4680-bda5-ea1506797474</code></li>
  <li>This option _is not_ user configurable. It must be set to the default.</li>
  <li>Valid Values:
  <ul>
    <li><tt>b9e4332e-b42b-4680-bda5-ea1506797474</tt> - Google Cloud Storage</li>
  </ul>
  </li>
</ul>


  </td>
</tr>

<tr>
  <td><tt>storage_class</tt></td>
  <td><i>string</i></td>
  <td>Storage Class</td>
  <td>
  The storage class of the bucket. See: https://cloud.google.com/storage/docs/storage-classes.


<ul>
  <li><b>Required</b></li>
</ul>


  </td>
</tr>

</table>



---------------------------------------

_Note: **Do not edit this file**, it was auto-generated by running <code>gcp-service-broker generate customization</code>. If you find an error, change the source code in <tt>customization-md.go</tt> or file a bug._
