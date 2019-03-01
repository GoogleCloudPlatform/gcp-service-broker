# Customizing Services

You can customize specific services by changing their defaults, disabling them, or creating custom plans.

The <tt>GSB_SERVICE_CONFIG</tt> environment variable holds all the customizations as a JSON map.
The keys of the map are the service's GUID, and the value is a configuration object.

Example:

	{
		"51b3e27e-d323-49ce-8c5f-1211e6409e82":{ /* Spanner Configuration Object */ },
		"628629e3-79f5-4255-b981-d14c6c7856be":{ /* Pub/Sub Configuration Object */ },
		...
	}

**Configuration Object**

| Property | Type | Description |
|----------|------|-------------|
| <tt>//</tt> | string | Space for your notes. |
| <tt>disabled</tt> | boolean | If set to true, this service will be hidden from the catalog. |
| <tt>provision_defaults</tt> | string:any map | A map of provision property/default pairs that are used to populate missing values in provision requests. |
| <tt>bind_defaults</tt> | string:any map | A map of bind property/default pairs that are used to populate missing values in bind requests. |
| <tt>custom_plans</tt> | array of custom plan objects | You can add custom service plans here. See below for the object structure. |

**Custom Plan Object**

| Property | Type | Description |
|----------|------|-------------|
| <tt>guid</tt> \* | string | A GUID for this plan, must be unique. Changing this value after services are using it WILL BREAK your instances. |
| <tt>name</tt> \* | string | A CLI friendly name for this plan. This can be changed without affecting existing instances, but may break scripts you've previously built referencing it. |
| <tt>display_name</tt> \* | string | A human readable name for this plan, this can be changed. |
| <tt>description</tt> \* | string | A human readable description for this plan, this can be changed. |
| <tt>properties</tt> \* | string:string map | Properties used to configure the plan. Each service has its own set of properties used to customize it. |

\* = Required



---------------------------------------

## Google BigQuery<a id="google-bigquery"></a>

A fast, economical and fully managed data warehouse for large-scale data analytics.

Configuration needs to be done under the GUID: <tt>f80c0a3e-bd4d-4809-a900-b4e33a6450f1</tt>.

#### Example

	{
	  "f80c0a3e-bd4d-4809-a900-b4e33a6450f1": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "See the 'provision defaults' section below for defaults you can change."
	    },
	    "bind_defaults": {
	      "//": "See the 'bind defaults' section below for defaults you can change."
	    },
	    "custom_plans": []
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `name` | string | The name of the BigQuery dataset. |
| `location` | string | The location of the BigQuery instance. |



#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `role` | string | The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. |



#### Custom Plan Properties

_There are no configurable properties for this object._



---------------------------------------

## Google Bigtable<a id="google-bigtable"></a>

A high performance NoSQL database service for large analytical and operational workloads.

Configuration needs to be done under the GUID: <tt>b8e19880-ac58-42ef-b033-f7cd9c94d1fe</tt>.

#### Example

	{
	  "b8e19880-ac58-42ef-b033-f7cd9c94d1fe": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "See the 'provision defaults' section below for defaults you can change."
	    },
	    "bind_defaults": {
	      "//": "See the 'bind defaults' section below for defaults you can change."
	    },
	    "custom_plans": [
	      {
	        "guid": "00000000-0000-0000-0000-000000000000",
	        "name": "a-cli-friendly-name",
	        "display_name": "A human-readable name",
	        "description": "What makes this plan different?",
	        "properties": {
	          "//": "See the custom plan properties section below for configurable properties."
	        },
	        "provision_overrides": {
	          "//": "You can override any user-settable provision variable here."
	        },
	        "bind_overrides": {
	          "//": "You can override any user-settable bind variable here."
	        }
	      }
	    ]
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `name` | string | The name of the Cloud Bigtable instance. |
| `cluster_id` | string | The ID of the Cloud Bigtable cluster. |
| `display_name` | string | The human-readable display name of the Bigtable instance. |
| `zone` | string | The zone to create the Cloud Bigtable cluster in. Zones that support Bigtable instances are noted on the Cloud Bigtable locations page: https://cloud.google.com/bigtable/docs/locations. |



