# Service Broker Test App

*sb-test* is a very naive [Golang](https://golang.org) demo application that demonstrates how to use bindings given by the [GCP Service Broker](https://github.com/GoogleCloudPlatform/gcp-service-broker).

This is not an official Google product.

## To run the application on Cloud Foundry

1. Log in to your Cloud Foundry using the `cf login` command.

1. From the main project directory, push your app to Cloud Foundry using the `cf push` command. Take note of the route to your app.
1. Create a new service or bind to an existing one, and hit the test-<servicename> endpoint to test your bindings. For example:

1. Create a Storage Bucket:
    ```
	cf create-service google-storage standard practice-storage -c '{"name": "my-practice-bucket"}'
    ```

1. Bind the bucket to your app and give the service account storage object admin permissions:
    ```
    cf bind-service sb-test practice-storage -c '{"role":"storage.objectAdmin"}'
    ```

1. Restage the app so the new environment variables take effect:
    ```
    cf restage sb-test
    ```

1. See your binding in action
    ```
    https://sb-test.<your app domain>.com/test-storage
    ```

    Your bucket should be in the list!