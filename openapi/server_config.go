package openapi

import (
	"errors"

	"github.com/ansonallard/go_utils/openapi/middleware/authz"
)

type serveOpenAPIConfig struct {
	openAPISpecFilePath string
	openAPISpecData     []byte
	authZMiddleware     authz.AuthZ
	isDevMode           bool
	port                uint16
	serviceController   interface{}
}

func NewServeOpenAPIConfig() *serveOpenAPIConfig {
	return &serveOpenAPIConfig{}
}

func (soc *serveOpenAPIConfig) WithOpenAPISpecFilePath(openAPISpecFilePath string) *serveOpenAPIConfig {
	soc.openAPISpecFilePath = openAPISpecFilePath
	return soc
}

func (soc *serveOpenAPIConfig) WithOpenAPISpecData(openAPISpecData []byte) *serveOpenAPIConfig {
	soc.openAPISpecData = openAPISpecData
	return soc
}

func (soc *serveOpenAPIConfig) WithAuthZMiddleware(middlewareFunction authz.AuthZ) *serveOpenAPIConfig {
	soc.authZMiddleware = middlewareFunction
	return soc
}

func (soc *serveOpenAPIConfig) WithIsDevMode(isDevMode bool) *serveOpenAPIConfig {
	soc.isDevMode = isDevMode
	return soc
}

func (soc *serveOpenAPIConfig) WithPort(port uint16) *serveOpenAPIConfig {
	soc.port = port
	return soc
}

func (soc *serveOpenAPIConfig) WithServiceController(serviceController interface{}) *serveOpenAPIConfig {
	soc.serviceController = serviceController
	return soc
}

// Build validates and returns the config or an error
func (soc *serveOpenAPIConfig) Build() (*serveOpenAPIConfig, error) {
	// Validate mutually exclusive options
	if soc.openAPISpecFilePath != "" && soc.openAPISpecData != nil {
		return nil, errors.New("cannot set both openAPISpecFilePath and openAPISpecData")
	}

	// Validate at least one is set
	if soc.openAPISpecFilePath == "" && soc.openAPISpecData == nil {
		return nil, errors.New("must set either openAPISpecFilePath or openAPISpecData")
	}

	if soc.port == 0 {
		return nil, errors.New("must provide valid port.")
	}

	if soc.serviceController == nil {
		return nil, errors.New("must provide valid service controller.")
	}

	return soc, nil
}
