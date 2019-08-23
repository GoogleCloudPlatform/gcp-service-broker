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