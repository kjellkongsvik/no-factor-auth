package controllers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/equinor/no-factor-auth/oidc"

	"github.com/equinor/no-factor-auth/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestJwks(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	if assert.NoError(t, Jwks(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		keys := oidc.JWKS{}
		err := json.Unmarshal(rec.Body.Bytes(), &keys)
		assert.NoError(t, err)
		assert.Equal(t, keys.Keys[0].Kid, "1")

		b64 := base64.RawURLEncoding.EncodeToString
		assert.Equal(t, keys.Keys[0].N, b64(config.PublicKey().N.Bytes()))
	}
}
