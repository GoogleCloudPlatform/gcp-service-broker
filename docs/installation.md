## Installation

Requires Go 1.8 and the associated buildpack.

Documentation for installing as a Pivotal Ops Manager tile is available [here](http://docs.pivotal.io/partners/gcp-sb/index.html)

### Prerequisites

#### Set up a GCP Project <a name="project"></a>

1. go to [Google Cloud Console](https://console.cloud.google.com) and sign up, walking through the setup wizard
1. next to the Google Cloud Platform logo in the upper left-hand corner, click the dropdown and select "Create Project"
1. give your project a name and click "Create"
1. when the project is created (a notification will show in the upper right), refresh the page.

#### Enable APIS <a name="apis"></a>

1. Navigate to **API Manager > Library**.
1. Enable the <a href="https://console.cloud.google.com/apis/api/cloudresourcemanager.googleapis.com/overview">Google Cloud Resource Manager API</a>
1. Enable the <a href="https://console.cloud.google.com/apis/api/iam.googleapis.com/overview">Google Identity and Access Management (IAM) API</a>
1. If you want to enable Cloud SQL as a service, enable the <a href="https://console.cloud.google.com/apis/api/sqladmin/overview">Cloud SQL API</a>
1. If you want to enable BigQuery as a service, enable the <a href="https://console.cloud.google.com/apis/api/bigquery/overview">BigQuery API</a>
1. If you want to enable Cloud Storage as a service, enable the <a href="https://console.cloud.google.com/apis/api/storage_component/overview">Cloud Storage API</a>
1. If you want to enable Pub/Sub as a service, enable the <a href="https://console.cloud.google.com/apis/api/pubsub/overview">Cloud Pub/Sub API</a>
1. If you want to enable Bigtable as a service, enable the <a href="https://console.cloud.google.com/apis/api/bigtableadmin/overview">Bigtable Admin API</a>

#### Create a root service account <a name="service-account"></a>

1. From the GCP console, navigate to **IAM & Admin > Service accounts** and click **Create Service Account**.
1. Enter a **Service account name**.
1. Select the checkbox to **Furnish a new Private Key**, and then click **Create**.
1. Save the automatically downloaded key file to a secure location.
1. Navigate to **IAM & Admin > IAM** and locate your service account.
1. From the dropdown on the right, choose **Project > Owner** and click **Save**.

#### Set up a backing database <a name="database"></a>

1. create new MySQL instance
1. run `CREATE DATABASE servicebroker;`
1. run `CREATE USER '<username>'@'%' IDENTIFIED BY '<password>';`
1. run `GRANT ALL PRIVILEGES ON servicebroker.* TO '<username>'@'%' WITH GRANT OPTION;`
1. (optional) create SSL certs for the database and save them somewhere secure

#### Set required env vars <a name="required-env"></a>

Add these to `manifest.yml`

* `ROOT_SERVICE_ACCOUNT_JSON` (the string version of the credentials file created for the Owner level Service Account)
* `SECURITY_USER_NAME` (a username to sign all service broker requests with - the same one used in cf create-service-broker)
* `SECURITY_USER_PASSWORD` (a password to sign all service broker requests with - the same one used in cf create-service-broker)
* `DB_HOST` (the host for the database to back the service broker)
* `DB_USERNAME` (the database username for the service broker to use)
* `DB_PASSWORD` (the database password for the service broker to use)

#### Optional env vars <a name="optional-env"></a>

optionally add these to `manifest.yml`

* `DB_PORT` (defaults to 3306)
* `DB_NAME` (defaults to "servicebroker")
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
* `BIGTABLE_CUSTOM_PLANS` (A map of plan names to string maps with fields `guid`, `name`, `description`,
`storage_type`, `num_nodes`, `display_name`, and `service` (Bigtable's service id)) - if unset, the service
will be disabled. e.g.,

```json
{
    "bt_plan": {
        "name": "bt_plan",
        "description": "Bigtable basic plan",
        "storage_type": "HDD",
        "num_nodes": "5",
        "display_name": "Bigtable Plan",
        "service": "b8e19880-ac58-42ef-b033-f7cd9c94d1fe"
    }
}
```
* `SPANNER_CUSTOM_PLANS` (A map of plan names to string maps with fields `guid`, `name`, `description`,
`num_nodes` `display_name`, and `service` (Spanner's service id)) - if unset, the service
will be disabled. e.g.,

```json
{
    "spannerplan": {
        "name": "spannerplan",
        "description": "Basic Spanner plan",
        "num_nodes": "15",
        "display_name": "Spanner Plan",
        "service": "51b3e27e-d323-49ce-8c5f-1211e6409e82"
    }
}
```

#### Push the service broker to CF and enable services <a name="push"></a>
1. `cf push gcp-service-broker`
1. `cf create-service-broker <service broker name> <username> <password> <service broker url>`
1. (for all applicable services, e.g.) `cf enable-service-access google-pubsub`

For more information, see the Cloud Foundry docs on [managing Service Brokers](https://docs.cloudfoundry.org/services/managing-service-brokers.html)


#### (Optional) Increase the default provision/bind timeout <a name="timeout"></a>
It is advisable, if you want to use CloudSQL, to increase the default timeout for provision and
bind operations to 90 seconds. CloudFoundry does not, at this point in time, support asynchronous
binding, and CloudSQL bind operations may exceed 60 seconds. To change this setting, set
`broker_client_timeout_seconds` = 90