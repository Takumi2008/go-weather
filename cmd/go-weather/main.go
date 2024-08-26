package main

import (
	"net/http"

	"github.com/takumi2008/go-weather/api"
	"github.com/takumi2008/go-weather/weather"
)

func main() {
	f := weather.NewFetcher(http.DefaultClient)
	a := api.New(f)
	http.HandleFunc("/forecast", a.HandleWeatherForecast)

	// TODO: Handle error.
	_ = http.ListenAndServe(":8080", nil)
}
