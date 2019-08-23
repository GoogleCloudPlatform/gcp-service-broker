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