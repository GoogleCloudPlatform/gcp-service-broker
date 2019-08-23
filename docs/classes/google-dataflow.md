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