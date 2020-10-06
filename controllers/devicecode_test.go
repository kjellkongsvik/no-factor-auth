package controllers

import (
	"bytes"
	"net/http"
	"mime/multipart"
	"net/http/httptest"
	"testing"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDeviceCode(t *testing.T) {
	body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)
	writer.WriteField("client_id", "123")
	myScope := "my_scope"
    writer.WriteField("scope", myScope)
	writer.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/devicecode", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, DeviceCode(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
