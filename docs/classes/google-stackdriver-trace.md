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