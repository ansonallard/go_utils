package authz

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3filter"
)

type AuthZ interface {
	AuthorizeCaller(ctx context.Context, ai *openapi3filter.AuthenticationInput) error
}

type AuthorizationFunction = func(ctx context.Context, ai *openapi3filter.AuthenticationInput) error
