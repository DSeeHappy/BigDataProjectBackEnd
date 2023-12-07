package repositories

import (
	"Backend/models"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		RETURNING id, name, address, city, state, zip_code, country, latitude, longitude,company_id, scheduled_date, scheduled, is_active`

	rows, err := rr.dbHandler.Query(query, job.Name, job.Address, job.City, job.State, job.ZipCode, job.Country, job.Latitude, job.Longitude, job.CompanyID, job.ScheduledDate)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var jobId string
	var name, address, city, state, zip_code, country, latitude, longitude, company_id, scheduled_date sql.NullString
	var scheduled, is_active sql.NullBool
	for rows.Next() {
		err := rows.Scan(&jobId, &name, &address, &city, &state, &zip_code, &country, &latitude, &longitude, &company_id, &scheduled_date, &scheduled, &is_active)
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
		Name:          name.String,
		Address:       address.String,
		City:          city.String,
		State:         state.String,
		ZipCode:       zip_code.String,
		Country:       country.String,
		Latitude:      latitude.String,
		Longitude:     longitude.String,
		CompanyID:     company_id.String,
		ScheduledDate: scheduled_date.String,
		Scheduled:     scheduled.Bool,
		IsActive:      is_active.Bool,
		Weathers:      nil,
	}, nil
}

func (rr JobsRepository) UpdateJob(job *models.JobUpdate) *models.ResponseError {
	query := "UPDATE jobs SET "
	var updates []string
	var params []interface{}
	if job.Name != nil {
		updates = append(updates, "name = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.Name)
	}
	if job.Address != nil {
		updates = append(updates, "address = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.Address)
	}
	if job.City != nil {
		updates = append(updates, "city = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.City)
	}

	if job.State != nil {
		updates = append(updates, "state = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.State)
	}

	if job.ZipCode != nil {
		updates = append(updates, "zip_code = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.ZipCode)
	}

	if job.Country != nil {
		updates = append(updates, "country = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.Country)
	}

	if job.Latitude != nil {
		updates = append(updates, "latitude = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.Latitude)
	}

	if job.Longitude != nil {
		updates = append(updates, "longitude = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.Longitude)
	}

	if job.ScheduledDate != nil {
		updates = append(updates, "scheduled_date = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.ScheduledDate)
	}

	if job.Scheduled != nil {
		updates = append(updates, "scheduled = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.Scheduled)
	}

	if job.IsActive != nil {
		updates = append(updates, "is_active = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.IsActive)
	}

	if job.CompanyID != nil {
		updates = append(updates, "company_id = $"+strconv.Itoa(len(updates)+1))
		params = append(params, job.CompanyID)
	}

	if len(updates) == 0 {
		return &models.ResponseError{
			Message: "No updates provided",
			Status:  http.StatusBadRequest,
		}
	}

	query += strings.Join(updates, ", ") + " WHERE id = $" + strconv.Itoa(len(updates)+1)
	params = append(params, job.ID)
	log.Printf("Query: %s\n", query)
	log.Printf("Params: %v\n", params)

	res, err := rr.dbHandler.Exec(query, params...)
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

func (rr JobsRepository) UpdateJobWeather(job *models.Job) *models.ResponseError {
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
	query := `SELECT * FROM jobs`

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
		rowErr := rows.Scan(&id, &name, &address, &city, &state, &zipCode, &country, &latitude, &longitude, &scheduledDate, &scheduled, &isActive, &companyId)
		if rowErr != nil {
			return nil, &models.ResponseError{
				Message: rowErr.Error(),
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
