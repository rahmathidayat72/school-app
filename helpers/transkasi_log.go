package helpers

import (
	"apk-sekolah/config"
	"apk-sekolah/database"
	"apk-sekolah/features/user"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	golangmodule "github.com/rahmathidayat72/golang-module"
	"gorm.io/gorm"
)

type (
	TransactionLog struct {
		ID           int       `json:"id"`
		Timestamp    time.Time `json:"timestamp"`
		UserID       string    `json:"user_id"`
		Perangkat    string    `json:"perangkat"`
		ServiceName  string    `json:"service_name"`
		RequestBody  string    `gorm:"type:jsonb" json:"request_body"`
		ResponseBody string    `gorm:"type:jsonb" json:"response_body"`
		RequestParam string    `gorm:"type:jsonb" json:"request_param"`
		Result       string    `json:"result"`
		Header       string    `gorm:"type:jsonb" json:"header"`
	}

	responseCapturer struct {
		http.ResponseWriter
		body *[]byte
	}
)

func (u *TransactionLog) TableName() string {
	return "school.transaksi_logs"
}

// Fungsi untuk mendapatkan nama layanan dari endpoint
func GetServiceNameFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 3 && parts[3] != "" {
		return parts[3]
	}
	return ""
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Membuat buffer untuk menampung request body
		var requestBody []byte
		var responseBody []byte

		db, err := database.InitPostgreSQL(config.InitConfig())
		if err != nil {
			c.Echo().Logger.Error("Gagal terhubung ke database:", err.Error())
			return err
		}
		//defer CloseDatabase(db)
		authorizationHeader := c.Request().Header.Get("Authorization")
		accessToken := golangmodule.GetTokenFromAuthorizationHeader(authorizationHeader)

		var userID string

		if accessToken != "" {
			token, err := golangmodule.VerifyToken(accessToken)
			if err != nil {
				return err
			}
			fmt.Println(token)

			metaToken, err := VerifyTokenHeader(accessToken)
			if err != nil {
				return err
			}

			sess := &MetaToken{
				ID: metaToken.ID,
			}

			userID = sess.ID
		} else {
			if c.Request().Body != nil {
				requestBodyBytes, err := ioutil.ReadAll(c.Request().Body)
				if err != nil {
					return err
				}

				c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))
				requestBody = requestBodyBytes

				var requestBodyMap map[string]interface{}
				if err := json.Unmarshal(requestBodyBytes, &requestBodyMap); err != nil {
					return err
				}
				username, ok := requestBodyMap["username"].(string)
				if !ok {
					return errors.New("username tidak ditemukan dalam request body")
				}

				userID, err = GetUserIDByUsername(db, username)
				if err != nil {
					userID = ""
				}
			}
		}
		if c.Request().Body != nil {
			requestBodyBytes, err := ioutil.ReadAll(c.Request().Body)
			if err != nil {
				return err
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))
			requestBody = requestBodyBytes
		}

		var requestBodyString string
		if len(requestBody) == 0 {

			requestBodyString = "{}"
		} else {
			requestBodyString = string(requestBody)
		}

		// Menghash password jika ada dalam request body
		requestBodyString = hashSensitiveData(requestBodyString)

		responseWriter := &responseCapturer{c.Response().Writer, &responseBody}
		c.Response().Writer = responseWriter

		if err := next(c); err != nil {
			c.Error(err)
		}

		go func() {
			perangkat, err := os.Hostname()
			if err != nil {
				errResponse := golangmodule.BuildResponse(nil, http.StatusInternalServerError, "Gagal mendapatkan perangkat", c)
				c.Error(errResponse)
				return
			}
			queryParam := c.QueryParams()
			paramJSON, err := json.Marshal(queryParam)
			if err != nil {
				errResponse := golangmodule.BuildResponse(nil, http.StatusInternalServerError, "Gagal mengurai query parameter", c)
				c.Error(errResponse)
				return
			}

			// Membuat header dalam format JSON
			headerJSON, err := json.Marshal(c.Request().Header)
			if err != nil {
				errResponse := golangmodule.BuildResponse(nil, http.StatusInternalServerError, "Gagal membuat header JSON", c)
				c.Error(errResponse)
				return
			}

			// Mencatat informasi transaksi log
			log := TransactionLog{
				Timestamp: time.Now(),
				UserID:    userID,
				//User:         sess,
				Perangkat:    perangkat,
				ServiceName:  GetServiceNameFromPath(c.Path()),
				RequestBody:  requestBodyString,
				ResponseBody: string(responseBody),
				RequestParam: string(paramJSON),
				Result:       "",
				Header:       string(headerJSON),
			}

			if c.Response().Status >= 400 {
				log.Result = "Failed"
			} else {
				log.Result = "Success"
			}

			if err := db.Create(&log).Error; err != nil {
				c.Echo().Logger.Error("Gagal menyimpan log transaksi:", err.Error())
				errResponse := golangmodule.BuildResponse(nil, http.StatusInternalServerError, "Gagal menyimpan log transaksi", c)
				c.Error(errResponse)
				return
			}
		}()

		return nil
	}
}

func (r *responseCapturer) Write(b []byte) (int, error) {
	*r.body = append(*r.body, b...)
	return r.ResponseWriter.Write(b)
}

func GetUserByUsername(db *gorm.DB, username string) (*user.UserCore, error) {
	var userInstance user.UserCore
	db = db.Debug().Table(userInstance.Nama)
	if err := db.Where("username = ?", username).First(&userInstance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &userInstance, nil
}

func GetUserIDByUsername(db *gorm.DB, username string) (string, error) {
	userInstance, err := GetUserByUsername(db, username)
	if err != nil {
		return "", err
	}
	return userInstance.ID, nil
}

// Fungsi untuk menghash data sensitif seperti password dalam request body
func hashSensitiveData(body string) string {
	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &bodyMap); err != nil {
		// Gagal mengurai JSON, kembalikan body asli
		return body
	}

	// Periksa apakah ada password dalam request body
	if password, ok := bodyMap["password"].(string); ok {
		// Hash password
		hashedPassword, err := golangmodule.HashPassword(password)
		if err != nil {
			// Penanganan kesalahan jika terjadi kesalahan saat menghash password
			// Misalnya, dapat mencatat kesalahan atau mengembalikan string kosong
			// tergantung pada kebutuhan aplikasi Anda.
			// Di sini, kita hanya mencatat kesalahan dalam log dan mengembalikan string kosong.
			log.Error("Gagal menghash password:", err)
			return ""
		}
		bodyMap["password"] = hashedPassword
	}

	// Mengembalikan JSON yang diubah
	hashedJSON, _ := json.Marshal(bodyMap)
	return string(hashedJSON)
}
