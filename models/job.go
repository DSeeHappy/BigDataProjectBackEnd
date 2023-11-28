package models

import "database/sql"

type Job struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Address       sql.NullString `json:"address,omitempty"`
	City          sql.NullString `json:"city,omitempty"`
	State         string         `json:"state,omitempty"`
	ZipCode       string         `json:"zipCode,omitempty"`
	Country       sql.NullString `json:"country,omitempty"`
	Latitude      sql.NullString `json:"latitude,omitempty"`
	Longitude     sql.NullString `json:"longitude,omitempty"`
	ScheduledDate sql.NullString `json:"scheduled_date,omitempty"`
	Scheduled     sql.NullBool   `json:"scheduled"`
	IsActive      bool           `json:"is_active"`
	Weathers      []*Weather     `json:"weather,omitempty"`
}
