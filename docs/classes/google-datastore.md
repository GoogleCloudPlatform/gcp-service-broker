# <a name="google-datastore"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/datastore.svg) Google Cloud Datastore
Google Cloud Datastore is a NoSQL document database service.

 * [Documentation](https://cloud.google.com/datastore/docs/)
 * [Support](https://cloud.google.com/datastore/docs/getting-support)
 * Catalog Metadata ID: `76d4abb2-fee7-4c8f-aee1-bcea2837f02b`
 * Tags: gcp, datastore
 * Service Name: `google-datastore`

## Provisioning

**Request Parameters**


 * `namespace` _string_ - A context for the identifiers in your entity’s dataset. This ensures that different systems can all interpret an entity's data the same way, based on the rules for the entity’s particular namespace. Blank means the default namespace will be used. Default: ``.
    * The string must have at most 100 characters.
    * The string must match the regular expression `^[A-Za-z0-9_-]*$`.


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
 * `namespace` _string_ - A context for the identifiers in your entity’s dataset.
    * The string must have at most 100 characters.
    * The string must match the regular expression `^[A-Za-z0-9_-]*$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `05f1fb6b-b5f0-48a2-9c2b-a5f236507a97`.
  * Description: Datastore default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Basic Configuration


Creates a datastore and a user with the permission `datastore.user`.
Uses plan: `05f1fb6b-b5f0-48a2-9c2b-a5f236507a97`.

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
$ cf create-service google-datastore default my-google-datastore-example -c `{}`
$ cf bind-service my-app my-google-datastore-example -c `{}`
</pre>


### Custom Namespace


Creates a datastore and returns the provided namespace along with bind calls.
Uses plan: `05f1fb6b-b5f0-48a2-9c2b-a5f236507a97`.

**Provision**

```javascript
{
    "namespace": "my-namespace"
}
```

**Bind**

```javascript
{}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-datastore default my-google-datastore-example -c `{"namespace":"my-namespace"}`
$ cf bind-service my-app my-google-datastore-example -c `{}`
</pre>