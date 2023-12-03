package controllers

import (
	"Backend/models"
	"Backend/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
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

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var weather models.WeatherInputDTO
	err = json.Unmarshal(body, &weather)

	//read response body

	weatherResponse, errWeather := rc.weatherService.RequestWeather(weather.Lat, weather.Lon, weather.JobID)
	if errWeather != nil {
		log.Fatalf("Error while reading weather response body %v", errWeather)
		return
	}

	ctx.JSON(http.StatusOK, weatherResponse)
	//ctx.Status(http.StatusNoContent)

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
