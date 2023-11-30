package repositories

import (
	"Backend/models"
	"database/sql"
	"net/http"
)

type JobsRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewJobsRepository(dbHandler *sql.DB) *JobsRepository {
	return &JobsRepository{
		dbHandler: dbHandler,
	}
}

func (rr JobsRepository) CreateJob(job *models.Job) (*models.Job, *models.ResponseError) {

	query := `
		INSERT INTO jobs(name,company_id,state, city, zip_code, scheduled, is_active)
		VALUES ($1, $2, $3, $4,$5, 'false', 'true')
		RETURNING id`

	rows, err := rr.dbHandler.Query(query, job.Name, job.CompanyID, job.City, job.State, job.ZipCode)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var jobId string
	for rows.Next() {
		err := rows.Scan(&jobId)
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

	return &models.Job{
		ID:        jobId,
		Name:      job.Name,
		Address:   job.Address,
		City:      job.City,
		State:     job.State,
		ZipCode:   job.ZipCode,
		Country:   job.Country,
		Latitude:  job.Latitude,
		Longitude: job.Longitude,
	}, nil
}

func (rr JobsRepository) UpdateJob(job *models.Job) *models.ResponseError {
	query := `
		UPDATE jobs
		SET
		    			name = $1,
		    			address = $2,
		    			city = $3,
		    			state = $4,
		    			zip_code = $5,
		    			country = $6,
		    			latitude = $7,
		    			longitude = $8,
		    			scheduled_date = $9,
		    			scheduled = $10,
		    			is_active = $11
		WHERE id = $12`
	res, err := rr.dbHandler.Exec(query, job.Name, job.Address, job.City, job.State, job.ZipCode, job.Country, job.Latitude, job.Longitude, job.ScheduledDate, job.Scheduled, job.IsActive, job.ID)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if rowsAffected == 0 {
		return &models.ResponseError{
			Message: "Job not found",
			Status:  http.StatusNotFound,
		}
	}

	return nil
}

func (rr JobsRepository) UpdateJobResults(job *models.Job) *models.ResponseError {
	query := `
		UPDATE jobs
		SET
			name = $1,
			address = $2,
			city = $3,
			state = $4,
			zip_code = $5,
			country = $6,
			latitude = $7,
			longitude = $8,
			scheduled_date = $9,
			scheduled = $10,
			is_active = $11
		WHERE id = $12`

	res, err := rr.transaction.Exec(query, job.Name, job.Address, job.City, job.State, job.ZipCode, job.Country, job.Latitude, job.Longitude, job.ScheduledDate, job.Scheduled, job.IsActive, job.ID)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if rowsAffected == 0 {
		return &models.ResponseError{
			Message: "Job not found",
			Status:  http.StatusNotFound,
		}
	}

	return nil
}

func (rr JobsRepository) DeleteJob(jobId string) *models.ResponseError {
	query := `UPDATE jobs SET is_active = 'false' WHERE id = $1`

	res, err := rr.dbHandler.Exec(query, jobId)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if rowsAffected == 0 {
		return &models.ResponseError{
			Message: "Job not found",
			Status:  http.StatusNotFound,
		}
	}

	return nil
}

func (rr JobsRepository) GetJob(jobId string) (*models.Job, *models.ResponseError) {
	query := `
		SELECT *
		FROM jobs
		WHERE id = $1`

	rows, err := rr.dbHandler.Query(query, jobId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var id, name, state, zipCode, address, city, country, latitude, longitude, scheduledDate, companyId sql.NullString
	var scheduled, isActive sql.NullBool
	for rows.Next() {
		err := rows.Scan(&id, &name, &address, &city, &state, &zipCode, &country, &latitude, &longitude, &scheduledDate, &scheduled, &isActive, &companyId)
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

	return &models.Job{
		ID:            id.String,
		Name:          name.String,
		CompanyID:     companyId.String,
		Address:       address.String,
		City:          city.String,
		State:         state.String,
		ZipCode:       zipCode.String,
		Country:       country.String,
		Latitude:      country.String,
		Longitude:     longitude.String,
		ScheduledDate: scheduledDate.String,
		Scheduled:     scheduled.Bool,
		IsActive:      isActive.Bool,
		Weathers:      nil,
	}, nil
}

func (rr JobsRepository) GetAllJobs() ([]*models.Job, *models.ResponseError) {
	query := `
	SELECT *
	FROM jobs`

	rows, err := rr.dbHandler.Query(query)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	jobs := make([]*models.Job, 0)
	var id, name, state, zipCode, address, city, country, latitude, longitude, scheduledDate, companyId sql.NullString
	var scheduled, isActive sql.NullBool

	for rows.Next() {
		err := rows.Scan(&id, &name, &address, &city, &state, &zipCode, &country, &latitude, &longitude, &scheduledDate, &scheduled, &isActive, &companyId)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		job := &models.Job{
			ID:            id.String,
			Name:          name.String,
			Address:       address.String,
			City:          city.String,
			State:         state.String,
			ZipCode:       zipCode.String,
			Country:       country.String,
			Latitude:      latitude.String,
			Longitude:     longitude.String,
			ScheduledDate: scheduledDate.String,
			Scheduled:     scheduled.Bool,
			CompanyID:     companyId.String,
			IsActive:      isActive.Bool,
			Weathers:      nil,
		}

		jobs = append(jobs, job)
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return jobs, nil
}

func (rr JobsRepository) GetJobsByCity(city string) ([]*models.Job, *models.ResponseError) {
	query := `
	SELECT id, address, city, state, zip_code, scheduled_date, scheduled
	FROM jobs
	ORDER BY city
	LIMIT 100`

	rows, err := rr.dbHandler.Query(query, city)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	jobs := make([]*models.Job, 0)
	var id, name, state, zipCode string
	var address, country, latitude, longitude, scheduledDate string
	var scheduled bool

	for rows.Next() {
		err := rows.Scan(&id, &address, &city, &state, &country, &latitude, &longitude, &scheduledDate, &scheduled)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		job := &models.Job{
			ID:            id,
			Name:          name,
			Address:       address,
			City:          city,
			State:         state,
			ZipCode:       zipCode,
			Country:       country,
			Latitude:      latitude,
			Longitude:     longitude,
			ScheduledDate: scheduledDate,
			Scheduled:     scheduled,
		}

		jobs = append(jobs, job)
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return jobs, nil
}

func (rr JobsRepository) GetJobsByZipCode(year int) ([]*models.Job, *models.ResponseError) {
	//query := `
	//SELECT jobs.id, jobs.address, jobs.city, jobs.state, jobs.zip_code, jobs.scheduled_date, jobs.scheduled, weathers.id
	//FROM jobs
	//INNER JOIN (
	//	SELECT job_id, MIN(weathers) as race_result
	//	FROM weathers
	//	WHERE year = $1
	//	GROUP BY job_id) results
	//ON jobs.id = results.job_id
	//ORDER BY results.race_result
	//LIMIT 10`
	//
	//rows, err := rr.dbHandler.Query(query, year)
	//if err != nil {
	//	return nil, &models.ResponseError{
	//		Message: err.Error(),
	//		Status:  http.StatusInternalServerError,
	//	}
	//}
	//
	//defer rows.Close()
	//
	//jobs := make([]*models.Job, 0)
	//var id, firstName, lastName, country string
	//var personalBest, seasonBest sql.NullString
	//var age int
	//var isActive bool
	//
	//for rows.Next() {
	//	err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
	//	if err != nil {
	//		return nil, &models.ResponseError{
	//			Message: err.Error(),
	//			Status:  http.StatusInternalServerError,
	//		}
	//	}
	//
	//	job := &models.Job{
	//		ID:           id,
	//		FirstName:    firstName,
	//		LastName:     lastName,
	//		Age:          age,
	//		IsActive:     isActive,
	//		Country:      country,
	//		PersonalBest: personalBest.String,
	//		SeasonBest:   seasonBest.String,
	//	}
	//
	//	jobs = append(jobs, job)
	//}
	//
	//if rows.Err() != nil {
	//	return nil, &models.ResponseError{
	//		Message: err.Error(),
	//		Status:  http.StatusInternalServerError,
	//	}
	//}

	var jobs []*models.Job

	return jobs, nil
}
