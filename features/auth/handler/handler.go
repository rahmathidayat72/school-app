package handler

import (
	"apk-sekolah/features/auth"
	"apk-sekolah/helpers"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	golangmodule "github.com/rahmathidayat72/golang-module"
)

type AuthHandler struct {
	userService auth.ServiceAuthInterface
}

func NewHandlerAuth(service auth.ServiceAuthInterface) *AuthHandler {
	return &AuthHandler{
		userService: service,
	}
}

func (handler *AuthHandler) Auth(c echo.Context) error {
	inputLogin := new(LoginRequest)
	err := c.Bind(inputLogin)
	if err != nil {
		return golangmodule.BuildResponse(err, http.StatusBadRequest, "error bind data, data not vali", c)
	}

	login, err := handler.userService.Login(inputLogin.Email, inputLogin.Password)
	if err != nil {
		if strings.Contains(err.Error(), "invalid Password") {
			return golangmodule.BuildResponse(nil, http.StatusUnauthorized, "Invalid credentials", c)

		}
		return golangmodule.BuildResponse(nil, http.StatusNotFound, "error: email or password is wrong", c)
	}
	data := map[string]interface{}{
		"id": login.ID,
	}
	token, expTime, err := helpers.SignToken(data)
	if err != nil {
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "Gagal menghasilkan token JWT", c)
	}
	var response = ResponseAuth{
		ID:         login.ID,
		Nama:       login.Nama,
		Email:      login.Email,
		Token:      token,
		Expiration: expTime,
	}
	return golangmodule.BuildResponse(response, http.StatusOK, "success login", c)

}
