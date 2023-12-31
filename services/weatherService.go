package services

import (
	"Backend/models"
	"Backend/repositories"
	"github.com/goccy/go-json"
	"io"
	"log"
	"net/http"
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

//func (rs WeatherService) RequestWeather(latLng *models.LatLng) (*[]models.Weather, *models.ResponseError) {
//	// request weather data from openweathermap.org
//
//	var lat = fmt.Sprintf("%f", latLng.Lat)
//	var lng = fmt.Sprintf("%f", latLng.Lon)
//	var url = "https://pro.openweathermap.org/data/2.5/forecast/climate?lat=" + lat + "&lon=" + lng + "&appid=fcc51394a211b5d91ede128ba9c971e5"
//
//	weather, err := http.Get(url)
//	if err != nil {
//		log.Fatalf("Error while requesting weather data: %v", err)
//		return nil, nil
//	}
//
//	// validation
//
//	weatherResponse, err := io.ReadAll(weather.Body)
//	if err != nil {
//		log.Println("Error while reading create result request body", err)
//		return nil, nil
//	}
//
//	var weatherDataResponse models.WeatherResponseDTO
//
//	err = json.Unmarshal(weatherResponse, &weatherDataResponse)
//	if err != nil {
//		log.Println("Error while unmarshalling create result request body", err)
//		return nil, nil
//	}
//	weatherData, weatherDataErr := models.MapDTOToWeatherModel(weatherDataResponse)
//	if weatherDataErr != nil {
//		log.Fatalf("Error while mapping weather data: %v", weatherDataErr)
//	}
//	return &weatherData, nil
//}

func (rs WeatherService) RequestWeather(lat, lon, jobId string) (*[]models.Weather, *models.ResponseError) {
	// request weather data from openweathermap.org
	validateLatLng(lat, lon)
	var url = "https://pro.openweathermap.org/data/2.5/forecast/climate?lat=" + lat + "&lon=" + lon + "&appid=fcc51394a211b5d91ede128ba9c971e5" + "&units=imperial"

	weather, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error while requesting weather data: %v", err)
		return nil, nil
	}

	transactionErr := repositories.BeginTransaction(rs.jobsRepository, rs.weatherRepository)
	if transactionErr != nil {
		return nil, &models.ResponseError{
			Message: "Failed to start transaction",
			Status:  http.StatusBadRequest,
		}
	}

	// validation

	weatherResponse, err := io.ReadAll(weather.Body)
	if err != nil {
		log.Println("Error while reading create result request body", err)
		return nil, nil
	}

	var weatherDataResponse models.WeatherResponseDTO

	err = json.Unmarshal(weatherResponse, &weatherDataResponse)
	if err != nil {
		log.Println("Error while unmarshalling create result request body", err)
		return nil, nil
	}

	weatherData, weatherDataErr := models.MapDTOToWeatherModel(weatherDataResponse, jobId)
	if weatherDataErr != nil {
		log.Fatalf("Error while mapping weather data: %v", weatherDataErr)
	}

	response, responseErr := rs.weatherRepository.CreateWeather(weatherData)
	if responseErr != nil {
		err := repositories.RollbackTransaction(rs.jobsRepository, rs.weatherRepository)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Rollback Error: Weather Data not saved" + err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		return nil, responseErr
	}

	err = repositories.CommitTransaction(rs.jobsRepository, rs.weatherRepository)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Commit Error: Weather Data not saved" + err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return response, nil
}

func (rs WeatherService) GetJobWithWeather(jobId string) (*models.Job, *models.ResponseError) {
	responseErr := ValidateJobId(jobId)
	if responseErr != nil {
		return nil, responseErr
	}

	job, responseErr := rs.jobsRepository.GetJob(jobId)
	if responseErr != nil {
		return nil, responseErr
	}

	weather, responseErr := rs.weatherRepository.GetAllJobsWeather(jobId)
	if responseErr != nil {
		return nil, responseErr
	}

	job.Weathers = weather

	return job, nil
}

func (rs WeatherService) DeleteWeather(weatherId string) *models.ResponseError {
	return nil
}

func validateLatLng(lat, lon string) *models.ResponseError {
	if lat != "" || lon != "" {
		return &models.ResponseError{
			Message: "Missing latitude or longitude",
			Status:  http.StatusBadRequest,
		}
	} else {
		return nil
	}
}
