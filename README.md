[![Go Report Card](https://goreportcard.com/badge/github.com/GoogleCloudPlatform/gcp-service-broker)](https://goreportcard.com/report/github.com/GoogleCloudPlatform/gcp-service-broker) [![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Cloud Foundry Service Broker for Google Cloud Platform

This is a service broker built to be used with [Cloud Foundry](https://docs.cloudfoundry.org/services/overview.html).
It adheres to the [Open Service Broker API v2.13](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md).

Service brokers provide a consistent way to create resources and accounts that can access those resources across a variety of different services.

The GCP Service Broker provides support for:

* [BigQuery](https://cloud.google.com/bigquery/)
* [Bigtable](https://cloud.google.com/bigtable/)
* [Cloud SQL](https://cloud.google.com/sql/)
* [Cloud Storage](https://cloud.google.com/storage/)
* [Datastore](https://cloud.google.com/datastore/)
* [Dialogflow](https://cloud.google.com/dialogflow-enterprise/)
* [Firestore](https://cloud.google.com/firestore/)
* [ML APIs](https://cloud.google.com/ml/)
* [PubSub](https://cloud.google.com/pubsub/)
* [Spanner](https://cloud.google.com/spanner/)
* [Stackdriver Debugger](https://cloud.google.com/debugger/)
* [Stackdriver Monitoring](https://cloud.google.com/monitoring/)
* [Stackdriver Trace](https://cloud.google.com/trace/)
* [Stackdriver Profiler](https://cloud.google.com/profiler/)

## Installation

This application can be installed as either a PCF Ops Man Tile _or_ deployed as a PCF application.
See the [installation instructions](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/installation.md) for a more detailed walkthrough.

## Upgrading

If you're upgrading, check the [upgrade guide](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/upgrading.md).

## Usage

For operators: see [docs/customization.md](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/customization.md) for details about configuring the service broker.

For developers: see [docs/use.md](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/use.md) for information about creating and binding specific GCP services with the broker.
Complete Spring Boot sample applications which use services can be found in the [service-broker-samples repository](https://github.com/GoogleCloudPlatform/service-broker-samples).


## Commands

The service broker can be run as both a server (the service broker) and as a general purpose command line utility.
It supports the following sub-commands:

 * `client` - A CLI client for the service broker.
 * `config` - Show and merge configuration options together.
 * `generate` - Generate documentation and tiles.
 * `help` - Help about any command.
 * `migrate` - Upgrade your database (you generally won't need this because the databases auto-upgrade).
 * `plan-info` - Dump plan information from the database.
 * `serve` - Start the service broker.
 * `show` - Show info about the provisioned resources.

## Testing

Production testing for the GCP Service Broker is administered via a private Concourse pipeline.

To run tests locally, use [Ginkgo](https://onsi.github.io/ginkgo/).

Integration tests require the `ROOT_SERVICE_ACCOUNT_JSON` environment variable to be set.

**Note: Integration tests create and destroy real project resources and therefore have associated costs to run**


## Support

[File a GitHub issue](https://github.com/GoogleCloudPlatform/gcp-service-broker/issues) for functional issues or feature requests.

Subscribe to the [gcp-service-broker Google group](https://groups.google.com/forum/#!forum/gcp-service-broker) for discussions and updates.


## Contributing

See [the contributing file](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/CONTRIBUTING.md) for more information.

This is not an officially supported Google product.
