package models

type Weather struct {
	ID          string    `json:"id"`
	JobID       string    `json:"job_id"`
	City        City      `json:"city"`
	Temp        Temp      `json:"temp"`
	FeelsLike   FeelsLike `json:"feels_like"`
	Pressure    int       `json:"pressure"`
	Humidity    int       `json:"humidity"`
	Sunrise     int       `json:"sunrise"`
	Sunset      int       `json:"sunset"`
	Speed       int       `json:"speed"`
	Deg         int       `json:"deg"`
	Clouds      int       `json:"clouds"`
	Rain        int       `json:"rain"`
	Snow        int       `json:"snow"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	Main        string    `json:"main"`
}

type City struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	LatLng     LatLng  `json:"coord"`
	Country    string  `json:"country"`
	Timezone   float64 `json:"timezone"`
	Population float64 `json:"population"`
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Temp struct {
	Day   float64 `json:"day"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type FeelsLike struct {
	Day   float64 `json:"day"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type WeatherResponseDTO struct {
	City    City             `json:"city"`
	Code    string           `json:"code"`
	Message float64          `json:"message"`
	Cnt     int              `json:"cnt"`
	List    []WeatherListDTO `json:"list"`
}

type WeatherListDTO struct {
	Dt        float64      `json:"dt"`
	Sunrise   float64      `json:"sunrise"`
	Sunset    float64      `json:"sunset"`
	Temp      Temp         `json:"temp"`
	FeelsLike FeelsLike    `json:"feels_like"`
	Pressure  float64      `json:"pressure"`
	Humidity  float64      `json:"humidity"`
	Weather   []WeatherDTO `json:"weather"`
	Speed     float64      `json:"speed"`
	Deg       float64      `json:"deg"`
	Clouds    float64      `json:"clouds"`
	Rain      float64      `json:"rain"`
	Snow      float64      `json:"snow"`
}

type WeatherDTO struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func MapDTOToWeatherModel(WeatherResponseDTO) Weather {
	return Weather{}
}

func MapWeatherModelToDTO(weather Weather, job Job) WeatherResponseDTO {
	return WeatherResponseDTO{}
}
