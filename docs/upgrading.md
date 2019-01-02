## Upgrading

### v3.X to 4.0

Version 4.0 of the broker contains significant improvements to security, ability to self-service, and documentation.
Specifically, new forms allow you to:

* Set role whitelists on your services
* Enable or disable services entirely
* Set the database name

Plans also come with a GUID so they'll be consistent across installs.

If you:

* Changed service details, description, names, GUID, or URLs via environment variables:
  * These will continue to work, but are deprecated.
* Rely on custom roles:
  * You will need to add the roles to the whitelists for the services they apply to.
* Used built-in plans prior to 4.0:
  * Enable **Compatibility with GCP Service Broker v3.X plans** in the **Compatibility** form.
  * You will also need to enable the new plans for your developers so they can upgrade.
  * If you are not using the tile, you can set the `GSB_COMPATIBILITY_THREE_TO_FOUR_LEGACY_PLANS`
    environment variable to true.
* Have scripts using the Cloud Storage `reduced_availability` plan:
  * Upgrade them to use `reduced-availability` instead (underscore to dash).
* Wanted to completely disable some services:
  * You can now do that using the enable/disable service properties.
* Set a custom environment variable to change the name of the database:
  * You can now set it in the database form of the PCF tile.
* Use a BigQuery billing export for chargebacks:
  * Read the [billing docs](https://github.com/GoogleCloudPlatform/gcp-service-broker/blob/master/docs/billing.md) to understand how labels are automatically applied to services now.

### v3.X to 5.0

You MUST upgrade your 3.X version to 4.x before upgrading to 5.x.
