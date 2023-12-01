package controllers

import (
	"Backend/models"
	"Backend/services"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WeatherController struct {
	weatherService *services.WeatherService
	usersService   *services.UsersService
}

func NewWeatherController(resultsService *services.WeatherService,
	userService *services.UsersService) *WeatherController {

	return &WeatherController{
		weatherService: resultsService,
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

	weather, err := http.Get("https://pro.openweathermap.org/data/2.5/forecast/climate?lat=39.63617&lon=-104.77700&appid=fcc51394a211b5d91ede128ba9c971e5")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	weatherResponse, err := io.ReadAll(weather.Body)
	if err != nil {
		log.Println("Error while reading create result request body", err)
		return
	}

	var weatherDataResponse models.WeatherResponseDTO

	err = json.Unmarshal(weatherResponse, &weatherDataResponse)
	if err != nil {
		log.Println("Error while unmarshalling create result request body", err)
		return
	}

	ctx.JSON(http.StatusOK, weatherDataResponse)

}

func (rc WeatherController) DeleteWeather(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	resultId := ctx.Param("id")

	responseErr = rc.weatherService.DeleteWeather(resultId)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}
