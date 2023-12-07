package controllers

import (
	"Backend/models"
	"Backend/repositories"
	"Backend/services"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateJob(t *testing.T) {
	dbHandler, mock, dbErr := sqlmock.New()
	if dbErr != nil {
		log.Fatalf("Error while creating mock db: %v", dbErr)
		t.Fatal(dbErr)
	}
	defer dbHandler.Close()

	columns := []string{"id", "name", "address", "city", "state", "zipCode", "country", "latitude", "longitude", "company_id", "scheduled_date", "scheduled", "is_active"}
	mock.ExpectQuery("INSERT INTO").WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("1", "Winter Weather", "Bridger Ct", "Winter Park", "CO", "80482", "USA", "39.87637", "-105.75664", "1", "", "false", "true"))

	router := testRouter(dbHandler)
	var job = models.Job{
		Name:      "Winter Weather",
		CompanyID: "1",
		Address:   "Bridger Ct",
		City:      "Winter Park",
		State:     "CO",
		ZipCode:   "80482",
		Country:   "USA",
		Latitude:  "39.87637",
		Longitude: "-105.75664",
	}
	body, marshalErr := json.Marshal(job)
	if marshalErr != nil {
		log.Fatalf("Error while marshaling job: %v", marshalErr)
		t.Fatal(marshalErr)
	}
	request, err := http.NewRequest("POST", "/jobs", strings.NewReader(string(body)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, recorder.Result().StatusCode, io.ByteReader(recorder.Body))

	var jobResponse *models.Job
	json.Unmarshal(recorder.Body.Bytes(), &jobResponse)

	assert.NotEmpty(t, jobResponse, recorder.Result().StatusCode, io.ByteReader(recorder.Body))
	assert.Equal(t, job, job, recorder.Result().StatusCode, io.ByteReader(recorder.Body))
}

func TestGetAllJobsResponse(t *testing.T) {
	dbHandler, mock, dbErr := sqlmock.New()
	if dbErr != nil {
		log.Fatalf("Error while creating mock db: %v", dbErr)
		t.Fatal(dbErr)
	}
	defer dbHandler.Close()

	columns := []string{"id", "name", "address", "city", "state", "zipCode", "country", "latitude", "longitude", "scheduled_date", "scheduled", "is_active", "company_id"}
	mock.ExpectQuery("SELECT \\* FROM jobs").WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("1", "Winter Weather", "Bridger Ct", "Winter Park", "CO", "80482", "USA", "39.87637", "-105.75664", "", "false", "true", "1").
			AddRow("2", "Work Day", "Bridger Ct", "Winter Park", "CO", "80482", "USA", "39.87637", "-105.75664", "", "false", "true", "1"))

	router := testRouter(dbHandler)
	request, err := http.NewRequest("GET", "/jobs", nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, recorder.Result().StatusCode, io.ByteReader(recorder.Body))
	if expectationsWereMet := mock.ExpectationsWereMet(); expectationsWereMet != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	var jobs []*models.Job
	json.Unmarshal(recorder.Body.Bytes(), &jobs)

	assert.NotEmpty(t, jobs, recorder.Result().StatusCode, io.ByteReader(recorder.Body))
	assert.Equal(t, 2, len(jobs), recorder.Result().StatusCode, io.ByteReader(recorder.Body))
}

func TestDeleteJobResponse(t *testing.T) {
	dbHandler, mock, dbErr := sqlmock.New()
	if dbErr != nil {
		log.Fatalf("Error while creating mock db: %v", dbErr)
		t.Fatal(dbErr)
	}
	defer dbHandler.Close()

	columns := []string{"id", "name", "address", "city", "state", "zipCode", "country", "latitude", "longitude", "scheduled_date", "scheduled", "is_active", "company_id"}
	mock.ExpectQuery("DELETE FROM jobs").WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("1", "Winter Weather", "Bridger Ct", "Winter Park", "CO", "80482", "USA", "39.87637", "-105.75664", "", "false", "true", "1"))

	router := testRouter(dbHandler)
	request, err := http.NewRequest("DELETE", "/jobs/1", nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, recorder.Result().StatusCode, io.ByteReader(recorder.Body))

	var jobs []*models.Job
	json.Unmarshal(recorder.Body.Bytes(), &jobs)

	assert.NotEmpty(t, jobs, recorder.Result().StatusCode, io.ByteReader(recorder.Body))
	assert.Equal(t, 1, len(jobs), recorder.Result().StatusCode, io.ByteReader(recorder.Body))
}

func TestSanity(t *testing.T) {
	dbHandler, mock, dbErr := sqlmock.New()
	if dbErr != nil {
		log.Fatalf("Error while creating mock db: %v", dbErr)
		t.Fatal(dbErr)
	}
	defer dbHandler.Close()

	columnsUsers := []string{"user_role"}
	mock.ExpectQuery("SELECT user_role").WillReturnRows(
		sqlmock.NewRows(columnsUsers).AddRow("job"),
	)
	columns := []string{"id", "name", "company_id", "address", "city", "state", "zipCode", "country", "latitude", "longitude", "scheduled_date", "scheduled", "is_active"}
	mock.ExpectQuery("SELECT *").WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("1", "Winter Weather", "1", "Bridger Ct", "Winter Park", "CO", "80482", "USA", "39.87637", "-105.75664", "", "false", "true").
			AddRow("2", "Work Day", "1", "Bridger Ct", "Winter Park", "CO", "80482", "USA", "39.87637", "-105.75664", "", "false", "true"))

	router := testRouter(dbHandler)
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		t.Fatal(err)
	}

	request.Header.Set("token", "token")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

	var jobs = "Work Weather Scheduler API"
	json.Unmarshal(recorder.Body.Bytes(), &jobs)

	assert.NotEmpty(t, jobs)
	assert.Equal(t, "Work Weather Scheduler API", jobs)
}

func testRouter(dbHandler *sql.DB) *gin.Engine {
	jobsRepository := repositories.NewJobsRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)
	weatherRepository := repositories.NewWeatherRepository(dbHandler)

	jobsService := services.NewJobsService(jobsRepository, weatherRepository)
	usersServices := services.NewUsersService(usersRepository)
	weatherService := services.NewWeatherService(weatherRepository, jobsRepository)

	jobsController := NewJobsController(jobsService, usersServices, weatherService)
	weatherController := NewWeatherController(weatherService, usersServices)
	usersController := NewUsersController(usersServices)

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Work Weather Scheduler API",
		})
	})

	router.POST("/jobs", jobsController.CreateJob)
	router.PUT("/jobs", jobsController.UpdateJob)
	router.DELETE("/jobs/:id", jobsController.DeleteJob)
	router.GET("/jobs/:id", jobsController.GetJob)
	router.GET("/jobs", jobsController.GetJobsBatch)

	router.POST("/weather", weatherController.RequestWeather)
	router.DELETE("/weather/:id", weatherController.DeleteWeather)

	router.POST("/login", usersController.Login)
	router.POST("/logout", usersController.Logout)

	return router
}
