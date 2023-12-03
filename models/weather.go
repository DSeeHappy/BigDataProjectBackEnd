package models

import (
	"log"
)

type Weather struct {
	ID          string    `json:"id"`
	JobID       string    `json:"job_id"`
	City        City      `json:"city"`
	Temp        Temp      `json:"temp"`
	FeelsLike   FeelsLike `json:"feels_like"`
	Pressure    float32   `json:"pressure"`
	Humidity    float32   `json:"humidity"`
	Sunrise     float32   `json:"sunrise"`
	Sunset      float32   `json:"sunset"`
	Speed       float32   `json:"speed"`
	Deg         float32   `json:"deg"`
	Clouds      float32   `json:"clouds"`
	Rain        float32   `json:"rain"`
	Snow        float32   `json:"snow"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	Main        string    `json:"main"`
}

type WeatherInputDTO struct {
	Lat   string `json:"lat"`
	Lon   string `json:"lon"`
	JobID string `json:"job_id"`
}

type City struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	LatLng     LatLng  `json:"coord"`
	Country    string  `json:"country"`
	Timezone   float32 `json:"timezone"`
	Population float32 `json:"population"`
}

type LatLng struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type Temp struct {
	Day   float32 `json:"day"`
	Min   float32 `json:"min"`
	Max   float32 `json:"max"`
	Night float32 `json:"night"`
	Eve   float32 `json:"eve"`
	Morn  float32 `json:"morn"`
}

type FeelsLike struct {
	Day   float32 `json:"day"`
	Night float32 `json:"night"`
	Eve   float32 `json:"eve"`
	Morn  float32 `json:"morn"`
}

type WeatherResponseDTO struct {
	City    City             `json:"city"`
	Code    string           `json:"code"`
	Message float32          `json:"message"`
	Cnt     int              `json:"cnt"`
	List    []WeatherListDTO `json:"list"`
}

type WeatherListDTO struct {
	Dt        float32      `json:"dt"`
	Sunrise   float32      `json:"sunrise"`
	Sunset    float32      `json:"sunset"`
	Temp      Temp         `json:"temp"`
	FeelsLike FeelsLike    `json:"feels_like"`
	Pressure  float32      `json:"pressure"`
	Humidity  float32      `json:"humidity"`
	Weather   []WeatherDTO `json:"weather"`
	Speed     float32      `json:"speed"`
	Deg       float32      `json:"deg"`
	Clouds    float32      `json:"clouds"`
	Rain      float32      `json:"rain"`
	Snow      float32      `json:"snow"`
}

type WeatherDTO struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func MapDTOToWeatherModel(w WeatherResponseDTO, jobId string) ([]Weather, error) {
	var list []Weather
	if w.List != nil || len(w.List) > 0 && len(w.List) > 1 {
		for index := range w.List {
			//fmt.Printf("Weather: %v\n", weather.Weather)
			//fmt.Printf("Temp: %v\n", weather.Temp)
			//fmt.Printf("Sunrise: %.4f\n", weather.Sunrise)
			//fmt.Printf("Sunset: %.4f\n", weather.Sunset)
			//fmt.Printf("index: %v\n", index)
			list = append(list, Weather{
				JobID:       jobId,
				Main:        w.List[index].Weather[0].Main,
				Description: w.List[index].Weather[0].Description,
				Icon:        w.List[index].Weather[0].Icon,
				Pressure:    w.List[index].Pressure,
				Humidity:    w.List[index].Humidity,
				Sunrise:     w.List[index].Sunrise,
				Sunset:      w.List[index].Sunset,
				Speed:       w.List[index].Speed,
				Deg:         w.List[index].Deg,
				Clouds:      w.List[index].Clouds,
				Rain:        w.List[index].Rain,
				Snow:        w.List[index].Snow,
				City: City{
					ID:         w.City.ID,
					Name:       w.City.Name,
					LatLng:     w.City.LatLng,
					Country:    w.City.Country,
					Timezone:   w.City.Timezone,
					Population: w.City.Population,
				},

				Temp: Temp{
					Day:   w.List[index].Temp.Day,
					Min:   w.List[index].Temp.Min,
					Max:   w.List[index].Temp.Max,
					Night: w.List[index].Temp.Night,
					Eve:   w.List[index].Temp.Eve,
					Morn:  w.List[index].Temp.Morn,
				},
				FeelsLike: FeelsLike{
					Day:   w.List[index].FeelsLike.Day,
					Night: w.List[index].FeelsLike.Night,
					Eve:   w.List[index].FeelsLike.Eve,
					Morn:  w.List[index].FeelsLike.Morn,
				},
			})
		}
	} else {
		log.Fatalf("Error while mapping weather data: %v", w)
		return nil, nil
	}

	return list, nil
}

func MapWeatherModelToDTO(weather Weather, job Job) WeatherResponseDTO {
	return WeatherResponseDTO{}
}
