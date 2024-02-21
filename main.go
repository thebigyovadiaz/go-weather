package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
	OpenWeatherULR       string `json:"OpenWeatherULR"`
}

type weatherData struct {
	Name  string `json:"name"`
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConf(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}

	return c, nil
}

func helloHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Weather with Go!\n"))
}

func weatherHandle(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	data, err := queryWeather(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func GetDataOpenWeather(url string) *http.Response {
	resp, errR := http.Get(url)

	if errR != nil {
		return &http.Response{}
	}

	return resp
}

func queryWeather(city string) (weatherData, error) {
	apiConfig, err := loadApiConf(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}

	url := apiConfig.OpenWeatherULR + apiConfig.OpenWeatherMapApiKey + "&q=" + city
	resp := GetDataOpenWeather(url)
	defer resp.Body.Close()

	var d weatherData
	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	d.Main.Kelvin -= math.Round(273.15)
	return d, nil
}

func main() {
	http.HandleFunc("/hello", helloHandle)
	http.HandleFunc("/weather/", weatherHandle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
