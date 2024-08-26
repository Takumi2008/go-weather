package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/takumi2008/go-weather/weather"
)

const (
	fieldLatitude  = "latitude"
	fieldLongitude = "longitude"
)

type T struct {
	weatherFetcher weather.Fetcher
}

func New(f weather.Fetcher) T {
	return T{
		weatherFetcher: f,
	}
}

func (t *T) HandleWeatherForecast(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	latitude, err := getFloatField(fieldLatitude, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	longitude, err := getFloatField(fieldLongitude, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	f, err := t.weatherFetcher.Forecast(latitude, longitude)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	// TODO: Handle error.
	_ = json.NewEncoder(w).Encode(f)
}

func getFloatField(fieldName string, req *http.Request) (float64, error) {
	value, ok := req.Form[fieldName]
	if !ok {
		return 0, fmt.Errorf("field %s not found", fieldName)
	}

	parsedValue, err := strconv.ParseFloat(value[0], 64)
	if err != nil {
		return 0, err
	}
	return parsedValue, nil

}
