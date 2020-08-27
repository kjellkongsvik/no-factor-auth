package controllers

import (
	"net/http"
	"strings"

	"github.com/equinor/no-factor-auth/oidc"

	"github.com/labstack/echo/v4"
)

// StdOidcConfigURI is the standard endpoint for oidc config
const StdOidcConfigURI = "/.well-known/openid-configuration"

func OpenIDConfigV2(c echo.Context) error {
	oidc := oidc.OidcV2(hostURLV2(c))
	return c.JSON(http.StatusOK, &oidc)
}

func hostURLV2(c echo.Context) string{
	suffix := strings.TrimSuffix(c.Request().URL.String(), StdOidcConfigURI)
	return "http://" + c.Request().Host + strings.TrimSuffix(suffix, "/v2.0")
}

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
