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

import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.net.URL;

import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Matchers;
import org.mockito.Mockito;
import org.mockito.Spy;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.context.junit4.SpringRunner;

import com.google.api.client.googleapis.auth.oauth2.GoogleCredential;
import com.google.api.services.vision.v1.Vision;
import com.google.api.services.vision.v1.model.AnnotateImageResponse;
import com.google.api.services.vision.v1.model.BatchAnnotateImagesRequest;
import com.google.api.services.vision.v1.model.BatchAnnotateImagesResponse;
import com.google.api.services.vision.v1.model.EntityAnnotation;
import com.google.cloud.servicebroker.awwvision.RedditResponse;
import com.google.cloud.servicebroker.awwvision.RedditScraper;
import com.google.cloud.servicebroker.awwvision.StorageAPI;
import com.google.cloud.servicebroker.awwvision.RedditResponse.Data;
import com.google.cloud.servicebroker.awwvision.RedditResponse.Image;
import com.google.cloud.servicebroker.awwvision.RedditResponse.Listing;
import com.google.cloud.servicebroker.awwvision.RedditResponse.ListingData;
import com.google.cloud.servicebroker.awwvision.RedditResponse.Preview;
import com.google.common.collect.ImmutableList;
import com.google.common.collect.ImmutableMap;

@RunWith(SpringRunner.class)
@AutoConfigureMockMvc
@SpringBootTest(properties = {"gcp-storage-bucket=fake-bucket"})
public class RedditScraperTest {

  @MockBean
  Vision vision;

  @MockBean
  StorageAPI storageAPI;

  // Even though this is not used directly in the test, mock it out so the application doesn't try
  // to read environment variables to set the credential.
  @MockBean
  GoogleCredential googleCredential;

  @Spy
  @Autowired
  RedditScraper scraper;

  @Before
  public void setUp() throws Exception {
    when(storageAPI.get(Matchers.anyString())).thenReturn(null);

    // Have the Vision API return "dog" for any request.
    Vision.Images images = Mockito.mock(Vision.Images.class);
    Vision.Images.Annotate annotate = Mockito.mock(Vision.Images.Annotate.class);
    when(vision.images()).thenReturn(images);
    when(images.annotate(Matchers.any(BatchAnnotateImagesRequest.class))).thenReturn(annotate);
    when(annotate.execute()).thenReturn(
        new BatchAnnotateImagesResponse().setResponses(ImmutableList.of(new AnnotateImageResponse()
            .setLabelAnnotations(ImmutableList.of(new EntityAnnotation().setDescription("dog"))))));

    doReturn("".getBytes()).when(scraper).download(Matchers.any(URL.class));
  }

  @Test
  public void testScrape() throws Exception {
    Image img1 = new Image(new RedditResponse.Source("http://url1"), "img1");
    Image img2 = new Image(new RedditResponse.Source("http://url2"), "img2");
    RedditResponse redditResponse = new RedditResponse(new Data(
        new Listing[] {new Listing(new ListingData(new Preview(new Image[] {img1, img2})))}));

    scraper.storeAndLabel(redditResponse);

    verify(storageAPI).uploadJpeg("img1.jpg", new URL("http://url1"),
        ImmutableMap.of("label", "dog"));
    verify(storageAPI).uploadJpeg("img2.jpg", new URL("http://url2"),
        ImmutableMap.of("label", "dog"));
  }
}