#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `role` | string | The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. |



#### Custom Plan Properties

| Property | Type | Description |
|----------|------|-------------|
| `storage_type` \* | string | Either HDD or SSD. See: https://cloud.google.com/bigtable/pricing for more information. |
| `num_nodes` \* | string | Number of nodes, between 3 and 30. See: https://cloud.google.com/bigtable/pricing for more information. |
\* = Required





---------------------------------------

## Google CloudSQL for MySQL<a id="google-cloudsql-mysql"></a>

Google CloudSQL for MySQL is a fully-managed MySQL database service.

Configuration needs to be done under the GUID: <tt>4bc59b9a-8520-409f-85da-1c7552315863</tt>.

#### Example

	{
	  "4bc59b9a-8520-409f-85da-1c7552315863": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "See the 'provision defaults' section below for defaults you can change."
	    },
	    "bind_defaults": {
	      "//": "See the 'bind defaults' section below for defaults you can change."
	    },
	    "custom_plans": [
	      {
	        "guid": "00000000-0000-0000-0000-000000000000",
	        "name": "a-cli-friendly-name",
	        "display_name": "A human-readable name",
	        "description": "What makes this plan different?",
	        "properties": {
	          "//": "See the custom plan properties section below for configurable properties."
	        },
	        "provision_overrides": {
	          "//": "You can override any user-settable provision variable here."
	        },
	        "bind_overrides": {
	          "//": "You can override any user-settable bind variable here."
	        }
	      }
	    ]
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `instance_name` | string | Name of the Cloud SQL instance. |
| `database_name` | string | Name of the database inside of the instance. Must be a valid identifier for your chosen database type. |
| `version` | string | The database engine type and version. Defaults to `MYSQL_5_6` for 1st gen MySQL instances or `MYSQL_5_7` for 2nd gen MySQL instances. |
| `failover_replica_name` | string | (only for 2nd generation instances) If specified, creates a failover replica with the given name. |
| `activation_policy` | string | The activation policy specifies when the instance is activated; it is applicable only when the instance state is RUNNABLE. |
| `binlog` | string | Whether binary log is enabled. If backup configuration is disabled, binary log must be disabled as well. Defaults: `false` for 1st gen, `true` for 2nd gen, set to `true` to use. |
| `disk_size` | string | In GB (only for 2nd generation instances). |
| `region` | string | The geographical region. See the instance locations list https://cloud.google.com/sql/docs/mysql/instance-locations for which regions support which databases. |
| `zone` | string | (only for 2nd generation instances) |
| `disk_type` | string | (only for 2nd generation instances) |
| `maintenance_window_day` | string | (only for 2nd generation instances) This specifies when a v2 CloudSQL instance should preferably be restarted for system maintenance purposes. Day of week (1-7), starting on Monday. |
| `maintenance_window_hour` | string | (only for 2nd generation instances) The hour of the day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Hour of day 0-23. |
| `backups_enabled` | string | Should daily backups be enabled for the service? |
| `backup_start_time` | string | Start time for the daily backup configuration in UTC timezone in the 24 hour format - HH:MM. |
| `authorized_networks` | string | A comma separated list without spaces. |
| `replication_type` | string | The type of replication this instance uses. This can be either ASYNCHRONOUS or SYNCHRONOUS. |
| `auto_resize` | string | (only for 2nd generation instances) Configuration to increase storage size automatically. |



#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `role` | string | The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. |
| `jdbc_uri_format` | string | If `true`, `uri` field will contain a JDBC formatted URI. |
| `username` | string | The SQL username for the account. |
| `password` | string | The SQL password for the account. |



#### Custom Plan Properties

