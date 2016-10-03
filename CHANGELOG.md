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