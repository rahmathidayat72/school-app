package handler

import (
	"apk-sekolah/user"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	golangmodule "github.com/rahmathidayat72/golang-module"
)

type UserHandler struct {
	userService user.ServiceInterface
}

func NewHandlerUser(service user.ServiceInterface) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

func (handler *UserHandler) CreatedUser(c echo.Context) error {
	dummyParam := c.QueryParam("dummy")
	if dummyParam == "true" {
		dummyData := GenerateDummyDataAdd()
		return golangmodule.BuildResponse(dummyData, http.StatusOK, "Success get dummy data", c)
	}

	userInput := new(RequestUser)
	// Binding data
	if err := c.Bind(userInput); err != nil {
		c.Echo().Logger.Error("Input error: ", err.Error())
		return golangmodule.BuildResponse(nil, http.StatusBadRequest, "Error input, invalid input data", c)
	}

	// userRequest := FormattingRequest(RequestUser{
	// 	Nama:     userInput.Nama,
	// 	Email:    userInput.Email,
	// 	Password: userInput.Password,
	// 	NoHp:     userInput.NoHp,
	// 	Alamat:   userInput.Alamat,
	// 	Role:     userInput.Role,
	// })

	// Formatting request
	userRequest := FormattingRequest(*userInput)

	err := handler.userService.Insert(userRequest)
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			return golangmodule.BuildResponse(nil, http.StatusBadRequest, err.Error(), c)
		}
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "Error", c)
	}

	responseUser := FormatterResponse(ResponseUser{
		ID:      userRequest.ID,
		Nama:    userRequest.Nama,
		Email:   userRequest.Email,
		Telepon: userRequest.Telepon,
		Alamat:  userRequest.Alamat,
	})

	return golangmodule.BuildResponse(responseUser, http.StatusCreated, "User created successfully", c)
}

func (handler *UserHandler) GetAllUsers(c echo.Context) error {
	// Check if dummy parameter is set to "true"
	dummyParam := c.QueryParam("dummy")
	if dummyParam == "true" {
		dummyData := GenerateDummyDataList()
		return golangmodule.BuildResponse(dummyData, http.StatusOK, "Success get dummy data", c)
	}

	// Filter role admin atau user
	roleParam := c.QueryParam("role")
	validRoles := []string{"admin", "user"} // Sesuaikan dengan role yang diperbolehkan

	if !isValidRole(roleParam, validRoles) {
		return golangmodule.BuildResponse(nil, http.StatusBadRequest, "Invalid role parameter", c)
	}
	if roleParam == "admin" || roleParam == "user" {
		// Jika parameter role adalah "admin" atau "user", maka kita filter data sesuai peran
		var userList []user.UserCore
		result, err := handler.userService.GetByRole(&userList, roleParam)
		if err != nil {
			log.Printf("Error in GetByRole: %s", err)
			return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "Internal Server Error", c)
		}
		var userResponse []ResponseUser

		// Menghindari membuat slice kosong jika tidak ada hasil
		if len(result) > 0 {
			userResponse = make([]ResponseUser, len(result))
			for i, v := range result {
				userResponse[i] = ResponseUser{
					ID:      v.ID,
					Nama:    v.Nama,
					Email:   v.Email,
					Telepon: v.Telepon,
					Alamat:  v.Alamat,
				}
			}
		}

		return golangmodule.BuildResponse(userResponse, http.StatusOK, "Success get data by role", c)
	}

	// Fitur pencarian berdasarkan nama, email, dan alamat
	searchParam := c.QueryParam("search")
	if searchParam != "" {
		var userList []user.UserCore
		result, err := handler.userService.SearchUsers(&userList, searchParam)
		if err != nil {
			log.Printf("Error in SearchUsers: %s", err)
			return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "Internal Server Error", c)
		}

		var userResponse []ResponseUser

		// Menghindari membuat slice kosong jika tidak ada hasil
		if len(result) > 0 {
			userResponse = make([]ResponseUser, len(result))
			for i, v := range result {
				userResponse[i] = ResponseUser{
					ID:      v.ID,
					Nama:    v.Nama,
					Email:   v.Email,
					Telepon: v.Telepon,
					Alamat:  v.Alamat,
				}
			}
		}

		return golangmodule.BuildResponse(userResponse, http.StatusOK, "Success get search data", c)
	}

	// Menggunakan tipe data slice kosong daripada nil jika tidak ada hasil
	result, err := handler.userService.GetAll()
	if err != nil {
		// Logging kesalahan untuk memudahkan debug
		log.Printf("Error in GetAllUsers (userService.GetAll): %s", err)
		// Menggunakan HTTP status code konstan dari paket net/http
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "Internal Server Error", c)
	}

	var userResponse []ResponseUser

	// Menghindari membuat slice kosong jika tidak ada hasil
	if len(result) > 0 {
		userResponse = make([]ResponseUser, len(result))
		for i, v := range result {
			userResponse[i] = ResponseUser{
				ID:      v.ID,
				Nama:    v.Nama,
				Email:   v.Email,
				Telepon: v.Telepon,
				Alamat:  v.Alamat,
			}
		}
	}

	// Logging informasi sukses
	log.Printf("Successfully fetched %d users", len(userResponse))

	return golangmodule.BuildResponse(userResponse, http.StatusOK, "Successfully get all users", c)
}
func isValidRole(role string, validRoles []string) bool {
	for _, validRole := range validRoles {
		if strings.ToLower(role) == strings.ToLower(validRole) {
			return true
		}
	}
	return false
}

