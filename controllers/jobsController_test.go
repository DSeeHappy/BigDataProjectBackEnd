package controllers

import (
	"Backend/models"
	"Backend/repositories"
	"Backend/services"
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJobController_CreateJob(t *testing.T) {
	dbHandler, mock, dbErr := sqlmock.New()
	if dbErr != nil {
		log.Fatalf("Error while creating mock db: %v", dbErr)
		t.Fatal(dbErr)
	}
	defer dbHandler.Close()

	columns := []string{"id", "name", "address", "city", "state", "zipCode", "country", "latitude", "longitude", "company_id", "scheduled_date", "scheduled", "is_active"}
	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO jobs(name, address, city, state, zip_code, country, latitude, longitude,company_id, scheduled_date, scheduled, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10, 'false', 'true')
		RETURNING id, name, address, city, state, zip_code, country, latitude, longitude,company_id, scheduled_date, scheduled, is_active`)).WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("59ee2dc2-963c-11ee-a48e-2f089e49b44b", "Winter Weather", "Bridger Ct", "Winter Park", "CO", "80482", "USA", "39.87637", "-105.75664", "1", "", "false", "true"))

	router := testRouter(dbHandler)
	var job = models.Job{
		Name:      "Winter Weather",
		CompanyID: "1",
		Address:   "Bridger Ct",
		City:      "Winter Park",
		State:     "CO",
		ZipCode:   "80482",
		Country:   "USA",
	}
	body, marshalErr := json.Marshal(job)
	if marshalErr != nil {
		log.Fatalf("Error while marshaling job: %v", marshalErr)
		t.Fatal(marshalErr)
	}
	request, err := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, recorder.Result().StatusCode, io.ByteReader(recorder.Body))

	var jobResponse *models.Job
	json.Unmarshal(recorder.Body.Bytes(), &jobResponse)

	assert.NotEmpty(t, jobResponse, recorder.Result().StatusCode, io.ByteReader(recorder.Body))
	assert.Equal(t, job, job, recorder.Result().StatusCode, io.ByteReader(recorder.Body))
}

func TestJobController_GetAllJobs(t *testing.T) {
	dbHandler, mock, dbErr := sqlmock.New()
	if dbErr != nil {
		log.Fatalf("Error while creating mock db: %v", dbErr)
		t.Fatal(dbErr)
	}
	defer dbHandler.Close()

	columns := []string{"id", "name", "address", "city", "state", "zipCode", "country", "latitude", "longitude", "scheduled_date", "scheduled", "is_active", "company_id"}
	mock.ExpectQuery("SELECT \\* FROM jobs").WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("59ee2dc2-963c-11ee-a48e-2f089e49b44b", "Winter Weather", "Bridger Ct", "Winter Park", "CO", "80482", "USA", "39.87637", "-105.75664", "", "false", "true", "1").
			AddRow("a4579258-9633-11ee-82fd-2b4e6228dbae", "Work Day", "Bridger Ct", "Winter Park", "CO", "80482", "USA", "39.87637", "-105.75664", "", "false", "true", "1"))

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

func TestJobsController_UpdateJob(t *testing.T) {

}

func TestJobsController_DeleteJob(t *testing.T) {
	dbHandler, mock, dbErr := sqlmock.New()
	if dbErr != nil {
		log.Fatalf("Error while creating mock db: %v", dbErr)
		t.Fatal(dbErr)
	}
	defer dbHandler.Close()

	//regexp.QuoteMeta() MOST IMPORTANT TOOL EVER
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE jobs SET is_active = 'false' WHERE id = $1`)).
		WithArgs("a4579258-9633-11ee-82fd-2b4e6228dbae")

	router := testRouter(dbHandler)
	request, err := http.NewRequest("DELETE", "/jobs/a4579258-9633-11ee-82fd-2b4e6228dbae", nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.NoError(t, mock.ExpectationsWereMet())
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
			AddRow("59ee2dc2-963c-11ee-a48e-2f089e49b44b", "Winter Weather", "1", "Bridger Ct", "Winter Park", "CO", "80482", "USA", "39.87637", "-105.75664", "", "false", "true").
			AddRow("a4579258-9633-11ee-82fd-2b4e6228dbae", "Work Day", "1", "Bridger Ct", "Winter Park", "CO", "80482", "USA", "39.87637", "-105.75664", "", "false", "true"))

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
	weatherController := NewWeatherController(weatherService, jobsService, usersServices)
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
