package main

import (
	"net/http"
	"encoding/json"
)

// Reflect openweathermap forecast JSON structure
type ExternalWeather struct {
	Main struct{
		Temp float64
	}
}

// Retreives external temperature based on openweathermap forecast
func GetExternalTemperature() (float64, error) {
	var temp float64
	forecastUrl := "http://api.openweathermap.org/data/2.5/weather?appid=78ac22a76a45e175dbb87e0fb0b38bd6&units=metric&q=Vinnitsya,ua"
	resp, err := http.Get(forecastUrl)
	if err != nil {
		return temp, err
	}

	defer resp.Body.Close()
	var weather ExternalWeather
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err == nil {
		temp = weather.Main.Temp
	}
	return temp, err
}

