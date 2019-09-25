# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [4.3.0] - 2019-08-26

### Fixed
- Fix bug that caused SQL users to never be deleted.

### Changed
- Built Brokerpaks now have a name of `{name}-{version}.brokerpak` as defined by the manifest rather than the name of the parent directory.
- Services inside Brokerpaks now have a file name that includes their CLI friendly name to help differentiate them.
- The `pak build` command now includes progress logs.
- An unbind of a service now fails as early as possible to prevent partial deletions.

### Added
 - Ability for plans to selectively override user variables.
 - Ability to get plan information from HIL execution environment on bind.

### Fixed

- "data too long" error for existing service broker installations when passing service provisioning configuration > 255 characters (#468)

## [4.2.3] - 2019-06-12

### Fixed
 - Added a workaround for an upstream CloudSQL issue that caused SQL user deletion to fail.
 - Delete replicas for CloudSQL instances before attempting to delete the instance.

### Added
 - A new option to MySQL that allows auto-generating a replica name from the
   name of the master.
 - A new option to Cloud Storage that allows buckets to be deleted even if they
   contain objects. This sets the label `sb-force-delete` to `true` on the
   bucket and will attempt to delete all contents before deleting the bucket.

### Changed
 - Removed requirement to have a 'default' network when using google-redis.

## [4.2.2] - 2019-02-06

### Fixed
 - The `pak run-examples` sub-command now returns a non-zero status code on failure.
 - JSONSchema validation no longer fails due to erroneous duplicate required fields.
 - Brokerpaks no longer use incorrect templates due to an invalid pointer.

## [4.2.0] - 2019-01-04

### Security
 - The broker uses a Pivotal library affected by [CVE-2018-15759](https://pivotal.io/security/cve-2018-15759). Until the library is updated, it's recommended that you not run the service broker on a public network. If you must run it on a public network, make it accessible through a proxy that supports fail2ban.

### Added
 - The ability to enable/disable services based on product lifecycle tags. See [#340](https://github.com/GoogleCloudPlatform/gcp-service-broker/pull/340) for context.
 - Preview support for Firestore.
 - Preview support for Dialogflow.
 - Preview support for Stackdriver Metrics.
 - Namespace support for Datastore.
 - Preview support for Dataflow.
 - Default roles for ML, BigQuery, BigTable, CloudSQL, Pub/Sub, Spanner, and Cloud Storage.
 - `/docs` endpoint that serves docs for your installation.
 - Preview support for Dataproc via brokerpaks.
 - Varcontext now supports casting computed HIL values.
 - New regional and multi-regional Cloud Storage plans.
 - Ability to expose JSONSchema in the service catalog by enabling the `enable-catalog-schemas` flag.

### Changed
 - Support links for services now point to service-specific pages where possible.
 - Feature flags are now handled through a generic toggles framework. Option labels and descriptions might change slightly in the tile.
 - Service definitions now get field-level validation to check for sanity before being registered.

### Removed
 - The `examples/` directory.

### Fixed
 - Fixed machine types for PostgreSQL to use custom but keep the old names for compatibility.

## [4.1.0] - 2018-11-05

### Added
- Pub/Sub now adds default labels to created topics and subscriptions.
- New validation documentation for Pub/Sub.
- Ability for operators to override the provision defaults with fixed values.
- New form to let operators set provision defaults.
- New `coldline` Cloud Storage plan.
- Ability to create custom Cloud Storage plans.
- New tile form for creating custom Cloud Storage plans.
- Examples of binding variables to the docs.
- Constraints/validation of the binding variables to the docs.
- New `version` sub-command to show the current version of the binary.
- New `generate` sub-commands to generate the `tile.yml` and `manifest.yml` files.

### Changed
- Role whitelists are now validated through JSON Schema checks.
- The `run-examples` sub-command now evaluates the credentials against the JSON Schema, improving robustness.

### Fixed
- Fixed issue where Cloud Datastore service accounts were getting the same name.

## [4.0.0] - 2018-10-01

### Added
- New sub-command `plan-info` to dump plan information to the console.
- New sub-command `client` to execute documentation examples and interact with the broker.
- New sub-command `help` which outputs help documentation.
- New sub-command `config` which can convert between configuration file formats.
- New sub-command `generate` to generate use, tile forms, and configuration documentation.
- New sub-command `serve` to run the service broker server.
- New sub-command `show` to dump database state.
- The ability to configure the system with YAML, TOML, properties, or JSON configuration files via the `--config` flag in conjunction with using environment variables.
- The ability to customize the database name in the Tile.
- The ability to turn on/off services via environment variable.
- Default plans for Spanner, BigTable, and CloudSQL.
- Whitelists for bindings so only certain "safe" roles can be chosen by end-users.
- Automatic labeling of resources with organization GUID, space GUID, and instance ID to BigQuery, CloudSQL, Spanner, and Cloud Storage.

### Deprecated
- Running the service by executing the main executable. Use the `serve` sub-command instead.

### Changed
- **Breaking** plan ids are now required and will not be generated if not supplied.
- **Breaking** changed custom plan id field name from `guid` to `id`.
- **Breaking** modified `"features"` plan configuration field name to `"service_properties"`.
- **Breaking** modified structure of all catalog-related environment variables - `plans` is now a sub-field of the Service object,
and Service objects are defined individually by setting environment variables like `GOOGLE_<SERVICE_NAME>`
- You no longer have to specify service information in the manifest.
- **Breaking** The Cloud Storage plan `reduced_availability` was changed to `reduced-availability` to be compliant with the spec and work with Kubernetes.
- Tables are created only if they do not exist on migration, fixing [#194](https://github.com/GoogleCloudPlatform/gcp-service-broker/issues/194).
- The broker now adheres to Open Service Broker API version 2.13.
- Improved ORM migrations and test coverage with SQLite3.

## [3.6.0] - 2017-01-03

- Changed default number of spanner nodes to 1 in `tile.yml`.
- Changed default authorized networks for CloudSQL instances from `0.0.0.0/0` to none.
- Added optional parameter to CloudSQL provision operation to specify authorized networks.
- Added service account and key provisioning on CloudSQL bind operations.
- Added optional parameter to CloudSQL bind operation to pass back JDBC formatted `uri` field.
- Updated PostgreSQL `uri` field to include SSL certificates.
- Removed waiting for SSL certs to finish being created in SQL account manager.

## [3.5.2] - 2017-10-17

- Added wait for SSL certs to finish being created in SQL account manager.
- Fixed PostgreSQL connection URI.

## [3.5.1] - 2017-09-06

- Added PostgreSQL support to CloudSQL (and migrated existing plans).
- Added Cloud Datastore support.

## [3.4.1] - 2017-05-22

- Fixed uninitialized security group.

## [3.4.0] - 2017-05-16

- Added Stackdriver services to `tile.yml`.

## [3.3.2] - 2017-05-11

- Security updates to address [CVE-2017-4975](https://nvd.nist.gov/vuln/detail/CVE-2017-4975).
- Fixed properly set buildpack.

## [3.3.1] - 2017-05-11

- Security updates to address [CVE-2017-4975](https://nvd.nist.gov/vuln/detail/CVE-2017-4975).
- Added environment variable `DB_NAME` to allow configuring the database name.

## [3.3.0] - 2017-03-31

- Added Stackdriver Debugger and Stackdriver Trace support.

## [3.2.1] - 2017-03-31

- Fixed a bug where Spanner instances could be deleted from Google but not deprovisioned in CF.
- Fixed a bug where CloudSQL instances were not being marked as deleted in the Service Broker database.

## [3.2.0] - 2017-03-22

- Added Spanner support.
- Added an example Golang application.
- Added Concourse CI pipeline.
- Added integration testing for asynchronous services.

## [3.1.2] - 2017-02-12

- Fixed a bug where supplying a custom name for Bigtable would cause an error.

## [3.1.1] - 2017-02-09

- Fixed a bug in reading in custom plans.

## [3.1.0] - 2017-02-07

- Updated vendor packages so that custom User-Agent string gets propagated for storage provision calls.
- Added Bigtable as a service.

## [3.0.1] - 2017-01-23

- Updated default user agent string so that bogus data will not be collected during testing.
- Updated service account bindings to include ProjectId.

## [3.0.0] - 2017-01-12

- Updated Pub/Sub library so that User-Agent string gets propagated correctly.
- Updated dependency management system.
- **Breaking**, changed org to system (Broker will need to be uninstalled and reinstalled for this change to take effect).

## [2.1.3] - 2017-01-05

- Fixed a bug where bind calls to ml-api service instances were failing because these service instances don't save extra access details.

## [2.1.2] - 2016-12-21

- Fixed a bug where anything that triggered an install repeat after installing version 2.1.0 or 2.1.1 would cause the installation to fail.

## [2.1.1] - 2016-12-02

- Added new URI parameter to CloudSQL bind credentials.

## [2.1.0] - 2016-12-02

- Removed need for service name for Pub/Sub (topic_name), BigQuery (name), Cloud Storage (name), and CloudSQL (instance_name).
- Instance details are now surfaced in bind requests for Pub/Sub (topic_name and subscription_name), Cloud Storage (bucket_name), BigQuery (dataset_id), and CloudSQL (instance_name, database_name, and host).

## [2.0.2] - 2016-11-16

- Fixed bug where CloudSQL was returning 400s for all 2nd gen instance provision requests.

## [2.0.1] - 2016-10-28

- Added CloudSQL will generate a username and password on bind if one is not provided.
- Fixed CloudSQL custom plans are now optional.
- Fixed username and password environment variable names in docs.
- Fixed CloudSQL custom plan names and/or features not updating.

## [2.0.0] - 2016-10-10

- Fixed CloudSQL docs link in `README.md`.
- Updated credentials type returned by bind call to be a map[string]string instead of a string.


## [1.0.1] - 2016-10-07

- Removed specified Stemcell version from `tile.yml` so the most recent Stemcell is used by default.

## [1.0.0] - 2016-10-03

- **Breaking**, switched from using `PLANS` environment variable to using `CLOUDSQL_CUSTOM_PLANS` to generate CloudSQL plans.
- Fixed bug where at least one CloudSQL plan was required.
- Changed database password type in tile configuration from string to secret.
- **Note** due to a migration issue in Ops Manager, you'll need to delete and re-install the broker if you are using it as a tile.