| Property | Type | Description |
|----------|------|-------------|
| `tier` \* | string | Case-sensitive tier/machine type name (see https://cloud.google.com/sql/pricing for more information). |
| `pricing_plan` \* | string | Select a pricing plan (only for 1st generation instances). |
| `max_disk_size` \* | string | Maximum disk size in GB (applicable only to Second Generation instances, 10 minimum/default). |
\* = Required





---------------------------------------

## Google CloudSQL for PostgreSQL<a id="google-cloudsql-postgres"></a>

Google CloudSQL for PostgreSQL is a fully-managed PostgreSQL database service.

Configuration needs to be done under the GUID: <tt>cbad6d78-a73c-432d-b8ff-b219a17a803a</tt>.

#### Example

	{
	  "cbad6d78-a73c-432d-b8ff-b219a17a803a": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "See the 'provision defaults' section below for defaults you can change."
	    },
	    "bind_defaults": {
	      "//": "See the 'bind defaults' section below for defaults you can change."
	    },
	    "custom_plans": [
	      {
	        "guid": "00000000-0000-0000-0000-000000000000",
	        "name": "a-cli-friendly-name",
	        "display_name": "A human-readable name",
	        "description": "What makes this plan different?",
	        "properties": {
	          "//": "See the custom plan properties section below for configurable properties."
	        },
	        "provision_overrides": {
	          "//": "You can override any user-settable provision variable here."
	        },
	        "bind_overrides": {
	          "//": "You can override any user-settable bind variable here."
	        }
	      }
	    ]
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `instance_name` | string | Name of the CloudSQL instance. |
| `database_name` | string | Name of the database inside of the instance. Must be a valid identifier for your chosen database type. |
| `version` | string | The database engine type and version. |
| `failover_replica_name` | string | (only for 2nd generation instances) If specified, creates a failover replica with the given name. |
| `activation_policy` | string | The activation policy specifies when the instance is activated; it is applicable only when the instance state is RUNNABLE. |
| `binlog` | string | Whether binary log is enabled. If backup configuration is disabled, binary log must be disabled as well. Defaults: `false` for 1st gen, `true` for 2nd gen, set to `true` to use. |
| `disk_size` | string | In GB (only for 2nd generation instances). |
| `region` | string | The geographical region. See the instance locations list https://cloud.google.com/sql/docs/mysql/instance-locations for which regions support which databases. |
| `zone` | string | (only for 2nd generation instances) |
| `disk_type` | string | (only for 2nd generation instances) |
| `maintenance_window_day` | string | (only for 2nd generation instances) This specifies when a v2 CloudSQL instance should preferably be restarted for system maintenance purposes. Day of week (1-7), starting on Monday. |
| `maintenance_window_hour` | string | (only for 2nd generation instances) The hour of the day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Hour of day 0-23. |
| `backups_enabled` | string | Should daily backups be enabled for the service? |
| `backup_start_time` | string | Start time for the daily backup configuration in UTC timezone in the 24 hour format - HH:MM. |
| `authorized_networks` | string | A comma separated list without spaces. |
| `replication_type` | string | The type of replication this instance uses. This can be either ASYNCHRONOUS or SYNCHRONOUS. |
| `auto_resize` | string | (only for 2nd generation instances) Configuration to increase storage size automatically. |



#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `role` | string | The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. |
| `jdbc_uri_format` | string | If `true`, `uri` field will contain a JDBC formatted URI. |
| `username` | string | The SQL username for the account. |
| `password` | string | The SQL password for the account. |



#### Custom Plan Properties

| Property | Type | Description |
|----------|------|-------------|
| `tier` \* | string | A string of the form db-custom-[CPUS]-[MEMORY_MBS], where memory is at least 3840. |
| `pricing_plan` \* | string | The pricing plan. |
| `max_disk_size` \* | string | Maximum disk size in GB, 10 is the minimum. |
\* = Required





---------------------------------------

## Google Cloud Dataflow<a id="google-dataflow"></a>

A managed service for executing a wide variety of data processing patterns built on Apache Beam.

Configuration needs to be done under the GUID: <tt>3e897eb3-9062-4966-bd4f-85bda0f73b3d</tt>.

#### Example

	{
	  "3e897eb3-9062-4966-bd4f-85bda0f73b3d": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "The provision action takes no params so it can't be overridden."
	    },
	    "bind_defaults": {
	      "//": "See the 'bind defaults' section below for defaults you can change."
	    },
	    "custom_plans": []
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `role` | string | The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. |



#### Custom Plan Properties

_There are no configurable properties for this object._



---------------------------------------

## Google Cloud Datastore<a id="google-datastore"></a>

Google Cloud Datastore is a NoSQL document database service.

Configuration needs to be done under the GUID: <tt>76d4abb2-fee7-4c8f-aee1-bcea2837f02b</tt>.

#### Example

	{
	  "76d4abb2-fee7-4c8f-aee1-bcea2837f02b": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "See the 'provision defaults' section below for defaults you can change."
	    },
	    "bind_defaults": {
	      "//": "The bind action takes no params so it can't be overridden."
	    },
	    "custom_plans": []
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `namespace` | string | A context for the identifiers in your entity’s dataset. This ensures that different systems can all interpret an entity's data the same way, based on the rules for the entity’s particular namespace. Blank means the default namespace will be used. |



#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Custom Plan Properties

_There are no configurable properties for this object._



---------------------------------------

## Google Cloud Dialogflow<a id="google-dialogflow"></a>

Dialogflow is an end-to-end, build-once deploy-everywhere development suite for creating conversational interfaces for websites, mobile applications, popular messaging platforms, and IoT devices.

Configuration needs to be done under the GUID: <tt>e84b69db-3de9-4688-8f5c-26b9d5b1f129</tt>.

#### Example

	{
	  "e84b69db-3de9-4688-8f5c-26b9d5b1f129": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "The provision action takes no params so it can't be overridden."
	    },
	    "bind_defaults": {
	      "//": "The bind action takes no params so it can't be overridden."
	    },
	    "custom_plans": []
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Custom Plan Properties

_There are no configurable properties for this object._



---------------------------------------

## Google Cloud Firestore<a id="google-firestore"></a>

Cloud Firestore is a fast, fully managed, serverless, cloud-native NoSQL document database that simplifies storing, syncing, and querying data for your mobile, web, and IoT apps at global scale.

Configuration needs to be done under the GUID: <tt>a2b7b873-1e34-4530-8a42-902ff7d66b43</tt>.

#### Example

	{
	  "a2b7b873-1e34-4530-8a42-902ff7d66b43": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "The provision action takes no params so it can't be overridden."
	    },
	    "bind_defaults": {
	      "//": "See the 'bind defaults' section below for defaults you can change."
	    },
	    "custom_plans": []
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `role` | string | The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. |



#### Custom Plan Properties

_There are no configurable properties for this object._



---------------------------------------

## Google Machine Learning APIs<a id="google-ml-apis"></a>

Machine Learning APIs including Vision, Translate, Speech, and Natural Language.

Configuration needs to be done under the GUID: <tt>5ad2dce0-51f7-4ede-8b46-293d6df1e8d4</tt>.

#### Example

	{
	  "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "The provision action takes no params so it can't be overridden."
	    },
	    "bind_defaults": {
	      "//": "See the 'bind defaults' section below for defaults you can change."
	    },
	    "custom_plans": []
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `role` | string | The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. |



#### Custom Plan Properties

_There are no configurable properties for this object._



---------------------------------------

## Google PubSub<a id="google-pubsub"></a>

A global service for real-time and reliable messaging and streaming data.

Configuration needs to be done under the GUID: <tt>628629e3-79f5-4255-b981-d14c6c7856be</tt>.

#### Example

	{
	  "628629e3-79f5-4255-b981-d14c6c7856be": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "See the 'provision defaults' section below for defaults you can change."
	    },
	    "bind_defaults": {
	      "//": "See the 'bind defaults' section below for defaults you can change."
	    },
	    "custom_plans": []
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `topic_name` | string | Name of the topic. Must not start with "goog". |
| `subscription_name` | string | Name of the subscription. Blank means no subscription will be created. Must not start with "goog". |
| `is_push` | string | Are events handled by POSTing to a URL? |
| `endpoint` | string | If `is_push` == 'true', then this is the URL that will be pushed to. |
| `ack_deadline` | string | Value is in seconds. Max: 600 This is the maximum time after a subscriber receives a message before the subscriber should acknowledge the message. After message delivery but before the ack deadline expires and before the message is acknowledged, it is an outstanding message and will not be delivered again during that time (on a best-effort basis).  |



#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `role` | string | The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. |



#### Custom Plan Properties

_There are no configurable properties for this object._



---------------------------------------

## Google Spanner<a id="google-spanner"></a>

The first horizontally scalable, globally consistent, relational database service.

Configuration needs to be done under the GUID: <tt>51b3e27e-d323-49ce-8c5f-1211e6409e82</tt>.

#### Example

	{
	  "51b3e27e-d323-49ce-8c5f-1211e6409e82": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "See the 'provision defaults' section below for defaults you can change."
	    },
	    "bind_defaults": {
	      "//": "See the 'bind defaults' section below for defaults you can change."
	    },
	    "custom_plans": [
	      {
	        "guid": "00000000-0000-0000-0000-000000000000",
	        "name": "a-cli-friendly-name",
	        "display_name": "A human-readable name",
	        "description": "What makes this plan different?",
	        "properties": {
	          "//": "See the custom plan properties section below for configurable properties."
	        },
	        "provision_overrides": {
	          "//": "You can override any user-settable provision variable here."
	        },
	        "bind_overrides": {
	          "//": "You can override any user-settable bind variable here."
	        }
	      }
	    ]
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `name` | string | A unique identifier for the instance, which cannot be changed after the instance is created. |
| `display_name` | string | The name of this instance configuration as it appears in UIs. |
| `location` | string | A configuration for a Cloud Spanner instance. Configurations define the geographic placement of nodes and their replication and are slightly different from zones. There are single region configurations, multi-region configurations, and multi-continent configurations. See the instance docs https://cloud.google.com/spanner/docs/instances for a list of configurations. |



#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `role` | string | The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. |



#### Custom Plan Properties

| Property | Type | Description |
|----------|------|-------------|
| `num_nodes` \* | string | Number of nodes, a minimum of 3 nodes is recommended for production environments. See: https://cloud.google.com/spanner/pricing for more information. |
\* = Required





---------------------------------------

## Stackdriver Debugger<a id="google-stackdriver-debugger"></a>

Stackdriver Debugger is a feature of the Google Cloud Platform that lets you inspect the state of an application at any code location without using logging statements and without stopping or slowing down your applications. Your users are not impacted during debugging. Using the production debugger you can capture the local variables and call stack and link it back to a specific line location in your source code.

Configuration needs to be done under the GUID: <tt>83837945-1547-41e0-b661-ea31d76eed11</tt>.

#### Example

	{
	  "83837945-1547-41e0-b661-ea31d76eed11": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "The provision action takes no params so it can't be overridden."
	    },
	    "bind_defaults": {
	      "//": "The bind action takes no params so it can't be overridden."
	    },
	    "custom_plans": []
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Custom Plan Properties

_There are no configurable properties for this object._



---------------------------------------

## Stackdriver Monitoring<a id="google-stackdriver-monitoring"></a>

Stackdriver Monitoring provides visibility into the performance, uptime, and overall health of cloud-powered applications.

Configuration needs to be done under the GUID: <tt>2bc0d9ed-3f68-4056-b842-4a85cfbc727f</tt>.

#### Example

	{
	  "2bc0d9ed-3f68-4056-b842-4a85cfbc727f": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "The provision action takes no params so it can't be overridden."
	    },
	    "bind_defaults": {
	      "//": "The bind action takes no params so it can't be overridden."
	    },
	    "custom_plans": []
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Custom Plan Properties

_There are no configurable properties for this object._



---------------------------------------

## Stackdriver Profiler<a id="google-stackdriver-profiler"></a>

Continuous CPU and heap profiling to improve performance and reduce costs.

Configuration needs to be done under the GUID: <tt>00b9ca4a-7cd6-406a-a5b7-2f43f41ade75</tt>.

#### Example

	{
	  "00b9ca4a-7cd6-406a-a5b7-2f43f41ade75": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "The provision action takes no params so it can't be overridden."
	    },
	    "bind_defaults": {
	      "//": "The bind action takes no params so it can't be overridden."
	    },
	    "custom_plans": []
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Custom Plan Properties

_There are no configurable properties for this object._



---------------------------------------

## Stackdriver Trace<a id="google-stackdriver-trace"></a>

Stackdriver Trace is a distributed tracing system that collects latency data from your applications and displays it in the Google Cloud Platform Console. You can track how requests propagate through your application and receive detailed near real-time performance insights.

Configuration needs to be done under the GUID: <tt>c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a</tt>.

#### Example

	{
	  "c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "The provision action takes no params so it can't be overridden."
	    },
	    "bind_defaults": {
	      "//": "The bind action takes no params so it can't be overridden."
	    },
	    "custom_plans": []
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

_There are no configurable properties for this object._

#### Custom Plan Properties

_There are no configurable properties for this object._



---------------------------------------

## Google Cloud Storage<a id="google-storage"></a>

Unified object storage for developers and enterprises. Cloud Storage allows world-wide storage and retrieval of any amount of data at any time.

Configuration needs to be done under the GUID: <tt>b9e4332e-b42b-4680-bda5-ea1506797474</tt>.

#### Example

	{
	  "b9e4332e-b42b-4680-bda5-ea1506797474": {
	    "disabled": false,
	    "provision_defaults": {
	      "//": "See the 'provision defaults' section below for defaults you can change."
	    },
	    "bind_defaults": {
	      "//": "See the 'bind defaults' section below for defaults you can change."
	    },
	    "custom_plans": [
	      {
	        "guid": "00000000-0000-0000-0000-000000000000",
	        "name": "a-cli-friendly-name",
	        "display_name": "A human-readable name",
	        "description": "What makes this plan different?",
	        "properties": {
	          "//": "See the custom plan properties section below for configurable properties."
	        },
	        "provision_overrides": {
	          "//": "You can override any user-settable provision variable here."
	        },
	        "bind_overrides": {
	          "//": "You can override any user-settable bind variable here."
	        }
	      }
	    ]
	  }
	}


_Note: the example includes the configuration and the GUID it should be nested under._

#### Provision Defaults

Setting a value for any of these in the <tt>provision_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `name` | string | The name of the bucket. There is a single global namespace shared by all buckets so it MUST be unique. |
| `location` | string | The location of the bucket. Object data for objects in the bucket resides in physical storage within this region. See: https://cloud.google.com/storage/docs/bucket-locations |



#### Bind Defaults

Setting a value for any of these in the <tt>bind_defaults</tt> map
will override the default value the provision call uses for the property.

| Property | Type | Description |
|----------|------|-------------|
| `role` | string | The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. |



#### Custom Plan Properties

| Property | Type | Description |
|----------|------|-------------|
| `storage_class` \* | string | The storage class of the bucket. See: https://cloud.google.com/storage/docs/storage-classes. |
\* = Required





---------------------------------------

_Note: **Do not edit this file**, it was auto-generated by <tt>service-config-md.go</tt>. If you find an error, change the source code or file a bug._ <nil>
