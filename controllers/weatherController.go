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

func (rc WeatherController) CreateWeather(ctx *gin.Context) {
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

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var result models.Weather
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Error while unmarshaling create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := rc.weatherService.CreateWeather(&result)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
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
