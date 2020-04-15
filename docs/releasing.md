# Release Process

## Before continuing

- [ ] Open a ticket with the name `release-vX.Y.Z` and copy the contents of this file into its description.
- [ ] Create a new pre-release branch in the GitHub repository labeled `release-vX.Y.Z`.

## Update the project in the branch

- [ ] Update version in `utils/version.go`.
- [ ] Update the stemcell and go version in `pkg/generator/pcf-artifacts.go`
- [ ] Run `./hack/build.sh`.
- [ ] Update the `CHANGELOG.md` to match the new version.
- [ ] Commit the changes on the new branch.
- [ ] Wait for the system to build and all Concourse tests to pass.

## Generate the OSDF

### Option 1: automatic

- [ ] Run `./hack/update-osdf.sh X.Y.Z`

### Option 2: fallback

- [ ] Get a list of licenses using the [license_finder](https://github.com/pivotal-legacy/LicenseFinder) tool.
- [ ] Fill in the [license template](https://docs.google.com/spreadsheets/d/1gqS1jwmpSIEdgTadQXhkbQqhm1hO3qOU1-AWwIYQqnw/edit#gid=0) with them.
- [ ] Download it as a CSV.
- [ ] Upload the CSV to the [OSDF Generator](http://osdf-generator.cfapps.io/static/index.html) and download the new OSDF file.
- [ ] Replace the current OSDF file in the root of the project with the OSDF for the release and commit the change.

## Our repository

- [ ] Draft a new release in GitHub with the tag `vX.Y.Z-rc`
- [ ] Include the version's changelog entry in the description.
- [ ] Upload the built tile and OSDF to the pre-release.
- [ ] Check the box labeled **This is a pre-release**.
- [ ] Publish the pre-release.

## Release on PivNet

- [ ] Validate that the name in the generated tile's `metadata.yml` matches the slug on PivNet.
- [ ] Ensure the release version is consistent on the tile and documentation.
- [ ] Create a [new release on Tanzu Network](network.pivotal.io) as an Admin Only release.
- [ ] Upload the tile and OSDF files that were staged to GitHub.
- [ ] Check that the tile passes the tests in the [build dashboard](https://tile-dashboard.cfapps.io/tiles/gcp-service-broker).

## Upgrade the documentation

- [ ] Submit a pull request to the [documentation repository](https://github.com/pivotal-cf/docs-google/tree/master/docs-content).
- [ ] Include the new release notes, changes, and the ERT/PAS and Ops Managers versions, as well as your Product Version and Release Date in the Product Snapshot on PivNet.

## File for release

- [ ] Fill out the [release form](https://docs.google.com/forms/d/e/1FAIpQLSctLGMU8iOuwq6NqDYI65aMhJ7widDQGo9SawDG0b8TFfq7Ag/viewform).
- [ ] An ISV Program Manager will make the release available to "All Users" after review. Partner Admins can make the release available to "Admin Users".
- [ ] Merge the release branch once done.
- [ ] Make a release announcement in the gcp-service-broker Google Group like [this one](https://groups.google.com/forum/#!topic/gcp-service-broker/7Ae9D2B1AzE).
- [ ] Submit an issue to https://github.com/cf-platform-eng/gcp-pcf-quickstart to update the GCP PCF quickstart.
