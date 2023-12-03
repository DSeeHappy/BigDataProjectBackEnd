package models

import "log"

type Job struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	CompanyID     string     `json:"company_id"`
	Address       string     `json:"address"`
	City          string     `json:"city"`
	State         string     `json:"state"`
	ZipCode       string     `json:"zipCode"`
	Country       string     `json:"country"`
	Latitude      string     `json:"latitude"`
	Longitude     string     `json:"longitude""`
	ScheduledDate string     `json:"scheduled_date,omitempty"`
	Scheduled     bool       `json:"scheduled"`
	IsActive      bool       `json:"is_active"`
	Weathers      []*Weather `json:"weather,omitempty"`
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
