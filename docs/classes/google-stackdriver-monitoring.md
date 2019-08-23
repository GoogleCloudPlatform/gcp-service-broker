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