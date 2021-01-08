package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/equinor/no-factor-auth/config"
	"github.com/labstack/echo/v4"
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
	Code         string `json:"code"`
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

	a, err := newTokenWithClaims("anon1", c.Request().Host, clientID, extraClaims)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TokenOKResponse{AccessToken: a, IDToken: a, TokenType: "Bearer"})
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

func newTokenV2(claims map[string]interface{}) (string, error) {
	defaultClaims := jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}

	for key, value := range claims {
		defaultClaims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, defaultClaims)

	token.Header = map[string]interface{}{
		"typ": "JWT",
		"alg": jwt.SigningMethodRS256.Name,
		"kid": "1",
	}

	tokenString, err := token.SignedString(config.PrivateKey())
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

type tokenQuery struct {
	GrantType string `query:"grant_type"`
}

func TokenV2(c echo.Context) error {
	r := new(tokenQuery)
	if err := c.Bind(r); err != nil {
		return err
	}
	claims := map[string]interface{}{
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}
	fmt.Printf("r: %v", r)
	fmt.Printf("grant_type: %v", r.GrantType)
	claims["iss"] = "http://no-factor-auth:8089/common/v2.0"
	claims["aud"] = "id"
	claims["email"] = "gpl@equinor.com"
	if r.GrantType == "authorization_code" {
		claims["iss"] = "http://no-factor-auth:8089/common/v2.0"
		claims["sub"] = "sub"

	} else if r.GrantType == "urn:ietf:params:oauth:grant-type:jwt-bearer" {
		claims["iss"] = os.Getenv("TOKEN_ENDPOINT_ISSUER")
		claims["sub"] = os.Getenv("TOKEN_ENDPOINT_SUBJECT")
		claims["aud"] = os.Getenv("TOKEN_ENDPOINT_AUDIENCE")
	}
	token, err := newTokenV2(claims)
	if err != nil {
		return err
	}
	exp := fmt.Sprintf("%v", claims["exp"])
	return c.JSON(http.StatusOK, TokenOKResponse{AccessToken: token, IDToken: token, ExpiresIn: exp})
}
