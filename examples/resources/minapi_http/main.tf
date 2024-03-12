terraform {
  required_providers {
    minapi = {
      source = "hashicorp.com/terraform-provider/minapi"
    }
  }
}

provider "minapi" {}

# resource "minapi_http" "cmdb" {
#   url     = "https://jsonplaceholder.typicode.com/posts"
#   method  = "POST"
#   payload = "{\"title\":\"bar\"}"

# }

# resource "minapi_http" "api2" {
#   url     = "https://reqres.in/api/"
#   method  = "POST"
#   payload = "{   \"name\": \"creazer\",    \"job\": \"terrafrom-gny\"}"
# }

resource "minapi_http" "cmdb_rds" {
  url    = "http://localhost:3000/test"
  method = "POST"
  request_headers = {
    Authorization = "fakeservertoken"
    Content-Type  = "application/json"
  }
  payload = "{\n    \"arn\": \"ar:aws:rds:us-east-1:721031640279:db:apgsql-apmtest100-00dev01\",\n    \"account_id\": \"721031640279\",\n    \"engine\": \"postgres\",\n    \"company\": \"ElevanceHealth\",\n    \"u_cloud_provider\": \"AWS (Amazon)\",\n    \"u_support_group\": \"Cloud Database Support\",\n    \"environment\": \"development\",\n    \"u_ip_address\": \"N/A\",\n    \"engine_version\": \"14.6\",\n    \"address\": \"N/A\",\n    \"tcp_port\": \"5432\",\n    \"u_connection_string\": \"N/A\",\n    \"storage_encrypted\": \"True\",\n    \"db_instance_identifier\": \"apgsql-apm1003207-00dev01\",\n    \"u_domain_name\": \"rds.amazonaws.com\",\n    \"u_instance_type\": \"DBMS Instance\",\n    \"db_cluster_identifier\": \"N/A\",\n    \"u_instance_readreplica\": \"False\",\n    \"db_name\": \"apgsqlapm1003207dev\",\n    \"apm_number\": \"apm1003206\",\n    \"u_support_model\": \"Outsourced\",\n    \"u_network_access_type\": \"Open\",\n    \"u_authentication_method\": \"DBMS\"\n}\n"
}


