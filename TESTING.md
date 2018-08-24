# Testing Locally

## Database Setup

You can set up a local MySQL database for testing using Docker:

```
$ docker run -p 3306:3306 --name test-mysql -e MYSQL_ROOT_PASSWORD=password -d mysql:5.7
$ docker exec -it test-mysql mysql -uroot -p
$ mysql> CREATE DATABASE servicebroker;
$ mysql> exit
```

Or, you can run the service broker using SQLite3 for development by sepecifying
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

  ./gcp-service-broker show bindings | jq --raw-output '.[] | "gcp-service-broker client unbind --bindingid \(.BindingId) --instanceid \(.ServiceInstanceId) --planid \(.PlanId) --serviceid \(.ServiceId)"'

Create deprovision commands for all bindings:

  ./gcp-service-broker show instances | jq --raw-output '.[] | "./gcp-service-broker client deprovision --instanceid \(.ID) --serviceid \(.ServiceId) --planid \(.PlanId)"'
