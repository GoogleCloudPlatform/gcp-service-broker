#!/usr/bin/env bash

mvn -DskipTests=true clean install 

service_name=awwvision-storage
app_name=awwvision

cf d -f $app_name
cf ds -f $service_name

cf a
cf s

cf create-service google-storage standard $service_name # -c '{"name": "awwvision-bucket"}'
cf push -f ./manifest.yml --no-start
cf bind-service $app_name awwvision-storage -c '{"role":"storage.objectAdmin"}'
cf restart $app_name