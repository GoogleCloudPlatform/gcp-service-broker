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