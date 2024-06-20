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
	// Mengambil path endpoint tanpa query parameter
	pathWithoutQuery := strings.Split(path, "?")[0]

	parts := strings.Split(pathWithoutQuery, "/")
	if len(parts) >= 3 && parts[len(parts)-2] != "" && parts[len(parts)-1] != "" {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	return ""
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var requestBody []byte
		var responseBody []byte

		db, err := database.InitPostgreSQL(config.InitConfig())
		if err != nil {
			c.Echo().Logger.Error("Gagal terhubung ke database:", err.Error())
			return err
		}

		authorizationHeader := c.Request().Header.Get("Authorization")
		accessToken := golangmodule.GetTokenFromAuthorizationHeader(authorizationHeader)

		var userID string

		if accessToken != "" {
			token, err := golangmodule.VerifyToken(accessToken)
			if err != nil {
				c.Echo().Logger.Error("Gagal memverifikasi token:", err)
				return err
			}
			fmt.Println(token)

			metaToken, err := VerifyTokenHeader(accessToken)
			if err != nil {
				c.Echo().Logger.Error("Gagal memverifikasi header token:", err)
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
					c.Echo().Logger.Error("Gagal membaca request body:", err)
					return err
				}
				c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))
				requestBody = requestBodyBytes

				var requestBodyMap map[string]interface{}
				if err := json.Unmarshal(requestBodyBytes, &requestBodyMap); err != nil {
					c.Echo().Logger.Error("Gagal mengurai request body:", err)
					return err
				}
				email, ok := requestBodyMap["email"].(string)
				if !ok {
					err := errors.New("username tidak ditemukan dalam request body")
					c.Echo().Logger.Error(err)
					return err
				}

				userID, err = GetUserIDByUsername(db, email)
				if err != nil {
					userID = ""
					c.Echo().Logger.Error("Gagal mendapatkan user ID:", err)
				}
			}
		}

		if c.Request().Body != nil {
			requestBodyBytes, err := ioutil.ReadAll(c.Request().Body)
			if err != nil {
				c.Echo().Logger.Error("Gagal membaca ulang request body:", err)
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
			c.Echo().Logger.Error("Error dari handler:", err)
			c.Error(err)
		}

		go func() {
			perangkat, err := os.Hostname()
			if err != nil {
				c.Echo().Logger.Error("Gagal mendapatkan perangkat:", err)
				return
			}
			queryParam := c.QueryParams()
			paramJSON, err := json.Marshal(queryParam)
			if err != nil {
				c.Echo().Logger.Error("Gagal mengurai query parameter:", err)
				return
			}

			headerJSON, err := json.Marshal(c.Request().Header)
			if err != nil {
				c.Echo().Logger.Error("Gagal membuat header JSON:", err)
				return
			}

			log := TransactionLog{
				Timestamp:    time.Now(),
				UserID:       userID,
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
				c.Echo().Logger.Error("Gagal menyimpan log transaksi:", err)
			}
		}()

		return nil
	}
}

func (r *responseCapturer) Write(b []byte) (int, error) {
	*r.body = append(*r.body, b...)
	return r.ResponseWriter.Write(b)
}

// Mengubah fungsi GetUserByUsername agar menggunakan raw SQL
func GetUserByUsername(db *gorm.DB, email string) (*user.UserCore, error) {
	var userInstance user.UserCore

	// Menggunakan raw SQL untuk mendapatkan user berdasarkan username
	sql := "SELECT * FROM user WHERE email = ? LIMIT 1"
	if err := db.Raw(sql, email).Scan(&userInstance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &userInstance, nil
}

// Fungsi untuk mendapatkan user ID berdasarkan username menggunakan raw SQL
func GetUserIDByUsername(db *gorm.DB, email string) (string, error) {
	userInstance, err := GetUserByUsername(db, email)
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
			log.Error("Gagal menghash password:", err)
			return ""
		}
		bodyMap["password"] = hashedPassword
	}

	// Mengembalikan JSON yang diubah
	hashedJSON, _ := json.Marshal(bodyMap)
	return string(hashedJSON)
}
