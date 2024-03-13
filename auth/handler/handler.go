package handler

import (
	auth "apk-sekolah/auth"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	golangmodule "github.com/rahmathidayat72/golang-module"
)

var ErrUserNotFound = errors.New("user not found")

type AuthHandler struct {
	authService auth.ServiceInterface
}

func NewHandlerAuth(service auth.ServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: service,
	}
}

func (handler *AuthHandler) LoginUser(c echo.Context) error {
	dummyParam := c.QueryParam("dummy")
	if dummyParam == "true" {
		dummyData := GenerateDummyDataLogin()
		return golangmodule.BuildResponse(dummyData, http.StatusOK, "Success get dummy data", c)
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	inputLogin := new(RequestAuth)
	if err := c.Bind(inputLogin); err != nil {
		return golangmodule.BuildResponse(nil, http.StatusBadRequest, "error bind data, data not valid", c)
	}
	// Authenticate user and generate token
	_, token, expTime, err := handler.authService.Login(email, password)
	if err != nil {
		return golangmodule.BuildResponse(nil, http.StatusBadRequest, err.Error(), c)
	}

	// Prepare response
	response := AuthResponse{
		Token:      token,
		Expiration: expTime,
	}

	return golangmodule.BuildResponse(response, http.StatusOK, "Login successful", c)
}
