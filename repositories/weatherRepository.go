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

func NewWeathersRepository(dbHAndler *sql.DB) *WeatherRepository {
	return &WeatherRepository{
		dbHandler: dbHAndler,
	}
}

func (rr WeatherRepository) CreateWeather(result *models.Weather) (*models.Weather, *models.ResponseError) {
	query := `
		INSERT INTO jobs(id, name, address, city, state, zip_code, country, latitude, longitude, scheduled_date, scheduled, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'false', 'true')
		RETURNING id`

	rows, err := rr.transaction.Query(query, result.ID)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

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

	return &models.Weather{
		ID:    resultId,
		JobID: result.JobID,
	}, nil
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
	SELECT id, job_id
	FROM weathers
	WHERE job_id = $1`

	rows, err := rr.dbHandler.Query(query, jobId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	results := make([]*models.Weather, 0)
	var id, raceWeather, location string
	var position, year int

	for rows.Next() {
		err := rows.Scan(&id, &raceWeather, &location, &position, &year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		result := &models.Weather{
			ID:    id,
			JobID: jobId,
		}

		results = append(results, result)
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return results, nil
}
