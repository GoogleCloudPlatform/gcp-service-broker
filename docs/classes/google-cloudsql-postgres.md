# <a name="google-cloudsql-postgres"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/sql.svg) Google CloudSQL for PostgreSQL
Google CloudSQL for PostgreSQL is a fully-managed PostgreSQL database service.

 * [Documentation](https://cloud.google.com/sql/docs/)
 * [Support](https://cloud.google.com/support/)
 * Catalog Metadata ID: `cbad6d78-a73c-432d-b8ff-b219a17a803a`
 * Tags: gcp, cloudsql, postgres
 * Service Name: `google-cloudsql-postgres`

## Provisioning

**Request Parameters**


 * `instance_name` _string_ - Name of the CloudSQL instance. Default: `pcf-sb-${counter.next()}-${time.nano()}`.
    * The string must have at most 86 characters.
    * The string must match the regular expression `^[a-z][a-z0-9-]+$`.
 * `database_name` _string_ - Name of the database inside of the instance. Must be a valid identifier for your chosen database type. Default: `pcf-sb-${counter.next()}-${time.nano()}`.
 * `version` _string_ - The database engine type and version. Default: `POSTGRES_9_6`.
    * The value must be one of: [POSTGRES_11 POSTGRES_9_6].
 * `activation_policy` _string_ - The activation policy specifies when the instance is activated; it is applicable only when the instance state is RUNNABLE. Default: `ALWAYS`.
    * The value must be one of: [ALWAYS NEVER].
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
 * `authorized_networks` _string_ - A comma separated list without spaces. Default: ``.
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
    * Examples: [pcf-sb-2-1540412407295372465].
 * `host` _string_ - **Required** The hostname or IP address of the database instance.
    * Examples: [127.0.0.1].
 * `instance_name` _string_ - **Required** The name of the database instance.
    * Examples: [pcf-sb-1-1540412407295273023].
    * The string must have at most 84 characters.
    * The string must match the regular expression `^[a-z][a-z0-9-]+$`.
 * `uri` _string_ - **Required** A database connection string.
    * Examples: [mysql://user:pass@127.0.0.1/pcf-sb-2-1540412407295372465?ssl_mode=required].
 * `last_master_operation_id` _string_ - (deprecated) The id of the last operation on the database.
    * Examples: [mysql://user:pass@127.0.0.1/pcf-sb-2-1540412407295372465?ssl_mode=required].
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