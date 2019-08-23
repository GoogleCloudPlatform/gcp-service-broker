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