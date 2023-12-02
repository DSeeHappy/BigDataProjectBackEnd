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
	weatherRepository := repositories.NewWeatherRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)
	jobsService := services.NewJobsService(jobsRepository, weatherRepository)
	weatherService := services.NewWeatherService(weatherRepository, jobsRepository)
	usersService := services.NewUsersService(usersRepository)
	jobsController := controllers.NewJobsController(jobsService, usersService, weatherService)
	weatherController := controllers.NewWeatherController(weatherService, usersService)
	usersController := controllers.NewUsersController(usersService)

	router := gin.Default()

	router.POST("/jobs", jobsController.CreateJob)
	router.PUT("/jobs", jobsController.UpdateJob)
	router.DELETE("/jobs/:id", jobsController.DeleteJob)
	router.GET("/jobs/:id", jobsController.GetJob)
	router.GET("/jobs", jobsController.GetJobsBatch)

	router.POST("/weather", weatherController.RequestWeather)
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
