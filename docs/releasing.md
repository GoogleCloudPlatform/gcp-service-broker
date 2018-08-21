# Release Process

## Pre-Check

- [ ] The system should build and all tests pass with the concourse pipeline.

## Generating the OSDF file

- [ ] Get a list of licenses using the [license_finder](https://github.com/pivotal-legacy/LicenseFinder) tool.
- [ ] Fill in the [license template](https://docs.google.com/spreadsheets/d/1gqS1jwmpSIEdgTadQXhkbQqhm1hO3qOU1-AWwIYQqnw/edit#gid=0) with them.
- [ ] Download it as a CSV.
- [ ] Upload the CSV to the [OSDF Generator](http://osdf-generator.cfapps.io/static/index.html) and download the new OSDF file.
- [ ] Create a new pre-release tag in the git repository.
