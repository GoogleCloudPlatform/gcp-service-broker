# <a name="google-bigquery"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/bigquery.svg) Google BigQuery
A fast, economical and fully managed data warehouse for large-scale data analytics.

 * [Documentation](https://cloud.google.com/bigquery/docs/)
 * [Support](https://cloud.google.com/bigquery/support)
 * Catalog Metadata ID: `f80c0a3e-bd4d-4809-a900-b4e33a6450f1`
 * Tags: gcp, bigquery
 * Service Name: `google-bigquery`

## Provisioning

**Request Parameters**


 * `name` _string_ - The name of the BigQuery dataset. Default: `pcf_sb_${counter.next()}_${time.nano()}`.
    * The string must have at most 1024 characters.
    * The string must match the regular expression `^[A-Za-z0-9_]+$`.
 * `location` _string_ - The location of the BigQuery instance. Default: `US`.
    * Examples: [US EU asia-northeast1].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `bigquery.user`.
    * The value must be one of: [bigquery.dataEditor bigquery.dataOwner bigquery.dataViewer bigquery.jobUser bigquery.user].

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].
 * `dataset_id` _string_ - **Required** The name of the BigQuery dataset.
    * The string must have at most 1024 characters.
    * The string must match the regular expression `^[A-Za-z0-9_]+$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `10ff4e72-6e84-44eb-851f-bdb38a791914`.
  * Description: BigQuery default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


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

**Cloud Foundry Example**

<pre>
$ cf create-service google-bigquery default my-google-bigquery-example -c `{"name":"orders_1997"}`
$ cf bind-service my-app my-google-bigquery-example -c `{"role":"bigquery.user"}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-bigtable"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/bigtable.svg) Google Bigtable
A high performance NoSQL database service for large analytical and operational workloads.

 * [Documentation](https://cloud.google.com/bigtable/)
 * [Support](https://cloud.google.com/bigtable/docs/support/getting-support)
 * Catalog Metadata ID: `b8e19880-ac58-42ef-b033-f7cd9c94d1fe`
 * Tags: gcp, bigtable
 * Service Name: `google-bigtable`

## Provisioning

**Request Parameters**


 * `name` _string_ - The name of the Cloud Bigtable instance. Default: `pcf-sb-${counter.next()}-${time.nano()}`.
    * The string must have at most 33 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z][-0-9a-z]+$`.
 * `cluster_id` _string_ - The ID of the Cloud Bigtable cluster. Default: `${str.truncate(20, name)}-cluster`.
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z][-0-9a-z]+[a-z]$`.
 * `display_name` _string_ - The human-readable display name of the Bigtable instance. Default: `${name}`.
    * The string must have at most 30 characters.
    * The string must have at least 4 characters.
 * `zone` _string_ - The zone to create the Cloud Bigtable cluster in. Zones that support Bigtable instances are noted on the Cloud Bigtable locations page: https://cloud.google.com/bigtable/docs/locations. Default: `us-east1-b`.
    * Examples: [us-central1-a europe-west2-b asia-northeast1-a australia-southeast1-c].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `bigtable.user`.
    * The value must be one of: [bigtable.reader bigtable.user bigtable.viewer].

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].
 * `instance_id` _string_ - **Required** The name of the BigTable dataset.
    * The string must have at most 33 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z][-0-9a-z]+$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`three-node-production-hdd`**
  * Plan ID: `65a49268-2c73-481e-80f3-9fde5bd5a654`.
  * Description: BigTable HDD basic production plan: Approx: Reads: 1,500 QPS @ 200ms or Writes: 30,000 QPS @ 50ms or Scans: 540 MB/s, 24TB storage.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`three-node-production-ssd`**
  * Plan ID: `38aa0e65-624b-4998-9c06-f9194b56d252`.
  * Description: BigTable SSD basic production plan: Approx: Reads: 30,000 QPS @ 6ms or Writes: 30,000 QPS @ 6ms or Scans: 660 MB/s, 7.5TB storage.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


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

**Cloud Foundry Example**

