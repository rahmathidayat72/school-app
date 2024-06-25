package handler

import (
	"apk-sekolah/features/mapel"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	golangmodule "github.com/rahmathidayat72/golang-module"
)

type MapelHandler struct {
	mapelService mapel.ServiceMapelInterface
}

func NewHandlerMapel(service mapel.ServiceMapelInterface) *MapelHandler {
	return &MapelHandler{
		mapelService: service,
	}
}

func (handler *MapelHandler) InsertMapel(c echo.Context) error {
	// dummyParam := c.QueryParam("dummy")
	// if dummyParam == "true" {
	// 	dummyData := GenerateDummyDataAdd()
	// 	return golangmodule.BuildResponse(dummyData, http.StatusOK, "Success get dummy data", c)
	// }

	mapelInput := new(RequestMapel)

	err := c.Bind(mapelInput)
	if err != nil {
		c.Echo().Logger.Error("Input error: ", err.Error())
		return golangmodule.BuildResponse(nil, http.StatusBadRequest, "Error input, invalid input data", c)
	}

	mapelRequest := FormatterRequest(*mapelInput)

	if err = handler.mapelService.Insert(mapelRequest); err != nil {
		if strings.Contains(err.Error(), "validation") {
			return golangmodule.BuildResponse(nil, http.StatusBadRequest, err.Error(), c)
		}
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "Error", c)
	}

	responseMapel := FormatterResponse(ResponseMapel{
		ID:     mapelRequest.ID,
		GuruID: mapelInput.GuruID,
		Mapel:  mapelRequest.Mapel,
	})

	return golangmodule.BuildResponse(responseMapel, http.StatusCreated, "Mata Pelajaran created successfully", c)
}

func (handler *MapelHandler) GetAllMapel(c echo.Context) error {
	guruIDParam := c.QueryParam("guru_id")
	searchParam := c.QueryParam("search")

	mapelCore, err := handler.getMapelByFilter(guruIDParam, searchParam)
	if err != nil {
		log.Printf("Error fetching mapel: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch mapel"})
	}

	var mapelResponse []ResponseMapel
	if len(mapelCore) > 0 {
		mapelResponse = make([]ResponseMapel, len(mapelCore))
		for i, v := range mapelCore {
			mapelResponse[i] = ResponseMapel{
				ID:       v.ID,
				GuruID:   v.GuruID,
				NamaGuru: v.NamaGuru,
				Mapel:    v.Mapel,
			}
		}
	}

	// Logging informasi sukses
	log.Printf("Successfully fetched %d mapel", len(mapelResponse))

	return golangmodule.BuildResponse(mapelResponse, http.StatusOK, "Successfully get mata pelajaran", c)

}
func (handler *MapelHandler) UpdateMapel(c echo.Context) error {
	idStr := c.QueryParam("id")

	err := handler.mapelService.Update(mapel.MapelCore{}, idStr)
	if err != nil {
		log.Printf("Error in Update MataPelajaran (mapelService.Update): %s", err)
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "error", c)
	}
	mapelUpdate := new(RequestMapel)
	err = c.Bind(&mapelUpdate)
	if err != nil {
		log.Printf("Error in Update Mata Pelajaran (mapelService.Update): %s", err)
		return golangmodule.BuildResponse(nil, http.StatusBadRequest, "error binding data", c)
	}
	updateMapel := FormatterRequest(*mapelUpdate)
	err = handler.mapelService.Update(updateMapel, idStr)
	if err != nil {
		// mengecek ada inputan sudah sesuai
		if strings.Contains(err.Error(), "validation") {
			log.Printf("Error in Update Mata Pelajaran (mapelService.Update): %s", err)
			return golangmodule.BuildResponse(nil, http.StatusBadRequest, err.Error(), c)
		}
		log.Printf("Error in Update Mata Pelajaran (mapelService.Update): %s", err)
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "error", c)
	}
	log.Printf("Successfully fetched mapel id %s ", updateMapel.ID)
	return golangmodule.BuildResponse(nil, http.StatusOK, "Mata pelajaran updated successfully", c)

}
func (handler *MapelHandler) DeleteMapel(c echo.Context) error {
	idStr := c.QueryParam("id")
	err := handler.mapelService.Delete(idStr)
	if err != nil {
		log.Printf("Error in Delete mapel (mapelService.Delete): %s", err)
		return golangmodule.BuildResponse(nil, http.StatusInternalServerError, "error", c)
	}
	log.Printf("Successfully delete mapel id %s ", idStr)
	return golangmodule.BuildResponse(nil, http.StatusOK, "Mapel delete successfully", c)

}
