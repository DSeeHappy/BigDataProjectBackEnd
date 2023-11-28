package services

import (
	"Backend/models"
	"Backend/repositories"
	"net/http"
	"strconv"
	"time"
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

func (rs JobsService) CreateJob(job *models.Job) (*models.Job, *models.ResponseError) {
	responseErr := validateJob(job)
	if responseErr != nil {
		return nil, responseErr
	}

	return rs.jobsRepository.CreateJob(job)
}

func (rs JobsService) UpdateJob(job *models.Job) *models.ResponseError {
	responseErr := validateJobId(job.ID)
	if responseErr != nil {
		return responseErr
	}

	responseErr = validateJob(job)
	if responseErr != nil {
		return responseErr
	}

	return rs.jobsRepository.UpdateJob(job)
}

func (rs JobsService) DeleteJob(jobId string) *models.ResponseError {
	responseErr := validateJobId(jobId)
	if responseErr != nil {
		return responseErr
	}

	return rs.jobsRepository.DeleteJob(jobId)
}

func (rs JobsService) GetJob(jobId string) (*models.Job, *models.ResponseError) {
	responseErr := validateJobId(jobId)
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

func (rs JobsService) GetJobsBatch(country string, year string) ([]*models.Job, *models.ResponseError) {
	if country != "" && year != "" {
		return nil, &models.ResponseError{
			Message: "Only one parameter, country or year, can be passed",
			Status:  http.StatusBadRequest,
		}
	}

	if country != "" {
		return rs.jobsRepository.GetJobsByCity(country)
	}

	if year != "" {
		intYear, err := strconv.Atoi(year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}

		currentYear := time.Now().Year()
		if intYear < 0 || intYear > currentYear {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}

		return rs.jobsRepository.GetJobsByZipCode(intYear)
	}

	return rs.jobsRepository.GetAllJobs()
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
