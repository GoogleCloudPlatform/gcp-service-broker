# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.14-alpine AS build

WORKDIR /go/src/github.com/GoogleCloudPlatform/gcp-service-broker
COPY . .

RUN CGO_ENABLED=0 go build -mod=mod -o /bin/gcp-service-broker

# Get latest CA certs
FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/GoogleCloudPlatform/gcp-service-broker /src
COPY --from=build /bin/gcp-service-broker /bin/gcp-service-broker

ENV PORT 8080
EXPOSE 8080/tcp

WORKDIR /tmp
WORKDIR /

ENTRYPOINT ["/bin/gcp-service-broker"]
CMD ["help"]
