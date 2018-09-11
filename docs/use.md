
--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/bigquery.svg) Google BigQuery

A fast, economical and fully managed data warehouse for large-scale data analytics.

 * [Documentation](https://cloud.google.com/bigquery/docs/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `f80c0a3e-bd4d-4809-a900-b4e33a6450f1`
 * Tags: gcp, bigquery

## Provisioning

**Request Parameters**


 * `name` _string_ - The name of the BigQuery dataset. Must be alphanumeric (plus underscores) and must be at most 1024 characters long. Default: `a generated value`.


## Binding

**Request Parameters**


 * `role` _string_ - **Required** The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. The following roles are available by default but may be overridden by your operator: 'bigquery.dataViewer', 'bigquery.dataEditor', 'bigquery.dataOwner', 'bigquery.user', 'bigquery.jobUser'.

**Response Parameters**

 * `Email` _string_ - Email address of the service account.
 * `Name` _string_ - The name of the service account.
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account.
 * `UniqueId` _string_ - Unique and stable id of the service account.
 * `dataset_id` _string_ - The name of the BigQuery dataset.

## Plans


  * **default**: BigQuery default plan. Plan ID: `10ff4e72-6e84-44eb-851f-bdb38a791914`.


## Examples




### Basic Configuration


Create a dataset and account that can manage and query the data.
Uses plan: `10ff4e72-6e84-44eb-851f-bdb38a791914`.

**Provision**

```javascript
{
    "name": "orders_1997"
}
```

**Bind**

```javascript
{
    "role": "bigquery.user"
}
```



--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/bigtable.svg) Google Bigtable

A high performance NoSQL database service for large analytical and operational workloads.

 * [Documentation](https://cloud.google.com/bigtable/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `b8e19880-ac58-42ef-b033-f7cd9c94d1fe`
 * Tags: gcp, bigtable

## Provisioning

**Request Parameters**


 * `name` _string_ - The name of the dataset. Should match `[a-z][a-z0-9\-]+[a-z0-9]`. Default: `a generated value`.
 * `cluster_id` _string_ - The name of the cluster. Default: `a generated value`.
 * `display_name` _string_ - The human-readable name of the dataset. Default: `a generated value`.
 * `zone` _string_ - The zone the data will reside in. Default: `us-east1-b`.


## Binding

**Request Parameters**


 * `role` _string_ - **Required** The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. The following roles are available by default but may be overridden by your operator: 'bigtable.user', 'bigtable.reader', 'bigtable.viewer'.

**Response Parameters**

 * `Email` _string_ - Email address of the service account.
 * `Name` _string_ - The name of the service account.
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account.
 * `UniqueId` _string_ - Unique and stable id of the service account.
 * `instance_id` _string_ - The name of the BigTable dataset.

## Plans


  * **three-node-production-hdd**: BigTable HDD basic production plan: Approx: Reads: 1,500 QPS @ 200ms or Writes: 30,000 QPS @ 50ms or Scans: 540 MB/s, 24TB storage. Plan ID: `65a49268-2c73-481e-80f3-9fde5bd5a654`.
  * **three-node-production-ssd**: BigTable SSD basic production plan: Approx: Reads: 30,000 QPS @ 6ms or Writes: 30,000 QPS @ 6ms or Scans: 660 MB/s, 7.5TB storage. Plan ID: `38aa0e65-624b-4998-9c06-f9194b56d252`.


## Examples




### Basic Production Configuration


Create an HDD production table and account that can manage and query the data.
Uses plan: `65a49268-2c73-481e-80f3-9fde5bd5a654`.

**Provision**

```javascript
{
    "name": "orders-table"
}
```

**Bind**

```javascript
{
    "role": "bigtable.user"
}
```



--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg) Google CloudSQL MySQL

