package oidc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOidcV2(t *testing.T){
	oidc := OidcV2("uri")
	assert.Equal(t, oidc.Issuer, "uri/v2.0")
	assert.Equal(t, oidc.AuthorizationEndpoint, "uri/oauth2/v2.0/authorize")
	assert.Equal(t, oidc.JwksURI, "uri/discovery/v2.0/keys")
}

