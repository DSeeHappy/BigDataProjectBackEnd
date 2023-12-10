package server

import (
	"Backend/controllers"
	"Backend/repositories"
	"Backend/services"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	weatherController := controllers.NewWeatherController(weatherService, jobsService, usersService)
	usersController := controllers.NewUsersController(usersService)

	router := gin.Default()

	router.Use(corsMiddleware())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Work Weather Scheduler API",
		})
	})

	router.POST("/jobs", jobsController.CreateJob)
	router.PUT("/jobs/:id", jobsController.UpdateJob)
	router.DELETE("/jobs/:id", jobsController.DeleteJob)
	router.GET("/jobs/:id", jobsController.GetJob)
	router.GET("/jobs", jobsController.GetJobsBatch)

	router.POST("/weather", weatherController.RequestWeather)
	router.GET("/weather/:id", weatherController.GetWeather)
	router.DELETE("/weather/:id", weatherController.DeleteWeather)

	router.POST("/login", usersController.Login)
	router.POST("/logout", usersController.Logout)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return HttpServer{
		config:             config,
		router:             router,
		jobsController:     jobsController,
		weathersController: weatherController,
		usersController:    usersController,
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Token")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (hs HttpServer) Start() {
	var err error

	gin.SetMode(gin.ReleaseMode)

	err = hs.router.Run(":" + hs.config.GetString("PORT"))

	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
