package controllers

import (
	"Backend/metrics"
	"Backend/models"
	"Backend/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
)

type WeatherController struct {
	weatherService *services.WeatherService
	jobService     *services.JobsService
	usersService   *services.UsersService
}

func NewWeatherController(weatherService *services.WeatherService, jobService *services.JobsService,
	userService *services.UsersService) *WeatherController {

	return &WeatherController{
		weatherService: weatherService,
		jobService:     jobService,
		usersService:   userService,
	}
}

func (rc WeatherController) RequestWeather(ctx *gin.Context) {
	//accessToken := ctx.Request.Header.Get("Token")
	//auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	//if responseErr != nil {
	//	ctx.JSON(responseErr.Status, responseErr)
	//	return
	//}
	//
	//if !auth {
	//	ctx.Status(http.StatusUnauthorized)
	//	return
	//}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var weather models.WeatherInputDTO
	var weatherString models.WeatherInputDTOString

	err = json.Unmarshal(body, &weather)
	if err != nil {
		err = json.Unmarshal(body, &weatherString)
		//read response body

		weatherResponse, errWeather := rc.weatherService.RequestWeather(weatherString.Lat, weatherString.Lon, weather.JobID)
		if errWeather != nil {
			log.Fatalf("Error while reading weather response body STRING TYPE %v", errWeather)
			return
		}
		ctx.JSON(http.StatusOK, weatherResponse)
	}

	if weather.Lat == 0 || weather.Lon == 0 {
		err = json.Unmarshal(body, &weatherString)
		//read response body

		weatherResponse, errWeather := rc.weatherService.RequestWeather(weatherString.Lat, weatherString.Lon, weather.JobID)
		if errWeather != nil {
			log.Fatalf("Error while reading weather response body STRING TYPE %v", errWeather)
			return
		}
		ctx.JSON(http.StatusOK, weatherResponse)
	} else {
		var latitude = strconv.FormatFloat(weather.Lat, 'f', 6, 64)
		var longitude = strconv.FormatFloat(weather.Lon, 'f', 6, 64)

		log.Printf("Latitude: %v", latitude)
		log.Printf("Longitude: %v", longitude)

		//read response body

		weatherResponse, errWeather := rc.weatherService.RequestWeather(latitude, longitude, weather.JobID)
		if errWeather != nil {
			log.Fatalf("Error while reading weather response body FLOAT TYPE %v", errWeather)
			return
		}
		ctx.JSON(http.StatusOK, weatherResponse)
	}
}

func (rc WeatherController) GetWeather(ctx *gin.Context) {
	metrics.HttpRequestsCounter.Inc()

	//accessToken := ctx.Request.Header.Get("Token")
	//auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_JOB})
	//if responseErr != nil {
	//	metrics.GetJobHttpResponsesCounter.WithLabelValues(
	//		strconv.Itoa(responseErr.Status)).Inc()
	//	ctx.JSON(responseErr.Status, responseErr)
	//	return
	//}
	//
	//if !auth {
	//	metrics.GetJobHttpResponsesCounter.WithLabelValues("401").Inc()
	//	ctx.Status(http.StatusUnauthorized)
	//	return
	//}

	jobId := ctx.Param("id")

	response, weatherResponseErr := rc.weatherService.GetJobWithWeather(jobId)
	if weatherResponseErr != nil {
		metrics.GetJobHttpResponsesCounter.WithLabelValues(
			strconv.Itoa(weatherResponseErr.Status)).Inc()
		ctx.JSON(weatherResponseErr.Status, weatherResponseErr)
		return
	}

	metrics.GetJobHttpResponsesCounter.WithLabelValues("200").Inc()
	ctx.JSON(http.StatusOK, response)
}

func (rc WeatherController) DeleteWeather(ctx *gin.Context) {
	//accessToken := ctx.Request.Header.Get("Token")
	//auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	//if responseErr != nil {
	//	ctx.JSON(responseErr.Status, responseErr)
	//	return
	//}
	//
	//if !auth {
	//	ctx.Status(http.StatusUnauthorized)
	//	return
	//}

	resultId := ctx.Param("id")

	responseErr := rc.weatherService.DeleteWeather(resultId)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}
