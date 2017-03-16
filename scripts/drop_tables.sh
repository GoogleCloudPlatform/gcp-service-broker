sql="drop table if exists   cloud_operations ; drop table if exists   migrations ; drop table if exists   plan_details ;  drop table if exists provision_request_details  ;  drop table if exists  service_binding_credentials;  drop table if exists    service_instance_details ;  " 
echo $sql | mysql -u gcp-service-broker -pqwerty servicebroker
