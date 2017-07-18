# Change Log
All releases of the GCP Service Broker will be documented in
this file. This project adheres to [Semantic Versioning](http://semver.org/).

## [1.0.0] - 2016-10-03

### Fixed
- Switched from using PLANS environment variable to using CLOUDSQL_CUSTOM_PLANS 
to generate CloudSQL plans. Fixed bug where at least one CloudSQL plan was required
and changed DB password type in tile config from string to secret. Note that due to
a migration issue in Ops Manager, you'll need to delete and re-install the broker 
if you are using it as a tile.

## [1.0.1] - 2016-10-07

### Fixed
- Removed specified stemcell version from tile.yml so that most recent stemcell is 
used by default.

## [2.0.0] - 2016-10-10

### Fixed
- fixed CloudSQL docs link in README
- updated credentials type returned by bind call to be a map[string]string instead
of a string.

## [2.0.1] - 2016-10-28

- CloudSQL will generate a username/password on bind if one is not provided.

### Fixed
- CloudSQL custom plans are now optional
- fixed username and password env var names in docs
- fixed CloudSQL custom plan names and/or features not updating

## [2.0.2] - 2016-11-16

### Fixed
- fixed bug where CloudSQL was returning 400s for all 2nd gen instance provision requests

## [2.1.0] - 2016-12-02

- Remove need for service name for PubSub (topic_name), BigQuery (name), Cloud Storage (name), and Cloud SQL (instance_name)
- Instance details are now surfaced in bind requests for Pubsub (topic_name and subscription_name),
Cloud Storage (bucket_name), BigQuery (dataset_id), and Cloud SQL (instance_name, database_name, and host)

## [2.1.1] - 2016-12-02

- added new uri parameter to Cloud SQL bind credentials

## [2.1.2] - 2016-12-21

### Fixed
- fixes a bug where anything that triggered an install repeat after installing version 2.1.0 or 2.1.1 would cause the
installation to fail.

## [2.1.3] - 2017-01-05

### Fixed
- fixes a bug where bind calls to ml-api service instances were failing because these service instances don't save
any extra access details

## [3.0.0] - 2017-01-12

- Updated pubsub library so that User Agent string gets propagated correctly
- Updated dependency management system
- Changed org to system (Broker will need to be uninstalled and reinstalled for this change to take effect)

## [3.0.1] - 2017-01-23

- Updated default user agent string so that bogus data will not be collected during testing
- Updated service account bindings to include ProjectId

## [3.1.0] - 2017-02-07

- Updated vendored packages so that custom UserAgent string gets propagated for storage provision calls
- Added Bigtable as a service

## [3.1.1] - 2017-02-09

- fixed a bug in reading in custom plans

## [3.1.2] - 2017-02-12

- fixed a bug where supplying a custom name for Bigtable would cause an error

## [3.2.0] - 2017-03-22

- Added Spanner support
- Added a Golang example application
- Added Concourse CI pipeline
- Added integration testing for async services

## [3.2.1] - 2017-03-31

- fixed a bug where Spanner instances could be deleted from Google but not deprovisioned in CF
- fixed a bug where Cloud SQL instances were not being marked as deleted in the Service Broker database

## [3.3.0] - 2017-03-31

- Added Stackdriver Debugger and Trace support

## [3.3.1] - 2017-05-11

- Security updates to address CVE-2017-4975
- database name is configurable by setting `DB_NAME`

## [3.3.2] - 2017-05-11

- Security updates to address CVE-2017-4975
- properly set buildpack

## [3.3.2] - 2017-05-11

- Security updates to address CVE-2017-4975
- properly set buildpack

## [3.4.0] - 2017-05-16

- Add Stackdriver services to tile.yml

## [3.4.1] - 2017-05-22

- fixes uninitialized security group

## [4.0.0] - 2017-07-XX

- added command `get_plan_info` to dump plan information to the console. Advised to run before updating.
- plan ids are now required and will not be generated if not supplied.
- changed custom plan id field name from `guid` to `id`
- modified `"features"` plan config field name to `"service_properties"`
- modified the formatting of custom plans to be consistent with that of preconfigured plans
