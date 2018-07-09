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

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.security.GeneralSecurityException;
import java.util.Base64;

import org.json.JSONObject;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.client.RestTemplate;

import com.google.api.client.googleapis.auth.oauth2.GoogleCredential;
import com.google.api.client.googleapis.javanet.GoogleNetHttpTransport;
import com.google.api.client.http.HttpTransport;
import com.google.api.client.json.JsonFactory;
import com.google.api.client.json.jackson2.JacksonFactory;
import com.google.api.services.storage.Storage;
import com.google.api.services.storage.StorageScopes;

/**
 * Sets up connections to client libraries and other injectable beans.
 */
@Configuration
public class StorageConfig {

  @Value("${gcp-application-name}")
  private String applicationName;

  @Bean
  JsonFactory jsonFactory() {
    return JacksonFactory.getDefaultInstance();
  }

  @Bean
  HttpTransport transport() throws GeneralSecurityException, IOException {
    return GoogleNetHttpTransport.newTrustedTransport();
  }

  @Bean
  GoogleCredential credential() throws IOException {
    String env = System.getenv("VCAP_SERVICES");
    
    String privateKeyData =
        new JSONObject(env)
          .getJSONArray("google-storage")
          .getJSONObject(0)
          .getJSONObject("credentials")
          .getString("PrivateKeyData");

    InputStream stream = new ByteArrayInputStream(Base64.getDecoder().decode(privateKeyData));
    return GoogleCredential.fromStream(stream);
  }

  @Bean
  Storage storage(HttpTransport transport, JsonFactory jsonFactory, GoogleCredential credential) {
    if (credential.createScopedRequired()) {
      credential = credential.createScoped(StorageScopes.all());
    }
    return new Storage.Builder(transport, jsonFactory, credential)
        .setApplicationName(applicationName).build();
  }

  @Bean
  RestTemplate restTemplate() {
    return new RestTemplate();
  }
}
