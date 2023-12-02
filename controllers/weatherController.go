package controllers

import (
	"Backend/services"
	"github.com/gin-gonic/gin"
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

	//var weather http.Response

	// read response body
	//
	//weatherResponse, errWeather := rc.weatherService.RequestWeather(&weather)
	//if errWeather != nil {
	//	log.Fatalf("Error while reading weather response body %v", errWeather)
	//	return
	//}
	//
	//ctx.JSON(http.StatusOK, weatherResponse)
	ctx.Status(http.StatusNoContent)

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
