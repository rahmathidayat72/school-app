package main

import (
	"apk-sekolah/config"
	"apk-sekolah/database"
	"apk-sekolah/user/data"
	"apk-sekolah/user/handler"
	"apk-sekolah/user/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	log.Println("App is running")

	cfg := config.InitConfig()
	db, err := database.InitPostgreSQL(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %s", err)
	}

	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	dataUser := data.NewDataUser(db)
	userService := service.NewServiceUser(dataUser)
	userHandlerAPI := handler.NewHandlerUser(userService)

	v1 := e.Group("/api/v1")
	v1.GET("/home", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"messages": "Hello, World!",
		})
	})
	user := v1.Group("/user")
	user.POST("/add", userHandlerAPI.CreatedUser)
	user.GET("/list", userHandlerAPI.GetAllUsers)
	user.GET("/detail-id/:id", userHandlerAPI.GetUsersById)
	user.GET("/detail-name/:nama", userHandlerAPI.DetailByName)
	user.POST("/update/:id", userHandlerAPI.UpdateUser)
	user.DELETE("/delete/:id", userHandlerAPI.DeleteUser)

	// Menambahkan pesan log untuk informasi port aplikasi
	port := ":8080"
	log.Printf("Server is listening on port %s", port)

	// Memastikan bahwa aplikasi berhenti jika terjadi kesalahan saat menjalankan server
	if err := e.Start(port); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
