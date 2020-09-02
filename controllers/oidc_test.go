package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/equinor/no-factor-auth/oidc"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHostURLV2(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, StdOidcConfigURI, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.Equal(t, hostURLV2(c), "http://example.com")
}

func TestOpenIDConfigV2(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, StdOidcConfigURI, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	var oidc oidc.OpenIDConfig
	if assert.NoError(t, OpenIDConfigV2(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &oidc))
	}
}

func TestOidc(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, StdOidcConfigURI, nil)
	rec := httptest.NewRecorder()
	authServer := "http://example.com"
	c := e.NewContext(req, rec)
	var oidc oidc.OpenIDConfig
	if assert.NoError(t, OidcConfig(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		json.Unmarshal(rec.Body.Bytes(), &oidc)
		assert.Equal(t, authServer+"/discovery/keys", oidc.JwksURI)
	}
}
func TestOidcTenant(t *testing.T) {
	e := echo.New()
	tenant := "/tenant"
	req := httptest.NewRequest(http.MethodGet, tenant+StdOidcConfigURI, nil)
	rec := httptest.NewRecorder()
	authServer := "http://example.com" + tenant
	c := e.NewContext(req, rec)
	var oidc oidc.OpenIDConfig
	if assert.NoError(t, OidcConfig(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		json.Unmarshal(rec.Body.Bytes(), &oidc)
		assert.Equal(t, authServer+"/discovery/keys", oidc.JwksURI)
	}
}
