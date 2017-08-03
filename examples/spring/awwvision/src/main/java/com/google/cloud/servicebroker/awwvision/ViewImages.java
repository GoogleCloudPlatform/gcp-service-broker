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
import java.security.GeneralSecurityException;
import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;

import com.google.api.services.storage.model.StorageObject;

/**
 * Provides request mappings for reading images from Cloud Storage.
 */
@Controller
public class ViewImages {

  @Autowired
  private StorageAPI storageAPI;

  @RequestMapping("/")
  public String view(Model model) throws IOException, GeneralSecurityException {
    List<StorageObject> objects = storageAPI.listAll();
    List<Image> images = new ArrayList<>();
    for (StorageObject obj : objects) {
      Image image = new Image(getPublicUrl(storageAPI.bucketName, obj.getName()), obj.getMetadata().get("label"));
      images.add(image);
    }
    model.addAttribute("images", images);
    return "index";
  }

  @RequestMapping("/label/{label}")
  String viewLabel(@PathVariable("label") String label, Model model)
      throws IOException, GeneralSecurityException {
    List<StorageObject> objects = storageAPI.listAll();
    List<Image> images = new ArrayList<>();
    for (StorageObject obj : objects) {
      Image image = new Image(getPublicUrl(storageAPI.bucketName, obj.getName()), obj.getMetadata().get("label"));
      if (image.label.equals(label)) {
        images.add(image);
      }
    }
    model.addAttribute("images", images);
    return "index";
  }

  static String getPublicUrl(String bucket, String object) {
    return String.format("http://storage.googleapis.com/%s/%s", bucket, object);
  }

  static class Image {
    private String URL;
    private String label;
    
    public Image(String URL, String label) {
      this.URL = URL;
      this.label = label;
    }
    
    public String getURL() {
      return URL;
    }
    
    public String getLabel() {
      return label;
    }
    
    @Override
    public boolean equals(Object object) {
      if (!(object instanceof Image)) {
        return false;
      }
      Image other = (Image) object;
      return other.getURL().equals(getURL()) && other.getLabel().equals(getLabel());
    }
  }
}
