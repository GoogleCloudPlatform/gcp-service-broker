# Testing Locally

## End to End Tests

The service broker has both unit and end-to-end tests.
End to end tests are generated from the documentation and examples and run outside the standard `go test` framework.
This ensures the auto-generated docs are always up-to-date and the examples work.
By executing the examples as an OSB client, it also ensures the service broker implements the OSB spec correctly.

To run the suite of end-to-end tests:

1. Start an instance of the broker `./gcp-service-broker serve`.
2. In a separate window, run the examples: `./gcp-service-broker client run-examples`
3. Wait for the examples to run and check the exit code. Exit codes other than 0 mean the end-to-end tests failed.

You can also target specific services in the end-to-end tests using the `--service-name` flag.
See `./gcp-service-broker client run-examples --help` for more details.

## Database Setup

You can set up a local MySQL database for testing using Docker:

```
$ docker run -p 3306:3306 --name test-mysql -e MYSQL_ROOT_PASSWORD=password -d mysql:5.7
$ docker exec -it test-mysql mysql -uroot -p
$ mysql> CREATE DATABASE servicebroker;
$ mysql> exit
```

Or, you can run the service broker using SQLite3 for development by specifying
the `db.type` and `db.path` fields:

```
db.type: sqlite3
db.path: service-broker-db.sqlite3
```

## Database Exploration

You can debug the database locally using the `show` sub-command.
It will dump a database table as JSON to stdout.
You can dump: `bindings`, `instances`, `migrations`, and `provisions`.

```
$ ./gcp-service-broker --config test.yaml show provisions
[
    {
        "ID": 1,
        "CreatedAt": "2018-07-17T10:08:07-07:00",
        "UpdatedAt": "2018-07-17T10:08:07-07:00",
        "DeletedAt": null,
        "ServiceInstanceId": "my-cloud-storage",
        "RequestDetails": ""
    }
]
```

## Configuration

Rather than setting environment variables to run the broker you can use the
settings file below and just set the service account JSON as an environment variable:

    ROOT_SERVICE_ACCOUNT_JSON=$(cat service-account.json) ./gcp-service-broker serve --config testconfig.yml



**testconfig.yml**

```
db:
  host: localhost
  name: servicebroker
  password: password
  port: "3306"
  user: root
api:
  user: user
  password: pass
  port: 8000
```


## Useful commands

Create unbind commands for all bindings:

  ./gcp-service-broker show bindings | jq --raw-output '.[] | "./gcp-service-broker client unbind --bindingid \(.BindingId) --instanceid \(.ServiceInstanceId) --planid \(.PlanId) --serviceid \(.ServiceId)"'

Create deprovision commands for all bindings:

  ./gcp-service-broker show instances | jq --raw-output '.[] | "./gcp-service-broker client deprovision --instanceid \(.ID) --serviceid \(.ServiceId) --planid \(.PlanId)"'
