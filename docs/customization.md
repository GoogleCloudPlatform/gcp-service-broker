# Installation Customization

Optionally add these to the env section of `manifest.yml`

## [Optional db vars](#optional-db)

* `DB_PORT` (defaults to 3306)
* `DB_NAME` (defaults to "servicebroker")
* `CA_CERT`
* `CLIENT_CERT`
* `CLIENT_KEY`

## [Optional plan vars](#optional-plan)

* `CLOUDSQL_MYSQL_CUSTOM_PLANS` (A map of plan names to string maps with fields `guid`, `name`, `description`, `authorized_networks`, 
`tier`, `pricing_plan`, `max_disk_size`, `display_name`, and `service` (CloudSQL MySQL's service id)) - if unset, the service
will be disabled. e.g.,

```json
{
    "test_plan": {
        "name": "test_plan",
        "description": "testplan",
	"authorized_networks": "[\"0.0.0.0/0\"]",
        "tier": "D8",
        "pricing_plan": "PER_USE",
        "max_disk_size": "15",
        "display_name": "FOOBAR",
        "service": "4bc59b9a-8520-409f-85da-1c7552315863"
    }
}
```

* `CLOUDSQL_POSTGRES_CUSTOM_PLANS` (A map of plan names to string maps with fields `guid`, `name`, `description`, `authorized_networks`, 
`tier`, `pricing_plan`, `max_disk_size`, `display_name`, and `service` (CloudSQL PostgreSQL's service id)) - if unset, the service
will be disabled. e.g.,

```json
{
    "test_plan": {
        "name": "test_plan",
        "description": "testplan",
	"authorized_networks": "[\"0.0.0.0/0\"]",
        "tier": "custom-db-4-4096",
        "pricing_plan": "PER_USE",
        "max_disk_size": "15",
        "display_name": "FOOBAR",
        "service": "cbad6d78-a73c-432d-b8ff-b219a17a803a"
    }
}
```

* `BIGTABLE_CUSTOM_PLANS` (A map of plan names to string maps with fields `guid`, `name`, `description`,
`storage_type`, `num_nodes`, `display_name`, and `service` (Bigtable's service id)) - if unset, the service
will be disabled. e.g.,

```json
{
    "bt_plan": {
        "name": "bt_plan",
        "description": "Bigtable basic plan",
        "storage_type": "HDD",
        "num_nodes": "5",
        "display_name": "Bigtable Plan",
        "service": "b8e19880-ac58-42ef-b033-f7cd9c94d1fe"
    }
}
```
* `SPANNER_CUSTOM_PLANS` (A map of plan names to string maps with fields `guid`, `name`, `description`,
`num_nodes` `display_name`, and `service` (Spanner's service id)) - if unset, the service
will be disabled. e.g.,

```json
{
    "spannerplan": {
        "name": "spannerplan",
        "description": "Basic Spanner plan",
        "num_nodes": "15",
        "display_name": "Spanner Plan",
        "service": "51b3e27e-d323-49ce-8c5f-1211e6409e82"
    }
}
```
