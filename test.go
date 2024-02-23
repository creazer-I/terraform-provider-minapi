internal\provider\minapi_resource_http.go:5:2: "encoding/json" imported and not used
internal\provider\minapi_resource_http.go:6:2: "fmt" imported and not used
internal\provider\minapi_resource_http.go:16:27: cannot use (*MinAPIHttpResource)(nil) (value of type *MinAPIHttpResource) as resource.Resource value in variable declaration: *MinAPIHttpResource does not implement resource.Resource (missing method Delete)
internal\provider\minapi_resource_http.go:19:9: cannot use &MinAPIHttpResource{} (value of type *MinAPIHttpResource) as resource.Resource value in return statement: *MinAPIHttpResource does not implement resource.Resource (missing method Delete)
internal\provider\minapi_resource_http.go:79:17: cannot use "id" (constant of type string) as context.Context value in argument to resp.State.Set: string does not implement context.Context (missing method Deadline)
internal\provider\minapi_resource_http.go:80:17: cannot use "response_body" (constant of type string) as context.Context value in argument to resp.State.Set: string does not implement context.Context (missing method Deadline)
