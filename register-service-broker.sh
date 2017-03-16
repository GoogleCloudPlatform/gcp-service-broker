#!/bin/bash

json_file=""

if [[ -z "$1" ]]; then
 json_file=${HOME}/Desktop/gcp.json
else
  json_file=$1
fi

echo $json_file

# first we need to deploy the go-based service broker as an app

function app_domain(){
    D=`cf apps | grep $1 | tr -s ' ' | cut -d' ' -f 6 | cut -d, -f1`
    echo $D
}


sb_an=gcp-service-broker
mysql_db=${sb_an}-db

# reset

cf d -f $sb_an
cf ds -f $mysql_db

echo "deleted app and service if they exist."

# deploy

cf cs p-mysql 100mb $mysql_db
cf bs $sb_an $mysql_db
cf push --random-route --no-start

cf env $sb_an > json

json=$(cat json |sed -e '1,/System-Provided:/d' | sed -e '/^}/q')

dbHostname=$( echo $json | jq -r .VCAP_SERVICES.\"p-mysql\"[0].credentials.hostname )
dbName=$( echo $json | jq -r .VCAP_SERVICES.\"p-mysql\"[0].credentials.name )
dbPw=$( echo $json | jq -r .VCAP_SERVICES.\"p-mysql\"[0].credentials.password )
dbPort=$( echo $json | jq -r .VCAP_SERVICES.\"p-mysql\"[0].credentials.port )
dbUser=$( echo $json | jq -r .VCAP_SERVICES.\"p-mysql\"[0].credentials.username )

rm json

cf set-env $sb_an DB_HOST $dbHostname
cf set-env $sb_an DB_NAME $dbName
cf set-env $sb_an DB_PORT $dbPort
cf set-env $sb_an DB_PASSWORD $dbPw
cf set-env $sb_an DB_USERNAME $dbUser
cf set-env $sb_an SECURITY_USER_NAME admin
cf set-env $sb_an SECURITY_USER_PASSWORD admin


echo "about to set JSON"
json_str=`cat $json_file | tr "\n" " "`
echo $json_str
cf set-env $sb_an ROOT_SERVICE_ACCOUNT_JSON '${json_str}'

route=`app_domain $sb_an`
echo $route


cf restart $sb_an
#
# sb_name=gcp-service-broker
# cf delete-service-broker $sb_name -f
# cf create-service-broker $sb_name admin admin http://host.pcfdev.io:8010 --space-scoped
