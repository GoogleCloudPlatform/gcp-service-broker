# <a name="google-dialogflow"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/dialogflow-enterprise.svg) Google Cloud Dialogflow
Dialogflow is an end-to-end, build-once deploy-everywhere development suite for creating conversational interfaces for websites, mobile applications, popular messaging platforms, and IoT devices.

 * [Documentation](https://cloud.google.com/dialogflow-enterprise/docs/)
 * [Support](https://cloud.google.com/dialogflow-enterprise/docs/support)
 * Catalog Metadata ID: `e84b69db-3de9-4688-8f5c-26b9d5b1f129`
 * Tags: gcp, dialogflow, preview
 * Service Name: `google-dialogflow`

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
  * Plan ID: `3ac4e1bd-b22d-4a99-864b-d3a3ac582348`.
  * Description: Dialogflow default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Reader


Creates a Dialogflow user and grants it permission to detect intent and read/write session properties (contexts, session entity types, etc.).
Uses plan: `3ac4e1bd-b22d-4a99-864b-d3a3ac582348`.

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
$ cf create-service google-dialogflow default my-google-dialogflow-example -c `{}`
$ cf bind-service my-app my-google-dialogflow-example -c `{}`
</pre>