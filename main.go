package main

import (
	"log"
	"os"

	"github.com/equinor/no-factor-auth/controllers"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func setup(e *echo.Echo) {

	com := e.Group("/common")
	com.GET("/.well-known/openid-configuration", controllers.OidcConfig)
	com.GET("/discovery/keys", controllers.Jwks)
	com.GET("/oauth2/authorize", controllers.Authorize)
	com.GET("/oauth2/token", controllers.Token)
}

func setupV2(e *echo.Echo, tokenClaims map[string]interface{}){
	com := e.Group("/common")
	com.GET("/v2.0/.well-known/openid-configuration", controllers.OpenIDConfigV2)
	com.GET("/discovery/v2.0/keys", controllers.Jwks)
	com.GET("/oauth2/v2.0/authorize", controllers.AuthorizeV2)
	com.POST("/oauth2/v2.0/token", controllers.TokenV2(tokenClaims))
}

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	tokenClaims := map[string]interface{}{
                "iss": os.Getenv("TOKEN_ENDPOINT_ISSUER"),
                "sub": os.Getenv("TOKEN_ENDPOINT_SUBJECT"),
                "aud": os.Getenv("TOKEN_ENDPOINT_AUDIENCE"),
        }

	setup(e)
	setupV2(e, tokenClaims)

	e.Logger.Fatal(e.Start("0.0.0.0:8089"))
}
