## This project is archived

This project is now deprecated. Please use the https://github.com/cloudfoundry/cloud-service-broker instead which is a fork of this service broker and is actively supported by the community.

----

## Open Service Broker for Google Cloud Platform

This is a service broker built to be used with [Cloud Foundry](https://docs.cloudfoundry.org/services/overview.html) and Kubernetes.
It adheres to the [Open Service Broker API v2.13](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md).

Service brokers provide a consistent way to create resources and accounts that can access those resources across a variety of different services.

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

## Support

**This is not an officially supported Google product.**

