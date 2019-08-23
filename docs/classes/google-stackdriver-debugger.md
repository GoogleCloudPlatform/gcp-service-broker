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