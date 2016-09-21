# Pivotal Cloud Foundry Service Broker for Google Cloud Platform

Depends on
[lager](https://github.com/pivotal-golang/lager) and
[gorilla/mux](https://github.com/gorilla/mux).

Requires go 1.6.

## Prerequisites

### GCP prereqs

1. create a new project
2. in the left nav, go to API Manager
3. search "google cloud resource manager api", click the option with no other modifiers, and enable.
4. search "Google Identity and Access Management (IAM) API", click the option with no other modifiers, and enable.
5. in the left nav, go to IAM and Admin
6. click Service Accounts
7. click create service account and set the role to Owner
8. check furnish new private key (leave json key type)
9. click create, save the downloaded file to a safe and accessible place.

### db prereqs

1. create new mysql instance
2. create "servicebroker" database
3. create a user for the service broker
4. grant the service broker user privileges on the servicebroker database
5. (optional) create ssl certs for the database

### required env vars

* ROOT_SERVICE_ACCOUNT_JSON (the string version of the credentials file created for the Owner level Service Account)
* SB_USERNAME (a username to sign all service broker requests with - the same one used in cf create-service-broker)
* SB_PASSWORD (a password to sign all service broker requests with - the same one used in cf create-service-broker)
* DB_HOST (the host for the database to back the service broker)
* DB_USERNAME (the database username for the service broker to use)
* DB_PASSWORD (the database password for the service broker to use)

### optional env vars

* DB_PORT (defaults to 3306)
* CA_CERT_B64 (base64 encoded version of your database ca cert)
* CLIENT_CERT_B64 (base64 encoded version of your database client cert)
* CLIENT_KEY_B64 (base64 encoded version of your database client key)


## Usage

### Install custom Go 1.6 buildpack if necessary
1. Download version 1.7.10 of the go buildpack from https://github.com/cloudfoundry/go-buildpack/releases
2. cf create-buildpack go_1_6_buildpack <path to buildpack file> <buildpack order> --enable

### (If using as an app) Update the manifest with ENV vars
1. replace any blank variables that are in manifest.yml with your own ENV vars

### Migrate the backing DB
1. cd cmd/migrate
2. ./migrate

### Push the service broker to CF and enable services
1. cf push gcp-service-broker -b go_1_6_buildpack
2. cf create-service-broker <service broker name> <username> <password> <service broker url>
3. cf enable-service-access pubsub

### Use!

e.g. cf create-service pubsub default foobar -c '{"topic_name": "foobar"}'
e.g. cf bind-service myapp foobar -c '{"role": "pubsub.admin"}'

### (Optional) Increase the default provision/bind timeout
It is advisable, if you want to use CloudSQL, to increase the default timeout for provision and
bind operations to 90 seconds. CloudFoundry does not, at this point in time, support asynchronous
binding, and CloudSQL bind operations may exceed 60 seconds. To change this setting, set
broker_client_timeout_seconds = 90


## Errors

`brokerapi` defines a handful of error types in `service_broker.go` 

The error types are:

```go
ErrInstanceAlreadyExists
ErrInstanceDoesNotExist
ErrInstanceLimitMet
ErrBindingAlreadyExists
ErrBindingDoesNotExist
ErrAsyncRequired
```

Additionally, some custom errors may be returned from the Google APIs.

## Change Notes


# This is not an official Google product.