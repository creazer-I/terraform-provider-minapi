---
page_title: "minapi_http Resource - terraform-provider-minapi"
subcategory: ""
description: |-
  The minapi resource makes an HTTP POST request to the given URL and exports
  information about the response.
---

# minapi_http (Resource)

The `minapi` resource makes an HTTP POST request to the given URL and exports
information about the response.


## Example Usage

```hcl

provider "minapi" {}

resource "minapi_http" "example" {
 method          = "POST"
 url             = "https://api.example.com/endpoint"
 payload         = jsonencode({ "key" = "value" })
 request_headers = {
   "Content-Type" = "application/json"
 }
}

output "response_body" {
 value = minapi_http.example.response_body
}
```

## Schema

### Required

- `method` (String) The HTTP method for the request.
- `url` (String) The URL for the request.

### Optional

- `payload` (String) The payload to be sent in the POST request.
- `request_headers` (Map of String) A map of request header field names and values.

### Read-Only

- `id` (String) The URL used for the request.
- `response_body` (String) The response body returned as a string.
- `status_code` (Number) The HTTP status code of the response.
