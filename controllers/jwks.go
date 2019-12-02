package controllers

import (
	"encoding/base64"
	"math/big"
	"net/http"

	"github.com/equinor/no-factor-auth/config"
	"github.com/equinor/no-factor-auth/oidc"
	"github.com/labstack/echo/v4"
)

func rsaKeyset() (*oidc.JWKS, error) {

	pub := config.PublicKey()
	b64 := base64.RawURLEncoding.EncodeToString

	e := big.Int{}
	e.SetUint64(uint64(pub.E))

	keys := oidc.JWKS{
		Keys: []oidc.JWK{
			{
				Alg: "RS256",
				Kty: "RSA",
				N:   b64(pub.N.Bytes()),
				E:   b64(e.Bytes()),
				Kid: "1",
				X5T: "1",
				Use: "sig",
			}}}
	return &keys, nil
}

func hmacKeyset() (*oidc.JWKS, error) {

	keys := oidc.JWKS{
		Keys: []oidc.JWK{
			{
				Alg: "HS256",
				Kty: "oct",
				Kid: "hmac",
				Use: "sig",
				K:   string(config.HMACKey()),
			}}}
	return &keys, nil
}

// Jwks provides oidc keyset
func Jwks(c echo.Context) error {

	jwks, err := rsaKeyset()
	if err != nil {
		return err
	}
	hmacJwks, err := hmacKeyset()
	if err != nil {
		return err
	}

	jwks.Keys = append(jwks.Keys, hmacJwks.Keys...)
	return c.JSON(http.StatusOK, jwks)
}
