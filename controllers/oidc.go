package controllers

import (
	"net/http"
	"strings"

	"github.com/equinor/no-factor-auth/oidc"

	"github.com/labstack/echo/v4"
)

// StdOidcConfigURI is the standard endpoint for oidc config
const StdOidcConfigURI = "/.well-known/openid-configuration"

// OidcConfig returns config for host
func OidcConfig(c echo.Context) error {

	// baseUrl := c

	hostURL := "http://" + c.Request().Host + strings.TrimSuffix(c.Request().URL.String(), StdOidcConfigURI)
	oidc := oidc.Default()
	oidc.JwksURI = hostURL + "/discovery/keys"
	oidc.Issuer = hostURL
	oidc.AuthorizationEndpoint = hostURL + "/oauth2/authorize"
	return c.JSON(http.StatusOK, oidc)
}
