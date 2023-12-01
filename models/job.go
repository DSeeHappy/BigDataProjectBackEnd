package models

type Job struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	CompanyID     string     `json:"company_id,default:1"`
	Address       string     `json:"address,default:123 Main St"`
	City          string     `json:"city,default:San Francisco"`
	State         string     `json:"state,default:CA"`
	ZipCode       string     `json:"zipCode,default:94105"`
	Country       string     `json:"country,default:USA"`
	Latitude      string     `json:"latitude,default:37.7917"`
	Longitude     string     `json:"longitude,default:-122.3933"`
	ScheduledDate string     `json:"scheduled_date,omitempty"`
	Scheduled     bool       `json:"scheduled"`
	IsActive      bool       `json:"is_active"`
	Weathers      []*Weather `json:"weather,omitempty"`
}
