package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/takumi2008/go-weather/api"
	"github.com/takumi2008/go-weather/weather"
)

func TestMissingLatitude(t *testing.T) {
	req := &http.Request{}
	req.Form = make(map[string][]string)
	req.Form.Set("longitude", "22.22")
	a := api.New(mockFetcher{})

	w := httptest.NewRecorder()
	a.HandleWeatherForecast(w, req)
	w.Flush()

	assert.Equal(t, 400, w.Code)
	require.Equal(t, "field latitude not found\n", w.Body.String())
}

func TestMissingLongitude(t *testing.T) {
	req := &http.Request{}
	req.Form = make(map[string][]string)
	req.Form.Set("latitude", "22.22")
	a := api.New(mockFetcher{})

	w := httptest.NewRecorder()
	a.HandleWeatherForecast(w, req)
	w.Flush()

	assert.Equal(t, 400, w.Code)
	require.Equal(t, "field longitude not found\n", w.Body.String())
}

func TestInvalidLatitude(t *testing.T) {
	req := &http.Request{}
	req.Form = make(map[string][]string)
	req.Form.Set("latitude", "invalid")
	req.Form.Set("longitude", "22.22")
	a := api.New(mockFetcher{})

	w := httptest.NewRecorder()
	a.HandleWeatherForecast(w, req)
	w.Flush()

	assert.Equal(t, 400, w.Code)
	require.Equal(t, "strconv.ParseFloat: parsing \"invalid\": invalid syntax\n", w.Body.String())

}

func TestInvalidLongitude(t *testing.T) {
	req := &http.Request{}
	req.Form = make(map[string][]string)
	req.Form.Set("latitude", "22.22")
	req.Form.Set("longitude", "invalid")
	a := api.New(mockFetcher{})

	w := httptest.NewRecorder()
	a.HandleWeatherForecast(w, req)
	w.Flush()

	assert.Equal(t, 400, w.Code)
	require.Equal(t, "strconv.ParseFloat: parsing \"invalid\": invalid syntax\n", w.Body.String())

}

func TestHappyPath(t *testing.T) {
	req := &http.Request{}
	req.Form = make(map[string][]string)
	req.Form.Set("latitude", "11.11")
	req.Form.Set("longitude", "22.22")
	a := api.New(mockFetcher{})

	w := httptest.NewRecorder()
	a.HandleWeatherForecast(w, req)
	w.Flush()

	assert.Equal(t, 200, w.Code)
	require.NotEmpty(t, w.Body.String())

	var got weather.Forecast
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, mockForecast.TemperatureFeels, got.TemperatureFeels)
	assert.Equal(t, mockForecast.Temperature, got.Temperature)
	assert.Equal(t, mockForecast.ShortTerm, got.ShortTerm)
}

type mockFetcher struct{}

func (m mockFetcher) Forecast(latitude, longitude float64) (weather.Forecast, error) {
	return mockForecast, nil
}

var mockForecast = weather.Forecast{
	ShortTerm:        "Cloudy",
	Temperature:      "35F",
	TemperatureFeels: "Cold",
}
