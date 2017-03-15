#!/bin/bash
sb_name=gcp-service-broker
cf delete-service-broker $sb_name -f
cf create-service-broker $sb_name admin admin http://host.pcfdev.io:8010 --space-scoped
