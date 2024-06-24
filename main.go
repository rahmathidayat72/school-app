package main

import (
	"apk-sekolah/config"
	"apk-sekolah/database"
	authdata "apk-sekolah/features/auth/data"
	authhandler "apk-sekolah/features/auth/handler"
	authservice "apk-sekolah/features/auth/service"
	mapeldata "apk-sekolah/features/mapel/data"
	mapelhandler "apk-sekolah/features/mapel/handler"
	mapelservice "apk-sekolah/features/mapel/service"
	"apk-sekolah/features/user/data"
	"apk-sekolah/features/user/handler"
	"apk-sekolah/features/user/service"
	"apk-sekolah/helpers"
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
	e.Use(helpers.LoggingMiddleware) // Panggil LoggingMiddleware

	dataUser := data.NewDataUser(db)
	userService := service.NewServiceUser(dataUser)
	userHandlerAPI := handler.NewHandlerUser(userService)

	authUser := authdata.NewDataAuth(db)
	authService := authservice.NewServiceAuth(authUser)
	authHandlerAPI := authhandler.NewHandlerAuth(authService)

	datamapel := mapeldata.NewDataMApel(db)
	mapelService := mapelservice.NewServiceMapel(datamapel)
	mapelHandlerAPI := mapelhandler.NewHandlerMapel(mapelService)

	v1 := e.Group("/api/v1")
	v1.GET("/home", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"messages": "Hello, World!",
		})
	})
	v1.POST("/register", userHandlerAPI.CreatedUser)
	user := v1.Group("/user")
	user.Use(middleware.JWT([]byte(cfg.JWT_SECRET)))
	user.GET("/listuser", userHandlerAPI.GetAllUsers)
	user.GET("/detail-id/:id", userHandlerAPI.GetUsersById)
	user.GET("/detail-user", userHandlerAPI.DetailByName)
	user.POST("/update", userHandlerAPI.UpdateUser)
	user.DELETE("/delete", userHandlerAPI.DeleteUser)

	auth := v1.Group("/auth")
	auth.POST("/login", authHandlerAPI.Auth)

	mapel := v1.Group("/mapel")
	mapel.Use(middleware.JWT([]byte(cfg.JWT_SECRET)))
	mapel.POST("/insert", mapelHandlerAPI.InsertMapel)
	mapel.GET("/list-mapel", mapelHandlerAPI.GetAllMapel)
	mapel.POST("/update", mapelHandlerAPI.UpdateMapel)
	mapel.DELETE("/delete", mapelHandlerAPI.DeleteMapel)

	// Menambahkan pesan log untuk informasi port aplikasi
	port := ":8080"
	log.Printf("Server is listening on port %s", port)

	// Memastikan bahwa aplikasi berhenti jika terjadi kesalahan saat menjalankan server
	if err := e.Start(port); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
