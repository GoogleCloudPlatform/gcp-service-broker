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

## Our repository

- [ ] Draft a new release in GitHub with the tag `vX.Y.Z-rc`
- [ ] Include the version's changelog entry in the description.
- [ ] Check the box labeled **This is a pre-release**.
- [ ] Publish the pre-release.

## File for release

- [ ] Merge the release branch once done.
