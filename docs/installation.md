## Installation

* [Installing as an Ops Manager tile](http://docs.pivotal.io/partners/gcp-sb/index.html)
* [Installing as a Cloud Foundry Application](#cf)
    * [Set up a GCP Project](#project)
    * [Enable APIs](#apis)
    * [Create a root service account](#service-account)
    * [Set up a backing database](#database)
    * [Set required env vars](#required-env)
    * [Optional env vars](#optional-env)
    * [Push the service broker to CF and enable services](#push)
    * [(Optional) Increase the default provision/bind timeout](#timeout)


### Installing as an Ops Manager tile

Documentation for installing as an Ops Manager tile is available [here](http://docs.pivotal.io/partners/gcp-sb/index.html).

### [Installing as a Cloud Foundry Application](#cf)

#### [Set up a GCP Project](#project)

1. Go to the [Google Cloud Console](https://console.cloud.google.com) and sign up, walking through the setup wizard.
1. A page then displays with a collection of options. Select "Create Project" option.
1. Give your project a name and click "Create".
1. The dashboard for the newly created project will be displayed.

#### [Enable APIs](#apis)

Enable the following services in **[APIs and services > Library](https://console.cloud.google.com/apis/library)**.

1. Enable the [Google Cloud Resource Manager API](https://console.cloud.google.com/apis/api/cloudresourcemanager.googleapis.com/overview)
1. Enable the [Google Identity and Access Management (IAM) API](https://console.cloud.google.com/apis/api/iam.googleapis.com/overview)
1. If you want to enable CloudSQL as a service, enable the [CloudSQL API]
(https://console.cloud.google.com/apis/library/sql-component.googleapis.com)
1. If you want to enable BigQuery as a service, enable the [BigQuery API](https://console.cloud.google.com/apis/api/bigquery/overview)
1. If you want to enable Cloud Storage as a service, enable the [Cloud Storage API](https://console.cloud.google.com/apis/api/storage_component/overview)
1. If you want to enable Pub/Sub as a service, enable the [Cloud Pub/Sub API](https://console.cloud.google.com/apis/api/pubsub/overview)
1. If you want to enable Bigtable as a service, enable the [Bigtable Admin API](https://console.cloud.google.com/apis/api/bigtableadmin/overview)
1. If you want to enable Datastore as a service, enable the [Datastore API](https://console.cloud.google.com/apis/api/datastore.googleapis.com/overview)

#### [Enable VPC](#vpc)

A Virtual Private Cloud (VPC) network is a virtual version of a physical
network, such as a data center network. It provides private network connectivity
between resources in your project.

The GCP Service Broker supports several services that attach to VPC networks.

To use VPC services, you have to enable [VPC peering](https://cloud.google.com/vpc/docs/vpc-peering).

Follow [these instructions to enable VPC peering](https://cloud.google.com/vpc/docs/using-vpc-peering#creating_a_peering_configuration)
for each network you wish to allow the broker to attach to.

#### [Create a root service account](#service-account)

1. From the GCP console, navigate to **IAM & Admin > Service accounts** and click **Create Service Account**.
1. Enter a **Service account name**.
1. In the **Project Role** dropdown, choose **Project > Owner**.
1. Select the checkbox to **Furnish a new Private Key**, make sure the **JSON** key type is specified.
1. Click **Save** to create the account, key and grant it the owner permission.
1. Save the automatically downloaded key file to a secure location.

#### [Set up a backing database](#database)

The GCP Service Broker stores the state of provisioned resources in a MySQL database.
You may use any database compatible with the MySQL protocol.
We recommend a second generation GCP CloudSQL instance with automatic backups, high availability and automatic maintenance.
The service broker does not require much disk space, but we do recommend an SSD for faster interactions with the broker.

1. Create new MySQL instance.
1. **CloudSQL Only** Make sure that the database can be accessed, add `0.0.0.0/0` as an authorized network.
1. Run `CREATE DATABASE servicebroker;`
1. Run `CREATE USER '<username>'@'%' IDENTIFIED BY '<password>';`
1. Run `GRANT ALL PRIVILEGES ON servicebroker.* TO '<username>'@'%' WITH GRANT OPTION;`
1. **CloudSQL Only** (Optional) create SSL certs for the database and save them somewhere secure.

#### [Set required environment variables](#required-env)

Add these to the `env` section of `manifest.yml`

* `ROOT_SERVICE_ACCOUNT_JSON` - the string version of the credentials file created for the Owner level Service Account.
* `SECURITY_USER_NAME` - the username to authenticate broker requests - the same one used in `cf create-service-broker`.
* `SECURITY_USER_PASSWORD` - the password to authenticate broker requests - the same one used in `cf create-service-broker`.
* `DB_HOST` - the host for the database to back the service broker.
* `DB_USERNAME` - the database username for the service broker to use.
* `DB_PASSWORD` - the database password for the service broker to use.

#### [Optional environment variables](#optional-env)

See [the customization documentation](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/customization.md)
for instructions about providing database name and port overrides, SSL certificates, custom service plans, and more.

#### [Push the service broker to CF and enable services](#push)
1. `cf push gcp-service-broker`
1. `cf create-service-broker <service broker name> <username> <password> <service broker url>`
1. (for all applicable services, e.g.) `cf enable-service-access google-pubsub`

For more information, see the Cloud Foundry docs on [managing Service Brokers](https://docs.cloudfoundry.org/services/managing-service-brokers.html).

#### [(Optional) Increase the default provision/bind timeout](#timeout)
If you want to use CloudSQL, we recommend increasing the default timeout for provision and bind operations to 90 seconds.
This is because CloudFoundry does not yet support asynchronous binding, and CloudSQL bind operations may exceed the default 60 second timeout.

Set `broker_client_timeout_seconds` = 90 in your deployment manifest to change this setting.
