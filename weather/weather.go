package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	// Temperature  below isCold is considered cold.
	isCold = 45

	// Temperature above isHot is considered hot.
	isHot = 75

	// Weather service http point api endpoint.
	// https://api.weather.gov/points/{latitude},{longitude}
	pointURL = "https://api.weather.gov/points/%v,%v"
)

type Fetcher interface {
	Forecast(latitude, longitude float64) (Forecast, error)
}

type Forecast struct {
	// Short term description of the weather.
	// e.g. `Sunny`
	ShortTerm string

	// The current temperature and unit.
	// e.g. 95F
	Temperature string

	// Description of how the temperature feels..
	// `Hot`, `Cold` or `Moderate`.
	TemperatureFeels string
}

type T struct {
	client *http.Client
}

func (t *T) Forecast(latitude, longitude float64) (Forecast, error) {
	p, err := t.getPointResponse(latitude, longitude)
	if err != nil {
		return Forecast{}, err
	}

	f, err := t.getForecast(p.HourlyForecastURL)
	if err != nil {
		return Forecast{}, err
	}

	// Return error if no forecast periods are found.
	if len(f.ForecastPeriods) == 0 {
		return Forecast{}, errors.New("no forecast periods found")
	}

	return Forecast{
		ShortTerm:        f.ForecastPeriods[0].ShortForecast,
		Temperature:      fmt.Sprintf("%v,%s", f.ForecastPeriods[0].Temperature, f.ForecastPeriods[0].TemperatureUnit),
		TemperatureFeels: temperatureFeelsLike(f.ForecastPeriods[0].Temperature),
	}, nil
}

func (t *T) getPointResponse(latitude, longitude float64) (pointResponse, error) {
	resp, err := t.client.Get(fmt.Sprintf(pointURL, latitude, longitude))
	if err != nil {
		return pointResponse{}, fmt.Errorf("while fetching data from point endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return pointResponse{}, fmt.Errorf("while reading error response from point endpoint: %w "+
				"status_code: %d", err, resp.StatusCode)
		}
		return pointResponse{}, fmt.Errorf("unexpected status code from point endpoint: %d, "+
			"resp_body: %s", resp.StatusCode, string(b))
	}

	var point pointResponse
	if err := json.NewDecoder(resp.Body).Decode(&point); err != nil {
		return pointResponse{}, fmt.Errorf("while decoding response from point endpoint: %w", err)
	}
	return point, nil
}
func (t *T) getForecast(url string) (forecastResponse, error) {
	resp, err := t.client.Get(url)
	if err != nil {
		return forecastResponse{}, fmt.Errorf("while fetching data from forecast endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return forecastResponse{}, fmt.Errorf("while reading error response from forecast endpoint: %w "+
				"status_code: %d", err, resp.StatusCode)
		}
		return forecastResponse{}, fmt.Errorf("unexpected status code from forecast endpoint: %d, "+
			"resp_body: %s", resp.StatusCode, string(b))
	}

	var forecast forecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return forecastResponse{}, fmt.Errorf("while decoding response from forecast endpoint: %w", err)
	}
	return forecast, nil
}

func temperatureFeelsLike(temp int) string {
	if temp < isCold {
		return "Cold"
	}
	if temp > isHot {
		return "Hot"
	}
	return "Moderate"
}

// To fetch the short term forcast first must fetch the wether
// using the `points` api endpoint Which, returns a json response
// containing the url to fetch the hourly forcast in the json
// attribute `properties.forecastHourly`.
// e.g. curl -G https://api.weather.gov/points/39.7456,-97.0892
//
// This struct represents the json response for this endpoint, exluding
// fields that are not being using.  This struct can be extended to
// include fields that may be needed in the future.
type pointResponse struct {
	HourlyForecastURL string `json:"properties.forecastHourly"`
}

type forecastResponse struct {
	// Slice of forcast periods fetched from
	// the National Weather Service API.
	ForecastPeriods []forecast `json:"properties.periods"`
}

type forecast struct {
	Temperature     int    `json:"temperature"`
	TemperatureUnit string `json:"temperatureUnit"`
	ShortForecast   string `json:"shortForecast"`
}
