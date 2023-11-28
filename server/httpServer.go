package server

import (
	"Backend/controllers"
	"Backend/repositories"
	"Backend/services"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

type HttpServer struct {
	config             *viper.Viper
	router             *gin.Engine
	jobsController     *controllers.JobsController
	weathersController *controllers.WeatherController
	usersController    *controllers.UsersController
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	jobsRepository := repositories.NewJobsRepository(dbHandler)
	weathersRepository := repositories.NewWeathersRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)
	jobsService := services.NewJobsService(jobsRepository, weathersRepository)
	weatherService := services.NewWeatherService(weathersRepository, jobsRepository)
	usersService := services.NewUsersService(usersRepository)
	jobsController := controllers.NewJobsController(jobsService, usersService)
	weatherController := controllers.NewWeatherController(weatherService, usersService)
	usersController := controllers.NewUsersController(usersService)

	router := gin.Default()

	router.POST("/job", jobsController.CreateJob)
	router.PUT("/job", jobsController.UpdateJob)
	router.DELETE("/job/:id", jobsController.DeleteJob)
	router.GET("/job/:id", jobsController.GetJob)
	router.GET("/job", jobsController.GetJobsBatch)

	router.POST("/weather", weatherController.CreateWeather)
	router.DELETE("/weather/:id", weatherController.DeleteWeather)

	router.POST("/login", usersController.Login)
	router.POST("/logout", usersController.Logout)

	return HttpServer{
		config:             config,
		router:             router,
		jobsController:     jobsController,
		weathersController: weatherController,
		usersController:    usersController,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.GetString("http.server_address"))
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
