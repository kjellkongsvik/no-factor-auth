package controllers

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type codeResponse struct {
	ClientID string `json:"device_code"`
	UserCode string `json:"user_code"`
	VerificationURI url.URL `json:"verification_uri"`
	ExpiresIn int `json:"expires_in"`
	Interval int `json:"interval"`
	Message string `json:"message"`
}


type codeRequest struct {
	ClientID string `json:"device_code"`
}

func DeviceCode(c echo.Context) error {
	req := new(codeRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	// proto := "http"
	// if c.IsTLS() {
	// 	proto = "https"
	// }

	var code codeResponse
	code.ClientID = req.ClientID
	// code.UserCode = ""
	// u, _ := url.Parse(proto + "://auth:")
	// code.VerificationURI = *u
	// code.ExpiresIn = 0
	// code.Interval = 0
	// code.Message = ""
	c.JSON(http.StatusOK, code)
	return nil
}
