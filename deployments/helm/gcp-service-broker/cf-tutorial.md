# Install the Service Broker into Kubernetes for use with CF

## Introduction

This tutorial will walk you through installing the GCP Service Broker into
Kubernetes using the Helm release and using it with Cloud Foundry.

**Time to complete**: About 10 minutes

**Prerequisites**:

* A Kubernetes cluster you want to install the broker into.

Click the **Continue** button to move to the next step.

## Create a Service Account for the broker

<walkthrough-watcher-constant key="service-account-name" value="gcp-service-broker">
</walkthrough-watcher-constant>


Create the service account:

    gcloud iam service-accounts create {{service-account-name}}

Create the environment variable for its email address:

    ACCOUNTEMAIL={{service-account-name}}@{{project-id}}.iam.gserviceaccount.com

Create new credentials to let the broker authenticate:

    gcloud iam service-accounts keys create key.json --iam-account $ACCOUNTEMAIL

Grant project owner permissions to the broker:

    gcloud projects add-iam-policy-binding {{project-id}} --member serviceAccount:$ACCOUNTEMAIL --role "roles/owner"

## Enable Required APIs

Now you need to enable APIs to let the broker provision those kind of resources.

The broker has a few APIs that are required for it to run, and a few that are
optional but must be enabled to provision resources of a particular type.

Enable the following services to allow the service broker to run:

<walkthrough-enable-apis apis="cloudresourcemanager.googleapis.com,iam.googleapis.com">
1. [Google Cloud Resource Manager API](https://console.cloud.google.com/apis/api/cloudresourcemanager.googleapis.com/overview)
1. [Google Identity and Access Management (IAM) API](https://console.cloud.google.com/apis/api/iam.googleapis.com/overview)
</walkthrough-enable-apis>


### Enable Service APIs

The following APIs must be enabled to use their respective services.
For example, you must enable the BigQuery API on the project if you want to
provision and use BigQuery instances.
It doesn't cost anything to enable them, so we recommend enabling them all unless you have a particular reason not to.

1. [BigQuery API](https://console.cloud.google.com/apis/api/bigquery/overview)
1. [BigTable API](https://console.cloud.google.com/apis/api/bigtableadmin/overview)
1. [CloudSQL API](https://console.cloud.google.com/apis/library/sql-component.googleapis.com)
1. [Datastore API](https://console.cloud.google.com/apis/api/datastore.googleapis.com/overview)
1. [Pub/Sub API](https://console.cloud.google.com/apis/api/pubsub/overview)
1. [Redis API](https://console.cloud.google.com/apis/api/redis.googleapis.com/overview)
1. [Storage API](https://console.cloud.google.com/apis/api/storage_component/overview)

## Install the Broker

First, update the dependencies of the helm chart:

    helm dependency update

Next, modify the `values.yaml` file.

<walkthrough-editor-open-file filePath="values.yaml" text="Open values.yaml">
</walkthrough-editor-open-file>

1. Set the value `broker.service_account_json` to the contents of `key.json`.
2. Set the value `svccat.register` to be `false` because you're using this
   installation with Cloud Foundry rather than the
3. **Optional:** read through the rest of the properties and change any you need
   to fit your environment.

Finally, install the broker:

    helm install .

## Set up CF

1. `cf create-service-broker <service broker name> <username> <password> <service broker url>`
1. (for all applicable services, e.g.) `cf enable-service-access google-pubsub`

For more information, see the Cloud Foundry docs on [managing Service Brokers](https://docs.cloudfoundry.org/services/managing-service-brokers.html).


See [the customization documentation](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/customization.md)
for instructions about providing database name and port overrides, SSL certificates, custom service plans, and more.

#### [(Optional) Increase the default provision/bind timeout](#timeout)
If you want to use CloudSQL, we recommend increasing the default timeout for provision and bind operations to 90 seconds.
This is because CloudFoundry does not yet support asynchronous binding, and CloudSQL bind operations may exceed the default 60 second timeout.

Set `broker_client_timeout_seconds` = 90 in your deployment manifest to change this setting.
