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
		INSERT INTO jobs(name, address, city, state, zip_code, country, latitude, longitude,company_id, scheduled_date, scheduled, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10, 'false', 'true')
		RETURNING id`

	rows, err := rr.dbHandler.Query(query, job.Name, job.Address, job.City, job.State, job.ZipCode, job.Country, job.Latitude, job.Longitude, job.CompanyID, job.ScheduledDate)
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
		ID:            jobId,
		Name:          job.Name,
		Address:       job.Address,
		City:          job.City,
		State:         job.State,
		ZipCode:       job.ZipCode,
		Country:       job.Country,
		Latitude:      job.Latitude,
		Longitude:     job.Longitude,
		CompanyID:     job.CompanyID,
		ScheduledDate: job.ScheduledDate,
		Scheduled:     job.Scheduled,
		IsActive:      job.IsActive,
		Weathers:      job.Weathers,
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
	SELECT id, name, address, city, state, zip_code, country, latitude, longitude, scheduled_date, scheduled, is_active, company_id
	FROM jobs
	ORDER BY city
	LIMIT 100`

	rows, err := rr.dbHandler.Query(query)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	jobs := make([]*models.Job, 0)
	var id, name, state, zipCode, address, country, latitude, longitude, scheduledDate, companyId sql.NullString
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
			City:          city,
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

func (rr JobsRepository) GetJobsByZipCode(zipCode string) ([]*models.Job, *models.ResponseError) {
	query := `
	SELECT id, name, address, city, state, zip_code, country, latitude, longitude, scheduled_date, scheduled, is_active, company_id
	FROM jobs
	INNER JOIN (
		SELECT job_id
		FROM weathers
		WHERE job_id = $1
		GROUP BY job_id) job
	ON jobs.id = job.job_id
	ORDER BY name 
	LIMIT 250`

	rows, err := rr.dbHandler.Query(query, zipCode)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	jobs := make([]*models.Job, 0)
	var id, name, state, city, address, country, latitude, longitude, scheduledDate, companyId sql.NullString
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
			ZipCode:       zipCode,
			Country:       country.String,
			Latitude:      latitude.String,
			Longitude:     longitude.String,
			ScheduledDate: scheduledDate.String,
			Scheduled:     scheduled.Bool,
			CompanyID:     companyId.String,
			IsActive:      isActive.Bool,
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
