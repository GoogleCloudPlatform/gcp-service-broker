## Cloud Build Configuration

### Substitutions
The `cloudbuild-release.yaml` template requires two substitutions:

1. `COMMIT_SHA`: The git commit (`git rev-parse HEAD`) of the repository being
   released
1. `GS_URL`: The URL of the GCS bucket, including path, that release artifacts
   will be uploaded to. For example, `gs://release-bucket/releases`

### Secrets
The JSON key for a service account with project owner permissions is required to
run service broker integration tests. That JSON key must be stored in Secrets
Manager in the project where the Cloud Build execution occurs. The secret should
be named `ROOT_SERVICE_ACCOUNT_JSON`. Configure it with these instructions: https://cloud.google.com/cloud-build/docs/securing-builds/use-encrypted-secrets-credentials

## Run

### Test only
To execute unit and integration tests - but not create a release - run this from
the root of the repository:

`gcloud builds submit --config=ci/cloudbuild.yaml .`

### Release 
To execute unit and integration tests and then publish release artifacts to GCS,
run this from the root of the repository:

```
gcloud builds submit                               \
  --config=ci/cloudbuild-release.yaml              \
  --substitutions=COMMIT_SHA=$(git rev-parse HEAD) .
```
