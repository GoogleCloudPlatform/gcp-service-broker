# Awwvision: Spring Boot edition

*Awwvision: Spring Boot edition* is a [Spring Boot](http://projects.spring.io/spring-boot/) demo application that uses the [Google Cloud Vision API](https://cloud.google.com/vision/) to classify (label) images from Reddit's [/r/aww](https://reddit.com/r/aww) subreddit, store the images and classifications in [Google Cloud Storage](https://cloud.google.com/storage/), and display the labeled results in a web app. It uses the GCP Service Broker to authenticate to the Vision and Storage APIs, and is based off of the [Python Awwvision sample app](https://github.com/GoogleCloudPlatform/cloud-vision/tree/master/python/awwvision).

This is not an official Google product.

Awwvision: Spring Boot edition has two endpoints:

1. A webapp that reads and displays the labels and associated images from GCS.
2. A scraper that downloads images from Reddit and classifies them using the Vision API.

## Prerequisites

1. Create a project in the [Google Cloud Platform Console](https://console.cloud.google.com).

1. [Enable billing](https://console.cloud.google.com/project/_/settings) for your project.

1. Enable the [Vision](https://console.cloud.google.com/apis/api/vision.googleapis.com) and [Storage](https://console.cloud.google.com/apis/api/storage_component) APIs. See the [Vision API Quickstart](https://cloud.google.com/vision/docs/quickstart) and [Storage API Quickstart](https://cloud.google.com/storage/docs/quickstart-console) for more information on using the two APIs.

## Run the application on Cloud Foundry

1. Log in to your Cloud Foundry using the `cf login` command.

1. From the main project directory, build an executable jar and push it to Cloud Foundry. This step will initially fail due to lack of credentials.
    ```
    mvn package -DskipTests && cf push -p target/awwvision-spring-0.0.1-SNAPSHOT.jar awwvision --no-start
    ```

1. Create a Storage Bucket:
    ```
    cf create-service google-storage standard awwvision-storage
    ```

1. Bind the bucket to your app and give the service account storage object admin permissions:
    ```
    cf bind-service awwvision awwvision-storage -c '{"role":"storage.objectAdmin"}'
    ```

1. Create a Machine Learning API Instance:
    ```
    cf create-service google-ml-apis default ml
    ```

1. Bind the machine learning API instance to your app:
    ```
    cf bind-service awwvision ml  -c '{"role":"ml.viewer"}'
    ```

1. Start the app so the new environment variables take effect:
    ```
    cf start awwvision
    ```

### Visit the application and start the crawler

Once your application is running, visit awwvision.\[your-cf-instance-url\]/reddit to start crawling. The page will display "Scrape completed." once it is done. From there, visit awwvision.\[your-cf-instance-url\] to view your images!
