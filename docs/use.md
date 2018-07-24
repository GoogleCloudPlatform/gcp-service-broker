# Usage

`cf create-service pubsub default foobar`

`cf bind-service myapp foobar -c '{"role": "pubsub.admin"}'`

Notes:

`bind-service` calls require a role, except for Cloud SQL, Stackdriver Debugger, and Stackdriver Trace.
`create-service` calls take the following optional custom parameters, all as strings.


* [google-pubsub](https://cloud.google.com/pubsub/docs/)
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
             "Email": "pcf-binding-abc123@projectid.iam.gserviceaccount.com",
             "Name": "pcf-binding-abc123",
             "PrivateKeyData": "redacted",
             "ProjectId": "projectid",
             "UniqueId": "12345",
             "subscription_name": "empty_if_not_set",
             "topic_name": "pcf_sb_1_123456",
        }
        ```

* [google-storage](https://cloud.google.com/storage/docs/)
    * Provision
        * `name` (defaults to a generated value)
        * `location` (for options, see https://cloud.google.com/storage/docs/bucket-locations. Defaults to `"us"`)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles)

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "pcf-binding-abc123@projectid.iam.gserviceaccount.com",
             "Name": "pcf-binding-abc123",
             "PrivateKeyData": "redacted",
             "ProjectId": "projectid",
             "UniqueId": "12345",
             "bucket_name": "pcf_sb_1_123456",
        }
        ```

* [google-bigquery](https://cloud.google.com/bigquery/docs/)
    * Provision
        * `name` (defaults to a generated value)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles), e.g. pubsub.admin

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "pcf-binding-abc123@projectid.iam.gserviceaccount.com",
             "Name": "pcf-binding-abc123",
             "PrivateKeyData": "redacted",
             "ProjectId": "projectid",
             "UniqueId": "12345",
             "dataset_id": "pcf_sb_1_123456",
        }
        ```

* [google-cloudsql-mysql and google-cloudsql-postgres](https://cloud.google.com/sql/docs/)
    * Provision
        * `instance_name` (defaults to a generated value)
        * `database_name` (defaults to a generated value)
        * `version` (defaults to `MYSQL_5_6` for 1st gen MySQL instances, `MYSQL_5_7` for 2nd gen MySQL instances, or `POSTGRES_9_6` for PostgreSQL instances)
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
        * `authorized_networks` (a comma separated list without spaces, defaults to none)
        * `replication_type` (defaults to `synchronous`)
        * `auto_resize` (2nd gen only, defaults to `false`, set to "true" to use)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles)
        * `username` (defaults to a generated value)
        * `password` (defaults to a generated value)
	* `jdbc_uri_format` (if `true`, `uri` field will contain a jdbc formatted uri, defaults to false)

        **Example Binding credentials**

        ```
        "credentials": {
             "CaCert": "-----BEGIN CERTIFICATE-----\nredacted\n-----END CERTIFICATE-----",
             "ClientCert": "-----BEGIN CERTIFICATE-----\nredacted\n-----END CERTIFICATE-----",
             "ClientKey": "-----BEGIN RSA PRIVATE KEY-----\redacted\n-----END RSA PRIVATE KEY-----",
             "Email": "pcf-binding-abc123@projectid.iam.gserviceaccount.com",
             "Password": "unencoded-redacted",
             "PrivateKeyData": "redacted",
             "ProjectId": "projectid",
             "Sha1Fingerprint": "redacted",
             "UniqueId": "12345",
             "UriPrefix": "empty_if_not_set",
             "Username": "aaa-bbb-c",
             "database_name": "pcf_sb_2_654321",
             "host": "255.255.255.255",
             "instance_name": "pcf_sb_1_123456",
             "last_master_operation_id": "some-guid",
             "region": "us-central",
             "uri": "mysql://username:encodedpassword@host/databasename?ssl_mode=required"
        }
        ```

* [google-ml-apis](https://cloud.google.com/ml/)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles)

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "pcf-binding-abc123@projectid.iam.gserviceaccount.com",
             "Name": "pcf-binding-abc123",
             "PrivateKeyData": "redacted",
             "ProjectId": "projectid",
             "UniqueId": "12345",
        }
        ```

