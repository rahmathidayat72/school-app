package handler

import (
	"apk-sekolah/features/user"
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
	searchParam := c.QueryParam("search")

	userList, err := handler.getUsersByFilter(roleParam, searchParam)
	if err != nil {
		log.Printf("Error in getUsersByFilter: %s", err)
		return golangmodule.BuildResponse(nil, http.StatusBadRequest, "Invalid input", c)
	}

	var userResponse []ResponseUser

	// Menghindari membuat slice kosong jika tidak ada hasil
	if len(userList) > 0 {
		userResponse = make([]ResponseUser, len(userList))
		for i, v := range userList {
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

	return golangmodule.BuildResponse(userResponse, http.StatusOK, "Successfully get users", c)
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
	name := c.QueryParam("nama")

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
	idStr := c.QueryParam("id")

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
	idStr := c.QueryParam("id")
	err := handler.userService.Delete(idStr)
	if err != nil {
		log.Printf("Error in Update user (userService.Delete): %s", err)
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "error", c)
	}
	log.Printf("Successfully fetched users id %s ", idStr)
	return golangmodule.BuildResponse(nil, http.StatusOK, "User delete successfully", c)
}
