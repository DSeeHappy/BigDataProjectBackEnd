package services

import (
	"Backend/models"
	"Backend/repositories"
)

type WeatherService struct {
	weatherRepository *repositories.WeatherRepository
	jobsRepository    *repositories.JobsRepository
}

func NewWeatherService(weathersRepository *repositories.WeatherRepository,
	runnersRepository *repositories.JobsRepository) *WeatherService {

	return &WeatherService{
		weatherRepository: weathersRepository,
		jobsRepository:    runnersRepository,
	}
}

func (rs WeatherService) CreateWeather(weather *models.Weather) (*models.Weather, *models.ResponseError) {
	// validation

	return nil, nil
}

func (rs WeatherService) DeleteWeather(weatherId string) *models.ResponseError {
	return nil
}
