package models

import "log"

type Job struct {
	ID            string     `json:"id" form:"id"`
	Name          string     `json:"name" form:"name"`
	CompanyID     string     `json:"company_id" form:"company_id"`
	Address       string     `json:"address" form:"address"`
	City          string     `json:"city" form:"city"`
	State         string     `json:"state" form:"state"`
	ZipCode       string     `json:"zip_code" form:"zip_code"`
	Country       string     `json:"country" form:"country"`
	Latitude      string     `json:"latitude" form:"latitude"`
	Longitude     string     `json:"longitude" form:"longitude" `
	ScheduledDate string     `json:"scheduled_date,omitempty" form:"scheduled_date"`
	Scheduled     bool       `json:"scheduled" form:"scheduled"`
	IsActive      bool       `json:"is_active" form:"is_active"`
	Weathers      []*Weather `json:"weather,omitempty" form:"weather,omitempty"`
}

type JobUpdate struct {
	ID            string     `json:"id" form:"id"`
	Name          *string    `json:"name" form:"name"`
	CompanyID     *string    `json:"company_id" form:"company_id"`
	Address       *string    `json:"address" form:"address"`
	City          *string    `json:"city" form:"city"`
	State         *string    `json:"state" form:"state"`
	ZipCode       *string    `json:"zip_code" form:"zip_code"`
	Country       *string    `json:"country" form:"country"`
	Latitude      *string    `json:"latitude" form:"latitude"`
	Longitude     *string    `json:"longitude" form:"longitude" `
	ScheduledDate *string    `json:"scheduled_date,omitempty" form:"scheduled_date"`
	Scheduled     *bool      `json:"scheduled" form:"scheduled"`
	IsActive      *bool      `json:"is_active" form:"is_active"`
	Weathers      []*Weather `json:"weather,omitempty" form:"weather,omitempty"`
}

func (j *Job) AddWeatherToJob(weather Weather) {
	log.Printf("Adding weather to job: %v", weather)
	j.Weathers = append(j.Weathers, &weather)
}

func (j *Job) RemoveWeatherFromJob(weather *Weather) {
	for i, w := range j.Weathers {
		if w.ID == weather.ID {
			j.Weathers = append(j.Weathers[:i], j.Weathers[i+1:]...)
		}
	}
}

func (j *Job) UpdateWeatherInJob(weather *Weather) {
	for i, w := range j.Weathers {
		if w.ID == weather.ID {
			j.Weathers[i] = weather
		}
	}
}

func (j *Job) AddWeatherListToJob(weatherList []Weather) {
	for _, w := range weatherList {
		j.AddWeatherToJob(w)
	}
}

func (j *JobUpdate) ConvertToJob() *Job {
	if j.Scheduled != nil {

	}

	return &Job{
		ID:            j.ID,
		Name:          *j.Name,
		CompanyID:     *j.CompanyID,
		Address:       *j.Address,
		City:          *j.City,
		State:         *j.State,
		ZipCode:       *j.ZipCode,
		Country:       *j.Country,
		Latitude:      *j.Latitude,
		Longitude:     *j.Longitude,
		ScheduledDate: *j.ScheduledDate,
		Scheduled:     *j.Scheduled,
		IsActive:      *j.IsActive,
		Weathers:      j.Weathers,
	}
}
