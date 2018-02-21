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

import org.apache.commons.io.IOUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.client.RestTemplate;

import com.google.api.services.storage.model.StorageObject;
import com.google.cloud.servicebroker.awwvision.RedditResponse.Listing;
import com.google.common.collect.ImmutableMap;

/**
 * Provides a request mapping for scraping images from reddit, labeling them with the Vision API,
 * and storing them in Cloud Storage.
 */
@Controller
public class RedditScraper {

  private static final String REDDIT_URL = "https://www.reddit.com/r/aww/hot.json";

  @Autowired
  private VisionAPI visionAPI;
  @Autowired
  private StorageAPI storageAPI;

  private final Log logger = LogFactory.getLog(getClass());

  @Value("${reddit-user-agent}")
  private String redditUserAgent;

  @RequestMapping("/reddit")
  String getRedditUrls(Model model, RestTemplate restTemplate) throws GeneralSecurityException {
    HttpHeaders headers = new HttpHeaders();
    headers.add(HttpHeaders.USER_AGENT, redditUserAgent);
    RedditResponse response = restTemplate
        .exchange(REDDIT_URL, HttpMethod.GET, new HttpEntity<String>(headers), RedditResponse.class)
        .getBody();

    storeAndLabel(response);

    return "reddit";
  }

  void storeAndLabel(RedditResponse response) throws GeneralSecurityException {
    for (Listing listing : response.data.children) {
      if (listing.data.preview != null) {
            URL url;
            byte[] raw;
            try {
              url = new URL(listing.data.url);
              raw = download(url);
            } catch (IOException e) {
              logger.warn("Issue in streaming image " + listing.data.url, e);
              continue;
            }
            try {
              // Only label and upload the image if it does not already exist in storage.
              StorageObject existing = storageAPI.get(listing.data.url);
              if (existing == null) {
                String label = visionAPI.labelImage(raw);
                if (label != null) {
                  storageAPI.uploadJpeg(listing.data.url, url, ImmutableMap.of("label", label));
                }
              }
            } catch (IOException e) {
              logger.error("Issue with labeling image " + listing.data.url, e);
            }
       }
    }
  }

  byte[] download(URL url) throws IOException {
    return IOUtils.toByteArray(url.openStream());
  }
}
