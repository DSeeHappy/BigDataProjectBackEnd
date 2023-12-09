package repositories

import (
	"Backend/models"
	"database/sql"
	"net/http"
)

type WeatherRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewWeatherRepository(dbHAndler *sql.DB) *WeatherRepository {
	return &WeatherRepository{
		dbHandler: dbHAndler,
	}
}

func (rr WeatherRepository) CreateWeather(result []models.Weather) (*[]models.Weather, *models.ResponseError) {
	var weather []models.Weather
	for index, w := range result {
		query := `
		INSERT INTO weathers(job_id, pressure, humidity, sunrise, sunset, speed, deg, clouds, rain, snow, icon, description, main, latitude, longitude, city_name, city_id, country, time_zone, population, temp_day, temp_min, temp_max, temp_night, temp_eve, temp_morn, feels_like_day, feels_like_night, feels_like_eve, feels_like_morn)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30)
		RETURNING id`

		rows, err := rr.transaction.Query(query, w.JobID, w.Pressure, w.Humidity, w.Sunrise, w.Sunset, w.Speed, w.Deg, w.Clouds, w.Rain, w.Snow, w.Icon, w.Description, w.Main, w.City.LatLng.Lat, w.City.LatLng.Lon, w.City.Name, w.City.ID, w.City.Country, w.City.Timezone, w.City.Population, w.Temp.Day, w.Temp.Min, w.Temp.Max, w.Temp.Night, w.Temp.Eve, w.Temp.Morn, w.FeelsLike.Day, w.FeelsLike.Night, w.FeelsLike.Eve, w.FeelsLike.Morn)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		var resultId string
		for rows.Next() {
			err := rows.Scan(&resultId)
			if err != nil {
				return nil, &models.ResponseError{
					Message: err.Error(),
					Status:  http.StatusInternalServerError,
				}
			}
		}

		if rows.Err() != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		weather = append(weather, w)
		weather[index].ID = resultId
	}

	return &weather, nil
}

func (rr WeatherRepository) DeleteWeather(resultId string) (*models.Weather, *models.ResponseError) {
	query := `
		DELETE FROM weathers
		WHERE id = $1
		RETURNING job_id`

	rows, err := rr.transaction.Query(query, resultId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var jobId, raceWeather string
	var year int
	for rows.Next() {
		err := rows.Scan(&jobId, &raceWeather, &year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Weather{
		ID:    resultId,
		JobID: jobId,
	}, nil
}

func (rr WeatherRepository) GetAllJobsWeather(jobId string) ([]*models.Weather, *models.ResponseError) {
	query := `
	SELECT id, job_id, pressure, humidity, sunrise, sunset, speed, deg, clouds, rain, snow, icon, description, main, latitude, longitude, city_name, city_id, country, time_zone, population, temp_day, temp_min, temp_max, temp_night, temp_eve, temp_morn, feels_like_day, feels_like_night, feels_like_eve, feels_like_morn
	FROM weathers
	WHERE job_id = $1`

	var pressure, humidity, sunrise, sunset, speed, deg, clouds, rain, snow sql.NullFloat64
	var icon, description, main sql.NullString
	var latitude, longitude sql.NullFloat64
	var city_name, country sql.NullString
	var city_id sql.NullInt64
	var population, time_zone sql.NullFloat64
	var temp_day, temp_min, temp_max, temp_night, temp_eve, temp_morn sql.NullFloat64
	var feels_like_day, feels_like_night, feels_like_eve, feels_like_morn sql.NullFloat64

	rows, err := rr.dbHandler.Query(query, jobId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "GetAllJobsWeather Repo Query: " + err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	results := make([]*models.Weather, 0)
	var id string

	for rows.Next() {
		err := rows.Scan(&id, &jobId, &pressure, &humidity, &sunrise, &sunset, &speed, &deg, &clouds, &rain, &snow, &icon, &description, &main, &latitude, &longitude, &city_name, &city_id, &country, &time_zone, &population, &temp_day, &temp_min, &temp_max, &temp_night, &temp_eve, &temp_morn, &feels_like_day, &feels_like_night, &feels_like_eve, &feels_like_morn)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "GetAllJobsWeather Repo Scan: " + err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		result := &models.Weather{
			ID:    id,
			JobID: jobId,
			City: models.City{
				ID:         int(city_id.Int64),
				Name:       city_name.String,
				LatLng:     models.LatLng{Lat: float32(latitude.Float64), Lon: float32(longitude.Float64)},
				Country:    country.String,
				Timezone:   float32(time_zone.Float64),
				Population: float32(population.Float64),
			},
			Temp: models.Temp{
				Day:   float32(temp_day.Float64),
				Min:   float32(temp_min.Float64),
				Max:   float32(temp_max.Float64),
				Night: float32(temp_night.Float64),
				Eve:   float32(temp_eve.Float64),
				Morn:  float32(temp_morn.Float64),
			},
			FeelsLike: models.FeelsLike{
				Day:   float32(feels_like_day.Float64),
				Night: float32(feels_like_night.Float64),
				Eve:   float32(feels_like_eve.Float64),
				Morn:  float32(feels_like_morn.Float64),
			},
			Pressure:    float32(pressure.Float64),
			Humidity:    float32(humidity.Float64),
			Sunrise:     float32(sunrise.Float64),
			Sunset:      float32(sunset.Float64),
			Speed:       float32(speed.Float64),
			Deg:         float32(deg.Float64),
			Clouds:      float32(clouds.Float64),
			Rain:        float32(rain.Float64),
			Snow:        float32(snow.Float64),
			Icon:        icon.String,
			Description: description.String,
			Main:        main.String,
		}

		results = append(results, result)
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: "GetAllJobsWeather Repo Rows.Err: " + err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return results, nil
}
