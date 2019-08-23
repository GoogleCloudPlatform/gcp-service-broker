# <a name="google-firestore"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/firestore.svg) Google Cloud Firestore
Cloud Firestore is a fast, fully managed, serverless, cloud-native NoSQL document database that simplifies storing, syncing, and querying data for your mobile, web, and IoT apps at global scale.

 * [Documentation](https://cloud.google.com/firestore/docs/)
 * [Support](https://cloud.google.com/firestore/docs/getting-support)
 * Catalog Metadata ID: `a2b7b873-1e34-4530-8a42-902ff7d66b43`
 * Tags: gcp, firestore, preview, beta
 * Service Name: `google-firestore`

## Provisioning

**Request Parameters**

_No parameters supported._


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `datastore.user`.
    * The value must be one of: [datastore.user datastore.viewer].

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
  * Plan ID: `64403af0-4413-4ef3-a813-37f0306ef498`.
  * Description: Firestore default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Reader Writer


Creates a general Firestore user and grants it permission to read and write entities.
Uses plan: `64403af0-4413-4ef3-a813-37f0306ef498`.

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
$ cf create-service google-firestore default my-google-firestore-example -c `{}`
$ cf bind-service my-app my-google-firestore-example -c `{}`
</pre>


### Read Only


Creates a Firestore user that can only view entities.
Uses plan: `64403af0-4413-4ef3-a813-37f0306ef498`.

**Provision**

```javascript
{}
```

**Bind**

```javascript
{
    "role": "datastore.viewer"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-firestore default my-google-firestore-example -c `{}`
$ cf bind-service my-app my-google-firestore-example -c `{"role":"datastore.viewer"}`
</pre>