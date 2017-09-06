# Usage

`cf create-service pubsub default foobar`

`cf bind-service myapp foobar -c '{"role": "pubsub.admin"}'`

Notes:

`bind-service` calls require a role, except for Cloud SQL, Stackdriver Debugger, and Stackdriver Trace.
`create-service` calls take the following optional custom parameters, all as strings.


* [PubSub](https://cloud.google.com/pubsub/docs/)
    * Provision
        * `topic_name` (defaults to a generated value)
        * `subscription_name`
        * `is_push` (defaults to `false`, to set use "true")
        * `endpoint` (for when is_push == "true", defaults to `nil`)
        * `ack_deadline` (in seconds, defaults to `10`, max 600)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles)

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "redacted",
             "Name": "redacted",
             "PrivateKeyData": "redacted",
             "UniqueId": "redacted",
             "topic_name": "foobar",
             "subscription_name": "empty_if_not_set",
        }
        ```

* [Cloud Storage](https://cloud.google.com/storage/docs/)
    * Provision
        * `name` (defaults to a generated value)
        * `location` (for options, see https://cloud.google.com/storage/docs/bucket-locations. Defaults to `"us"`)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles)

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "redacted",
             "Name": "redacted",
             "PrivateKeyData": "redacted",
             "UniqueId": "redacted",
             "bucket_name": "foobar",
        }
        ```

* [BigQuery](https://cloud.google.com/bigquery/docs/)
    * Provision
        * `name` (defaults to a generated value)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles), e.g. pubsub.admin

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "redacted",
             "Name": "redacted",
             "PrivateKeyData": "redacted",
             "UniqueId": "redacted",
             "dataset_id": "foobar",
        }
        ```

* [CloudSQL](https://cloud.google.com/sql/docs/)
    * Provision
        * `instance_name` (defaults to a generated value)
        * `database_name` (defaults to a generated value)
        * `version` (defaults to `5.6`)
        * `disk_size`in GB (only for 2nd gen, defaults to `10`)
        * `region` (defaults to `"us-central"`)
        * `zone` (for 2nd gen)
        * `disk_type` (for 2nd gen, defaults to `ssd`)
        * `failover_replica_name` (only for 2nd gen, if specified creates a failover replica, defaults to `""`)
        * `maintenance_window_day` (for 2nd gen only, defaults to `1` (Sunday))
        * `maintenance_window_hour` (for 2nd gen only, defaults to `0`)
        * `backups_enabled` (defaults to `true`, set to "false" to disable)
        * `backup_start_time` (defaults to `"06:00"`)
        * `binlog` (defaults to `false` for 1st gen, true for 2nd gen, set to "true" to use)
        * `activation_policy` (defaults to `on demand`)
        * `replication_type` (defaults to `synchronous`)
        * `auto_resize` (2nd gen only, defaults to `false`, set to "true" to use)
    * Bind
        * `username` (defaults to a generated value)
        * `password` (defaults to a generated value)

        **Example Binding credentials**

        ```
        "credentials": {
             "CaCert": "-----BEGIN CERTIFICATE-----\nredacted\n-----END CERTIFICATE-----",
             "ClientCert": "-----BEGIN CERTIFICATE-----\nredacted\n-----END CERTIFICATE-----",
             "ClientKey": "-----BEGIN RSA PRIVATE KEY-----\redacted\n-----END RSA PRIVATE KEY-----",
             "Password": "unencoded-redacted",
             "Sha1Fingerprint": "redacted",
             "Username": "redacted",
             "database_name": "redacted",
             "host": "255.255.255.255",
             "instance_name": "redacted",
             "last_master_operation_id": "some-guid",
             "uri": "mysql://username:encodedpassword@host/databasename?ssl_mode=required"
        }
        ```

* [ML APIs](https://cloud.google.com/ml/)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles)

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "redacted",
             "Name": "redacted",
             "PrivateKeyData": "redacted",
             "UniqueId": "redacted",
        }
        ```

* [Bigtable](https://cloud.google.com/bigtable/docs/)
    * Provison
        * `name` (defaults to a generated value)
        * `cluster_id` (defaults to a generated value)
        * `display_name` (defaults to a generated value)
        * `storage_type` (one of "SSD" or "HDD", defaults to `"SSD"`)
        * `zone` (defaults to `"us-east1-b"`)
        * `num_nodes` (defaults to `3`)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles), e.g. editor

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "redacted",
             "Name": "redacted",
             "PrivateKeyData": "redacted",
             "UniqueId": "redacted",
             "instance_id": "foobar",
        }
        ```

* [Spanner](https://cloud.google.com/spanner/docs/) (BETA Google Service)
    * Provison
        * `name` (defaults to a generated value)
        * `display_name` (defaults to a generated value)
        * `location` (defaults to `"regional-us-central1"`)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles), e.g. spanner.admin

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "redacted",
             "Name": "redacted",
             "PrivateKeyData": "redacted",
             "UniqueId": "redacted",
             "instance_id": "foobar",
        }
        ```

* [Stackdriver Debugger](https://cloud.google.com/debugger/)
    * Provison (none)
    * Bind (none)
	* provided credentials will have the role of `clouddebugger.agent`

* [Stackdriver Trace](https://cloud.google.com/trace/)
    * Provison (none)
    * Bind (none)
	* provided credentials will have the role of `cloudtrace.agent`
	
* [Datastore](https://cloud.google.com/datastore/)
    * Provison (none)
    * Bind (none)
	* provided credentials will have the role of `datastore.user`	
