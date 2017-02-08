# Cloud Foundry Service Broker for Google Cloud Platform

Depends on [lager](https://github.com/pivotal-golang/lager) and [gorilla/mux](https://github.com/gorilla/mux).

Requires Go 1.6 and the associated buildpack.

## Examples

See the [examples](https://github.com/GoogleCloudPlatform/gcp-service-broker/tree/master/examples/) folder.

## Prerequisites

### Set up a GCP Project

1. go to [Google Cloud Console](https://console.cloud.google.com) and sign up, walking through the setup wizard
1. next to the Google Cloud Platform logo in the upper left-hand corner, click the dropdown and select "Create Project"
1. give your project a name and click "Create"
1. when the project is created (a notification will show in the upper right), refresh the page.

### Enable APIS

1. Navigate to **API Manager > Library**.
1. Enable the <a href="https://console.cloud.google.com/apis/api/cloudresourcemanager.googleapis.com/overview">Google Cloud Resource Manager API</a>
1. Enable the <a href="https://console.cloud.google.com/apis/api/iam.googleapis.com/overview">Google Identity and Access Management (IAM) API</a>
1. If you want to enable Cloud SQL as a service, enable the <a href="https://console.cloud.google.com/apis/api/sqladmin/overview">Cloud SQL API</a>
1. If you want to enable BigQuery as a service, enable the <a href="https://console.cloud.google.com/apis/api/bigquery/overview">BigQuery API</a>
1. If you want to enable Cloud Storage as a service, enable the <a href="https://console.cloud.google.com/apis/api/storage_component/overview">Cloud Storage API</a>
1. If you want to enable Pub/Sub as a service, enable the <a href="https://console.cloud.google.com/apis/api/pubsub/overview">Cloud Pub/Sub API</a>
1. If you want to enable Bigtable as a service, enable the <a href="https://console.cloud.google.com/apis/api/bigtableadmin/overview">Bigtable Admin API</a>

### Create a root service account

1. From the GCP console, navigate to **IAM & Admin > Service accounts** and click **Create Service Account**.
1. Enter a **Service account name**.
1. Select the checkbox to **Furnish a new Private Key**, and then click **Create**.
1. Save the automatically downloaded key file to a secure location.
1. Navigate to **IAM & Admin > IAM** and locate your service account.
1. From the dropdown on the right, choose **Project > Owner** and click **Save**.

### Set up a backing database

1. create new MySQL instance
1. run `CREATE DATABASE servicebroker;`
1. run `CREATE USER '<username>'@'%' IDENTIFIED BY '<password>';`
1. run `GRANT ALL PRIVILEGES ON servicebroker.* TO '<username>'@'%' WITH GRANT OPTION;`
1. (optional) create SSL certs for the database and save them somewhere secure

### Set required env vars - if deploying as an app, add these to missing-properties.yml

* `ROOT_SERVICE_ACCOUNT_JSON` (the string version of the credentials file created for the Owner level Service Account)
* `SECURITY_USER_NAME` (a username to sign all service broker requests with - the same one used in cf create-service-broker)
* `SECURITY_USER_PASSWORD` (a password to sign all service broker requests with - the same one used in cf create-service-broker)
* `DB_HOST` (the host for the database to back the service broker)
* `DB_USERNAME` (the database username for the service broker to use)
* `DB_PASSWORD` (the database password for the service broker to use)

### optional env vars - if deploying as an app, optionally add these to missing-properties.yml

* `DB_PORT` (defaults to 3306)
* `CA_CERT`
* `CLIENT_CERT`
* `CLIENT_KEY`
* `CLOUDSQL_CUSTOM_PLANS` (A map of plan names to string maps with fields `guid`, `name`, `description`, `tier`,
`pricing_plan`, `max_disk_size`, `display_name`, and `service` (Cloud SQL's service id)) - if unset, the service
will be disabled. e.g.,

```json
{
    "test_plan": {
        "name": "test_plan",
        "description": "testplan",
        "tier": "D8",
        "pricing_plan": "PER_USE",
        "max_disk_size": "15",
        "display_name": "FOOBAR",
        "service": "4bc59b9a-8520-409f-85da-1c7552315863"
    }
}
```


## Usage

### As an App

#### Update the manifest with ENV vars
1. replace any blank variables that are in `manifest.yml` with your own ENV vars

#### Push the service broker to CF and enable services
1. `cf push gcp-service-broker`
1. `cf create-service-broker <service broker name> <username> <password> <service broker url>`
1. (for all applicable services, e.g.) `cf enable-service-access google-pubsub`

### As a Tile

#### Import the product into Ops Manager
1. Click "Import a Product" and upload the .pivotal file from the product directory

#### Add the product to your Dashboard
1. Click the plus icon next to the uploaded product

#### Configure the Service Broker
1. Click on the tile and fill in any required fields (tabs will be orange if updates are needed)
1. Once the tile is green and updates are applied, review the service/plan access and
update if necessary using cf disable-service-access. By default, all services and plans
are enabled except CloudSQL (unless plans have been saved for it). If you wish to change this,
you'll need to use the cf cli's service-access commands.

### (Optional) Increase the default provision/bind timeout
It is advisable, if you want to use CloudSQL, to increase the default timeout for provision and
bind operations to 90 seconds. CloudFoundry does not, at this point in time, support asynchronous
binding, and CloudSQL bind operations may exceed 60 seconds. To change this setting, set
`broker_client_timeout_seconds` = 90

### Use!

For example:

* `cf create-service pubsub default foobar`
* `cf bind-service myapp foobar -c '{"role": "pubsub.admin"}'`

Notes:

* `create-service` calls take the following optional custom parameters, all as strings.
* `bind-service` calls require a role, except for Cloud SQL.

* [PubSub](https://cloud.google.com/pubsub/docs/)
    * Provision
        * `topic_name` (defaults to a generated value)
        * `subscription_name`
        * `is_push` (defaults to false, to set use "true")
        * `endpoint` (for when is_push == "true", defaults to nil)
        * `ack_deadline` (in seconds, defaults to 10, max 600)
    * Bind
        * `role`without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles)

        **Example Binding credentials**

        ```json
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
        * `location` (for options, see https://cloud.google.com/storage/docs/bucket-locations. Defaults to us)
    * Bind
        * `role`without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles)

        **Example Binding credentials**

        ```json
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
        * `role`without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles), e.g. pubsub.admin

        **Example Binding credentials**

        ```json
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
        * `version` (defaults to 5.6)
        * `disk_size`in GB (only for 2nd gen, defaults to 10)
        * `region` (defaults to us-central)
        * `zone` (for 2nd gen)
        * `disk_type` (for 2nd gen, defaults to ssd)
        * `failover_replica_name` (only for 2nd gen, if specified creates a failover replica, defaults to "")
        * `maintenance_window_day` (for 2nd gen only, defaults to 1 (Sunday))
        * `maintenance_window_hour` (for 2nd gen only, defaults to 0)
        * `backups_enabled` (defaults to true, set to "false" to disable)
        * `backup_start_time` (defaults to 06:00)
        * `binlog` (defaults to false for 1st gen, true for 2nd gen, set to "true" to use)
        * `activation_policy` (defaults to on demand)
        * `replication_type` (defaults to synchronous)
        * `auto_resize` (2nd gen only, defaults to false, set to "true" to use)
    * Bind
        * `username` (defaults to a generated value)
        * `password` (defaults to a generated value)

        **Example Binding credentials**

        ```json
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

        ```json
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
        * `storage_type` (one of "SSD" or "HDD", defaults to "SSD")
        * `zone` (defaults to us-east1-b)
        * `num_nodes` (defaults to 3)
    * Bind
        * `role` without "roles/" prefix (see https://cloud.google.com/iam/docs/understanding-roles for available roles), e.g. editor

        **Example Binding credentials**

        ```json
        "credentials": {
             "Email": "redacted",
             "Name": "redacted",
             "PrivateKeyData": "redacted",
             "UniqueId": "redacted",
             "instance_id": "foobar",
        }
        ```

## Change Notes

see https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/CHANGELOG.md

## Support

For functional issues with the service broker or feature requests, please file a github issue here:

https://github.com/GoogleCloudPlatform/gcp-service-broker/issues

They will be prioritized and updated here:

https://github.com/GoogleCloudPlatform/gcp-service-broker/projects/1

For discussions and updates, please subscribe to this group:

https://groups.google.com/forum/#!forum/gcp-service-broker

# This is not an official Google product.
