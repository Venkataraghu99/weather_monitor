package handlers

import (
	"encoding/json"
	"net/http"
	"weather_monitor/clients"
	"weather_monitor/models"
	"weather_monitor/utils"
)

type WeatherMonitor interface {
	GetWeatherStatus(res http.ResponseWriter, req *http.Request)
}

type DefaultWeatherMonitor struct {
	env           *utils.Env
	weatherClient clients.OpenWeatherClient
}

func NewDefaultWeatherMonitor(env *utils.Env, weatherClient clients.OpenWeatherClient) WeatherMonitor {
	return &DefaultWeatherMonitor{env: env, weatherClient: weatherClient}
}

func (w DefaultWeatherMonitor) GetWeatherStatus(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	latitude := req.URL.Query().Get("latitude")
	longitude := req.URL.Query().Get("longitude")

	if len(latitude) == 0 || len(longitude) == 0 {
		err := models.AppError{ErrorMessage: "please provide latitude and longitude"}
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(err)
		return
	}

	weatherResponse, weatherResponseErr := w.weatherClient.GetCurrentWeather(ctx, latitude, longitude)
	if weatherResponseErr != nil {
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(weatherResponseErr)
		return
	}

	weatherCondition := ""
	for _, weather := range weatherResponse.Weather {
		weatherCondition = weather.Main
	}

	celsius := weatherResponse.Main.Temp - 273.15

	temperatureRange := ""
	if celsius < 15 {
		temperatureRange = "Cold"
	} else if celsius < 20 {
		temperatureRange = "Warm"
	} else if celsius < 30 {
		temperatureRange = "Warm to hot"
	} else {
		temperatureRange = "Feeling hot"
	}

	response := models.WeatherMonitorResponse{
		WeatherCondition: weatherCondition,
		TemperatureRange: temperatureRange,
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(res).Encode(response)
	return

}
