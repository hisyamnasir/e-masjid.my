package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Dev4w4n/e-masjid.my/api/core/config"
	"github.com/Dev4w4n/e-masjid.my/api/core/env"
	"github.com/Dev4w4n/e-masjid.my/api/tabung-api/controller"
	_ "github.com/Dev4w4n/e-masjid.my/api/tabung-api/docs"
	"github.com/Dev4w4n/e-masjid.my/api/tabung-api/helper"
	"github.com/Dev4w4n/e-masjid.my/api/tabung-api/repository"
	"github.com/Dev4w4n/e-masjid.my/api/tabung-api/router"
	"github.com/Dev4w4n/e-masjid.my/api/tabung-api/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			Tabung Service API
// @version		1.0
// @description	A Tabung service API in Go using Gin framework
func main() {

	log.Println("Starting server ...")

	// Database
	env, err := env.GetEnvironment()
	if err != nil {
		log.Fatalf("Error getting environment: %v", err)
	}

	db, err := config.DatabaseConnection(env)
	if err != nil {
		log.Fatalf("Error getting database connection: %v", err)
	}

	// Repository
	tabungRepository := repository.NewTabungRepository(db)
	tabungService := service.NewTabungService(tabungRepository)

	tabungTypeRepository := repository.NewTabungTypeRepository(db)
	tabungTypeService := service.NewTabungTypeService(tabungTypeRepository)

	kutipanRepository := repository.NewKutipanRepository(db)
	kutipanService := service.NewKutipanService(kutipanRepository)

	// Controller
	tabungController := controller.NewTabungController(tabungService)
	tabungTypeController := controller.NewTabungTypeController(tabungTypeService)
	kutipanController := controller.NewKutipanController(kutipanService)

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	config.MaxAge = 3600
	config.AllowOrigins = []string{env.AllowOrigins}
	config.AllowMethods = []string{"GET", "POST", "DELETE", "PUT"}

	// Router
	gin.SetMode(gin.ReleaseMode)
	_router := gin.Default()
	_router.Use(cors.New(config))

	// enable swagger for dev env
	isLocalEnv := os.Getenv("GO_ENV")
	if isLocalEnv == "local" || isLocalEnv == "dev" {
		_router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	var routes *gin.Engine = _router
	routes = router.NewTabungRouter(tabungController, routes, env)
	routes = router.NewTabungTypeRouter(tabungTypeController, routes, env)
	routes = router.NewKutipanRouter(kutipanController, routes, env)

	server := &http.Server{
		Addr:    ":" + env.ServerPort,
		Handler: routes,
	}

	log.Println("Server listening on port ", env.ServerPort)

	err = server.ListenAndServe()
	helper.ErrorPanic(err)

}
