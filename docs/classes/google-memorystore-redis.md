# <a name="google-memorystore-redis"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/cache.svg) Google Cloud Memorystore for Redis API
Creates and manages Redis instances on the Google Cloud Platform.

 * [Documentation](https://cloud.google.com/memorystore/docs/redis)
 * [Support](https://cloud.google.com/memorystore/docs/redis/support)
 * Catalog Metadata ID: `3ea92b54-838c-4fe1-b75d-9bda513380aa`
 * Tags: gcp, memorystore, redis
 * Service Name: `google-memorystore-redis`

## Provisioning

**Request Parameters**


 * `instance_id` _string_ - The name of the instance. The name must be unique per project. Default: `gsb-${counter.next()}-${time.nano()}`.
    * The string must have at most 40 characters.
    * The string must have at least 1 characters.
    * The string must match the regular expression `^[a-z]([-0-9a-z]*[a-z0-9]$)*`.
 * `authorized_network` _string_ - The name of the VPC network to attach the instance to. Default: `default`.
    * Examples: [default projects/MYPROJECT/global/networks/MYNETWORK].
 * `region` _string_ - The region to create the instance in. Supported regions can be found here: https://cloud.google.com/memorystore/docs/redis/regions. Default: `us-east1`.
    * The string must match the regular expression `^[A-Za-z][-a-z0-9A-Z]+$`.
 * `memory_size_gb` _integer_ - Redis memory size in GiB. Default: `4`.
 * `tier` _string_ - The performance tier. Default: `BASIC`.
    * The value must be one of: [BASIC STANDARD_HA].


## Binding

**Request Parameters**

_No parameters supported._

**Response Parameters**

 * `authorized_network` _string_ - Name of the VPC network the instance is attached to.
 * `reserved_ip_range` _string_ - Range of IP addresses reserved for the instance.
 * `redis_version` _string_ - The version of Redis software.
 * `memory_size_gb` _integer_ - Redis memory size in GiB.
 * `host` _string_ - Hostname or IP address of the exposed Redis endpoint used by clients to connect to the service.
 * `port` _integer_ - The port number of the exposed Redis endpoint.
 * `uri` _string_ - URI of the instance.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `df10762e-6ef1-44e3-84c2-07e9358ceb1f`.
  * Description: Lets you chose your own values for all properties.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.
* **`basic`**
  * Plan ID: `dd1923b6-ac26-4697-83d6-b3a0c05c2c94`.
  * Description: Provides a standalone Redis instance. Use this tier for applications that require a simple Redis cache.
  * This plan overrides the following user variables on provision.
    * `service_tier` = `BASIC`
  * This plan doesn't override user variables on bind.
* **`standard_ha`**
  * Plan ID: `41771881-b456-4940-9081-34b6424744c6`.
  * Description: Provides a highly available Redis instance.
  * This plan overrides the following user variables on provision.
    * `service_tier` = `STANDARD_HA`
  * This plan doesn't override user variables on bind.


## Examples




### Standard Redis Configuration


Create a Redis instance with standard service tier.
Uses plan: `dd1923b6-ac26-4697-83d6-b3a0c05c2c94`.

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
$ cf create-service google-memorystore-redis basic my-google-memorystore-redis-example -c `{}`
$ cf bind-service my-app my-google-memorystore-redis-example -c `{}`
</pre>


### HA Redis Configuration


Create a Redis instance with high availability.
Uses plan: `41771881-b456-4940-9081-34b6424744c6`.

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
$ cf create-service google-memorystore-redis standard_ha my-google-memorystore-redis-example -c `{}`
$ cf bind-service my-app my-google-memorystore-redis-example -c `{}`
</pre>