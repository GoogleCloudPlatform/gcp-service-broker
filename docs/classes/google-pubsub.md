# <a name="google-pubsub"></a> ![](https://cloud.google.com/_static/images/cloud/products/logos/svg/pubsub.svg) Google PubSub
A global service for real-time and reliable messaging and streaming data.

 * [Documentation](https://cloud.google.com/pubsub/docs/)
 * [Support](https://cloud.google.com/pubsub/docs/support)
 * Catalog Metadata ID: `628629e3-79f5-4255-b981-d14c6c7856be`
 * Tags: gcp, pubsub
 * Service Name: `google-pubsub`

## Provisioning

**Request Parameters**


 * `topic_name` _string_ - Name of the topic. Must not start with "goog". Default: `pcf_sb_${counter.next()}_${time.nano()}`.
    * The string must have at most 255 characters.
    * The string must have at least 3 characters.
    * The string must match the regular expression `^[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+$`.
 * `subscription_name` _string_ - Name of the subscription. Blank means no subscription will be created. Must not start with "goog". Default: ``.
    * The string must have at most 255 characters.
    * The string must have at least 0 characters.
    * The string must match the regular expression `^(|[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+)`.
 * `is_push` _string_ - Are events handled by POSTing to a URL? Default: `false`.
    * The value must be one of: [false true].
 * `endpoint` _string_ - If `is_push` == 'true', then this is the URL that will be pushed to. Default: ``.
 * `ack_deadline` _string_ - Value is in seconds. Max: 600 This is the maximum time after a subscriber receives a message before the subscriber should acknowledge the message. After message delivery but before the ack deadline expires and before the message is acknowledged, it is an outstanding message and will not be delivered again during that time (on a best-effort basis).  Default: `10`.


## Binding

**Request Parameters**


 * `role` _string_ - The role for the account without the "roles/" prefix. See: https://cloud.google.com/iam/docs/understanding-roles for more details. Default: `pubsub.editor`.
    * The value must be one of: [pubsub.editor pubsub.publisher pubsub.subscriber pubsub.viewer].

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
 * `subscription_name` _string_ - Name of the subscription.
    * The string must have at most 255 characters.
    * The string must have at least 0 characters.
    * The string must match the regular expression `^(|[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+)`.
 * `topic_name` _string_ - **Required** Name of the topic.
    * The string must have at most 255 characters.
    * The string must have at least 3 characters.
    * The string must match the regular expression `^[a-zA-Z][a-zA-Z0-9\d\-_~%\.\+]+$`.

## Plans

The following plans are built-in to the GCP Service Broker and may be overridden
or disabled by the broker administrator.


* **`default`**
  * Plan ID: `622f4da3-8731-492a-af29-66a9146f8333`.
  * Description: PubSub Default plan.
  * This plan doesn't override user variables on provision.
  * This plan doesn't override user variables on bind.


## Examples




### Basic Configuration


Create a topic and a publisher to it.
Uses plan: `622f4da3-8731-492a-af29-66a9146f8333`.

**Provision**

```javascript
{
    "subscription_name": "example_topic_subscription",
    "topic_name": "example_topic"
}
```

**Bind**

```javascript
{
    "role": "pubsub.publisher"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-pubsub default my-google-pubsub-example -c `{"subscription_name":"example_topic_subscription","topic_name":"example_topic"}`
$ cf bind-service my-app my-google-pubsub-example -c `{"role":"pubsub.publisher"}`
</pre>


### No Subscription


Create a topic without a subscription.
Uses plan: `622f4da3-8731-492a-af29-66a9146f8333`.

**Provision**

```javascript
{
    "topic_name": "example_topic"
}
```

**Bind**

```javascript
{
    "role": "pubsub.publisher"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-pubsub default my-google-pubsub-example -c `{"topic_name":"example_topic"}`
$ cf bind-service my-app my-google-pubsub-example -c `{"role":"pubsub.publisher"}`
</pre>


### Custom Timeout


Create a subscription with a custom deadline for long processess.
Uses plan: `622f4da3-8731-492a-af29-66a9146f8333`.

**Provision**

```javascript
{
    "ack_deadline": "200",
    "subscription_name": "long_deadline_subscription",
    "topic_name": "long_deadline_topic"
}
```

**Bind**

```javascript
{
    "role": "pubsub.publisher"
}
```

**Cloud Foundry Example**

<pre>
$ cf create-service google-pubsub default my-google-pubsub-example -c `{"ack_deadline":"200","subscription_name":"long_deadline_subscription","topic_name":"long_deadline_topic"}`
$ cf bind-service my-app my-google-pubsub-example -c `{"role":"pubsub.publisher"}`
</pre>