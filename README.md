# go_utils

Shared Go utilities

## openapi

Package for running a REST API defined by an OpenAPI spec. Callers must provide:

- HTTP Port
- File path or go embedded byte data to OpenAPI spec

Clients can use the builder pattern to provide config:

```go
config, err := NewServeOpenAPIConfig().
    WithOpenAPISpecData(data).
    WithPort(8080).
    WithServiceController(serviceController).
    Build()
if err != nil {
    return err
}

if err := ServeOpenAPI(ctx, config); err != nil {
    log.Fatal().Err(err).Msg("Failed to run OpenAPI server")
}
```

### Controllers

The application running this server must implement all controllers in the following manner:

- Method name must match OpenAPI spec `operationId`.
- Method params must be `(ctx context.Context, request request.Request)`, where the request interface is:

```go
type Request interface {
	GetQueryParams() map[string][]string
	GetHeaders() map[string][]string
	GetPathParams() map[string]string
	GetRequestBody() io.ReadCloser
}
```

- Methods should return either `error` or `T, error`.

Example:

```go
type ServiceController interface {
    CreateService(ctx context.Context, request request.Request) (*api.CreateServiceResponse, error)
}
```

### Logging

This package requires the use of `zerolog`. Callers can pass in a zerolog instance via the `ctx` or the library will instantiate a logging instance for you.

### Middleware

#### AuthZ Middleware

Callers can implement the following interface to authenticate API operations:

```go
AuthorizeCaller(ctx context.Context, ai *openapi3filter.AuthenticationInput) error
```
