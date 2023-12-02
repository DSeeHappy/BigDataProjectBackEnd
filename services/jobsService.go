package services

import (
	"Backend/models"
	"Backend/repositories"
	"net/http"
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

func (js JobsService) CreateJob(job *models.Job, weather []models.Weather) (*models.Job, *models.ResponseError) {
	responseErr := validateJob(job)
	if responseErr != nil {
		return nil, responseErr
	}

	job.AddWeatherListToJob(weather)

	return js.jobsRepository.CreateJob(job)
}

func (js JobsService) UpdateJob(job *models.Job) *models.ResponseError {
	responseErr := validateJobId(job.ID)
	if responseErr != nil {
		return responseErr
	}

	responseErr = validateJob(job)
	if responseErr != nil {
		return responseErr
	}

	return js.jobsRepository.UpdateJob(job)
}

func (js JobsService) DeleteJob(jobId string) *models.ResponseError {
	responseErr := validateJobId(jobId)
	if responseErr != nil {
		return responseErr
	}

	return js.jobsRepository.DeleteJob(jobId)
}

func (js JobsService) GetJob(jobId string) (*models.Job, *models.ResponseError) {
	responseErr := validateJobId(jobId)
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

func (js JobsService) GetJobsBatch(city string, state string) ([]*models.Job, *models.ResponseError) {
	if city != "" && state != "" {
		return nil, &models.ResponseError{
			Message: "Only one parameter, city or state, can be passed",
			Status:  http.StatusBadRequest,
		}
	}

	if city != "" {
		return js.jobsRepository.GetJobsByCity(city)
	}

	if state != "" {

		return js.jobsRepository.GetJobsByZipCode(state)
	}

	return js.jobsRepository.GetAllJobs()
}

func validateJob(job *models.Job) *models.ResponseError {
	if job.Name == "" {
		return &models.ResponseError{
			Message: "Invalid name",
			Status:  http.StatusBadRequest,
		}
	}

	if job.State == "" {
		return &models.ResponseError{
			Message: "Invalid state",
			Status:  http.StatusBadRequest}
	}

	if job.ZipCode == "" {
		return &models.ResponseError{
			Message: "Invalid Zip Code",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}

func validateJobId(jobId string) *models.ResponseError {
	if jobId == "" {
		return &models.ResponseError{
			Message: "Invalid job ID",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}
