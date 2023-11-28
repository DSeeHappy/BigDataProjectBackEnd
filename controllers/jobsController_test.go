package controllers

import (
	"Backend/models"
	"Backend/repositories"
	"Backend/services"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetJobsResponse(t *testing.T) {
	dbHandler, mock, _ := sqlmock.New()
	defer dbHandler.Close()

	columnsUsers := []string{"user_role"}
	mock.ExpectQuery("SELECT user_role").WillReturnRows(
		sqlmock.NewRows(columnsUsers).AddRow("runner"),
	)

	columns := []string{"id", "first_name", "last_name", "age", "is_active", "country", "personal_best", "season_best"}
	mock.ExpectQuery("SELECT *").WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("1", "John", "Smith", 30, true, "United States", "02:00:41", "02:13:13").
			AddRow("2", "Marijana", "Komatinovic", 30, true, "Serbia", "01:18:28", "01:18:28"))

	router := initTestRouter(dbHandler)
	request, _ := http.NewRequest("GET", "/runner", nil)
	request.Header.Set("token", "token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

	var jobs []*models.Job
	json.Unmarshal(recorder.Body.Bytes(), &jobs)

	assert.NotEmpty(t, jobs)
	assert.Equal(t, 2, len(jobs))
}

func initTestRouter(dbHandler *sql.DB) *gin.Engine {
	jobsRepository := repositories.NewJobsRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)
	jobsService := services.NewJobsService(jobsRepository, nil)
	usersServices := services.NewUsersService(usersRepository)
	jobsController := NewJobsController(jobsService, usersServices)

	router := gin.Default()

	router.GET("/runner", jobsController.GetJobsBatch)

	return router
}
