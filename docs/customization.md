# Installation Customization

Optionally add these to the env section of `manifest.yml`

## [Optional db vars](#optional-db)

* `DB_PORT` (defaults to 3306)
* `DB_NAME` (defaults to "servicebroker")
* `CA_CERT`
* `CLIENT_CERT`
* `CLIENT_KEY`

## [Optional service vars](#optional-plan)

update the following variables in `manifest.yml` if you wish to enable these services

* `GOOGLE_CLOUDSQL.plans` (A list of json objects with fields `id`, `name`, `description`, `
service_properties` (containing `tier`, `pricing_plan`, `max_disk_size`), `display_name`, and `service_id` 
(Cloud SQL's service id)) - if unset, the service will be disabled. 

e.g.,

```json
[
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
```
* `GOOGLE_BIGTABLE.plans` (A list of json objects with fields `id`, `name`, `description`,
`service_properties` (containing `storage_type`, `num_nodes`), `display_name`, and `service_id` (Bigtable's service id)) 
- if unset, the service will be disabled. 

e.g.,

```json
[
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
```
* `GOOGLE_SPANNER.plans` (A list of json objects with fields `id`, `name`, `description`, `service_properties` (containing 
`num_nodes`), `display_name`, and `service_id` (Spanner's service id)) - if unset, the service will be disabled. 

e.g.,

```json
[
    {
        "id": "test-spanner-plan",
        "name": "spannerplan",
        "description": "Basic Spanner plan",
        "service_properties": {
            "num_nodes": "15"
        },
        "display_name": "Spanner Plan",
        "service_id": "51b3e27e-d323-49ce-8c5f-1211e6409e82"
    }
]
```
