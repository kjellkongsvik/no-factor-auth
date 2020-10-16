package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// TokenOKResponse ok type
type TokenOKResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int `json:"expires_in"`
	ExpiresOn    int `json:"expires_on"`
	Resource     string `json:"resource"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

// TokenErrorResponse error type
type TokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCodes       []int  `json:"error_codes"`
	Timestamp        string `json:"timestamp"`
	TraceID          string `json:"trace_id"`
	CorrelationID    string `json:"correlation_id"`
}

// Token provides id_token and access_token to anyone who asks
func Token(c echo.Context) error {
	redirectURI := c.QueryParam("redirect_uri")
	if redirectURI == "" {
		return c.JSON(http.StatusBadRequest, TokenErrorResponse{Error: "No redirect_uri"})

	}
	clientID := c.QueryParam("client_id")
	if clientID == "" {
		return c.JSON(http.StatusBadRequest, TokenErrorResponse{Error: "No client_id"})
	}
	grantType := c.QueryParam("grant_type")
	if grantType == "" {
		return c.JSON(http.StatusBadRequest, TokenErrorResponse{Error: "No grant_type"})
	}
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, TokenErrorResponse{Error: "No code"})
	}
	clientSecret := c.QueryParam("client_secret")
	if clientSecret == "" {
		return c.JSON(http.StatusBadRequest, TokenErrorResponse{Error: "No client_secret"})
	}

	extraClaimsBytes := []byte(c.QueryParam("extra_claims"))
	var extraClaims map[string]interface{}
	var err error
	if len(extraClaimsBytes) > 0 {
		extraClaims, err = ParseExtraClaims(extraClaimsBytes)
		if err != nil {
			return c.JSON(http.StatusBadRequest,
				TokenErrorResponse{Error: fmt.Sprintf("Unable to parse extra_claims: %s", err.Error())})
		}
	}

	a, err := newTokenWithClaims("anon1", c.Request().Host, clientID, "Foo", "Jane Doe", extraClaims)
	if err != nil {
		return err
	}

	var t TokenOKResponse
	t.AccessToken = a
	t.IDToken = a
	t.TokenType = "Bearer"
	t.ExpiresIn = 0
	t.ExpiresOn = 0

	return c.JSON(http.StatusOK, t)
}

func ParseExtraClaims(addClaims []byte) (map[string]interface{}, error) {
	var f interface{}
	var res map[string]interface{}
	var err error

	if addClaims != nil {
		err = json.Unmarshal(addClaims, &f)
	}
	if f != nil {
		res = f.(map[string]interface{})
	}

	return res, err
}

type req struct {
	GrantType string `query:"grant_type"`
	ClientID string `query:"client_id"`
}

func TokenV2(claims map[string]interface{}) func(c echo.Context) error {
	return func(c echo.Context) error {

		r := new(req)
		if err := c.Bind(r); err != nil {
			return err
		}
		var ccc map[string]interface{}
		log.Println(r)
		claims["aud"] = r.ClientID
		if r.GrantType == "urn:ietf:params:oauth:grant-type:jwt-bearer" {
			ccc = claims
		} else if r.GrantType == "urn:ietf:params:oauth:grant-type:device_code" {
			claims["client_id"] = r.ClientID
			ccc = claims
		} else {

		}

		a, err := newTokenV2(ccc)
		if err != nil {
			return err
		}
		var t TokenOKResponse
		t.AccessToken = a
		t.IDToken = a
		t.TokenType = "Bearer"
		t.ExpiresIn = 0
		t.ExpiresOn = 0

		log.Printf("%v\n", t)
		return c.JSON(http.StatusOK, t)

	}
}
