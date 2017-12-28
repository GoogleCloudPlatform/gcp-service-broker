# Cloud Foundry Service Broker for Google Cloud Platform

This is the home of the Cloud Foundry Service Broker for Google Cloud Platform. For a demo of installing and using the broker,
see [here](https://www.youtube.com/watch?v=8nc4624K91A&list=PLIivdWyY5sqKJ48ycao632rEDuVbFm8yJ&index=3)

## Background

### Service Brokers

This product is a [Cloud Foundry Service Broker](https://docs.cloudfoundry.org/services/overview.html). It adheres
to [v2.8](https://docs.pivotal.io/pivotalcf/1-7/services/api.html) of the Service Broker API.

### Google Cloud Platform (GCP)

[GCP](cloud.google.com) is a cloud service provider. In addition to VMs and networking, many other useful services are available. The ones
available through this Service Broker are:

* [BigQuery](https://cloud.google.com/bigquery/)
* [Bigtable](https://cloud.google.com/bigtable/)
* [Cloud SQL](https://cloud.google.com/sql/)
* [Cloud Storage](https://cloud.google.com/storage/)
* [Datastore](https://cloud.google.com/datastore/)
* [ML APIs](https://cloud.google.com/ml/)
* [PubSub](https://cloud.google.com/pubsub/)
* [Spanner](https://cloud.google.com/spanner/)
* [Stackdriver Debugger](https://cloud.google.com/debugger/)
* [Stackdriver Trace](https://cloud.google.com/trace/)


## Installation

Requires Go 1.8 and the associated buildpack.

* [Installing as a Pivotal Ops Manager tile](http://docs.pivotal.io/partners/gcp-sb/index.html)
* [Installing as a Cloud Foundry Application](#cf)
    * [Set up a GCP Project](#project)
    * [Enable APIs](#apis)
    * [Create a root service account](#service-account)
    * [Set up a backing database](#database)
    * [Set required env vars](#required-env)
    * [Optional env vars](#optional-env)
    * [Push the service broker to CF and enable services](#push)
    * [(Optional) Increase the default provision/bind timeout](#timeout)


### Installing as a Pivotal Ops Manager tile

Documentation for installing as a Pivotal Ops Manager tile is available [here](http://docs.pivotal.io/partners/gcp-sb/index.html)

### [Installing as a Cloud Foundry Application](#cf)

#### [Set up a GCP Project](#project)

1. Go to [Google Cloud Console](https://console.cloud.google.com) and sign up, walking through the setup wizard
1. Next to the Google Cloud Platform logo in the upper left-hand corner, click the dropdown and select "Create Project"
1. Give your project a name and click "Create"
1. When the project is created (a notification will show in the upper right), refresh the page.

#### [Enable APIs](#apis)

Enable the following services in **[API Manager > Library](https://console.cloud.google.com/apis/library)**.

1. Enable the [Google Cloud Resource Manager API](https://console.cloud.google.com/apis/api/cloudresourcemanager.googleapis.com/overview)
1. Enable the [Google Identity and Access Management (IAM) API](https://console.cloud.google.com/apis/api/iam.googleapis.com/overview)
1. If you want to enable Cloud SQL as a service, enable the [Cloud SQL API](https://console.cloud.google.com/apis/api/sqladmin/overview)
1. If you want to enable BigQuery as a service, enable the [BigQuery API](https://console.cloud.google.com/apis/api/bigquery/overview)
1. If you want to enable Cloud Storage as a service, enable the [Cloud Storage API](https://console.cloud.google.com/apis/api/storage_component/overview)
1. If you want to enable Pub/Sub as a service, enable the [Cloud Pub/Sub API](https://console.cloud.google.com/apis/api/pubsub/overview)
1. If you want to enable Bigtable as a service, enable the [Bigtable Admin API](https://console.cloud.google.com/apis/api/bigtableadmin/overview)
1. If you want to enable Datastore as a service, enable the [Datastore API](https://console.cloud.google.com/apis/api/datastore.googleapis.com/overview)

#### [Create a root service account](#service-account)

1. From the GCP console, navigate to **IAM & Admin > Service accounts** and click **Create Service Account**.
1. Enter a **Service account name**.
1. Select the checkbox to **Furnish a new Private Key**, and then click **Create**.
1. Save the automatically downloaded key file to a secure location.
1. Navigate to **IAM & Admin > IAM** and locate your service account.
1. From the dropdown on the right, choose **Project > Owner** and click **Save**.

#### [Set up a backing database](#database)

1. Create new MySQL instance
1. Make sure that the database can be accessed, if you are using GCP cloudsql, add `0.0.0.0/0` as an authorized network.
1. Run `CREATE DATABASE servicebroker;`
1. Run `CREATE USER '<username>'@'%' IDENTIFIED BY '<password>';`
1. Run `GRANT ALL PRIVILEGES ON servicebroker.* TO '<username>'@'%' WITH GRANT OPTION;`
1. (Optional) create SSL certs for the database and save them somewhere secure

#### [Set required env vars](#required-env)

Add these to the env section of `manifest.yml`

* `ROOT_SERVICE_ACCOUNT_JSON` (the string version of the credentials file created for the Owner level Service Account)
* `SECURITY_USER_NAME` (a username to sign all service broker requests with - the same one used in cf create-service-broker)
* `SECURITY_USER_PASSWORD` (a password to sign all service broker requests with - the same one used in cf create-service-broker)
* `DB_HOST` (the host for the database to back the service broker)
* `DB_USERNAME` (the database username for the service broker to use)
* `DB_PASSWORD` (the database password for the service broker to use)

#### [Optional env vars](#optional-env)

See https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/customization.md 
for instructions on providing database name and port overrides, ssl certs, and custom service plans for Cloud SQL, Bigtable, and Spanner.

#### [Push the service broker to CF and enable services](#push)
1. `cf push gcp-service-broker`
1. `cf create-service-broker <service broker name> <username> <password> <service broker url>`
1. (for all applicable services, e.g.) `cf enable-service-access google-pubsub`

For more information, see the Cloud Foundry docs on [managing Service Brokers](https://docs.cloudfoundry.org/services/managing-service-brokers.html)

#### [(Optional) Increase the default provision/bind timeout](#timeout)
It is advisable, if you want to use CloudSQL, to increase the default timeout for provision and
bind operations to 90 seconds. CloudFoundry does not, at this point in time, support asynchronous
binding, and CloudSQL bind operations may exceed 60 seconds. To change this setting, set
`broker_client_timeout_seconds` = 90 in your deployment manifest.

## Usage

See https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/use.md for instructions on creating and binding to GCP Services
 
See the [examples](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/examples/) folder to understand how to use services once they are created and bound.

## Commands

The [cmd](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/cmd/) folder contains commands that can be run independent of the broker.

* `migrate`: migrates the database to the latest schema

## Testing

Production testing for the GCP Service Broker is administered via a private Concourse pipeline.

To run tests locally, use [Ginkgo](https://onsi.github.io/ginkgo/). 

Integration tests require the `ROOT_SERVICE_ACCOUNT_JSON` environment variable to be set.
 
**Note: Integration tests create and destroy real project resources and therefore have associated costs to run**


## Change Notes

see https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/CHANGELOG.md

## Support

For functional issues with the service broker or feature requests, please file a github issue here:

https://github.com/GoogleCloudPlatform/gcp-service-broker/issues

They will be prioritized and updated here:

https://github.com/GoogleCloudPlatform/gcp-service-broker/projects/1

For discussions and updates, please subscribe to this group:

https://groups.google.com/forum/#!forum/gcp-service-broker

## Contributing

see https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/CONTRIBUTING

# This is not an official Google product.
