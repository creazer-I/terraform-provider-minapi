
internal\provider\minapi_resource_http.go:14:27: cannot use (*MinAPIHttpResource)(nil) (value of type *MinAPIHttpResource) as resource.Resource value in variable declaration: *MinAPIHttpResource does not implement resource.Resource (missing method Delete)
internal\provider\minapi_resource_http.go:17:9: cannot use &MinAPIHttpResource{} (value of type *MinAPIHttpResource) as resource.Resource value in return statement: *MinAPIHttpResource does not implement resource.Resource (missing method Delete)
internal\provider\minapi_resource_http.go:77:28: too many arguments in call to resp.State.Set
        have (context.Context, string, string)
        want (context.Context, interface{})
internal\provider\minapi_resource_http.go:78:39: too many arguments in call to resp.State.Set
        have (context.Context, string, basetypes.StringValue)
        want (context.Context, interface{})
