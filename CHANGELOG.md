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
