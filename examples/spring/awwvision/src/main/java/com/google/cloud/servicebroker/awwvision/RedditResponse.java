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

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

/**
 * JSON class structure for Reddit API.
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public class RedditResponse {
  public Data data;
  
  public RedditResponse() {}
  public RedditResponse(Data data) {
    this.data = data;
  }

  @JsonIgnoreProperties(ignoreUnknown = true)
  public static class Data {
    public Listing[] children;
    
    public Data() {}
    public Data(Listing[] children) {
      this.children = children;
    }
  }

  @JsonIgnoreProperties(ignoreUnknown = true)
  public static class Listing {
    public ListingData data;

    public Listing() {}
    public Listing(ListingData data) {
      this.data = data;
    }
  }

  @JsonIgnoreProperties(ignoreUnknown = true)
  public static class ListingData {
    public Preview preview;
    public String url;

    public ListingData() {}
    public ListingData(Preview preview) {
      this.preview = preview;
    }
  }

  @JsonIgnoreProperties(ignoreUnknown = true)
  public static class Preview {
    public Image[] images;

    public Preview() {}
    public Preview(Image[] images) {
      this.images = images;
    }
  }

  @JsonIgnoreProperties(ignoreUnknown = true)
  public static class Image {
    public Source source;
    public String id;

    public Image() {}
    public Image(Source source, String id) {
      this.source = source;
      this.id = id;
    }
  }

  @JsonIgnoreProperties(ignoreUnknown = true)
  public static class Source {
    public String url;

    public Source() {}
    public Source(String url) {
      this.url = url;
    }
  }
}
