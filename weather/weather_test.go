package weather_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/takumi2008/go-weather/weather"
)

func TestForecast(t *testing.T) {
	latitude := 39.7456
	longtitude := -97.0892

	f := weather.NewFetcher(&http.Client{})

	got, err := f.Forecast(latitude, longtitude)
	require.NoError(t, err)

	assert.NotEmpty(t, got.Temperature)
	assert.NotEmpty(t, got.TemperatureFeels)
	assert.NotEmpty(t, got.ShortTerm)
}
