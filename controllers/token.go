package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// TokenOKResponse ok type
type TokenOKResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"`
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
		if err != nil{
			return c.JSON(http.StatusBadRequest,
				TokenErrorResponse{Error: fmt.Sprintf("Unable to parse extra_claims: %s",err.Error())})
		}
	}

	a, err := newTokenWithClaims("anon1", c.Request().Host, clientID, "Foo", "Jane Doe", extraClaims)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TokenOKResponse{AccessToken: a, IDToken: a, TokenType: "Bearer"})
}

func ParseExtraClaims(addClaims []byte)(map[string]interface{}, error){
	var f interface{}
	var res map[string]interface{}
	var err error

	if addClaims != nil{
		err = json.Unmarshal(addClaims, &f)
	}
	if f != nil{
		res = f.(map[string]interface{})
	}

	return res, err
}
