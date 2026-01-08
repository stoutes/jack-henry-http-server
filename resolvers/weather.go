package resolvers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const userAgent = "(JackHenryWeatherServer, stoutes1242@gmail.com)"

type NWSPointResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

type NWSForecastResponse struct {
	Properties struct {
		Periods []struct {
			Temperature   int    `json:"temperature"`
			ShortForecast string `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

func characterizeTemperature(tempF int) string {
	switch {
	case tempF >= 85:
		return "hot"
	case tempF >= 60:
		return "moderate"
	default:
		return "cold"
	}
}

func makeNWSRequest(url, userAgent string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	return http.DefaultClient.Do(req)
}

func GetWeatherReport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lat, long := ps.ByName("lat"), ps.ByName("lon")
	if len(lat) == 0 || len(long) == 0 {
		http.Error(w, "Missing lat or lon parameter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	// Get grid point
	pointURL := fmt.Sprintf("https://api.weather.gov/points/%s,%s", lat, long)
	resp, err := makeNWSRequest(pointURL, userAgent)
	if err != nil {
		http.Error(w, "Failed to fetch weather grid point", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Location may be outside US", http.StatusNotFound)
		return
	}

	var pointData NWSPointResponse
	if err := json.NewDecoder(resp.Body).Decode(&pointData); err != nil {
		http.Error(w, "Failed to parse grid point data", http.StatusInternalServerError)
		return
	}

	// Get forecast
	forecastResp, err := makeNWSRequest(pointData.Properties.Forecast, userAgent)
	if err != nil {
		http.Error(w, "Failed to fetch forecast", http.StatusInternalServerError)
		return
	}
	defer forecastResp.Body.Close()

	var forecast NWSForecastResponse
	if err := json.NewDecoder(forecastResp.Body).Decode(&forecast); err != nil {
		http.Error(w, "Failed to parse forecast data", http.StatusInternalServerError)
		return
	}

	if len(forecast.Properties.Periods) == 0 {
		http.Error(w, "No forecast data available", http.StatusNotFound)
		return
	}

	// Get current period data
	currentPeriod := forecast.Properties.Periods[0]
	tempCharacterization := characterizeTemperature(currentPeriod.Temperature)

	// Format: "hot, Partly Cloudy"
	output := fmt.Sprintf("%s, %s", tempCharacterization, currentPeriod.ShortForecast)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}
