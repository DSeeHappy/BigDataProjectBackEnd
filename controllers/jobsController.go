package controllers

import (
	"Backend/metrics"
	"Backend/models"
	"Backend/services"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

const ROLE_ADMIN = "admin"
const ROLE_JOB = "job"

type JobsController struct {
	jobsService    *services.JobsService
	usersService   *services.UsersService
	weatherService *services.WeatherService
}

func NewJobsController(jobsService *services.JobsService, usersService *services.UsersService, weatherService *services.WeatherService) *JobsController {
	return &JobsController{
		jobsService:    jobsService,
		usersService:   usersService,
		weatherService: weatherService,
	}
}

func (rc JobsController) CreateJob(ctx *gin.Context) {
	metrics.HttpRequestsCounter.Inc()

	//accessToken := ctx.Request.Header.Get("Token")
	//auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	//if responseErr != nil {
	//	ctx.JSON(responseErr.Status, responseErr)
	//	return
	//}

	//if !auth {
	//	ctx.Status(http.StatusUnauthorized)
	//	return
	//}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create job request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var job models.Job
	err = json.Unmarshal(body, &job)
	if err != nil {
		log.Println("Error while unmarshalling create job request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := rc.jobsService.CreateJob(&job)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (rc JobsController) UpdateJob(ctx *gin.Context) {
	metrics.HttpRequestsCounter.Inc()

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
		log.Println("Error while reading update job request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var job models.Job
	err = json.Unmarshal(body, &job)
	if err != nil {
		log.Println("Error while unmarshalling update job request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	responseErr := rc.jobsService.UpdateJob(&job)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (rc JobsController) DeleteJob(ctx *gin.Context) {
	metrics.HttpRequestsCounter.Inc()

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

	jobId := ctx.Param("id")

	responseErr := rc.jobsService.DeleteJob(jobId)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (rc JobsController) GetJob(ctx *gin.Context) {
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

	response, responseErr := rc.jobsService.GetJob(jobId)
	if responseErr != nil {
		metrics.GetJobHttpResponsesCounter.WithLabelValues(
			strconv.Itoa(responseErr.Status)).Inc()
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	metrics.GetJobHttpResponsesCounter.WithLabelValues("200").Inc()
	ctx.JSON(http.StatusOK, response)
}

func (rc JobsController) GetJobsBatch(ctx *gin.Context) {
	metrics.HttpRequestsCounter.Inc()
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		metrics.GetAllJobsTimer.Observe(f)
	}))

	defer func() {
		timer.ObserveDuration()
	}()

	//accessToken := ctx.Request.Header.Get("Token")
	//auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_JOB})
	//fmt.Println("Response error", responseErr)
	//if responseErr != nil {
	//	ctx.JSON(responseErr.Status, responseErr)
	//	return
	//}
	//
	//if !auth {
	//	ctx.Status(http.StatusUnauthorized)
	//	return
	//}

	params := ctx.Request.URL.Query()
	city := params.Get("city")
	state := params.Get("zipCode")

	response, responseErr := rc.jobsService.GetJobsBatch(city, state)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
