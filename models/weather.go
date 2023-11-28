package models

type Weather struct {
	ID         string `json:"id"`
	JobID      string `json:"job_id"`
	JobWeather string `json:"job_weather"`
	Location   string `json:"location"`
	Position   int    `json:"position,omitempty"`
	Year       int    `json:"year"`
}
