# <a name="google-filestore"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/storage.svg) Google Cloud Filestore
Fully managed NFS file storage with predictable performance.

 * [Documentation](https://cloud.google.com/filestore/docs/)
 * [Support](https://cloud.google.com/filestore/docs/getting-support)
 * Catalog Metadata ID: `494eb82e-c4ca-4bed-871d-9c3f02f66e01`
 * Tags: gcp, filestore, nfs
 * Service Name: `google-filestore`

## Provisioning

**Request Parameters**


 * `instance_id` _string_ - The name of the instance. The name must be unique per zone. Default: `gsb-${counter.next()}-${time.nano()}`.
    * The string must have at most 63 characters.
    * The string must have at least 1 characters.
    * The string must match the regular expression `^[a-z]([-0-9a-z]*[a-z0-9]$)*`.
 * `zone` _string_ - The zone to create the instance in. Supported zones can be found here: https://cloud.google.com/filestore/docs/regions. Default: `us-west1-a`.
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.
 * `tier` _string_ - The performance tier. Default: `STANDARD`.
    * The value must be one of: [PREMIUM STANDARD].
 * `authorized_network` _string_ - The name of the network to attach the instance to. Default: `default`.
 * `address_mode` _string_ - The address mode of the service. Default: `MODE_IPV4`.
    * The value must be one of: [MODE_IPV4].
 * `capacity_gb` _integer_ - The capacity of the Filestore. Standard minimum is 1TiB and Premium is minimum 2.5TiB. Default: `1024`.


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `authorized_network` _string_ - Name of the VPC network the instance is attached to.
 * `reserved_ip_range` _string_ - Range of IP addresses reserved for the instance.
 * `ip_address` _string_ - IP address of the service.
 * `file_share_name` _string_ - Name of the share.
 * `capacity_gb` _integer_ - Capacity of the share in GiB.
 * `uri` _string_ - URI of the instance.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `e4c83975-e60f-43cf-afde-ebec573c6c2e`.
  * Description: Filestore default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Standard


Creates a standard Filestore.
Uses plan: `e4c83975-e60f-43cf-afde-ebec573c6c2e`.

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
$ cf create-service google-filestore default my-google-filestore-example -c `{}`
$ cf bind-service my-app my-google-filestore-example -c `{}`
</pre>


### Premium


Creates a premium Filestore.
Uses plan: `e4c83975-e60f-43cf-afde-ebec573c6c2e`.

**Provision**

```javascript
{
    "capacity_gb": 2560,
    "tier": "PREMIUM"
}
```

**Bind**

```javascript
{}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-filestore default my-google-filestore-example -c `{"capacity_gb":2560,"tier":"PREMIUM"}`
$ cf bind-service my-app my-google-filestore-example -c `{}`
</pre>