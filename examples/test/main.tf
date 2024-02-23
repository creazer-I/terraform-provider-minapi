terraform {
  required_providers {
    minapi = {
      source = "hashicorp.com/test/minapi"
  }
  }
}

provider "minapi" {}

#data "minapi" "example" {}

# resource "minapi_http" "api1" {
#   url = "https://jsonplaceholder.typicode.com/posts"
#   method = "POST"
#   payload = "{\"title\":\"foo\"}"

# }


resource "minapi_http" "api2" {
  url = "https://reqres.in/api/users"
  method = "POST"
  payload = "{   \"name\": \"creazer\",    \"job\": \"terrafrom-developer-guy\"}"

}
# output "name" {
#   value = minapi_http.api1.response_body
# }