func (handler *UserHandler) GetUsersById(c echo.Context) error {
	dummyParam := c.QueryParam("dummy")
	if dummyParam == "true" {
		dummyData := GenerateDummyGetUsersBy()
		return golangmodule.BuildResponse(dummyData, http.StatusOK, "Success get user dummy by id data", c)
	}
	idStr := c.Param("id")

	user, err := handler.userService.SelectById(idStr)
	if err != nil {
		log.Printf("Error in GetAllUsers (userService.GetAll): %s", err)
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "error", c)
	}
	userById := ResponseUser{
		ID:      user.ID,
		Nama:    user.Nama,
		Email:   user.Email,
		Telepon: user.Telepon,
		Alamat:  user.Alamat,
	}
	log.Printf("Successfully fetched users id %s ", userById.ID)
	return golangmodule.BuildResponse(userById, http.StatusOK, "Successfully get user by id", c)

}

func (handler *UserHandler) DetailByName(c echo.Context) error {
	dummyParam := c.QueryParam("dummy")
	if dummyParam == "true" {
		dummyData := GenerateDummyGetUsersBy()
		return golangmodule.BuildResponse(dummyData, http.StatusOK, "Success get user dummy by name data", c)
	}
	name := c.Param("nama")

	user, err := handler.userService.DetailByName(name)
	if err != nil {
		log.Printf("Error in GetAllUsers (userService.GetAll): %s", err)
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "error", c)
	}
	userByName := FormatterSimpleResponse(ResponseUserSimple{
		ID:    user.ID,
		Nama:  user.Nama,
		Email: user.Email,
	})
	log.Printf("Successfully fetched users id %s ", userByName.Nama)
	return golangmodule.BuildResponse(userByName, http.StatusOK, "Successfully get user by name", c)
}

func (handler *UserHandler) UpdateUser(c echo.Context) error {
	dummyParam := c.QueryParam("dummy")
	if dummyParam == "true" {
		return golangmodule.BuildResponse(nil, http.StatusOK, "Success update dummy data", c)
	}
	idStr := c.Param("id")

	err := handler.userService.Update(user.UserCore{}, idStr)
	if err != nil {
		log.Printf("Error in Update user (userService.Update): %s", err)
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "error", c)
	}

	userUpdate := new(RequestUser)
	err = c.Bind(&userUpdate)

	if err != nil {
		log.Printf("Error in Update user (userService.Update): %s", err)
		return golangmodule.BuildResponse(nil, http.StatusBadRequest, "error binding data", c)
	}

	updateUser := FormattingRequest(*userUpdate)
	err = handler.userService.Update(updateUser, idStr)
	if err != nil {
		// mengecek ada inputan sudah sesuai
		if strings.Contains(err.Error(), "validation") {
			log.Printf("Error in Update user (userService.Update): %s", err)
			return golangmodule.BuildResponse(nil, http.StatusBadRequest, err.Error(), c)
		}
		log.Printf("Error in Update user (userService.Update): %s", err)
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "error", c)
	}
	log.Printf("Successfully fetched users id %s ", updateUser.ID)
	return golangmodule.BuildResponse(nil, http.StatusOK, "User updated successfully", c)

}

func (handler *UserHandler) DeleteUser(c echo.Context) error {
	dummyParam := c.QueryParam("dummy")
	if dummyParam == "true" {
		return golangmodule.BuildResponse(nil, http.StatusOK, "Success delete dummy data", c)
	}
	idStr := c.Param("id")
	err := handler.userService.Delete(idStr)
	if err != nil {
		log.Printf("Error in Update user (userService.Delete): %s", err)
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "error", c)
	}
	log.Printf("Successfully fetched users id %s ", idStr)
	return golangmodule.BuildResponse(nil, http.StatusOK, "User delete successfully", c)
}
