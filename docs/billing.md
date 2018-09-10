# Billing

The GCP Service Broker automatically labels supported resources with organization GUID, space GUID and instance ID.

When these supported services are provisioned, they will have the following labels populated with information from the request:

 * `pcf-organization-guid`
 * `pcf-space-guid`
 * `pcf-instance-id`

GCP labels have a more restricted character set than the Service Broker so unsupported characters will be mapped to the underscore character (`_`).

## Support

The following resources support these automatically generated labels:

 * BigQuery
 * CloudSQL (PostgreSQL and MySQL)
 * Cloud Storage
 * Spanner

## Usage

You can use these labels with the [BigQuery Billing Export](https://cloud.google.com/billing/docs/how-to/bq-examples)
to create reports about which organizations and spaces are incurring cost in your GCP project.