* [google-bigtable](https://cloud.google.com/bigtable/docs/)
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
             "Email": "pcf-binding-abc123@projectid.iam.gserviceaccount.com",
             "Name": "pcf-binding-abc123",
             "PrivateKeyData": "redacted",
             "ProjectId": "projectid",
             "UniqueId": "12345",
             "instance_id": "pcf_sb_1_123456",
        }
        ```

* [google-spanner](https://cloud.google.com/spanner/docs/) (BETA Google Service)
    * Provison
        * `name` (defaults to a generated value)
        * `display_name` (defaults to a generated value)
        * `location` (defaults to `"regional-us-central1"`)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles), e.g. spanner.admin

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "pcf-binding-abc123@projectid.iam.gserviceaccount.com",
             "Name": "pcf-binding-abc123",
             "PrivateKeyData": "redacted",
             "ProjectId": "projectid",
             "UniqueId": "12345",
             "instance_id": "pcf_sb_1_123456",
        }
        ```

* [google-stackdriver-debugger](https://cloud.google.com/debugger/)
    * Provison (none)
    * Bind (none)
	* provided credentials will have the role of `clouddebugger.agent`

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "pcf-binding-abc123@projectid.iam.gserviceaccount.com",
             "Name": "pcf-binding-abc123",
             "PrivateKeyData": "redacted",
             "ProjectId": "projectid",
             "UniqueId": "12345",
        }
        ```

* [google-stackdriver-trace](https://cloud.google.com/trace/)
    * Provison (none)
    * Bind (none)
	* provided credentials will have the role of `cloudtrace.agent`

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "pcf-binding-abc123@projectid.iam.gserviceaccount.com",
             "Name": "pcf-binding-abc123",
             "PrivateKeyData": "redacted",
             "ProjectId": "projectid",
             "UniqueId": "12345",
        }
        ```

* [google-datastore](https://cloud.google.com/datastore/)
    * Provison (none)
    * Bind (none)
	* provided credentials will have the role of `datastore.user`

        **Example Binding credentials**

        ```
        "credentials": {
             "Email": "pcf-binding-abc123@projectid.iam.gserviceaccount.com",
             "Name": "pcf-binding-abc123",
             "PrivateKeyData": "redacted",
             "ProjectId": "projectid",
             "UniqueId": "12345",
        }
        ```

------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/storage.svg) Google Cloud Storage

A Powerful, Simple and Cost Effective Object Storage Service

 * [Documentation](https://cloud.google.com/storage/docs/overview)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `b9e4332e-b42b-4680-bda5-ea1506797474`
 * Tags: gcp, storage

## Provisioning

* Request Parameters
    * `name` _string_ - The name of the bucket. There is a single global namespace shared by all buckets so it MUST be unique. Default: `a generated value`
    * `location` _string_ - The location of the bucket. Object data for objects in the bucket resides in physical storage within this region. See https://cloud.google.com/storage/docs/bucket-locations Default: `us`


## Binding

 * Request Parameters
    * `role` _string_ - **Required** The role for the account without the "roles/" prefix. See https://cloud.google.com/iam/docs/understanding-roles for available roles.

 * Response Parameters
    * `Email` _string_ - Email address of the service account
    * `Name` _string_ - The name of the service account
    * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
    * `ProjectId` _string_ - ID of the project that owns the service account
    * `UniqueId` _string_ - Unique and stable id of the service account
    * `bucket_name` _string_ - Name of the bucket this binding is for


## Plans

  * **standard**: Standard storage class - Plan ID: `e1d11f65-da66-46ad-977c-6d56513baf43`
  * **nearline**: Nearline storage class - Plan ID: `a42c1182-d1a0-4d40-82c1-28220518b360`
  * **reduced_availability**: Durable Reduced Availability storage class - Plan ID: `1a1f4fe6-1904-44d0-838c-4c87a9490a6b`


## Examples

### Basic Configuration

Create a nearline bucket with a service account that can create/read/delete the objects in it.
Uses plan: `a42c1182-d1a0-4d40-82c1-28220518b360`

**Provision**
```.json
{
    "location": "us"
}
```


**Bind**

```.json
{
    "role": "storage.objectAdmin"
}
```
