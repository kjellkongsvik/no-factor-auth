package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/equinor/no-factor-auth/config"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizeV2(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet,
		"http://auth:8089/common/oauth2/v2.0/authorize?client_id=ci", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	if !assert.NoError(t, AuthorizeV2(c)) {
		return
	}

	if !assert.Equal(t, http.StatusFound, rec.Code) {
		return
	}

	loc, err := rec.Result().Location()

	if !assert.NoError(t, err) {
		return
	}

	fragments := strings.Split(loc.Fragment, "&")
	var accessToken string
	prefix := "access_token="
	for _, frag := range fragments {
		if strings.HasPrefix(frag, prefix) {
			accessToken = strings.TrimPrefix(frag, prefix)
			break
		}
	}

	token, _, _ := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	claims, ok := token.Claims.(jwt.MapClaims)

	assert.True(t, ok)
	assert.Equal(t, "ci", claims["aud"])
	assert.Equal(t, "http://auth:8089/common/v2.0", claims["iss"])

}

func TestAuthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	if !assert.NoError(t, Authorize(c)) {
		return
	}

	if !assert.Equal(t, http.StatusFound, rec.Code) {
		return
	}

	loc, err := rec.Result().Location()

	if !assert.NoError(t, err) {
		return
	}

	fragments := strings.Split(loc.Fragment, "&")
	accessToken := ""
	tokenPrefix := "id_token="
	for _, frag := range fragments {
		if strings.HasPrefix(frag, tokenPrefix) {
			accessToken = strings.TrimPrefix(frag, tokenPrefix)
			break
		}
	}

	if len(accessToken) == 0 {
		t.Errorf("No access/idtoken token")
		return
	}

	p := config.PublicKey()
	t.Run("Check token", checktoken(accessToken, p))

}

func checktoken(tokenString string, pubKey interface{}) func(t *testing.T) {

	return func(t *testing.T) {

		assert.NotEqual(t, len(tokenString), 0)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return pubKey, nil
		})

		if !assert.NoError(t, err) {
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {

			assert.Error(t, fmt.Errorf("No claims"), "No claims found")
		}

		l := len(claims)
		if !assert.Equal(t, true, l > 0) {
			return
		}

		v, ok := claims["sub"]
		if !ok {

			assert.Error(t, fmt.Errorf("No subject claim found"), "No subject claim found")
		}
		sub, ok := v.(string)
		if !ok {

			t.Errorf("Subject can't be parsed")
		}
		fmt.Println(sub)
		return

	}

}
