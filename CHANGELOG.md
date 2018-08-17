# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unresolved]

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

### Deprecated
- Running the service by executing the main executable. Use the `serve` sub-command instead.

### Changed  
- **Breaking** plan ids are now required and will not be generated if not supplied.
- **Breaking** changed custom plan id field name from `guid` to `id`.
- **Breaking** modified `"features"` plan config field name to `"service_properties"`.
- **Breaking** modified structure of all catalog-related environment variables - `plans` is now a sub-field of the Service object,
and Service objects are defined individually by setting env variables like `GOOGLE_<SERVICE_NAME>`
- You no longer have to specify service information in the manifest.
- **Breaking** The GCS plan `reduced_availability` was changed to `reduced-availability` to be compliant with the spec and work with Kubernetes.
- Tables are created only if they do not exist on migration, fixing [#194](https://github.com/GoogleCloudPlatform/gcp-service-broker/issues/194).
- The broker now adheres to OSB version 2.13.

## [3.6.0] - 2017-01-03

- changed default authorized networks for cloudsql instances from `0.0.0.0/0` to none
- added optional parameter to cloudsql provision operation to specify authorized networks
- added service account and key provisioning on cloudsql bind operations
- updated postgres `uri` field to include ssl certificates
- added optional parameter to cloudsql bind operation to pass back jdbc formatted `uri` field
- removed waiting for ssl certs to finish being created in sql account manager
- changed default number of spanner nodes to 1 in `tile.yml`

## [3.5.2] - 2017-10-17

- fixed Postgres connection uri
- added wait for ssl certs to finish being created in sql account manager

## [3.5.1] - 2017-09-06

- added PostgreSQL support to CloudSQL (and migrated existing plans)
- added Datastore support

## [3.4.1] - 2017-05-22

- fixes uninitialized security group

## [3.4.0] - 2017-05-16

- Add Stackdriver services to tile.yml

## [3.3.2] - 2017-05-11

- Security updates to address CVE-2017-4975
- properly set buildpack

## [3.3.2] - 2017-05-11

- Security updates to address CVE-2017-4975
- properly set buildpack

## [3.3.1] - 2017-05-11

- Security updates to address CVE-2017-4975
- database name is configurable by setting `DB_NAME`

## [3.3.0] - 2017-03-31

- Added Stackdriver Debugger and Trace support

## [3.2.1] - 2017-03-31

- fixed a bug where Spanner instances could be deleted from Google but not deprovisioned in CF
- fixed a bug where Cloud SQL instances were not being marked as deleted in the Service Broker database

## [3.2.0] - 2017-03-22

- Added Spanner support
- Added a Golang example application
- Added Concourse CI pipeline
- Added integration testing for async services

## [3.1.2] - 2017-02-12

- fixed a bug where supplying a custom name for Bigtable would cause an error

## [3.1.1] - 2017-02-09

- fixed a bug in reading in custom plans

## [3.1.0] - 2017-02-07

- Updated vendored packages so that custom UserAgent string gets propagated for storage provision calls
- Added Bigtable as a service

## [3.0.1] - 2017-01-23

- Updated default user agent string so that bogus data will not be collected during testing
- Updated service account bindings to include ProjectId

## [3.0.0] - 2017-01-12

- Updated pubsub library so that User Agent string gets propagated correctly
- Updated dependency management system
- Changed org to system (Broker will need to be uninstalled and reinstalled for this change to take effect)

## [2.1.3] - 2017-01-05

### Fixed
- fixes a bug where bind calls to ml-api service instances were failing because these service instances don't save
any extra access details

## [2.1.2] - 2016-12-21

### Fixed
- fixes a bug where anything that triggered an install repeat after installing version 2.1.0 or 2.1.1 would cause the
installation to fail.

## [2.1.1] - 2016-12-02

- added new uri parameter to Cloud SQL bind credentials

## [2.1.0] - 2016-12-02

- Remove need for service name for PubSub (topic_name), BigQuery (name), Cloud Storage (name), and Cloud SQL (instance_name)
- Instance details are now surfaced in bind requests for Pubsub (topic_name and subscription_name),
Cloud Storage (bucket_name), BigQuery (dataset_id), and Cloud SQL (instance_name, database_name, and host)

## [2.0.2] - 2016-11-16
### Fixed
- fixed bug where CloudSQL was returning 400s for all 2nd gen instance provision requests

## [2.0.1] - 2016-10-28

- CloudSQL will generate a username/password on bind if one is not provided.

### Fixed
- CloudSQL custom plans are now optional
- fixed username and password env var names in docs
- fixed CloudSQL custom plan names and/or features not updating

## [2.0.0] - 2016-10-10
### Fixed
- fixed CloudSQL docs link in README
- updated credentials type returned by bind call to be a map[string]string instead
of a string.


## [1.0.1] - 2016-10-07
### Fixed
- Removed specified stemcell version from tile.yml so that most recent stemcell is
used by default.

## [1.0.0] - 2016-10-03
### Fixed
- Switched from using PLANS environment variable to using CLOUDSQL_CUSTOM_PLANS
to generate CloudSQL plans. Fixed bug where at least one CloudSQL plan was required
and changed DB password type in tile config from string to secret. Note that due to
a migration issue in Ops Manager, you'll need to delete and re-install the broker
if you are using it as a tile.
