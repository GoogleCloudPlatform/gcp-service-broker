[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Open Service Broker for Google Cloud Platform

This is a service broker built to be used with [Cloud Foundry](https://docs.cloudfoundry.org/services/overview.html) and Kubernetes.
It adheres to the [Open Service Broker API v2.13](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md).

Service brokers provide a consistent way to create resources and accounts that can access those resources across a variety of different services.

The GCP Service Broker provides support for:

* [BigQuery](https://cloud.google.com/bigquery/)
* [Bigtable](https://cloud.google.com/bigtable/)
* [Cloud SQL](https://cloud.google.com/sql/)
* [Cloud Storage](https://cloud.google.com/storage/)
* [Dataflow](https://cloud.google.com/dataflow/) (preview)
* [Dataproc](https://cloud.google.com/dataproc/docs/overview) (preview)
* [Datastore](https://cloud.google.com/datastore/)
* [Dialogflow](https://cloud.google.com/dialogflow-enterprise/) (preview)
* [Firestore](https://cloud.google.com/firestore/) (preview)
* [Memorystore for Redis](https://cloud.google.com/memorystore/docs/redis/) (preview)
* [ML APIs](https://cloud.google.com/ml/)
* [PubSub](https://cloud.google.com/pubsub/)
* [Spanner](https://cloud.google.com/spanner/)
* [Stackdriver Debugger](https://cloud.google.com/debugger/)
* [Stackdriver Monitoring](https://cloud.google.com/monitoring/) (preview)
* [Stackdriver Trace](https://cloud.google.com/trace/)
* [Stackdriver Profiler](https://cloud.google.com/profiler/)

## Installation

This application can be installed as either a Tanzu Ops Man Tile _or_ deployed as a Cloud Foundry application.
See the [installation instructions](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/installation.md) for a more detailed walkthrough.

## Upgrading

If you're upgrading, check the [upgrade guide](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/upgrading.md).

## Usage

For operators: see [docs/customization.md](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/customization.md) for details about configuring the service broker.

For developers: see [docs/use.md](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/use.md) for information about creating and binding specific GCP services with the broker.
Complete Spring Boot sample applications which use services can be found in the [service-broker-samples repository](https://github.com/GoogleCloudPlatform/service-broker-samples).

You can get documentation specific to your install from the `/docs` endpoint of your deployment.

## Commands

The service broker can be run as both a server (the service broker) and as a general purpose command line utility.
It supports the following sub-commands:

 * `client` - A CLI client for the service broker.
 * `config` - Show and merge configuration options together.
 * `generate` - Generate documentation and tiles.
 * `help` - Help about any command.
 * `serve` - Start the service broker.

## Testing

Pull requests are unit-tested with Travis. You can run the same tests Travis does using `go test ./...`.

Unit and integration tests may be run with Google Cloud Build. See the `ci/`
directory in this repository for instructions.

## Support

**This is not an officially supported Google product.**

[File a GitHub issue](https://github.com/GoogleCloudPlatform/gcp-service-broker/issues) for functional issues or feature requests and use GitHub's notification settings to watch the repository for new releases.

## Contributing

See [the contributing file](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/CONTRIBUTING.md) for more information.
