#!/bin/bash

script_dir=`dirname $0`
upload_dir=$script_dir/upload
mkdir -p $upload_dir

function stage(){

  pushd $script_dir/../../.. > /dev/null
  export GOPATH=`pwd`
  popd > /dev/null

  binary=$GOPATH/bin/gcp-service-broker
  # rm -rf $binary

  cd $script_dir
  [[ -f $binary  ]] || ./build_for_linux.sh

  cp $binary $upload_dir
  echo "web: ./gcp-service-broker" > $upload_dir/Procfile
  cp $script_dir/../manifest-linux-binary.yml $upload_dir/manifest.yml
  ls -la
  echo $upload_dir
}


function app_domain(){
    D=`cf apps | grep $1 | tr -s ' ' | cut -d' ' -f 6 | cut -d, -f1`
    echo $D
}


function  deploy() {
  json_file=$1
  cd $upload_dir
  echo $json_file

  sb_an=gcp-service-broker
  mysql_db=${sb_an}-db
  #db_service_name=cleardb
  db_service_name=p-mysql
  db_service_plan=512mb

  # reset

  cf d -f $sb_an
  cf ds -f $mysql_db

  echo "deleted app and service if they exist."

  # deploy
  cf cs $db_service_name $db_service_plan $mysql_db
  cf push  --random-route --no-start
  cf bs $sb_an $mysql_db

  cf env $sb_an > json

  json=$(cat json |sed -e '1,/System-Provided:/d' | sed -e '/^}/q')


  dbHostname=$( echo $json | jq -r .VCAP_SERVICES.\"${db_service_name}\"[0].credentials.hostname )
  dbName=$( echo $json | jq -r .VCAP_SERVICES.\"${db_service_name}\"[0].credentials.name )
  dbPw=$( echo $json | jq -r .VCAP_SERVICES.\"${db_service_name}\"[0].credentials.password )
  dbPort=$( echo $json | jq -r .VCAP_SERVICES.\"${db_service_name}\"[0].credentials.port )
  dbUser=$( echo $json | jq -r .VCAP_SERVICES.\"${db_service_name}\"[0].credentials.username )

  rm json

  cf set-env $sb_an DB_HOST $dbHostname
  cf set-env $sb_an DB_NAME $dbName
  cf set-env $sb_an DB_PORT $dbPort
  cf set-env $sb_an DB_PASSWORD $dbPw
  cf set-env $sb_an DB_USERNAME $dbUser
  cf set-env $sb_an SECURITY_USER_NAME admin
  cf set-env $sb_an SECURITY_USER_PASSWORD admin


  echo "about to set JSON"
  cf set-env $sb_an ROOT_SERVICE_ACCOUNT_JSON "$(< $json_file)"

  route=`app_domain $sb_an`
  echo $route


  cf restart $sb_an

  cf delete-service-broker $sb_an -f
  cf create-service-broker $sb_an admin admin http://`app_domain $sb_an` --space-scoped

}

stage
deploy $HOME/Desktop/gcp.json
