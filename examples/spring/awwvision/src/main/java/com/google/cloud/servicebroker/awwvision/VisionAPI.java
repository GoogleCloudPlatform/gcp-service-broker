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

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import com.google.api.services.vision.v1.Vision;
import com.google.api.services.vision.v1.model.AnnotateImageRequest;
import com.google.api.services.vision.v1.model.AnnotateImageResponse;
import com.google.api.services.vision.v1.model.BatchAnnotateImagesRequest;
import com.google.api.services.vision.v1.model.BatchAnnotateImagesResponse;
import com.google.api.services.vision.v1.model.Feature;
import com.google.api.services.vision.v1.model.Image;
import com.google.common.collect.ImmutableList;

/**
 * Helper methods for interacting with the Cloud Vision API.
 */
@Component
public class VisionAPI {

  @Autowired
  private Vision vision;

  /**
   * Calls the Vision API to get a single label for the given image.
   * @param bytes The image, in bytes
   * @return The label given to the image, or null if the request could not successfully label the image
   * @throws IOException
   */
  public String labelImage(byte[] bytes) throws IOException {
    AnnotateImageRequest request =
        new AnnotateImageRequest().setImage(new Image().encodeContent(bytes)).setFeatures(
            ImmutableList.of(new Feature().setType("LABEL_DETECTION").setMaxResults(1)));
    return sendAndParseRequest(request);
  }

  private String sendAndParseRequest(AnnotateImageRequest request) throws IOException {
    AnnotateImageResponse response = sendRequest(request);
    if (response == null) {
      return null;
    }
    if (response.getLabelAnnotations() == null) {
      throw new IOException(response.getError() != null ? response.getError().getMessage()
          : "Unknown error getting image annotations");
    }
    return response.getLabelAnnotations().get(0).getDescription();
  }

  private AnnotateImageResponse sendRequest(AnnotateImageRequest request) throws IOException {
    Vision.Images.Annotate annotate = vision.images()
        .annotate(new BatchAnnotateImagesRequest().setRequests(ImmutableList.of(request)));

    BatchAnnotateImagesResponse batchResponse = annotate.execute();
    if (batchResponse == null || batchResponse.getResponses() == null) {
      return null;
    }
    return batchResponse.getResponses().get(0);
  }
}
