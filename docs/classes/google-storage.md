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