<pre>
$ cf create-service google-bigtable three-node-production-hdd my-google-bigtable-example -c `{"name":"orders-table"}`
$ cf bind-service my-app my-google-bigtable-example -c `{"role":"bigtable.user"}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-cloudsql-mysql"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg) Google CloudSQL for MySQL
Google CloudSQL for MySQL is a fully-managed MySQL database service.

 * [Documentation](https://cloud.google.com/sql/docs/)
 * [Support](https://cloud.google.com/sql/docs/getting-support/)
 * Catalog Metadata ID: `4bc59b9a-8520-409f-85da-1c7552315863`
 * Tags: gcp, cloudsql, mysql
 * Service Name: `google-cloudsql-mysql`

## Provisioning

**Request Parameters**


 * `instance_name` _string_ - Name of the CloudSQL instance. Default: `sb-${counter.next()}-${time.nano()}`.
    * The string must have at most 84 characters.
    * The string must match the regular expression `^[a-z][a-z0-9-]+$`.
 * `database_name` _string_ - Name of the database inside of the instance. Must be a valid identifier for your chosen database type. Default: `sb-${counter.next()}-${time.nano()}`.
 * `version` _string_ - The database engine type and version. Default: `MYSQL_5_7`.
    * The value must be one of: [MYSQL_5_6 MYSQL_5_7].
 * `activation_policy` _string_ - The activation policy specifies when the instance is activated; it is applicable only when the instance state is RUNNABLE. Default: `ALWAYS`.
    * The value must be one of: [ALWAYS NEVER].
 * `binlog` _string_ - Whether binary log is enabled. Must be enabled for high availability. Default: `true`.
    * The value must be one of: [false true].
 * `authorized_networks` _string_ - A comma separated list without spaces. Default: ``.
 * `region` _string_ - The geographical region. See the instance locations list https://cloud.google.com/sql/docs/mysql/instance-locations for which regions support which databases. Default: `us-central`.
    * Examples: [northamerica-northeast1 southamerica-east1 us-east1].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.
 * `disk_size` _string_ - In GB. Default: `10`.
    * Examples: [10 500 10230].
    * The string must have at most 5 characters.
    * The string must match the regular expression `^[1-9][0-9]+$`.
 * `database_flags` _string_ - The database flags passed to the instance at startup (comma separated list of flags, e.g. general_log=on,skip_show_database=off). Default: ``.
    * Examples: [long_query_time=10 general_log=on,skip_show_database=off].
    * The string must match the regular expression `^(|([a-z_]+=[a-zA-Z0-9\.\+\:-]+)(,[a-z_]+=[a-zA-Z0-9\.\+\:-]+)*)$`.
 * `zone` _string_ - Optional, the specific zone in the region to run the instance. Default: ``.
    * The string must match the regular expression `^(|[A-Za-z][-a-z0-9A-Z]+)$`.
 * `disk_type` _string_ - The type of disk backing the database. Default: `PD_SSD`.
    * The value must be one of: [PD_HDD PD_SSD].
 * `maintenance_window_day` _string_ - The day of week a CloudSQL instance should preferably be restarted for system maintenance purposes. (1-7), starting on Monday. Default: `1`.
    * The value must be one of: [1 2 3 4 5 6 7].
 * `maintenance_window_hour` _string_ - The hour of the day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Hour of day 0-23. Default: `0`.
    * The string must match the regular expression `^([0-9]|1[0-9]|2[0-3])$`.
 * `backups_enabled` _string_ - Should daily backups be enabled for the service? Default: `true`.
    * The value must be one of: [false true].
 * `backup_start_time` _string_ - Start time for the daily backup configuration in UTC timezone in the 24 hour format - HH:MM. Default: `06:00`.
    * The string must match the regular expression `^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$`.
 * `replication_type` _string_ - The type of replication this instance uses. This can be either ASYNCHRONOUS or SYNCHRONOUS. Default: `SYNCHRONOUS`.
    * The value must be one of: [ASYNCHRONOUS SYNCHRONOUS].
 * `auto_resize` _string_ - Configuration to increase storage size automatically. Default: `false`.
    * The value must be one of: [false true].
 * `auto_resize_limit` _string_ - The maximum size to which storage capacity can be automatically increased. Default: `0`.
    * Examples: [10 500 10230].
    * The string must have at most 5 characters.
    * The string must match the regular expression `^[0-9][0-9]*$`.
 * `availability_type` _string_ - Availability type specifies whether the instance serves data from multiple zones. Default: `ZONAL`.
    * The value must be one of: [REGIONAL ZONAL].


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `cloudsql.client`.
    * The value must be one of: [cloudsql.client cloudsql.editor cloudsql.viewer].
 * `jdbc_uri_format` _string_ - If `true`, `uri` field will contain a JDBC formatted URI. Default: `false`.
    * The value must be one of: [false true].
 * `username` _string_ - The SQL username for the account. Default: `sb${str.truncate(14, time.nano())}`.
 * `password` _string_ - The SQL password for the account. Default: `${rand.base64(32)}`.

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].
 * `CaCert` _string_ - **Required** The server Certificate Authority's certificate.
    * Examples: [-----BEGIN CERTIFICATE-----BASE64 Certificate Text-----END CERTIFICATE-----].
 * `ClientCert` _string_ - **Required** The client certificate.
    * Examples: [-----BEGIN CERTIFICATE-----BASE64 Certificate Text-----END CERTIFICATE-----].
 * `ClientKey` _string_ - **Required** The client certificate key.
    * Examples: [-----BEGIN RSA PRIVATE KEY-----BASE64 Key Text-----END RSA PRIVATE KEY-----].
 * `Sha1Fingerprint` _string_ - **Required** The SHA1 fingerprint of the client certificate.
    * Examples: [e6d0c68f35032c6c2132217d1f1fb06b12ed32e2].
    * The string must match the regular expression `^[0-9a-f]{40}$`.
 * `UriPrefix` _string_ - The connection prefix.
    * Examples: [jdbc: ].
 * `Username` _string_ - **Required** The name of the SQL user provisioned.
    * Examples: [sb15404128767777].
 * `Password` _string_ - **Required** The database password for the SQL user.
    * Examples: [N-JPz7h2RHPZ81jB5gDHdnluddnIFMWG4nd5rKjR_8A=].
 * `database_name` _string_ - **Required** The name of the database on the instance.
    * Examples: [sb-2-1540412407295372465].
 * `host` _string_ - **Required** The hostname or IP address of the database instance.
    * Examples: [127.0.0.1].
 * `instance_name` _string_ - **Required** The name of the database instance.
    * Examples: [sb-1-1540412407295273023].
    * The string must have at most 84 characters.
    * The string must match the regular expression `^[a-z][a-z0-9-]+$`.
 * `uri` _string_ - **Required** A database connection string.
    * Examples: [mysql://user:pass@127.0.0.1/sb-2-1540412407295372465?ssl_mode=required].
 * `last_master_operation_id` _string_ - (deprecated) The id of the last operation on the database.
    * Examples: [mysql://user:pass@127.0.0.1/sb-2-1540412407295372465?ssl_mode=required].
 * `region` _string_ - **Required** The region the database is in.
    * Examples: [northamerica-northeast1 southamerica-east1 us-east1].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`mysql-db-f1-micro`**
  * Plan ID: `7d8f9ade-30c1-4c96-b622-ea0205cc5f0b`.
  * Description: MySQL on a db-f1-micro (Shared CPUs, 0.6 GB/RAM, 3062 GB/disk, 250 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-g1-small`**
  * Plan ID: `b68bf4d8-1636-4121-af2f-087e46189929`.
  * Description: MySQL on a db-g1-small (Shared CPUs, 1.7 GB/RAM, 3062 GB/disk, 1,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-standard-1`**
  * Plan ID: `bdfd8033-c2b9-46e9-9b37-1f3a5889eef4`.
  * Description: MySQL on a db-n1-standard-1 (1 CPUs, 3.75 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-standard-2`**
  * Plan ID: `2c99e938-4c1e-4da7-810a-94c9f5b71b57`.
  * Description: MySQL on a db-n1-standard-2 (2 CPUs, 7.5 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-standard-4`**
  * Plan ID: `d520a5f5-7485-4a83-849b-5439f911fe26`.
  * Description: MySQL on a db-n1-standard-4 (4 CPUs, 15 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-standard-8`**
  * Plan ID: `7ef42bb4-87e3-4ead-8118-4e88c98ed2e6`.
  * Description: MySQL on a db-n1-standard-8 (8 CPUs, 30 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-standard-16`**
  * Plan ID: `200bd90a-4323-46d8-8aa5-afd4601498d0`.
  * Description: MySQL on a db-n1-standard-16 (16 CPUs, 60 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-standard-32`**
  * Plan ID: `52305df2-1e64-4cdb-a4c9-bb5dddb33c3e`.
  * Description: MySQL on a db-n1-standard-32 (32 CPUs, 120 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-standard-64`**
  * Plan ID: `e45d7c44-4990-4dac-a14d-c5127e9ae0c5`.
  * Description: MySQL on a db-n1-standard-64 (64 CPUs, 240 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-highmem-2`**
  * Plan ID: `07b8a04c-0efe-42d3-8b2c-2c23f7c79583`.
  * Description: MySQL on a db-n1-highmem-2 (2 CPUs, 13 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-highmem-4`**
  * Plan ID: `50fa4baa-e36f-41c3-bbe9-c986d9fbe3c8`.
  * Description: MySQL on a db-n1-highmem-4 (4 CPUs, 26 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-highmem-8`**
  * Plan ID: `6e8e5bc3-bf68-4e57-bda1-d9c9a67faee0`.
  * Description: MySQL on a db-n1-highmem-8 (8 CPUs, 52 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-highmem-16`**
  * Plan ID: `3c83ff6b-165e-47bf-9bba-f4801390d0ff`.
  * Description: MySQL on a db-n1-highmem-16 (16 CPUs, 104 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-highmem-32`**
  * Plan ID: `cbc6d376-8fd3-4a34-9ab5-324311f038f6`.
  * Description: MySQL on a db-n1-highmem-32 (32 CPUs, 208 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`mysql-db-n1-highmem-64`**
  * Plan ID: `b0742cc5-caba-4b8d-98e0-03380ae9522b`.
  * Description: MySQL on a db-n1-highmem-64 (64 CPUs, 416 GB/RAM, 10230 GB/disk, 4,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### HA Instance


A regionally available database with automatic failover.
Uses plan: `7d8f9ade-30c1-4c96-b622-ea0205cc5f0b`.

**Provision**

```javascript
{
    "availability_type": "REGIONAL",
    "backups_enabled": "true",
    "binlog": "true"
}
```

**Bind**

```javascript
{
    "role": "cloudsql.editor"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-cloudsql-mysql mysql-db-f1-micro my-google-cloudsql-mysql-example -c `{"availability_type":"REGIONAL","backups_enabled":"true","binlog":"true"}`
$ cf bind-service my-app my-google-cloudsql-mysql-example -c `{"role":"cloudsql.editor"}`
</pre>


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

**Cloud Foundry Example**

<pre>
$ cf create-service google-cloudsql-mysql mysql-db-f1-micro my-google-cloudsql-mysql-example -c `{"backups_enabled":"false","binlog":"false","disk_size":"10"}`
$ cf bind-service my-app my-google-cloudsql-mysql-example -c `{"role":"cloudsql.editor"}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-cloudsql-mysql-vpc"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg) Google CloudSQL for MySQL
Google CloudSQL for MySQL is a fully-managed MySQL database service.

 * [Documentation](https://cloud.google.com/sql/docs/)
 * [Support](https://cloud.google.com/sql/docs/getting-support/)
 * Catalog Metadata ID: `b48d2a6b-b1b0-499f-8389-57ba33bfbb19`
 * Tags: gcp, cloudsql, mysql
 * Service Name: `google-cloudsql-mysql-vpc`

## Provisioning

**Request Parameters**


 * `instance_name` _string_ - Name of the CloudSQL instance. Default: `sb-${counter.next()}-${time.nano()}`.
    * The string must have at most 84 characters.
    * The string must match the regular expression `^[a-z][a-z0-9-]+$`.
 * `database_name` _string_ - Name of the database inside of the instance. Must be a valid identifier for your chosen database type. Default: `sb-${counter.next()}-${time.nano()}`.
 * `version` _string_ - The database engine type and version. Default: `MYSQL_5_7`.
    * The value must be one of: [MYSQL_5_6 MYSQL_5_7].
 * `binlog` _string_ - Whether binary log is enabled. Must be enabled for high availability. Default: `true`.
    * The value must be one of: [false true].
 * `tier` _string_ - The machine type the database will run on. MySQL has predefined tiers, other databases use the a string of the form db-custom-[CPUS]-[MEMORY_MBS], where memory is at least 3840.
    * Examples: [db-n1-standard-1 db-custom-1-3840].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.
 * `private_network` _string_ - The private network to attach to. If specified the instance will only be accessible on the VPC. Default: `default`.
    * Examples: [projects/my-project/global/networks/default].
 * `region` _string_ - The geographical region. See the instance locations list https://cloud.google.com/sql/docs/mysql/instance-locations for which regions support which databases. Default: `us-central`.
    * Examples: [northamerica-northeast1 southamerica-east1 us-east1].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.
 * `disk_size` _string_ - In GB. Default: `10`.
    * Examples: [10 500 10230].
    * The string must have at most 5 characters.
    * The string must match the regular expression `^[1-9][0-9]+$`.
 * `database_flags` _string_ - The database flags passed to the instance at startup (comma separated list of flags, e.g. general_log=on,skip_show_database=off). Default: ``.
    * Examples: [long_query_time=10 general_log=on,skip_show_database=off].
    * The string must match the regular expression `^(|([a-z_]+=[a-zA-Z0-9\.\+\:-]+)(,[a-z_]+=[a-zA-Z0-9\.\+\:-]+)*)$`.
 * `zone` _string_ - Optional, the specific zone in the region to run the instance. Default: ``.
    * The string must match the regular expression `^(|[A-Za-z][-a-z0-9A-Z]+)$`.
 * `disk_type` _string_ - The type of disk backing the database. Default: `PD_SSD`.
    * The value must be one of: [PD_HDD PD_SSD].
 * `maintenance_window_day` _string_ - The day of week a CloudSQL instance should preferably be restarted for system maintenance purposes. (1-7), starting on Monday. Default: `1`.
    * The value must be one of: [1 2 3 4 5 6 7].
 * `maintenance_window_hour` _string_ - The hour of the day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Hour of day 0-23. Default: `0`.
    * The string must match the regular expression `^([0-9]|1[0-9]|2[0-3])$`.
 * `backups_enabled` _string_ - Should daily backups be enabled for the service? Default: `true`.
    * The value must be one of: [false true].
 * `backup_start_time` _string_ - Start time for the daily backup configuration in UTC timezone in the 24 hour format - HH:MM. Default: `06:00`.
    * The string must match the regular expression `^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$`.
 * `replication_type` _string_ - The type of replication this instance uses. This can be either ASYNCHRONOUS or SYNCHRONOUS. Default: `SYNCHRONOUS`.
    * The value must be one of: [ASYNCHRONOUS SYNCHRONOUS].
 * `auto_resize` _string_ - Configuration to increase storage size automatically. Default: `false`.
    * The value must be one of: [false true].
 * `auto_resize_limit` _string_ - The maximum size to which storage capacity can be automatically increased. Default: `0`.
    * Examples: [10 500 10230].
    * The string must have at most 5 characters.
    * The string must match the regular expression `^[0-9][0-9]*$`.
 * `availability_type` _string_ - Availability type specifies whether the instance serves data from multiple zones. Default: `ZONAL`.
    * The value must be one of: [REGIONAL ZONAL].


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `cloudsql.client`.
    * The value must be one of: [cloudsql.client cloudsql.editor cloudsql.viewer].
 * `jdbc_uri_format` _string_ - If `true`, `uri` field will contain a JDBC formatted URI. Default: `false`.
    * The value must be one of: [false true].
 * `username` _string_ - The SQL username for the account. Default: `sb${str.truncate(14, time.nano())}`.
 * `password` _string_ - The SQL password for the account. Default: `${rand.base64(32)}`.

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].
 * `CaCert` _string_ - **Required** The server Certificate Authority's certificate.
    * Examples: [-----BEGIN CERTIFICATE-----BASE64 Certificate Text-----END CERTIFICATE-----].
 * `ClientCert` _string_ - **Required** The client certificate.
    * Examples: [-----BEGIN CERTIFICATE-----BASE64 Certificate Text-----END CERTIFICATE-----].
 * `ClientKey` _string_ - **Required** The client certificate key.
    * Examples: [-----BEGIN RSA PRIVATE KEY-----BASE64 Key Text-----END RSA PRIVATE KEY-----].
 * `Sha1Fingerprint` _string_ - **Required** The SHA1 fingerprint of the client certificate.
    * Examples: [e6d0c68f35032c6c2132217d1f1fb06b12ed32e2].
    * The string must match the regular expression `^[0-9a-f]{40}$`.
 * `UriPrefix` _string_ - The connection prefix.
    * Examples: [jdbc: ].
 * `Username` _string_ - **Required** The name of the SQL user provisioned.
    * Examples: [sb15404128767777].
 * `Password` _string_ - **Required** The database password for the SQL user.
    * Examples: [N-JPz7h2RHPZ81jB5gDHdnluddnIFMWG4nd5rKjR_8A=].
 * `database_name` _string_ - **Required** The name of the database on the instance.
    * Examples: [sb-2-1540412407295372465].
 * `host` _string_ - **Required** The hostname or IP address of the database instance.
    * Examples: [127.0.0.1].
 * `instance_name` _string_ - **Required** The name of the database instance.
    * Examples: [sb-1-1540412407295273023].
    * The string must have at most 84 characters.
    * The string must match the regular expression `^[a-z][a-z0-9-]+$`.
 * `uri` _string_ - **Required** A database connection string.
    * Examples: [mysql://user:pass@127.0.0.1/sb-2-1540412407295372465?ssl_mode=required].
 * `last_master_operation_id` _string_ - (deprecated) The id of the last operation on the database.
    * Examples: [mysql://user:pass@127.0.0.1/sb-2-1540412407295372465?ssl_mode=required].
 * `region` _string_ - **Required** The region the database is in.
    * Examples: [northamerica-northeast1 southamerica-east1 us-east1].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `89e2c84e-4d5c-457c-ad14-329dcf44b806`.
  * Description: MySQL attached to a VPC
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### HA Instance


A regionally available database with automatic failover.
Uses plan: `89e2c84e-4d5c-457c-ad14-329dcf44b806`.

**Provision**

```javascript
{
    "availability_type": "REGIONAL",
    "backups_enabled": "true",
    "binlog": "true",
    "tier": "db-n1-standard-1"
}
```

**Bind**

```javascript
{
    "role": "cloudsql.editor"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-cloudsql-mysql-vpc default my-google-cloudsql-mysql-vpc-example -c `{"availability_type":"REGIONAL","backups_enabled":"true","binlog":"true","tier":"db-n1-standard-1"}`
$ cf bind-service my-app my-google-cloudsql-mysql-vpc-example -c `{"role":"cloudsql.editor"}`
</pre>


### Development Sandbox


An inexpensive MySQL sandbox for developing with no backups.
Uses plan: `89e2c84e-4d5c-457c-ad14-329dcf44b806`.

**Provision**

```javascript
{
    "backups_enabled": "false",
    "binlog": "false",
    "disk_size": "10",
    "tier": "db-n1-standard-1"
}
```

**Bind**

```javascript
{
    "role": "cloudsql.editor"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-cloudsql-mysql-vpc default my-google-cloudsql-mysql-vpc-example -c `{"backups_enabled":"false","binlog":"false","disk_size":"10","tier":"db-n1-standard-1"}`
$ cf bind-service my-app my-google-cloudsql-mysql-vpc-example -c `{"role":"cloudsql.editor"}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-cloudsql-postgres"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg) Google CloudSQL for PostgreSQL
Google CloudSQL for PostgreSQL is a fully-managed PostgreSQL database service.

 * [Documentation](https://cloud.google.com/sql/docs/)
 * [Support](https://cloud.google.com/sql/docs/getting-support/)
 * Catalog Metadata ID: `cbad6d78-a73c-432d-b8ff-b219a17a803a`
 * Tags: gcp, cloudsql, postgres
 * Service Name: `google-cloudsql-postgres`

## Provisioning

**Request Parameters**


 * `instance_name` _string_ - Name of the CloudSQL instance. Default: `sb-${counter.next()}-${time.nano()}`.
    * The string must have at most 86 characters.
    * The string must match the regular expression `^[a-z][a-z0-9-]+$`.
 * `database_name` _string_ - Name of the database inside of the instance. Must be a valid identifier for your chosen database type. Default: `sb-${counter.next()}-${time.nano()}`.
 * `version` _string_ - The database engine type and version. Default: `POSTGRES_11`.
    * The value must be one of: [POSTGRES_10 POSTGRES_11 POSTGRES_12 POSTGRES_9_6].
 * `activation_policy` _string_ - The activation policy specifies when the instance is activated; it is applicable only when the instance state is RUNNABLE. Default: `ALWAYS`.
    * The value must be one of: [ALWAYS NEVER].
 * `authorized_networks` _string_ - A comma separated list without spaces. Default: ``.
 * `region` _string_ - The geographical region. See the instance locations list https://cloud.google.com/sql/docs/mysql/instance-locations for which regions support which databases. Default: `us-central`.
    * Examples: [northamerica-northeast1 southamerica-east1 us-east1].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.
 * `disk_size` _string_ - In GB. Default: `10`.
    * Examples: [10 500 10230].
    * The string must have at most 5 characters.
    * The string must match the regular expression `^[1-9][0-9]+$`.
 * `database_flags` _string_ - The database flags passed to the instance at startup (comma separated list of flags, e.g. general_log=on,skip_show_database=off). Default: ``.
    * Examples: [long_query_time=10 general_log=on,skip_show_database=off].
    * The string must match the regular expression `^(|([a-z_]+=[a-zA-Z0-9\.\+\:-]+)(,[a-z_]+=[a-zA-Z0-9\.\+\:-]+)*)$`.
 * `zone` _string_ - Optional, the specific zone in the region to run the instance. Default: ``.
    * The string must match the regular expression `^(|[A-Za-z][-a-z0-9A-Z]+)$`.
 * `disk_type` _string_ - The type of disk backing the database. Default: `PD_SSD`.
    * The value must be one of: [PD_HDD PD_SSD].
 * `maintenance_window_day` _string_ - The day of week a CloudSQL instance should preferably be restarted for system maintenance purposes. (1-7), starting on Monday. Default: `1`.
    * The value must be one of: [1 2 3 4 5 6 7].
 * `maintenance_window_hour` _string_ - The hour of the day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Hour of day 0-23. Default: `0`.
    * The string must match the regular expression `^([0-9]|1[0-9]|2[0-3])$`.
 * `backups_enabled` _string_ - Should daily backups be enabled for the service? Default: `true`.
    * The value must be one of: [false true].
 * `backup_start_time` _string_ - Start time for the daily backup configuration in UTC timezone in the 24 hour format - HH:MM. Default: `06:00`.
    * The string must match the regular expression `^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$`.
 * `replication_type` _string_ - The type of replication this instance uses. This can be either ASYNCHRONOUS or SYNCHRONOUS. Default: `SYNCHRONOUS`.
    * The value must be one of: [ASYNCHRONOUS SYNCHRONOUS].
 * `auto_resize` _string_ - Configuration to increase storage size automatically. Default: `false`.
    * The value must be one of: [false true].
 * `auto_resize_limit` _string_ - The maximum size to which storage capacity can be automatically increased. Default: `0`.
    * Examples: [10 500 10230].
    * The string must have at most 5 characters.
    * The string must match the regular expression `^[0-9][0-9]*$`.
 * `availability_type` _string_ - Availability type specifies whether the instance serves data from multiple zones. Default: `ZONAL`.
    * The value must be one of: [REGIONAL ZONAL].


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `cloudsql.client`.
    * The value must be one of: [cloudsql.client cloudsql.editor cloudsql.viewer].
 * `jdbc_uri_format` _string_ - If `true`, `uri` field will contain a JDBC formatted URI. Default: `false`.
    * The value must be one of: [false true].
 * `username` _string_ - The SQL username for the account. Default: `sb${str.truncate(14, time.nano())}`.
 * `password` _string_ - The SQL password for the account. Default: `${rand.base64(32)}`.

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].
 * `CaCert` _string_ - **Required** The server Certificate Authority's certificate.
    * Examples: [-----BEGIN CERTIFICATE-----BASE64 Certificate Text-----END CERTIFICATE-----].
 * `ClientCert` _string_ - **Required** The client certificate.
    * Examples: [-----BEGIN CERTIFICATE-----BASE64 Certificate Text-----END CERTIFICATE-----].
 * `ClientKey` _string_ - **Required** The client certificate key.
    * Examples: [-----BEGIN RSA PRIVATE KEY-----BASE64 Key Text-----END RSA PRIVATE KEY-----].
 * `Sha1Fingerprint` _string_ - **Required** The SHA1 fingerprint of the client certificate.
    * Examples: [e6d0c68f35032c6c2132217d1f1fb06b12ed32e2].
    * The string must match the regular expression `^[0-9a-f]{40}$`.
 * `UriPrefix` _string_ - The connection prefix.
    * Examples: [jdbc: ].
 * `Username` _string_ - **Required** The name of the SQL user provisioned.
    * Examples: [sb15404128767777].
 * `Password` _string_ - **Required** The database password for the SQL user.
    * Examples: [N-JPz7h2RHPZ81jB5gDHdnluddnIFMWG4nd5rKjR_8A=].
 * `database_name` _string_ - **Required** The name of the database on the instance.
    * Examples: [sb-2-1540412407295372465].
 * `host` _string_ - **Required** The hostname or IP address of the database instance.
    * Examples: [127.0.0.1].
 * `instance_name` _string_ - **Required** The name of the database instance.
    * Examples: [sb-1-1540412407295273023].
    * The string must have at most 84 characters.
    * The string must match the regular expression `^[a-z][a-z0-9-]+$`.
 * `uri` _string_ - **Required** A database connection string.
    * Examples: [mysql://user:pass@127.0.0.1/sb-2-1540412407295372465?ssl_mode=required].
 * `last_master_operation_id` _string_ - (deprecated) The id of the last operation on the database.
    * Examples: [mysql://user:pass@127.0.0.1/sb-2-1540412407295372465?ssl_mode=required].
 * `region` _string_ - **Required** The region the database is in.
    * Examples: [northamerica-northeast1 southamerica-east1 us-east1].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`postgres-db-f1-micro`**
  * Plan ID: `2513d4d9-684b-4c3c-add4-6404969006de`.
  * Description: PostgreSQL on a db-f1-micro (Shared CPUs, 0.6 GB/RAM, 3062 GB/disk, 250 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-g1-small`**
  * Plan ID: `6c1174d8-243c-44d1-b7a8-e94a779f67f5`.
  * Description: PostgreSQL on a db-g1-small (Shared CPUs, 1.7 GB/RAM, 3062 GB/disk, 1,000 Connections)
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-standard-1`**
  * Plan ID: `c4e68ab5-34ca-4d02-857d-3e6b3ab079a7`.
  * Description: PostgreSQL with 1 CPU, 3.75 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-standard-2`**
  * Plan ID: `3f578ecf-885c-4b60-b38b-60272f34e00f`.
  * Description: PostgreSQL with 2 CPUs, 7.5 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-standard-4`**
  * Plan ID: `b7fcab5d-d66d-4e82-af16-565e84cef7f9`.
  * Description: PostgreSQL with 4 CPUs, 15 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-standard-8`**
  * Plan ID: `4b2fa14a-caf1-42e0-bd8c-3342502008a8`.
  * Description: PostgreSQL with 8 CPUs, 30 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-standard-16`**
  * Plan ID: `ca2e770f-bfa5-4fb7-a249-8b943c3474ca`.
  * Description: PostgreSQL with 16 CPUs, 60 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-standard-32`**
  * Plan ID: `b44f8294-b003-4a50-80c2-706858073f44`.
  * Description: PostgreSQL with 32 CPUs, 120 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-standard-64`**
  * Plan ID: `d97326e0-5af2-4da5-b970-b4772d59cded`.
  * Description: PostgreSQL with 64 CPUs, 240 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-highmem-2`**
  * Plan ID: `c10f8691-02f5-44eb-989f-7217393012ca`.
  * Description: PostgreSQL with 2 CPUs, 13 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-highmem-4`**
  * Plan ID: `610cc78d-d26a-41a9-90b7-547a44517f03`.
  * Description: PostgreSQL with 4 CPUs, 26 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-highmem-8`**
  * Plan ID: `2a351e8d-958d-4c4f-ae46-c984fec18740`.
  * Description: PostgreSQL with 8 CPUs, 52 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-highmem-16`**
  * Plan ID: `51d3ca0c-9d21-447d-a395-3e0dc0659775`.
  * Description: PostgreSQL with 16 CPUs, 104 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-highmem-32`**
  * Plan ID: `2e72b386-f7ce-4f0d-a149-9f9a851337d4`.
  * Description: PostgreSQL with 32 CPUs, 208 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`postgres-db-n1-highmem-64`**
  * Plan ID: `82602649-e4ac-4a2f-b80d-dacd745aed6a`.
  * Description: PostgreSQL with 64 CPUs, 416 GB/RAM, 10230 GB/disk, supporting 4,000 connections.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Dedicated Machine Sandbox


A low end PostgreSQL sandbox that uses a dedicated machine.
Uses plan: `c4e68ab5-34ca-4d02-857d-3e6b3ab079a7`.

**Provision**

```javascript
{
    "backups_enabled": "false",
    "disk_size": "25"
}
```

**Bind**

```javascript
{
    "role": "cloudsql.editor"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-cloudsql-postgres postgres-db-n1-standard-1 my-google-cloudsql-postgres-example -c `{"backups_enabled":"false","disk_size":"25"}`
$ cf bind-service my-app my-google-cloudsql-postgres-example -c `{"role":"cloudsql.editor"}`
</pre>


### Development Sandbox


An inexpensive PostgreSQL sandbox for developing with no backups.
Uses plan: `2513d4d9-684b-4c3c-add4-6404969006de`.

**Provision**

```javascript
{
    "backups_enabled": "false",
    "disk_size": "10"
}
```

**Bind**

```javascript
{
    "role": "cloudsql.editor"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-cloudsql-postgres postgres-db-f1-micro my-google-cloudsql-postgres-example -c `{"backups_enabled":"false","disk_size":"10"}`
$ cf bind-service my-app my-google-cloudsql-postgres-example -c `{"role":"cloudsql.editor"}`
</pre>


### HA Instance


A regionally available database with automatic failover.
Uses plan: `c4e68ab5-34ca-4d02-857d-3e6b3ab079a7`.

**Provision**

```javascript
{
    "availability_type": "REGIONAL",
    "backups_enabled": "false",
    "disk_size": "25"
}
```

**Bind**

```javascript
{
    "role": "cloudsql.editor"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-cloudsql-postgres postgres-db-n1-standard-1 my-google-cloudsql-postgres-example -c `{"availability_type":"REGIONAL","backups_enabled":"false","disk_size":"25"}`
$ cf bind-service my-app my-google-cloudsql-postgres-example -c `{"role":"cloudsql.editor"}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-cloudsql-postgres-vpc"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg) Google CloudSQL for PostgreSQL
Google CloudSQL for PostgreSQL is a fully-managed PostgreSQL database service.

 * [Documentation](https://cloud.google.com/sql/docs/)
 * [Support](https://cloud.google.com/sql/docs/getting-support/)
 * Catalog Metadata ID: `c90ea118-605a-47e8-8f63-57fc09c113f1`
 * Tags: gcp, cloudsql, postgres
 * Service Name: `google-cloudsql-postgres-vpc`

## Provisioning

**Request Parameters**


 * `instance_name` _string_ - Name of the CloudSQL instance. Default: `sb-${counter.next()}-${time.nano()}`.
    * The string must have at most 86 characters.
    * The string must match the regular expression `^[a-z][a-z0-9-]+$`.
 * `database_name` _string_ - Name of the database inside of the instance. Must be a valid identifier for your chosen database type. Default: `sb-${counter.next()}-${time.nano()}`.
 * `version` _string_ - The database engine type and version. Default: `POSTGRES_11`.
    * The value must be one of: [POSTGRES_10 POSTGRES_11 POSTGRES_12 POSTGRES_9_6].
 * `tier` _string_ - The machine type the database will run on. MySQL has predefined tiers, other databases use the a string of the form db-custom-[CPUS]-[MEMORY_MBS], where memory is at least 3840.
    * Examples: [db-n1-standard-1 db-custom-1-3840].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.
 * `private_network` _string_ - The private network to attach to. If specified the instance will only be accessible on the VPC. Default: `default`.
    * Examples: [projects/my-project/global/networks/default].
 * `region` _string_ - The geographical region. See the instance locations list https://cloud.google.com/sql/docs/mysql/instance-locations for which regions support which databases. Default: `us-central`.
    * Examples: [northamerica-northeast1 southamerica-east1 us-east1].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.
 * `disk_size` _string_ - In GB. Default: `10`.
    * Examples: [10 500 10230].
    * The string must have at most 5 characters.
    * The string must match the regular expression `^[1-9][0-9]+$`.
 * `database_flags` _string_ - The database flags passed to the instance at startup (comma separated list of flags, e.g. general_log=on,skip_show_database=off). Default: ``.
    * Examples: [long_query_time=10 general_log=on,skip_show_database=off].
    * The string must match the regular expression `^(|([a-z_]+=[a-zA-Z0-9\.\+\:-]+)(,[a-z_]+=[a-zA-Z0-9\.\+\:-]+)*)$`.
 * `zone` _string_ - Optional, the specific zone in the region to run the instance. Default: ``.
    * The string must match the regular expression `^(|[A-Za-z][-a-z0-9A-Z]+)$`.
 * `disk_type` _string_ - The type of disk backing the database. Default: `PD_SSD`.
    * The value must be one of: [PD_HDD PD_SSD].
 * `maintenance_window_day` _string_ - The day of week a CloudSQL instance should preferably be restarted for system maintenance purposes. (1-7), starting on Monday. Default: `1`.
    * The value must be one of: [1 2 3 4 5 6 7].
 * `maintenance_window_hour` _string_ - The hour of the day when disruptive updates (updates that require an instance restart) to this CloudSQL instance can be made. Hour of day 0-23. Default: `0`.
    * The string must match the regular expression `^([0-9]|1[0-9]|2[0-3])$`.
 * `backups_enabled` _string_ - Should daily backups be enabled for the service? Default: `true`.
    * The value must be one of: [false true].
 * `backup_start_time` _string_ - Start time for the daily backup configuration in UTC timezone in the 24 hour format - HH:MM. Default: `06:00`.
    * The string must match the regular expression `^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$`.
 * `replication_type` _string_ - The type of replication this instance uses. This can be either ASYNCHRONOUS or SYNCHRONOUS. Default: `SYNCHRONOUS`.
    * The value must be one of: [ASYNCHRONOUS SYNCHRONOUS].
 * `auto_resize` _string_ - Configuration to increase storage size automatically. Default: `false`.
    * The value must be one of: [false true].
 * `auto_resize_limit` _string_ - The maximum size to which storage capacity can be automatically increased. Default: `0`.
    * Examples: [10 500 10230].
    * The string must have at most 5 characters.
    * The string must match the regular expression `^[0-9][0-9]*$`.
 * `availability_type` _string_ - Availability type specifies whether the instance serves data from multiple zones. Default: `ZONAL`.
    * The value must be one of: [REGIONAL ZONAL].


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `cloudsql.client`.
    * The value must be one of: [cloudsql.client cloudsql.editor cloudsql.viewer].
 * `jdbc_uri_format` _string_ - If `true`, `uri` field will contain a JDBC formatted URI. Default: `false`.
    * The value must be one of: [false true].
 * `username` _string_ - The SQL username for the account. Default: `sb${str.truncate(14, time.nano())}`.
 * `password` _string_ - The SQL password for the account. Default: `${rand.base64(32)}`.

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].
 * `CaCert` _string_ - **Required** The server Certificate Authority's certificate.
    * Examples: [-----BEGIN CERTIFICATE-----BASE64 Certificate Text-----END CERTIFICATE-----].
 * `ClientCert` _string_ - **Required** The client certificate.
    * Examples: [-----BEGIN CERTIFICATE-----BASE64 Certificate Text-----END CERTIFICATE-----].
 * `ClientKey` _string_ - **Required** The client certificate key.
    * Examples: [-----BEGIN RSA PRIVATE KEY-----BASE64 Key Text-----END RSA PRIVATE KEY-----].
 * `Sha1Fingerprint` _string_ - **Required** The SHA1 fingerprint of the client certificate.
    * Examples: [e6d0c68f35032c6c2132217d1f1fb06b12ed32e2].
    * The string must match the regular expression `^[0-9a-f]{40}$`.
 * `UriPrefix` _string_ - The connection prefix.
    * Examples: [jdbc: ].
 * `Username` _string_ - **Required** The name of the SQL user provisioned.
    * Examples: [sb15404128767777].
 * `Password` _string_ - **Required** The database password for the SQL user.
    * Examples: [N-JPz7h2RHPZ81jB5gDHdnluddnIFMWG4nd5rKjR_8A=].
 * `database_name` _string_ - **Required** The name of the database on the instance.
    * Examples: [sb-2-1540412407295372465].
 * `host` _string_ - **Required** The hostname or IP address of the database instance.
    * Examples: [127.0.0.1].
 * `instance_name` _string_ - **Required** The name of the database instance.
    * Examples: [sb-1-1540412407295273023].
    * The string must have at most 84 characters.
    * The string must match the regular expression `^[a-z][a-z0-9-]+$`.
 * `uri` _string_ - **Required** A database connection string.
    * Examples: [mysql://user:pass@127.0.0.1/sb-2-1540412407295372465?ssl_mode=required].
 * `last_master_operation_id` _string_ - (deprecated) The id of the last operation on the database.
    * Examples: [mysql://user:pass@127.0.0.1/sb-2-1540412407295372465?ssl_mode=required].
 * `region` _string_ - **Required** The region the database is in.
    * Examples: [northamerica-northeast1 southamerica-east1 us-east1].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `60f0b6c0-c48f-4f84-baab-57836611e013`.
  * Description: PostgreSQL attached to a VPC
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Dedicated Machine Sandbox


A low end PostgreSQL sandbox that uses a dedicated machine.
Uses plan: `60f0b6c0-c48f-4f84-baab-57836611e013`.

**Provision**

```javascript
{
    "backups_enabled": "false",
    "disk_size": "25",
    "tier": "db-custom-1-3840"
}
```

**Bind**

```javascript
{
    "role": "cloudsql.editor"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-cloudsql-postgres-vpc default my-google-cloudsql-postgres-vpc-example -c `{"backups_enabled":"false","disk_size":"25","tier":"db-custom-1-3840"}`
$ cf bind-service my-app my-google-cloudsql-postgres-vpc-example -c `{"role":"cloudsql.editor"}`
</pre>


### HA Instance


A regionally available database with automatic failover.
Uses plan: `60f0b6c0-c48f-4f84-baab-57836611e013`.

**Provision**

```javascript
{
    "availability_type": "REGIONAL",
    "backups_enabled": "true",
    "disk_size": "25",
    "tier": "db-custom-1-3840"
}
```

**Bind**

```javascript
{
    "role": "cloudsql.editor"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-cloudsql-postgres-vpc default my-google-cloudsql-postgres-vpc-example -c `{"availability_type":"REGIONAL","backups_enabled":"true","disk_size":"25","tier":"db-custom-1-3840"}`
$ cf bind-service my-app my-google-cloudsql-postgres-vpc-example -c `{"role":"cloudsql.editor"}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-dataflow"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/dataflow.svg) Google Cloud Dataflow
A managed service for executing a wide variety of data processing patterns built on Apache Beam.

 * [Documentation](https://cloud.google.com/dataflow/docs/)
 * [Support](https://cloud.google.com/dataflow/docs/support)
 * Catalog Metadata ID: `3e897eb3-9062-4966-bd4f-85bda0f73b3d`
 * Tags: gcp, dataflow, preview
 * Service Name: `google-dataflow`

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `dataflow.developer`.
    * The value must be one of: [dataflow.developer dataflow.viewer].

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `8e956dd6-8c0f-470c-9a11-065537d81872`.
  * Description: Dataflow default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Developer


Creates a Dataflow user and grants it permission to create, drain and cancel jobs.
Uses plan: `8e956dd6-8c0f-470c-9a11-065537d81872`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-dataflow default my-google-dataflow-example -c `{}`
$ cf bind-service my-app my-google-dataflow-example -c `{}`
</pre>


### Viewer


Creates a Dataflow user and grants it permission to create, drain and cancel jobs.
Uses plan: `8e956dd6-8c0f-470c-9a11-065537d81872`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{
    "role": "dataflow.viewer"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-dataflow default my-google-dataflow-example -c `{}`
$ cf bind-service my-app my-google-dataflow-example -c `{"role":"dataflow.viewer"}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-datastore"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/datastore.svg) Google Cloud Datastore
Google Cloud Datastore is a NoSQL document database service.

 * [Documentation](https://cloud.google.com/datastore/docs/)
 * [Support](https://cloud.google.com/datastore/docs/getting-support)
 * Catalog Metadata ID: `76d4abb2-fee7-4c8f-aee1-bcea2837f02b`
 * Tags: gcp, datastore
 * Service Name: `google-datastore`

## Provisioning

**Request Parameters**


 * `namespace` _string_ - A context for the identifiers in your entity’s dataset. This ensures that different systems can all interpret an entity's data the same way, based on the rules for the entity’s particular namespace. Blank means the default namespace will be used. Default: ``.
    * The string must have at most 100 characters.
    * The string must match the regular expression `^[A-Za-z0-9_-]*$`.


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].
 * `namespace` _string_ - A context for the identifiers in your entity’s dataset.
    * The string must have at most 100 characters.
    * The string must match the regular expression `^[A-Za-z0-9_-]*$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `05f1fb6b-b5f0-48a2-9c2b-a5f236507a97`.
  * Description: Datastore default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


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

**Cloud Foundry Example**

<pre>
$ cf create-service google-datastore default my-google-datastore-example -c `{}`
$ cf bind-service my-app my-google-datastore-example -c `{}`
</pre>


### Custom Namespace


Creates a datastore and returns the provided namespace along with bind calls.
Uses plan: `05f1fb6b-b5f0-48a2-9c2b-a5f236507a97`.

**Provision**

```javascript
{
    "namespace": "my-namespace"
}
```

**Bind**

```javascript
{}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-datastore default my-google-datastore-example -c `{"namespace":"my-namespace"}`
$ cf bind-service my-app my-google-datastore-example -c `{}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-dialogflow"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/dialogflow-enterprise.svg) Google Cloud Dialogflow
Dialogflow is an end-to-end, build-once deploy-everywhere development suite for creating conversational interfaces for websites, mobile applications, popular messaging platforms, and IoT devices.

 * [Documentation](https://cloud.google.com/dialogflow-enterprise/docs/)
 * [Support](https://cloud.google.com/dialogflow-enterprise/docs/support)
 * Catalog Metadata ID: `e84b69db-3de9-4688-8f5c-26b9d5b1f129`
 * Tags: gcp, dialogflow, preview
 * Service Name: `google-dialogflow`

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `3ac4e1bd-b22d-4a99-864b-d3a3ac582348`.
  * Description: Dialogflow default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Reader


Creates a Dialogflow user and grants it permission to detect intent and read/write session properties (contexts, session entity types, etc.).
Uses plan: `3ac4e1bd-b22d-4a99-864b-d3a3ac582348`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-dialogflow default my-google-dialogflow-example -c `{}`
$ cf bind-service my-app my-google-dialogflow-example -c `{}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-filestore"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/storage.svg) Google Cloud Filestore
Fully managed NFS file storage with predictable performance.

 * [Documentation](https://cloud.google.com/filestore/docs/)
 * [Support](https://cloud.google.com/filestore/docs/getting-support)
 * Catalog Metadata ID: `494eb82e-c4ca-4bed-871d-9c3f02f66e01`
 * Tags: gcp, filestore, nfs
 * Service Name: `google-filestore`

## Provisioning

**Request Parameters**


 * `instance_id` _string_ - The name of the instance. The name must be unique per zone. Default: `gsb-${counter.next()}-${time.nano()}`.
    * The string must have at most 63 characters.
    * The string must have at least 1 characters.
    * The string must match the regular expression `^[a-z]([-0-9a-z]*[a-z0-9]$)*`.
 * `zone` _string_ - The zone to create the instance in. Supported zones can be found here: https://cloud.google.com/filestore/docs/regions. Default: `us-west1-a`.
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.
 * `tier` _string_ - The performance tier. Default: `STANDARD`.
    * The value must be one of: [PREMIUM STANDARD].
 * `authorized_network` _string_ - The name of the network to attach the instance to. Default: `default`.
 * `address_mode` _string_ - The address mode of the service. Default: `MODE_IPV4`.
    * The value must be one of: [MODE_IPV4].
 * `capacity_gb` _integer_ - The capacity of the Filestore. Standard minimum is 1TiB and Premium is minimum 2.5TiB. Default: `1024`.


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `authorized_network` _string_ - Name of the VPC network the instance is attached to.
 * `reserved_ip_range` _string_ - Range of IP addresses reserved for the instance.
 * `ip_address` _string_ - IP address of the service.
 * `file_share_name` _string_ - Name of the share.
 * `capacity_gb` _integer_ - Capacity of the share in GiB.
 * `uri` _string_ - URI of the instance.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `e4c83975-e60f-43cf-afde-ebec573c6c2e`.
  * Description: Filestore default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Standard


Creates a standard Filestore.
Uses plan: `e4c83975-e60f-43cf-afde-ebec573c6c2e`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-filestore default my-google-filestore-example -c `{}`
$ cf bind-service my-app my-google-filestore-example -c `{}`
</pre>


### Premium


Creates a premium Filestore.
Uses plan: `e4c83975-e60f-43cf-afde-ebec573c6c2e`.

**Provision**

```javascript
{
    "capacity_gb": 2560,
    "tier": "PREMIUM"
}
```

**Bind**

```javascript
{}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-filestore default my-google-filestore-example -c `{"capacity_gb":2560,"tier":"PREMIUM"}`
$ cf bind-service my-app my-google-filestore-example -c `{}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-firestore"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/firestore.svg) Google Cloud Firestore
Cloud Firestore is a fast, fully managed, serverless, cloud-native NoSQL document database that simplifies storing, syncing, and querying data for your mobile, web, and IoT apps at global scale.

 * [Documentation](https://cloud.google.com/firestore/docs/)
 * [Support](https://cloud.google.com/firestore/docs/getting-support)
 * Catalog Metadata ID: `a2b7b873-1e34-4530-8a42-902ff7d66b43`
 * Tags: gcp, firestore, preview, beta
 * Service Name: `google-firestore`

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `datastore.user`.
    * The value must be one of: [datastore.user datastore.viewer].

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `64403af0-4413-4ef3-a813-37f0306ef498`.
  * Description: Firestore default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Reader Writer


Creates a general Firestore user and grants it permission to read and write entities.
Uses plan: `64403af0-4413-4ef3-a813-37f0306ef498`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-firestore default my-google-firestore-example -c `{}`
$ cf bind-service my-app my-google-firestore-example -c `{}`
</pre>


### Read Only


Creates a Firestore user that can only view entities.
Uses plan: `64403af0-4413-4ef3-a813-37f0306ef498`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{
    "role": "datastore.viewer"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-firestore default my-google-firestore-example -c `{}`
$ cf bind-service my-app my-google-firestore-example -c `{"role":"datastore.viewer"}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-memorystore-redis"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/cache.svg) Google Cloud Memorystore for Redis API
Creates and manages Redis instances on the Google Cloud Platform.

 * [Documentation](https://cloud.google.com/memorystore/docs/redis)
 * [Support](https://cloud.google.com/memorystore/docs/redis/support)
 * Catalog Metadata ID: `3ea92b54-838c-4fe1-b75d-9bda513380aa`
 * Tags: gcp, memorystore, redis
 * Service Name: `google-memorystore-redis`

## Provisioning

**Request Parameters**


 * `instance_id` _string_ - The name of the instance. The name must be unique per project. Default: `gsb-${counter.next()}-${time.nano()}`.
    * The string must have at most 40 characters.
    * The string must have at least 1 characters.
    * The string must match the regular expression `^[a-z]([-0-9a-z]*[a-z0-9]$)*`.
 * `authorized_network` _string_ - The name of the VPC network to attach the instance to. Default: `default`.
    * Examples: [default projects/MYPROJECT/global/networks/MYNETWORK].
 * `region` _string_ - The region to create the instance in. Supported regions can be found here: https://cloud.google.com/memorystore/docs/redis/regions. Default: `us-east1`.
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.
 * `memory_size_gb` _integer_ - Redis memory size in GiB. Default: `4`.
 * `tier` _string_ - The performance tier. Default: `BASIC`.
    * The value must be one of: [BASIC STANDARD_HA].


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `authorized_network` _string_ - Name of the VPC network the instance is attached to.
 * `reserved_ip_range` _string_ - Range of IP addresses reserved for the instance.
 * `redis_version` _string_ - The version of Redis software.
 * `memory_size_gb` _integer_ - Redis memory size in GiB.
 * `host` _string_ - Hostname or IP address of the exposed Redis endpoint used by clients to connect to the service.
 * `port` _integer_ - The port number of the exposed Redis endpoint.
 * `uri` _string_ - URI of the instance.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `df10762e-6ef1-44e3-84c2-07e9358ceb1f`.
  * Description: Lets you chose your own values for all properties.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`basic`**
  * Plan ID: `dd1923b6-ac26-4697-83d6-b3a0c05c2c94`.
  * Description: Provides a standalone Redis instance. Use this tier for applications that require a simple Redis cache.
  * This plan overrides the following user variables on provision.
    * `service_tier` = `BASIC`
  * This plan doesn't override user variables on bind.
* **`standard_ha`**
  * Plan ID: `41771881-b456-4940-9081-34b6424744c6`.
  * Description: Provides a highly available Redis instance.
  * This plan overrides the following user variables on provision.
    * `service_tier` = `STANDARD_HA`
  * This plan doesn't override user variables on bind.


## Examples




### Standard Redis Configuration


Create a Redis instance with standard service tier.
Uses plan: `dd1923b6-ac26-4697-83d6-b3a0c05c2c94`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-memorystore-redis basic my-google-memorystore-redis-example -c `{}`
$ cf bind-service my-app my-google-memorystore-redis-example -c `{}`
</pre>


### HA Redis Configuration


Create a Redis instance with high availability.
Uses plan: `41771881-b456-4940-9081-34b6424744c6`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-memorystore-redis standard_ha my-google-memorystore-redis-example -c `{}`
$ cf bind-service my-app my-google-memorystore-redis-example -c `{}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-ml-apis"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/machine-learning.svg) Google Machine Learning APIs
Machine Learning APIs including Vision, Translate, Speech, and Natural Language.

 * [Documentation](https://cloud.google.com/ml/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `5ad2dce0-51f7-4ede-8b46-293d6df1e8d4`
 * Tags: gcp, ml
 * Service Name: `google-ml-apis`

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `ml.modelUser`.
    * The value must be one of: [ml.developer ml.jobOwner ml.modelOwner ml.modelUser ml.operationOwner ml.viewer].

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `be7954e1-ecfb-4936-a0b6-db35e6424c7a`.
  * Description: Machine Learning API default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


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

**Cloud Foundry Example**

<pre>
$ cf create-service google-ml-apis default my-google-ml-apis-example -c `{}`
$ cf bind-service my-app my-google-ml-apis-example -c `{"role":"ml.developer"}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-pubsub"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/pubsub.svg) Google PubSub
A global service for real-time and reliable messaging and streaming data.

 * [Documentation](https://cloud.google.com/pubsub/docs/)
 * [Support](https://cloud.google.com/pubsub/docs/support)
 * Catalog Metadata ID: `628629e3-79f5-4255-b981-d14c6c7856be`
 * Tags: gcp, pubsub
 * Service Name: `google-pubsub`

## Provisioning

**Request Parameters**


 * `topic_name` _string_ - Name of the topic. Must not start with "goog". Default: `pcf_sb_${counter.next()}_${time.nano()}`.
    * The string must have at most 255 characters.
    * The string must have at least 3 characters.
    * The string must match the regular expression `^[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+$`.
 * `subscription_name` _string_ - Name of the subscription. Blank means no subscription will be created. Must not start with "goog". Default: ``.
    * The string must have at most 255 characters.
    * The string must have at least 0 characters.
    * The string must match the regular expression `^(|[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+)`.
 * `is_push` _string_ - Are events handled by POSTing to a URL? Default: `false`.
    * The value must be one of: [false true].
 * `endpoint` _string_ - If `is_push` == 'true', then this is the URL that will be pushed to. Default: ``.
 * `ack_deadline` _string_ - Value is in seconds. Max: 600 This is the maximum time after a subscriber receives a message before the subscriber should acknowledge the message. After message delivery but before the ack deadline expires and before the message is acknowledged, it is an outstanding message and will not be delivered again during that time (on a best-effort basis).  Default: `10`.


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `pubsub.editor`.
    * The value must be one of: [pubsub.editor pubsub.publisher pubsub.subscriber pubsub.viewer].

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].
 * `subscription_name` _string_ - Name of the subscription.
    * The string must have at most 255 characters.
    * The string must have at least 0 characters.
    * The string must match the regular expression `^(|[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+)`.
 * `topic_name` _string_ - **Required** Name of the topic.
    * The string must have at most 255 characters.
    * The string must have at least 3 characters.
    * The string must match the regular expression `^[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `622f4da3-8731-492a-af29-66a9146f8333`.
  * Description: PubSub Default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


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

**Cloud Foundry Example**

<pre>
$ cf create-service google-pubsub default my-google-pubsub-example -c `{"subscription_name":"example_topic_subscription","topic_name":"example_topic"}`
$ cf bind-service my-app my-google-pubsub-example -c `{"role":"pubsub.publisher"}`
</pre>


### No Subscription


Create a topic without a subscription.
Uses plan: `622f4da3-8731-492a-af29-66a9146f8333`.

**Provision**

```javascript
{
    "topic_name": "example_topic"
}
```

**Bind**

```javascript
{
    "role": "pubsub.publisher"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-pubsub default my-google-pubsub-example -c `{"topic_name":"example_topic"}`
$ cf bind-service my-app my-google-pubsub-example -c `{"role":"pubsub.publisher"}`
</pre>


### Custom Timeout


Create a subscription with a custom deadline for long processess.
Uses plan: `622f4da3-8731-492a-af29-66a9146f8333`.

**Provision**

```javascript
{
    "ack_deadline": "200",
    "subscription_name": "long_deadline_subscription",
    "topic_name": "long_deadline_topic"
}
```

**Bind**

```javascript
{
    "role": "pubsub.publisher"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-pubsub default my-google-pubsub-example -c `{"ack_deadline":"200","subscription_name":"long_deadline_subscription","topic_name":"long_deadline_topic"}`
$ cf bind-service my-app my-google-pubsub-example -c `{"role":"pubsub.publisher"}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-spanner"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/spanner.svg) Google Spanner
The first horizontally scalable, globally consistent, relational database service.

 * [Documentation](https://cloud.google.com/spanner/)
 * [Support](https://cloud.google.com/spanner/docs/support)
 * Catalog Metadata ID: `51b3e27e-d323-49ce-8c5f-1211e6409e82`
 * Tags: gcp, spanner
 * Service Name: `google-spanner`

## Provisioning

**Request Parameters**


 * `name` _string_ - A unique identifier for the instance, which cannot be changed after the instance is created. Default: `pcf-sb-${counter.next()}-${time.nano()}`.
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z][-a-z0-9]*[a-z0-9]$`.
 * `display_name` _string_ - The name of this instance configuration as it appears in UIs. Default: `${name}`.
    * The string must have at most 30 characters.
    * The string must have at least 4 characters.
 * `location` _string_ - A configuration for a Cloud Spanner instance. Configurations define the geographic placement of nodes and their replication and are slightly different from zones. There are single region configurations, multi-region configurations, and multi-continent configurations. See the instance docs https://cloud.google.com/spanner/docs/instances for a list of configurations. Default: `regional-us-central1`.
    * Examples: [regional-asia-east1 nam3 nam-eur-asia1].
    * The string must match the regular expression `^[a-z][-a-z0-9]*[a-z0-9]$`.


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `spanner.databaseUser`.
    * The value must be one of: [spanner.databaseAdmin spanner.databaseReader spanner.databaseUser spanner.viewer].

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].
 * `instance_id` _string_ - **Required** Name of the Spanner instance the account can connect to.
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z][-a-z0-9]*[a-z0-9]$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`sandbox`**
  * Plan ID: `44828436-cfbd-47ae-b4bc-48854564347b`.
  * Description: Useful for testing, not eligible for SLA.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`minimal-production`**
  * Plan ID: `0752b1ad-a784-4dcc-96eb-64149089a1c9`.
  * Description: A minimal production level Spanner setup eligible for 99.99% SLA. Each node can provide up to 10,000 QPS of reads or 2,000 QPS of writes (writing single rows at 1KB data per row), and 2 TiB storage.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


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

**Cloud Foundry Example**

<pre>
$ cf create-service google-spanner sandbox my-google-spanner-example -c `{"name":"auth-database"}`
$ cf bind-service my-app my-google-spanner-example -c `{"role":"spanner.databaseAdmin"}`
</pre>


### 99.999% availability


Create a spanner instance spanning North America.
Uses plan: `44828436-cfbd-47ae-b4bc-48854564347b`.

**Provision**

```javascript
{
    "location": "nam3",
    "name": "auth-database"
}
```

**Bind**

```javascript
{
    "role": "spanner.databaseAdmin"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-spanner sandbox my-google-spanner-example -c `{"location":"nam3","name":"auth-database"}`
$ cf bind-service my-app my-google-spanner-example -c `{"role":"spanner.databaseAdmin"}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-stackdriver-debugger"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/debugger.svg) Stackdriver Debugger
Inspect the state of an app, at any code location, without stopping or slowing it down.

 * [Documentation](https://cloud.google.com/debugger/docs/)
 * [Support](https://cloud.google.com/stackdriver/docs/getting-support)
 * Catalog Metadata ID: `83837945-1547-41e0-b661-ea31d76eed11`
 * Tags: gcp, stackdriver, debugger
 * Service Name: `google-stackdriver-debugger`

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `10866183-a775-49e8-96e3-4e7a901e4a79`.
  * Description: Stackdriver Debugger default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


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

**Cloud Foundry Example**

<pre>
$ cf create-service google-stackdriver-debugger default my-google-stackdriver-debugger-example -c `{}`
$ cf bind-service my-app my-google-stackdriver-debugger-example -c `{}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-stackdriver-monitoring"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/stackdriver.svg) Stackdriver Monitoring
Stackdriver Monitoring provides visibility into the performance, uptime, and overall health of cloud-powered applications.

 * [Documentation](https://cloud.google.com/monitoring/docs/)
 * [Support](https://cloud.google.com/stackdriver/docs/getting-support)
 * Catalog Metadata ID: `2bc0d9ed-3f68-4056-b842-4a85cfbc727f`
 * Tags: gcp, stackdriver, monitoring, preview
 * Service Name: `google-stackdriver-monitoring`

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `2e4b85c1-0ce6-46e4-91f5-eebeb373e3f5`.
  * Description: Stackdriver Monitoring default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Basic Configuration


Creates an account with the permission `monitoring.metricWriter` for writing metrics.
Uses plan: `2e4b85c1-0ce6-46e4-91f5-eebeb373e3f5`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-stackdriver-monitoring default my-google-stackdriver-monitoring-example -c `{}`
$ cf bind-service my-app my-google-stackdriver-monitoring-example -c `{}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-stackdriver-profiler"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/stackdriver.svg) Stackdriver Profiler
Continuous CPU and heap profiling to improve performance and reduce costs.

 * [Documentation](https://cloud.google.com/profiler/docs/)
 * [Support](https://cloud.google.com/stackdriver/docs/getting-support)
 * Catalog Metadata ID: `00b9ca4a-7cd6-406a-a5b7-2f43f41ade75`
 * Tags: gcp, stackdriver, profiler
 * Service Name: `google-stackdriver-profiler`

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `594627f6-35f5-462f-9074-10fb033fb18a`.
  * Description: Stackdriver Profiler default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


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

**Cloud Foundry Example**

<pre>
$ cf create-service google-stackdriver-profiler default my-google-stackdriver-profiler-example -c `{}`
$ cf bind-service my-app my-google-stackdriver-profiler-example -c `{}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-stackdriver-trace"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/trace.svg) Stackdriver Trace
A real-time distributed tracing system.

 * [Documentation](https://cloud.google.com/trace/docs/)
 * [Support](https://cloud.google.com/stackdriver/docs/getting-support)
 * Catalog Metadata ID: `c5ddfe15-24d9-47f8-8ffe-f6b7daa9cf4a`
 * Tags: gcp, stackdriver, trace
 * Service Name: `google-stackdriver-trace`

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `ab6c2287-b4bc-4ff4-a36a-0575e7910164`.
  * Description: Stackdriver Trace default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


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

**Cloud Foundry Example**

<pre>
$ cf create-service google-stackdriver-trace default my-google-stackdriver-trace-example -c `{}`
$ cf bind-service my-app my-google-stackdriver-trace-example -c `{}`
</pre>



--------------------------------------------------------------------------------

# <a name="google-storage"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/storage.svg) Google Cloud Storage
Unified object storage for developers and enterprises. Cloud Storage allows world-wide storage and retrieval of any amount of data at any time.

 * [Documentation](https://cloud.google.com/storage/docs/overview)
 * [Support](https://cloud.google.com/storage/docs/getting-support)
 * Catalog Metadata ID: `b9e4332e-b42b-4680-bda5-ea1506797474`
 * Tags: gcp, storage
 * Service Name: `google-storage`

## Provisioning

**Request Parameters**


 * `name` _string_ - The name of the bucket. There is a single global namespace shared by all buckets so it MUST be unique. Default: `pcf_sb_${counter.next()}_${time.nano()}`.
    * The string must have at most 222 characters.
    * The string must have at least 3 characters.
    * The string must match the regular expression `^[a-z0-9_.-]+$`.
 * `location` _string_ - The location of the bucket. Object data for objects in the bucket resides in physical storage within this region. See: https://cloud.google.com/storage/docs/bucket-locations Default: `US`.
    * Examples: [US EU southamerica-east1].
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.
 * `force_delete` _string_ - Attempt to erase bucket contents before deleting bucket on deprovision. Default: `false`.
    * The value must be one of: [true false].


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `storage.objectAdmin`.
    * The value must be one of: [storage.objectAdmin storage.objectCreator storage.objectViewer].

**Response Parameters**

 * `Email` _string_ - **Required** Email address of the service account.
    * Examples: [pcf-binding-ex312029@my-project.iam.gserviceaccount.com].
    * The string must match the regular expression `^pcf-binding-[a-z0-9-]+@.+\.gserviceaccount\.com$`.
 * `Name` _string_ - **Required** The name of the service account.
    * Examples: [pcf-binding-ex312029].
 * `PrivateKeyData` _string_ - **Required** Service account private key data. Base64 encoded JSON.
    * The string must have at least 512 characters.
    * The string must match the regular expression `^[A-Za-z0-9+/]*=*$`.
 * `ProjectId` _string_ - **Required** ID of the project that owns the service account.
    * Examples: [my-project].
    * The string must have at most 30 characters.
    * The string must have at least 6 characters.
    * The string must match the regular expression `^[a-z0-9-]+$`.
 * `UniqueId` _string_ - **Required** Unique and stable ID of the service account.
    * Examples: [112447814736626230844].
 * `bucket_name` _string_ - **Required** Name of the bucket this binding is for.
    * The string must have at most 222 characters.
    * The string must have at least 3 characters.
    * The string must match the regular expression `^[A-Za-z0-9_\.]+$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`standard`**
  * Plan ID: `e1d11f65-da66-46ad-977c-6d56513baf43`.
  * Description: Standard storage class. Auto-selects either regional or multi-regional based on the location.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`nearline`**
  * Plan ID: `a42c1182-d1a0-4d40-82c1-28220518b360`.
  * Description: Nearline storage class.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`reduced-availability`**
  * Plan ID: `1a1f4fe6-1904-44d0-838c-4c87a9490a6b`.
  * Description: Durable Reduced Availability storage class.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`coldline`**
  * Plan ID: `c8538397-8f15-45e3-a229-8bb349c3a98f`.
  * Description: Google Cloud Storage Coldline is a very-low-cost, highly durable storage service for data archiving, online backup, and disaster recovery.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`regional`**
  * Plan ID: `5e6161d2-0202-48be-80c4-1006cce19b9d`.
  * Description: Data is stored in a narrow geographic region, redundant across availability zones with a 99.99% typical monthly availability.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`multiregional`**
  * Plan ID: `a5e8dfb5-e5ec-472a-8d36-33afcaff2fdb`.
  * Description: Data is stored geo-redundantly with >99.99% typical monthly availability.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Basic Configuration


Create a nearline bucket with a service account that can create/read/list/delete the objects in it.
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

**Cloud Foundry Example**

<pre>
$ cf create-service google-storage nearline my-google-storage-example -c `{"location":"us"}`
$ cf bind-service my-app my-google-storage-example -c `{"role":"storage.objectAdmin"}`
</pre>


### Cold Storage


Create a coldline bucket with a service account that can create/read/list/delete the objects in it.
Uses plan: `c8538397-8f15-45e3-a229-8bb349c3a98f`.

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

**Cloud Foundry Example**

<pre>
$ cf create-service google-storage coldline my-google-storage-example -c `{"location":"us"}`
$ cf bind-service my-app my-google-storage-example -c `{"role":"storage.objectAdmin"}`
</pre>


### Regional Storage


Create a regional bucket with a service account that can create/read/list/delete the objects in it.
Uses plan: `5e6161d2-0202-48be-80c4-1006cce19b9d`.

**Provision**

```javascript
{
    "location": "us-west1"
}
```

**Bind**

```javascript
{
    "role": "storage.objectAdmin"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-storage regional my-google-storage-example -c `{"location":"us-west1"}`
$ cf bind-service my-app my-google-storage-example -c `{"role":"storage.objectAdmin"}`
</pre>


### Multi-Regional Storage


Create a multi-regional bucket with a service account that can create/read/list/delete the objects in it.
Uses plan: `a5e8dfb5-e5ec-472a-8d36-33afcaff2fdb`.

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

**Cloud Foundry Example**

<pre>
$ cf create-service google-storage multiregional my-google-storage-example -c `{"location":"us"}`
$ cf bind-service my-app my-google-storage-example -c `{"role":"storage.objectAdmin"}`
</pre>


### Delete even if not empty


Sets the label sb-force-delete=true on the bucket. The broker will try to erase all contents before deleting the bucket.
Uses plan: `5e6161d2-0202-48be-80c4-1006cce19b9d`.

**Provision**

```javascript
{
    "force_delete": "true",
    "location": "us-west1"
}
```

**Bind**

```javascript
{
    "role": "storage.objectAdmin"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-storage regional my-google-storage-example -c `{"force_delete":"true","location":"us-west1"}`
$ cf bind-service my-app my-google-storage-example -c `{"role":"storage.objectAdmin"}`
</pre>



--------------------------------------------------------------------------------
