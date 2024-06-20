package handler

import (
	"apk-sekolah/features/auth"
	"apk-sekolah/helpers"
	"log"
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
		log.Printf("Error binding data: %v", err)
		return golangmodule.BuildResponse(err, http.StatusBadRequest, "error bind data, data not valid", c)
	}

	log.Printf("Attempting to log in user with email: %s", inputLogin.Email)
	login, err := handler.userService.Login(inputLogin.Email, inputLogin.Password)
	if err != nil {
		if strings.Contains(err.Error(), "invalid Password") {
			log.Printf("Invalid credentials for email: %s", inputLogin.Email)
			return golangmodule.BuildResponse(nil, http.StatusUnauthorized, "Invalid credentials", c)
		}
		log.Printf("Error during login for email %s: %v", inputLogin.Email, err)
		return golangmodule.BuildResponse(nil, http.StatusNotFound, "error: email or password is wrong", c)
	}

	data := map[string]interface{}{
		"id": login.ID,
	}
	log.Printf("Generating JWT token for user ID: %d", login.ID)
	token, expTime, err := helpers.SignToken(data)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "Gagal menghasilkan token JWT", c)
	}

	var response = ResponseAuth{
		ID:         login.ID,
		Nama:       login.Nama,
		Email:      login.Email,
		Token:      token,
		Expiration: expTime,
	}
	log.Printf("Login successful for user with email: %s", inputLogin.Email)
	return golangmodule.BuildResponse(response, http.StatusOK, "success login", c)
}
