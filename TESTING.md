# Testing Locally

## Database Setup

You can set up a local MySQL database for testing using Docker:

```
$ docker run -p 3306:3306 --name test-mysql -e MYSQL_ROOT_PASSWORD=password -d mysql:5.7
$ docker exec -it test-mysql mysql -uroot -p
$ mysql> CREATE DATABASE servicebroker;
$ mysql> exit
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

Rather than setting ~20 environment variables to run the broker you can use the
settings file below and just set the service account JSON as an environment variable:

    ROOT_SERVICE_ACCOUNT_JSON=$(cat service-account.json) ./gcp-service-broker serve --config testconfig.yml



**testconfig.yml**

```
service:
  google-storage: |-
    {
      "id": "b9e4332e-b42b-4680-bda5-ea1506797474",
      "description": "A Powerful, Simple and Cost Effective Object Storage Service",
      "name": "google-storage",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google Cloud Storage",
        "longDescription": "A Powerful, Simple and Cost Effective Object Storage Service",
        "documentationUrl": "https://cloud.google.com/storage/docs/overview",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/storage.svg"
      },
      "tags": ["gcp", "storage"],
      "plans": [
        {
          "id": "e1d11f65-da66-46ad-977c-6d56513baf43",
          "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
          "name": "standard",
          "display_name": "Standard",
          "description": "Standard storage class",
          "service_properties": {"storage_class": "STANDARD"}
        },
        {
          "id": "a42c1182-d1a0-4d40-82c1-28220518b360",
          "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
          "name": "nearline",
          "display_name": "Nearline",
          "description": "Nearline storage class",
          "service_properties": {"storage_class": "NEARLINE"}
        },
        {
          "id": "1a1f4fe6-1904-44d0-838c-4c87a9490a6b",
          "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
          "name": "reduced_availability",
          "display_name": "Durable Reduced Availability",
          "description": "Durable Reduced Availability storage class",
          "service_properties": {"storage_class": "DURABLE_REDUCED_AVAILABILITY"}
        }
      ]
    }
  google-pubsub: |-
    {
      "id": "628629e3-79f5-4255-b981-d14c6c7856be",
      "description": "A global service for real-time and reliable messaging and streaming data",
      "name": "google-pubsub",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google PubSub",
        "longDescription": "A global service for real-time and reliable messaging and streaming data",
        "documentationUrl": "https://cloud.google.com/pubsub/docs/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/pubsub.svg"
      },
      "tags": ["gcp", "pubsub"],
      "plans": [
        {
          "id": "622f4da3-8731-492a-af29-66a9146f8333",
          "service_id": "628629e3-79f5-4255-b981-d14c6c7856be",
          "name": "default",
          "display_name": "Default",
          "description": "PubSub Default plan",
          "service_properties": {}
        }
      ]
    }
  google-bigquery: |-
    {
      "id": "f80c0a3e-bd4d-4809-a900-b4e33a6450f1",
      "description": "A fast, economical and fully managed data warehouse for large-scale data analytics",
      "name": "google-bigquery",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google BigQuery",
        "longDescription": "A fast, economical and fully managed data warehouse for large-scale data analytics",
        "documentationUrl": "https://cloud.google.com/bigquery/docs/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/bigquery.svg"
      },
      "tags": ["gcp", "bigquery"],
      "plans": [
        {
          "id": "10ff4e72-6e84-44eb-851f-bdb38a791914",
          "service_id": "f80c0a3e-bd4d-4809-a900-b4e33a6450f1",
          "name": "default",
          "display_name": "Default",
          "description": "BigQuery default plan",
          "service_properties": {}
        }
      ]
    }
  google-cloudsql-mysql: |-
    {
        "id": "4bc59b9a-8520-409f-85da-1c7552315863",
        "description": "Google Cloud SQL is a fully-managed MySQL database service",
        "name": "google-cloudsql-mysql",
        "bindable": true,
        "plan_updateable": false,
        "metadata": {
          "displayName": "Google CloudSQL MySQL",
          "longDescription": "Google Cloud SQL is a fully-managed MySQL database service",
          "documentationUrl": "https://cloud.google.com/sql/docs/",
          "supportUrl": "https://cloud.google.com/support/",
          "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg"
        },
        "tags": ["gcp", "cloudsql", "mysql"],
        "plans": [
          {
              "id": "test-cloudsql-plan",
              "name": "test_plan",
              "description": "testplan",
              "service_properties": {
                "tier": "D8",
                "pricing_plan": "PER_USE",
                "max_disk_size": "15"
              },
              "display_name": "FOOBAR",
              "service_id": "4bc59b9a-8520-409f-85da-1c7552315863"
            }
        ]
    }
  google-cloudsql-postgres: |-
    {
        "id": "cbad6d78-a73c-432d-b8ff-b219a17a803a",
          "description": "Google Cloud SQL is a fully-managed PostgreSQL database service",
          "name": "google-cloudsql-postgres",
          "bindable": true,
          "plan_updateable": false,
          "metadata": {
          "displayName": "Google CloudSQL PostgreSQL",
          "longDescription": "Google Cloud SQL is a fully-managed MySQL database service",
          "documentationUrl": "https://cloud.google.com/sql/docs/",
          "supportUrl": "https://cloud.google.com/support/",
          "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg"
          },
          "tags": ["gcp", "cloudsql", "postgres"],
        "plans": [
          {
            "id": "test_cloudsql_postgres_plan",
            "name": "test_cloudsql_postgres_plan",
            "description": "test-postgres-plan",
            "service_properties": {
              "tier": "db-custom-2-4096",
              "pricing_plan": "PER_USE",
              "max_disk_size": "20"
            },
            "display_name": "test_cloudsql_postgres_plan",
            "service_id": "cbad6d78-a73c-432d-b8ff-b219a17a803a"
          }
        ]
    }
  google-ml-apis: |-
    {
      "id": "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4",
      "description": "Machine Learning Apis including Vision, Translate, Speech, and Natural Language",
      "name": "google-ml-apis",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google Machine Learning APIs",
        "longDescription": "Machine Learning Apis including Vision, Translate, Speech, and Natural Language",
        "documentationUrl": "https://cloud.google.com/ml/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/machine-learning.svg"
      },
      "tags": ["gcp", "ml"],
      "plans":  [
        {
         "id": "be7954e1-ecfb-4936-a0b6-db35e6424c7a",
         "service_id": "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4",
         "name": "default",
         "display_name": "Default",
         "description": "Machine Learning api default plan",
         "service_properties": {}
        }
      ]
    }
  google-stackdriver-debugger: |-
    {
      "id": "83837945-1547-41e0-b661-ea31d76eed11",
      "description": "Stackdriver Debugger",
      "name": "google-stackdriver-debugger",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Stackdriver Debugger",
        "longDescription": "Stackdriver Debugger is a feature of the Google Cloud Platform that lets you inspect the state of an application at any code location without using logging statements and without stopping or slowing down your applications. Your users are not impacted during debugging. Using the production debugger you can capture the local variables and call stack and link it back to a specific line location in your source code.",
        "documentationUrl": "https://cloud.google.com/debugger/docs/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/debugger.svg"
      },
      "tags": ["gcp", "stackdriver", "debugger"],
      "plans": [
        {
          "id": "10866183-a775-49e8-96e3-4e7a901e4a79",
          "service_id": "83837945-1547-41e0-b661-ea31d76eed11",
          "name": "default",
          "display_name": "Default",
          "description": "Stackdriver Debugger default plan",
          "service_properties": {}
        }
      ]
    }
  google-stackdriver-trace: |-
    {
      "id": "c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a",
      "description": "Stackdriver Trace",
      "name": "google-stackdriver-trace",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Stackdriver Trace",
        "longDescription": "Stackdriver Trace is a distributed tracing system that collects latency data from your applications and displays it in the Google Cloud Platform Console. You can track how requests propagate through your application and receive detailed near real-time performance insights.",
        "documentationUrl": "https://cloud.google.com/trace/docs/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/trace.svg"
      },
      "tags": ["gcp", "stackdriver", "trace"],
      "plans": [
        {
          "id": "ab6c2287-b4bc-4ff4-a36a-0575e7910164",
          "service_id": "c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a",
          "name": "default",
          "display_name": "Default",
          "description": "Stackdriver Trace default plan",
          "service_properties": {}
        }
      ]
    }
  google-spanner: |-
    {
      "id": "51b3e27e-d323-49ce-8c5f-1211e6409e82",
      "description": "The first horizontally scalable, globally consistent, relational database service",
      "name": "google-spanner",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google Spanner",
        "longDescription": "The first horizontally scalable, globally consistent, relational database service",
        "documentationUrl": "https://cloud.google.com/spanner/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/spanner.svg"
      },
      "tags": ["gcp", "spanner"],
      "plans": [
        {
            "id": "test-spanner-plan",
            "name": "spannerplan",
            "description": "Basic Spanner plan",
            "service_properties": {
                "num_nodes": "3"
            },
            "display_name": "Spanner Plan",
            "service_id": "51b3e27e-d323-49ce-8c5f-1211e6409e82"
        }
      ]
    }
  google-bigtable: |-
    {
      "id": "b8e19880-ac58-42ef-b033-f7cd9c94d1fe",
      "description": "A high performance NoSQL database service for large analytical and operational workloads",
      "name": "google-bigtable",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
          "displayName": "Google Bigtable",
          "longDescription": "A high performance NoSQL database service for large analytical and operational workloads",
          "documentationUrl": "https://cloud.google.com/bigtable/",
          "supportUrl": "https://cloud.google.com/support/",
          "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/bigtable.svg"
      },
      "tags": ["gcp", "bigtable"],
      "plans": [
        {
          "id": "test-bigtable-plan",
          "name": "bt_plan",
          "description": "Bigtable basic plan",
          "service_properties": {
            "storage_type": "HDD",
            "num_nodes": "5"
          },
          "display_name": "Bigtable Plan",
          "service_id": "b8e19880-ac58-42ef-b033-f7cd9c94d1fe"
        }
      ]
    }
  google-datastore: |-
    {
      "id": "76d4abb2-fee7-4c8f-aee1-bcea2837f02b",
      "description": "Google Cloud Datastore is a NoSQL document database service",
      "name": "google-datastore",
      "bindable": true,
      "plan_updateable": false,
      "metadata": {
        "displayName": "Google Cloud Datastore",
        "longDescription": "Google Cloud Datastore is a NoSQL document database built for automatic scaling, high performance, and ease of application development.",
        "documentationUrl": "https://cloud.google.com/datastore/docs/",
        "supportUrl": "https://cloud.google.com/support/",
        "imageUrl": "https://cloud.google.com/_static/images/cloud/products/logos/svg/datastore.svg"
      },
      "tags": ["gcp", "datastore"],
      "plans": [
        {
         "id": "05f1fb6b-b5f0-48a2-9c2b-a5f236507a97",
         "service_id": "76d4abb2-fee7-4c8f-aee1-bcea2837f02b",
         "name": "default",
         "display_name": "Default",
         "description": "Datastore default plan",
         "service_properties": {}
        }
      ]
    }
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
