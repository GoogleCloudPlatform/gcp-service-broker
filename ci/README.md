## Credhub configuration

The following parameters should be set in the `/concourse/cf` keyspace of the
Credhub instance associated with the Concourse installation that runs pipelines
for the service broker:

| Name                                  | Credhub command                                                                        |
| ---                                   | ---                                                                                    |
| code_branch                           | `credhub set --type value --name /concourse/cf/code_branch --value ...`                              |
| ci_branch                             | `credhub set --type value --name /concourse/cf/ci_branch --value ...`                                |
| artifacts_bucket_name                 | `credhub set --type value --name /concourse/cf/artifacts_bucket_name --value ...`                    |
| artifacts_json_key                    | `credhub set --type value --name /concourse/cf/artifacts_json_key --value ...`                        |
| integration_test_service_account_json | `credhub set --type value --name /concourse/cf/integration_test_service_account_json --value ...`     |
| integration_test_db_username          | `credhub set --type value --name /concourse/cf/integration_test_db_username --value ...`             |
| integration_test_db_password          | `credhub set --type value --name /concourse/cf/integration_test_db_password --value ...`             |
| integration_test_db_host              | `credhub set --type value --name /concourse/cf/integration_test_db_host --value ...`                 |
| integration_test_ca_cert              | `credhub set --type value --name /concourse/cf/integration_test_ca_cert --value ...`     |
| integration_test_client_cert          | `credhub set --type value --name /concourse/cf/integration_test_client_cert --value ...` |
| integration_test_client_key           | `credhub set --type value --name /concourse/cf/integration_test_client_key --value ...`              |
