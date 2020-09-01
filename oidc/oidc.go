package oidc

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
)

// OpenIDConfig is the expected return from the wellk-known endpoint
type OpenIDConfig struct {
	Issuer                                     string   `json:"issuer"`
	AuthorizationEndpoint                      string   `json:"authorization_endpoint"`
	TokenEndpoint                              string   `json:"token_endpoint"`
	TokenEndpointAuthMethodsSupported          []string `json:"token_endpoint_auth_methods_supported"`
	TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported"`
	UserinfoEndpoint                           string   `json:"userinfo_endpoint"`
	CheckSessionIframe                         string   `json:"check_session_iframe"`
	EndSessionEndpoint                         string   `json:"end_session_endpoint"`
	JwksURI                                    string   `json:"jwks_uri"`
	RegistrationEndpoint                       string   `json:"registration_endpoint"`
	ScopesSupported                            []string `json:"scopes_supported"`
	ResponseTypesSupported                     []string `json:"response_types_supported"`
	AcrValuesSupported                         []string `json:"acr_values_supported"`
	SubjectTypesSupported                      []string `json:"subject_types_supported"`
	UserinfoSigningAlgValuesSupported          []string `json:"userinfo_signing_alg_values_supported"`
	UserinfoEncryptionAlgValuesSupported       []string `json:"userinfo_encryption_alg_values_supported"`
	UserinfoEncryptionEncValuesSupported       []string `json:"userinfo_encryption_enc_values_supported"`
	IDTokenSigningAlgValuesSupported           []string `json:"id_token_signing_alg_values_supported"`
	IDTokenEncryptionAlgValuesSupported        []string `json:"id_token_encryption_alg_values_supported"`
	IDTokenEncryptionEncValuesSupported        []string `json:"id_token_encryption_enc_values_supported"`
	RequestObjectSigningAlgValuesSupported     []string `json:"request_object_signing_alg_values_supported"`
	DisplayValuesSupported                     []string `json:"display_values_supported"`
	ClaimTypesSupported                        []string `json:"claim_types_supported"`
	ClaimsSupported                            []string `json:"claims_supported"`
	ClaimsParameterSupported                   bool     `json:"claims_parameter_supported"`
	ServiceDocumentation                       string   `json:"service_documentation"`
	UILocalesSupported                         []string `json:"ui_locales_supported"`
}

// JWK JSON Web Key
type JWK struct {
	Alg string   `json:"alg,omitempty"`
	Kty string   `json:"kty,omitempty"`
	Use string   `json:"use,omitempty"`
	Kid string   `json:"kid,omitempty"`
	X5T string   `json:"x5t,omitempty"`
	K   string   `json:"k,omitempty"`
	N   string   `json:"n,omitempty"`
	E   string   `json:"e,omitempty"`
	X5C []string `json:"x5c,omitempty"`
}

// JWKS keyset from openID
type JWKS struct {
	Keys []JWK `json:"keys"`
}

var configClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(url string, target interface{}) error {
	r, err := configClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// GetKey gets the authservers signing key
func GetKey(authserver string, kid string) (interface{}, error) {
	oidcConf := OpenIDConfig{}
	err := getJSON(authserver+"/.well-known/openid-configuration", &oidcConf)
	if err != nil {
		return nil, err
	}

	jwksURI := oidcConf.JwksURI

	jwks := JWKS{}
	err = getJSON(jwksURI, &jwks)
	if err != nil {
		return nil, err
	}

	if len(jwks.Keys) == 0 {
		return nil, fmt.Errorf("No keys in key set: %s", jwksURI)
	}
	fromB64 := base64.RawURLEncoding.DecodeString
	jwk := jwks.Keys[0]
	if jwk.Kty == "RSA" {

		b, err := fromB64(jwk.E)
		if err != nil {
			return nil, err
		}
		e := big.Int{}
		e.SetBytes(b)

		b, err = fromB64(jwk.N)
		if err != nil {
			return nil, err
		}
		n := big.Int{}
		n.SetBytes(b)

		return &rsa.PublicKey{N: &n, E: int(e.Int64())}, nil

	}
	fmt.Println("uri", jwks.Keys[0])
	// var ks jose.JsonWebKeySet

	// resp, err := httpClient.Get(origin)
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()

	// if err := json.NewDecoder(resp.Body).Decode(&ks); err != nil {
	// 	return nil, err
	// }

	// return ks.Keys, nil
	log.Panic("foo")
	return nil, nil
}

func Default() *OpenIDConfig {
	oidc := &OpenIDConfig{
		// JwksURI:               authServer + "/" + tenantID + "/discovery/keys",
		SubjectTypesSupported:               []string{"pairwise"},
		IDTokenEncryptionAlgValuesSupported: []string{"RS256"},
		IDTokenSigningAlgValuesSupported:    []string{"RS256", "HS256"},
	}

	return oidc
}

func OidcV2(url string) *OpenIDConfig {
	oidc := Default()
	oidc.JwksURI = url + "/discovery/v2.0/keys"
	oidc.Issuer = url + "/v2.0"
	oidc.AuthorizationEndpoint = url + "/oauth2/v2.0/authorize"
	oidc.TokenEndpoint = url + "/oauth2/v2.0/token"
	return oidc
}
