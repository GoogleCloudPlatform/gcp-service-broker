/*
 * Copyright 2016 Google Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */
package com.google.cloud.servicebroker.awwvision;

import java.io.IOException;
import java.net.URL;
import java.security.GeneralSecurityException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Map;
import org.json.JSONObject;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Controller;

import com.google.api.client.http.InputStreamContent;
import com.google.api.services.storage.Storage;
import com.google.api.services.storage.model.ObjectAccessControl;
import com.google.api.services.storage.model.Objects;
import com.google.api.services.storage.model.StorageObject;

/**
 * Helper methods for interacting with the Cloud Storage API.
 * 
 * Uses the Cloud Storage Bucket configured in the application properties.
 */
@Controller
public class StorageAPI {

  @Autowired
  private Storage storageService;

  public final String bucketName;

  public StorageAPI(){
    String env = System.getenv("VCAP_SERVICES");

    this.bucketName =
        new JSONObject(env)
          .getJSONArray("google-storage")
          .getJSONObject(0)
          .getJSONObject("credentials")
          .getString("bucket_name");
  }

  /**
   * Uploads a JPEG image to Cloud Storage.
   * @param name The name of the image
   * @param url A URL pointing to the image
   * @param metadata Metadata about the image
   * @throws IOException
   * @throws GeneralSecurityException
   */
  public void uploadJpeg(String name, URL url, Map<String, String> metadata)
      throws IOException, GeneralSecurityException {
    InputStreamContent contentStream = new InputStreamContent("image/jpeg", url.openStream());
    StorageObject objectMetadata = new StorageObject().setName(name)
        .setAcl(Arrays.asList(new ObjectAccessControl().setEntity("allUsers").setRole("READER")))
        .setMetadata(metadata);

    storageService.objects().insert(this.bucketName, objectMetadata, contentStream).execute();
  }

  /**
   * Returns a List of all objects in the configured Cloud Storage bucket.
   * @throws IOException
   * @throws GeneralSecurityException
   */
  public List<StorageObject> listAll() throws IOException, GeneralSecurityException {
    Storage.Objects.List listRequest = storageService.objects().list(this.bucketName);

    List<StorageObject> results = new ArrayList<StorageObject>();
    Objects objects;

    // Iterate through each page of results, and add them to our results list.
    do {
      objects = listRequest.execute();
      if (objects.getItems() == null) {
        break;
      }
      // Add the items in this page of results to the list we'll return.
      results.addAll(objects.getItems());

      // Get the next page, in the next iteration of this loop.
      listRequest.setPageToken(objects.getNextPageToken());
    } while (null != objects.getNextPageToken());

    return results;
  }

  /**
   * Gets a specific object in the configured Cloud Storage bucket.
   * @param name The name of the object.
   * @return The StorageObject with the specified name, or null if one does not exist.
   */
  public StorageObject get(String name) {
    try {
      return storageService.objects().get(this.bucketName, name).execute();
    } catch (IOException e) {
      return null;
    }
  }
}
