package controllers

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type Code struct {
	ClientID  string `json:"device_code"`
	UserCode string `json:"user_code"`
	VerificationURI url.URL `json:"verification_uri"`
	ExpiresIn int `json:"expires_in"`
	Interval int `json:"interval"`
	Message string `json:"message"`
}

func DeviceCode(c echo.Context) error {
	proto := "http"
	if c.IsTLS() {
		proto = "https"
	}

	code := new(Code)
	code.ClientID = "id"
	code.UserCode = ""
	u, _ := url.Parse("https://auth:")
	code.VerificationURI = *u
	code.ExpiresIn = 0
	code.Interval = 0
	code.Message = ""
	c.JSON(http.StatusOK, code)
	return nil
}
