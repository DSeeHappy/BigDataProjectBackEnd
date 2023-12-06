package tests

import (
	"Backend/models"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateWeather(t *testing.T) {
	dbHandler, mock, dbErr := sqlmock.New()
	if dbErr != nil {
		log.Fatalf("Error while creating mock db: %v", dbErr)
		t.Fatal(dbErr)
	}
	defer dbHandler.Close()

	columns := []string{"id", "job_id", "pressure", "humidity", "sunrise", "sunset", "speed", "deg", "clouds", "rain", "snow", "icon", "description", "main", "city_id", "city_name", "country", "time_zone", "population", "latitude", "longitude", "temp_day", "temp_min", "temp_max", "temp_night", "temp_eve", "temp_morn", "feels_like_day", "feels_like_night", "feels_like_eve", "feels_like_morn"}
	mock.ExpectQuery("INSERT INTO").WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("1", "1", "1000", "100", "1600000000", "1600000000", "10", "10", "10", "10", "10", "10d", "description", "main", "1", "Winter Park", "USA", "America/Denver", "1000", "39.87637", "-105.75664", "10", "10", "10", "10", "10", "10", "10", "10", "10", "10"))

	router := testRouter(dbHandler)
	var weather = models.Weather{
		ID:          "1",
		JobID:       "1",
		City:        models.City{},
		Temp:        models.Temp{},
		FeelsLike:   models.FeelsLike{},
		Pressure:    0,
		Humidity:    0,
		Sunrise:     0,
		Sunset:      0,
		Speed:       0,
		Deg:         0,
		Clouds:      0,
		Rain:        0,
		Snow:        0,
		Icon:        "",
		Description: "",
		Main:        "",
	}

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
		Weathers:  make([]*models.Weather, 0),
	}

	job.Weathers = append(job.Weathers, &weather)
	body, marshalErr := json.Marshal(weather)
	if marshalErr != nil {
		log.Fatalf("Error while marshaling weather: %v", marshalErr)
		t.Fatal(marshalErr)
	}
	request, err := http.NewRequest("POST", "/weathers", strings.NewReader(string(body)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	var jobResponse *models.Job
	json.Unmarshal(recorder.Body.Bytes(), &jobResponse)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, recorder.Result().StatusCode, io.ByteReader(recorder.Body))

	assert.Equal(t, job, jobResponse, recorder.Result().StatusCode, io.ByteReader(recorder.Body))
}
