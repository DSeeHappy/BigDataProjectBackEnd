package services

import (
	"Backend/models"
	"Backend/repositories"
	"github.com/kelvins/geocoder"
	"log"
	"net/http"
	"strconv"
)

type JobsService struct {
	jobsRepository    *repositories.JobsRepository
	weatherRepository *repositories.WeatherRepository
}

func NewJobsService(jobsRepository *repositories.JobsRepository, weatherRepository *repositories.WeatherRepository) *JobsService {
	return &JobsService{
		jobsRepository:    jobsRepository,
		weatherRepository: weatherRepository,
	}
}

func (js JobsService) CreateJob(job *models.Job) (*models.Job, *models.ResponseError) {
	responseErr := ValidateJob(job)
	if responseErr != nil {
		return nil, responseErr
	}

	if job.Latitude == "" || job.Longitude == "" {
		//var number int
		//var convErr error
		var address geocoder.Address

		geocoder.ApiKey = "AIzaSyCdZebGh7LnvVq5cINvbSlYupdykRlANw4"
		// Geocoding for lat/lon values of location
		address = geocoder.Address{
			Street:     job.Address,
			City:       job.City,
			State:      job.State,
			PostalCode: job.ZipCode,
			Country:    job.Country,
		}
		log.Printf("Address: %v", address)

		location, geocodeErr := geocoder.Geocoding(address)
		if geocodeErr != nil {
			return nil, &models.ResponseError{
				Message: "Error while geocoding address" + geocodeErr.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		job.Latitude = strconv.FormatFloat(location.Latitude, 'f', 6, 64)
		job.Longitude = strconv.FormatFloat(location.Longitude, 'f', 6, 64)

	}

	jobWithWeather, err := js.jobsRepository.CreateJob(job)
	if err != nil {
		return nil, err
	}

	return jobWithWeather, responseErr
}

func (js JobsService) UpdateJob(job *models.JobUpdate) *models.ResponseError {
	responseErr := ValidateJobId(job.ID)
	scheduled, validationErr := ValidateJobScheduledDate(*job.ScheduledDate)
	if validationErr != nil {
		return validationErr
	}
	if scheduled {
		job.Scheduled = &scheduled
	}
	if responseErr != nil {
		return responseErr
	}

	//responseErr = ValidateJobUpdate(job)
	if responseErr != nil {
		return responseErr
	}

	return js.jobsRepository.UpdateJob(job)
}

func (js JobsService) UpdateJobWeather(job *models.Job, weather []models.Weather) (*models.Job, *models.ResponseError) {
	jobWithWeather, err := js.jobsRepository.CreateJob(job)
	jobWithWeather.AddWeatherListToJob(weather)

	if err != nil {
		repositories.RollbackTransaction(js.jobsRepository, js.weatherRepository)
		return nil, err
	}

	commitErr := repositories.CommitTransaction(js.jobsRepository, js.weatherRepository)
	if commitErr != nil {
		return nil, &models.ResponseError{
			Message: "Commit Error: Job Data not saved" + commitErr.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return jobWithWeather, nil
}

func (js JobsService) DeleteJob(jobId string) *models.ResponseError {
	responseErr := ValidateJobId(jobId)
	if responseErr != nil {
		return responseErr
	}

	return js.jobsRepository.DeleteJob(jobId)
}

func (js JobsService) GetJob(jobId string) (*models.Job, *models.ResponseError) {
	responseErr := ValidateJobId(jobId)
	if responseErr != nil {
		return nil, responseErr
	}

	job, responseErr := js.jobsRepository.GetJob(jobId)
	if responseErr != nil {
		return nil, responseErr
	}

	weather, responseErr := js.weatherRepository.GetAllJobsWeather(jobId)
	if responseErr != nil {
		return nil, responseErr
	}

	job.Weathers = weather

	return job, nil
}

func (js JobsService) GetJobsBatch(city string, zipCode string) ([]*models.Job, *models.ResponseError) {
	if city != "" && zipCode != "" {
		return nil, &models.ResponseError{
			Message: "Only one parameter, city or zipCode, can be passed",
			Status:  http.StatusBadRequest,
		}
	}

	if city != "" {
		return js.jobsRepository.GetJobsByCity(city)
	}

	if zipCode != "" {

		return js.jobsRepository.GetJobsByZipCode(zipCode)
	}

	return js.jobsRepository.GetAllJobs()
}

func ValidateJob(job *models.Job) *models.ResponseError {

	if job.Name == "" {
		return &models.ResponseError{
			Message: "Invalid name: " + job.Name,
			Status:  http.StatusBadRequest,
		}
	}

	if job.State == "" {
		return &models.ResponseError{
			Message: "Invalid state: " + job.State,
			Status:  http.StatusBadRequest}
	}

	if job.ZipCode == "" {
		return &models.ResponseError{
			Message: "Invalid Zip Code: " + job.ZipCode,
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}
func ValidateJobUpdate(job *models.JobUpdate) *models.ResponseError {

	if job.Name == nil {
		return &models.ResponseError{
			Message: "Invalid name",
			Status:  http.StatusBadRequest,
		}
	}

	if job.State == nil {
		return &models.ResponseError{
			Message: "Invalid state",
			Status:  http.StatusBadRequest}
	}

	if job.ZipCode == nil {
		return &models.ResponseError{
			Message: "Invalid Zip Code",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}

func ValidateJobId(jobId string) *models.ResponseError {
	if jobId == "" {
		return &models.ResponseError{
			Message: "Invalid job ID",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}

func ValidateJobScheduledDate(jobScheduledDate string) (bool, *models.ResponseError) {
	if jobScheduledDate != "" {
		return true, nil
	}

	return false, nil
}