Google Cloud SQL is a fully-managed MySQL database service.

 * [Documentation](https://cloud.google.com/sql/docs/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `4bc59b9a-8520-409f-85da-1c7552315863`
 * Tags: gcp, cloudsql, mysql

## Provisioning

**Request Parameters**


 * `instance_name` _string_ - Name of the Cloud SQL instance. Default: `a generated value`.
 * `database_name` _string_ - Name of the database inside of the instance. Default: `a generated value`.
 * `version` _string_ - The database engine type and version. Defaults to `MYSQL_5_6` for 1st gen MySQL instances, `MYSQL_5_7` for 2nd gen MySQL instances, or `POSTGRES_9_6` for PostgreSQL instances.
 * `disk_size` _string_ - In GB (only for 2nd generation instances). Default: `10`.
 * `region` _string_ - The geographical region. Default: `us-central`.
 * `zone` _string_ - (only for 2nd generation instances)
 * `disk_type` _string_ - (only for 2nd generation instances) Default: `ssd`.
 * `failover_replica_name` _string_ - (only for 2nd generation instances) If specified, creates a failover replica with the given name. Default: ``.
 * `maintenance_window_day` _string_ - (only for 2nd generation instances) The day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Day of week (1-7), starting on Monday. Default: `1`.
 * `maintenance_window_hour` _string_ - (only for 2nd generation instances) The hour of the day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Hour of day 0-23. Default: `0`.
 * `backups_enabled` _string_ - Should daily backups be enabled for the service? Default: `true`.
 * `backup_start_time` _string_ - Start time for the daily backup configuration in UTC timezone in the 24 hour format - HH:MM. Default: `06:00`.
 * `binlog` _string_ - Whether binary log is enabled. If backup configuration is disabled, binary log must be disabled as well. Defaults: `false` for 1st gen, `true` for 2nd gen, set to `true` to use.
 * `activation_policy` _string_ - The activation policy specifies when the instance is activated; it is applicable only when the instance state is RUNNABLE. Default: `ON_DEMAND`.
 * `authorized_networks` _string_ - A comma separated list without spaces. Default: `none`.
 * `replication_type` _string_ - The type of replication this instance uses. This can be either ASYNCHRONOUS or SYNCHRONOUS. Default: `SYNCHRONOUS`.
 * `auto_resize` _string_ - (only for 2nd generation instances) Configuration to increase storage size automatically. Default: `false`.


## Binding

**Request Parameters**


 * `role` _string_ - **Required** The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. The following roles are available by default but may be overridden by your operator: 'cloudsql.editor', 'cloudsql.viewer', 'cloudsql.client'.
 * `jdbc_uri_format` _string_ - If `true`, `uri` field will contain a JDBC formatted URI. Default: `false`.
 * `username` _string_ - The SQL username for the account. Default: `a generated value`.
 * `password` _string_ - The SQL password for the account. Default: `a generated value`.

**Response Parameters**

 * `Email` _string_ - Email address of the service account
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account
 * `UniqueId` _string_ - Unique and stable id of the service account
 * `CaCert` _string_ - The server Certificate Authority's certificate.
 * `ClientCert` _string_ - The client certificate. For First Generation instances, the new certificate does not take effect until the instance is restarted.
 * `ClientKey` _string_ - The client certificate key.
 * `Sha1Fingerprint` _string_ - The SHA1 fingerprint of the client certificate.
 * `UriPrefix` _string_ - The connection prefix e.g. `mysql` or `postgres`
 * `Username` _string_ - The name of the SQL user provisioned
 * `database_name` _string_ - The name of the database on the instance
 * `host` _string_ - The hostname or ip of the database instance
 * `instance_name` _string_ - The name of the database instance
 * `uri` _string_ - A database connection string
 * `last_master_operation_id` _string_ - (GCP internals) The id of the last operation on the database.
 * `region` _string_ - The region the database is in.

## Plans


  * **mysql-db-f1-micro**: MySQL on a db-f1-micro (Shared CPUs, 0.6 GB/RAM, 3062 GB/disk, 250 Connections) Plan ID: `7d8f9ade-30c1-4c96-b622-ea0205cc5f0b`.
  * **mysql-db-g1-small**: MySQL on a db-g1-small (Shared CPUs, 1.7 GB/RAM, 3062 GB/disk, 1,000 Connections) Plan ID: `b68bf4d8-1636-4121-af2f-087e46189929`.
  * **mysql-db-n1-standard-1**: MySQL on a db-n1-standard-1 (1 CPUs, 3.75 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `bdfd8033-c2b9-46e9-9b37-1f3a5889eef4`.
  * **mysql-db-n1-standard-2**: MySQL on a db-n1-standard-2 (2 CPUs, 7.5 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `2c99e938-4c1e-4da7-810a-94c9f5b71b57`.
  * **mysql-db-n1-standard-4**: MySQL on a db-n1-standard-4 (4 CPUs, 15 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `d520a5f5-7485-4a83-849b-5439f911fe26`.
  * **mysql-db-n1-standard-8**: MySQL on a db-n1-standard-8 (8 CPUs, 30 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `7ef42bb4-87e3-4ead-8118-4e88c98ed2e6`.
  * **mysql-db-n1-standard-16**: MySQL on a db-n1-standard-16 (16 CPUs, 60 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `200bd90a-4323-46d8-8aa5-afd4601498d0`.
  * **mysql-db-n1-standard-32**: MySQL on a db-n1-standard-32 (32 CPUs, 120 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `52305df2-1e64-4cdb-a4c9-bb5dddb33c3e`.
  * **mysql-db-n1-standard-64**: MySQL on a db-n1-standard-64 (64 CPUs, 240 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `e45d7c44-4990-4dac-a14d-c5127e9ae0c5`.
  * **mysql-db-n1-highmem-2**: MySQL on a db-n1-highmem-2 (2 CPUs, 13 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `07b8a04c-0efe-42d3-8b2c-2c23f7c79583`.
  * **mysql-db-n1-highmem-4**: MySQL on a db-n1-highmem-4 (4 CPUs, 26 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `50fa4baa-e36f-41c3-bbe9-c986d9fbe3c8`.
  * **mysql-db-n1-highmem-8**: MySQL on a db-n1-highmem-8 (8 CPUs, 52 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `6e8e5bc3-bf68-4e57-bda1-d9c9a67faee0`.
  * **mysql-db-n1-highmem-16**: MySQL on a db-n1-highmem-16 (16 CPUs, 104 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `3c83ff6b-165e-47bf-9bba-f4801390d0ff`.
  * **mysql-db-n1-highmem-32**: MySQL on a db-n1-highmem-32 (32 CPUs, 208 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `cbc6d376-8fd3-4a34-9ab5-324311f038f6`.
  * **mysql-db-n1-highmem-64**: MySQL on a db-n1-highmem-64 (64 CPUs, 416 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `b0742cc5-caba-4b8d-98e0-03380ae9522b`.


## Examples




### Development Sandbox


An inexpensive MySQL sandbox for developing with no backups.
Uses plan: `7d8f9ade-30c1-4c96-b622-ea0205cc5f0b`.

**Provision**

```javascript
{
    "backups_enabled": "false",
    "binlog": "false",
    "disk_size": "10"
}
```

**Bind**

```javascript
{
    "role": "cloudsql.editor"
}
```



--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg) Google CloudSQL PostgreSQL

Google Cloud SQL is a fully-managed MySQL database service.

 * [Documentation](https://cloud.google.com/sql/docs/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `cbad6d78-a73c-432d-b8ff-b219a17a803a`
 * Tags: gcp, cloudsql, postgres

## Provisioning

**Request Parameters**


 * `instance_name` _string_ - Name of the Cloud SQL instance. Default: `a generated value`.
 * `database_name` _string_ - Name of the database inside of the instance. Default: `a generated value`.
 * `version` _string_ - The database engine type and version. Defaults to `MYSQL_5_6` for 1st gen MySQL instances, `MYSQL_5_7` for 2nd gen MySQL instances, or `POSTGRES_9_6` for PostgreSQL instances.
 * `disk_size` _string_ - In GB (only for 2nd generation instances). Default: `10`.
 * `region` _string_ - The geographical region. Default: `us-central`.
 * `zone` _string_ - (only for 2nd generation instances)
 * `disk_type` _string_ - (only for 2nd generation instances) Default: `ssd`.
 * `failover_replica_name` _string_ - (only for 2nd generation instances) If specified, creates a failover replica with the given name. Default: ``.
 * `maintenance_window_day` _string_ - (only for 2nd generation instances) The day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Day of week (1-7), starting on Monday. Default: `1`.
 * `maintenance_window_hour` _string_ - (only for 2nd generation instances) The hour of the day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Hour of day 0-23. Default: `0`.
 * `backups_enabled` _string_ - Should daily backups be enabled for the service? Default: `true`.
 * `backup_start_time` _string_ - Start time for the daily backup configuration in UTC timezone in the 24 hour format - HH:MM. Default: `06:00`.
 * `binlog` _string_ - Whether binary log is enabled. If backup configuration is disabled, binary log must be disabled as well. Defaults: `false` for 1st gen, `true` for 2nd gen, set to `true` to use.
 * `activation_policy` _string_ - The activation policy specifies when the instance is activated; it is applicable only when the instance state is RUNNABLE. Default: `ON_DEMAND`.
 * `authorized_networks` _string_ - A comma separated list without spaces. Default: `none`.
 * `replication_type` _string_ - The type of replication this instance uses. This can be either ASYNCHRONOUS or SYNCHRONOUS. Default: `SYNCHRONOUS`.
 * `auto_resize` _string_ - (only for 2nd generation instances) Configuration to increase storage size automatically. Default: `false`.


## Binding

**Request Parameters**


 * `role` _string_ - **Required** The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. The following roles are available by default but may be overridden by your operator: 'cloudsql.editor', 'cloudsql.viewer', 'cloudsql.client'.
 * `jdbc_uri_format` _string_ - If `true`, `uri` field will contain a JDBC formatted URI. Default: `false`.
 * `username` _string_ - The SQL username for the account. Default: `a generated value`.
 * `password` _string_ - The SQL password for the account. Default: `a generated value`.

**Response Parameters**

 * `Email` _string_ - Email address of the service account
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account
 * `UniqueId` _string_ - Unique and stable id of the service account
 * `CaCert` _string_ - The server Certificate Authority's certificate.
 * `ClientCert` _string_ - The client certificate. For First Generation instances, the new certificate does not take effect until the instance is restarted.
 * `ClientKey` _string_ - The client certificate key.
 * `Sha1Fingerprint` _string_ - The SHA1 fingerprint of the client certificate.
 * `UriPrefix` _string_ - The connection prefix e.g. `mysql` or `postgres`
 * `Username` _string_ - The name of the SQL user provisioned
 * `database_name` _string_ - The name of the database on the instance
 * `host` _string_ - The hostname or ip of the database instance
 * `instance_name` _string_ - The name of the database instance
 * `uri` _string_ - A database connection string
 * `last_master_operation_id` _string_ - (GCP internals) The id of the last operation on the database.
 * `region` _string_ - The region the database is in.

## Plans


  * **postgres-db-f1-micro**: PostgreSQL on a db-f1-micro (Shared CPUs, 0.6 GB/RAM, 3062 GB/disk, 250 Connections) Plan ID: `2513d4d9-684b-4c3c-add4-6404969006de`.
  * **postgres-db-g1-small**: PostgreSQL on a db-g1-small (Shared CPUs, 1.7 GB/RAM, 3062 GB/disk, 1,000 Connections) Plan ID: `6c1174d8-243c-44d1-b7a8-e94a779f67f5`.
  * **postgres-db-n1-standard-1**: PostgreSQL on a db-n1-standard-1 (1 CPUs, 3.75 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `c4e68ab5-34ca-4d02-857d-3e6b3ab079a7`.
  * **postgres-db-n1-standard-2**: PostgreSQL on a db-n1-standard-2 (2 CPUs, 7.5 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `3f578ecf-885c-4b60-b38b-60272f34e00f`.
  * **postgres-db-n1-standard-4**: PostgreSQL on a db-n1-standard-4 (4 CPUs, 15 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `b7fcab5d-d66d-4e82-af16-565e84cef7f9`.
  * **postgres-db-n1-standard-8**: PostgreSQL on a db-n1-standard-8 (8 CPUs, 30 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `4b2fa14a-caf1-42e0-bd8c-3342502008a8`.
  * **postgres-db-n1-standard-16**: PostgreSQL on a db-n1-standard-16 (16 CPUs, 60 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `ca2e770f-bfa5-4fb7-a249-8b943c3474ca`.
  * **postgres-db-n1-standard-32**: PostgreSQL on a db-n1-standard-32 (32 CPUs, 120 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `b44f8294-b003-4a50-80c2-706858073f44`.
  * **postgres-db-n1-standard-64**: PostgreSQL on a db-n1-standard-64 (64 CPUs, 240 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `d97326e0-5af2-4da5-b970-b4772d59cded`.
  * **postgres-db-n1-highmem-2**: PostgreSQL on a db-n1-highmem-2 (2 CPUs, 13 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `c10f8691-02f5-44eb-989f-7217393012ca`.
  * **postgres-db-n1-highmem-4**: PostgreSQL on a db-n1-highmem-4 (4 CPUs, 26 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `610cc78d-d26a-41a9-90b7-547a44517f03`.
  * **postgres-db-n1-highmem-8**: PostgreSQL on a db-n1-highmem-8 (8 CPUs, 52 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `2a351e8d-958d-4c4f-ae46-c984fec18740`.
  * **postgres-db-n1-highmem-16**: PostgreSQL on a db-n1-highmem-16 (16 CPUs, 104 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `51d3ca0c-9d21-447d-a395-3e0dc0659775`.
  * **postgres-db-n1-highmem-32**: PostgreSQL on a db-n1-highmem-32 (32 CPUs, 208 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `2e72b386-f7ce-4f0d-a149-9f9a851337d4`.
  * **postgres-db-n1-highmem-64**: PostgreSQL on a db-n1-highmem-64 (64 CPUs, 416 GB/RAM, 10230 GB/disk, 4,000 Connections) Plan ID: `82602649-e4ac-4a2f-b80d-dacd745aed6a`.


## Examples




### Development Sandbox


An inexpensive PostgreSQL sandbox for developing with no backups.
Uses plan: `2513d4d9-684b-4c3c-add4-6404969006de`.

**Provision**

```javascript
{
    "backups_enabled": "false",
    "binlog": "false",
    "disk_size": "10"
}
```

**Bind**

```javascript
{
    "role": "cloudsql.editor"
}
```



--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/datastore.svg) Google Cloud Datastore

Google Cloud Datastore is a NoSQL document database built for automatic scaling, high performance, and ease of application development.

 * [Documentation](https://cloud.google.com/datastore/docs/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `76d4abb2-fee7-4c8f-aee1-bcea2837f02b`
 * Tags: gcp, datastore

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `Email` _string_ - Email address of the service account.
 * `Name` _string_ - The name of the service account.
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account.
 * `UniqueId` _string_ - Unique and stable id of the service account.

## Plans


  * **default**: Datastore default plan. Plan ID: `05f1fb6b-b5f0-48a2-9c2b-a5f236507a97`.


## Examples




### Basic Configuration


Creates a datastore and a user with the permission `datastore.user`.
Uses plan: `05f1fb6b-b5f0-48a2-9c2b-a5f236507a97`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{}
```



--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/machine-learning.svg) Google Machine Learning APIs

Machine Learning APIs including Vision, Translate, Speech, and Natural Language.

 * [Documentation](https://cloud.google.com/ml/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `5ad2dce0-51f7-4ede-8b46-293d6df1e8d4`
 * Tags: gcp, ml

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**


 * `role` _string_ - **Required** The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. The following roles are available by default but may be overridden by your operator: 'ml.developer', 'ml.viewer', 'ml.modelOwner', 'ml.modelUser', 'ml.jobOwner', 'ml.operationOwner'.

**Response Parameters**

 * `Email` _string_ - Email address of the service account.
 * `Name` _string_ - The name of the service account.
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account.
 * `UniqueId` _string_ - Unique and stable id of the service account.

## Plans


  * **default**: Machine Learning API default plan. Plan ID: `be7954e1-ecfb-4936-a0b6-db35e6424c7a`.


## Examples




### Basic Configuration


Create an account with developer access to your ML models.
Uses plan: `be7954e1-ecfb-4936-a0b6-db35e6424c7a`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{
    "role": "ml.developer"
}
```



--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/pubsub.svg) Google PubSub

A global service for real-time and reliable messaging and streaming data.

 * [Documentation](https://cloud.google.com/pubsub/docs/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `628629e3-79f5-4255-b981-d14c6c7856be`
 * Tags: gcp, pubsub

## Provisioning

**Request Parameters**


 * `topic_name` _string_ - Name of the topic. Default: `a generated value`.
 * `subscription_name` _string_ - **Required** Name of the subscription.
 * `is_push` _string_ - Are events handled by POSTing to a URL? Default: `false`.
 * `endpoint` _string_ - If `is_push` == 'true', then this is the URL that will be pused to. Default: ``.
 * `ack_deadline` _string_ - Value is in seconds. Max: 600 This is the maximum time after a subscriber receives a message before the subscriber should acknowledge the message. After message delivery but before the ack deadline expires and before the message is acknowledged, it is an outstanding message and will not be delivered again during that time (on a best-effort basis).  Default: `10`.


## Binding

**Request Parameters**


 * `role` _string_ - **Required** The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. The following roles are available by default but may be overridden by your operator: 'pubsub.publisher', 'pubsub.subscriber', 'pubsub.viewer', 'pubsub.editor'.

**Response Parameters**

 * `Email` _string_ - Email address of the service account.
 * `Name` _string_ - The name of the service account.
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account.
 * `UniqueId` _string_ - Unique and stable id of the service account.
 * `subscription_name` _string_ - Name of the subscription.
 * `topic_name` _string_ - Name of the topic.

## Plans


  * **default**: PubSub Default plan. Plan ID: `622f4da3-8731-492a-af29-66a9146f8333`.


## Examples




### Basic Configuration


Create a topic and a publisher to it.
Uses plan: `622f4da3-8731-492a-af29-66a9146f8333`.

**Provision**

```javascript
{
    "subscription_name": "example_topic_subscription",
    "topic_name": "example_topic"
}
```

**Bind**

```javascript
{
    "role": "pubsub.publisher"
}
```



--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/spanner.svg) Google Spanner

The first horizontally scalable, globally consistent, relational database service.

 * [Documentation](https://cloud.google.com/spanner/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `51b3e27e-d323-49ce-8c5f-1211e6409e82`
 * Tags: gcp, spanner

## Provisioning

**Request Parameters**


 * `name` _string_ - The name of the instance. Default: `a generated value`.
 * `display_name` _string_ - A human-readable name for the instance. Default: `a generated value`.
 * `location` _string_ - The location of the Spanner instance. Default: `regional-us-central1`.


## Binding

**Request Parameters**


 * `role` _string_ - **Required** The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. The following roles are available by default but may be overridden by your operator: 'spanner.databaseAdmin', 'spanner.databaseReader', 'spanner.databaseUser', 'spanner.viewer'.

**Response Parameters**

 * `Email` _string_ - Email address of the service account.
 * `Name` _string_ - The name of the service account.
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account.
 * `UniqueId` _string_ - Unique and stable id of the service account.
 * `instance_id` _string_ - Name of the spanner instance the account can connect to.

## Plans


  * **sandbox**: Useful for testing, not eligible for SLA. Plan ID: `44828436-cfbd-47ae-b4bc-48854564347b`.
  * **minimal-production**: A minimal production level Spanner setup eligible for 99.99% SLA. Each node can provide up to 10,000 QPS of reads or 2,000 QPS of writes (writing single rows at 1KB data per row), and 2 TiB storage. Plan ID: `0752b1ad-a784-4dcc-96eb-64149089a1c9`.


## Examples




### Basic Configuration


Create a sandbox environment with a database admin account.
Uses plan: `44828436-cfbd-47ae-b4bc-48854564347b`.

**Provision**

```javascript
{
    "name": "auth-database"
}
```

**Bind**

```javascript
{
    "role": "spanner.databaseAdmin"
}
```



--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/debugger.svg) Stackdriver Debugger

Stackdriver Debugger is a feature of the Google Cloud Platform that lets you inspect the state of an application at any code location without using logging statements and without stopping or slowing down your applications. Your users are not impacted during debugging. Using the production debugger you can capture the local variables and call stack and link it back to a specific line location in your source code.

 * [Documentation](https://cloud.google.com/debugger/docs/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `83837945-1547-41e0-b661-ea31d76eed11`
 * Tags: gcp, stackdriver, debugger

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `Email` _string_ - Email address of the service account.
 * `Name` _string_ - The name of the service account.
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account.
 * `UniqueId` _string_ - Unique and stable id of the service account.

## Plans


  * **default**: Stackdriver Debugger default plan. Plan ID: `10866183-a775-49e8-96e3-4e7a901e4a79`.


## Examples




### Basic Configuration


Creates an account with the permission `clouddebugger.agent`.
Uses plan: `10866183-a775-49e8-96e3-4e7a901e4a79`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{}
```



--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/stackdriver.svg) Stackdriver Profiler

Continuous CPU and heap profiling to improve performance and reduce costs.

 * [Documentation](https://cloud.google.com/profiler/docs/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `00b9ca4a-7cd6-406a-a5b7-2f43f41ade75`
 * Tags: gcp, stackdriver, profiler

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `Email` _string_ - Email address of the service account.
 * `Name` _string_ - The name of the service account.
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account.
 * `UniqueId` _string_ - Unique and stable id of the service account.

## Plans


  * **default**: Stackdriver Profiler default plan. Plan ID: `594627f6-35f5-462f-9074-10fb033fb18a`.


## Examples




### Basic Configuration


Creates an account with the permission `cloudprofiler.agent`.
Uses plan: `594627f6-35f5-462f-9074-10fb033fb18a`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{}
```



--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/trace.svg) Stackdriver Trace

Stackdriver Trace is a distributed tracing system that collects latency data from your applications and displays it in the Google Cloud Platform Console. You can track how requests propagate through your application and receive detailed near real-time performance insights.

 * [Documentation](https://cloud.google.com/trace/docs/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a`
 * Tags: gcp, stackdriver, trace

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `Email` _string_ - Email address of the service account.
 * `Name` _string_ - The name of the service account.
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account.
 * `UniqueId` _string_ - Unique and stable id of the service account.

## Plans


  * **default**: Stackdriver Trace default plan. Plan ID: `ab6c2287-b4bc-4ff4-a36a-0575e7910164`.


## Examples




### Basic Configuration


Creates an account with the permission `cloudtrace.agent`.
Uses plan: `ab6c2287-b4bc-4ff4-a36a-0575e7910164`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{}
```



--------------------------------------------------------------------------------

# ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/storage.svg) Google Cloud Storage

Unified object storage for developers and enterprises. Cloud Storage allows world-wide storage and retrieval of any amount of data at any time.

 * [Documentation](https://cloud.google.com/storage/docs/overview)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `b9e4332e-b42b-4680-bda5-ea1506797474`
 * Tags: gcp, storage

## Provisioning

**Request Parameters**


 * `name` _string_ - The name of the bucket. There is a single global namespace shared by all buckets so it MUST be unique. Default: `a generated value`.
 * `location` _string_ - The location of the bucket. Object data for objects in the bucket resides in physical storage within this region. See: https://cloud.google.com/storage/docs/bucket-locations Default: `US`.


## Binding

**Request Parameters**


 * `role` _string_ - **Required** The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. The following roles are available by default but may be overridden by your operator: 'storage.objectCreator', 'storage.objectViewer', 'storage.objectAdmin'.

**Response Parameters**

 * `Email` _string_ - Email address of the service account.
 * `Name` _string_ - The name of the service account.
 * `PrivateKeyData` _string_ - Service account private key data. Base-64 encoded JSON.
 * `ProjectId` _string_ - ID of the project that owns the service account.
 * `UniqueId` _string_ - Unique and stable id of the service account.
 * `bucket_name` _string_ - Name of the bucket this binding is for

## Plans


  * **standard**: Standard storage class. Plan ID: `e1d11f65-da66-46ad-977c-6d56513baf43`.
  * **nearline**: Nearline storage class. Plan ID: `a42c1182-d1a0-4d40-82c1-28220518b360`.
  * **reduced-availability**: Durable Reduced Availability storage class. Plan ID: `1a1f4fe6-1904-44d0-838c-4c87a9490a6b`.


## Examples




### Basic Configuration


Create a nearline bucket with a service account that can create/read/delete the objects in it.
Uses plan: `a42c1182-d1a0-4d40-82c1-28220518b360`.

**Provision**

```javascript
{
    "location": "us"
}
```

**Bind**

```javascript
{
    "role": "storage.objectAdmin"
}
```



