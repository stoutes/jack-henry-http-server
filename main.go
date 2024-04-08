package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

type CurrentWeather struct {
	LocationData []LocalWeather  `json:"weather"`
	Main         TemperatureData `json:"main"`
}
type LocalWeather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}
type TemperatureData struct {
	Temp        float64 `json:"temp"`
	TempIndex   float64 `json:"feels_like"`
	TempMin     float64 `json:"temp_min"`
	TempMax     float64 `json:"temp_max"`
	Pressure    int     `json:"pressure"`
	Humidity    int     `json:"humidity"`
	SeaLevel    int     `json:"sea_level"`
	GroundLevel int     `json:"grnd_level"`
}

func getWeatherReport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	lat, long, apiKey := ps.ByName("lat"), ps.ByName("lon"), ps.ByName("apiKey")
	if len(lat) == 0 || len(long) == 0 || len(apiKey) == 0 {
		fmt.Println("empty input! error!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s", lat, long, apiKey)

	//set return types
	w.Header().Set("Content-Type", "application/text-plain")

	//GET op
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error!", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	//read response
	response, err := io.ReadAll(resp.Body)
	fmt.Println("response: ", string(response))
	if err != nil {
		http.Error(w, "failed reading response from weather api: ", 422)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// parse and return weather conditions
	var weather CurrentWeather

	parseErr := json.Unmarshal(response, &weather)
	if parseErr != nil {
		fmt.Println("PARSE ERROR: ", parseErr)
		http.Error(w, "Parse error unmarshalling json!: ", 422)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if weather.LocationData == nil {
		fmt.Println("PARSE ERROR: ", parseErr)
		http.Error(w, "Bad response from weather API!", 422)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("err when closing body: ", err)
			http.Error(w, "Error when closing body: ", 422)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}(resp.Body)

	//final temp output should be some general take on temp (aka its hot)
	//the temp output from the weather api is in kelvin btw
	//just range this bad boy
	temperatureOutput := ""
	if weather.Main.Temp >= 310 {
		temperatureOutput = "ITS HOT HOT HOT!, with "
	} else if weather.Main.Temp >= 290 && weather.Main.Temp < 310 {
		temperatureOutput = "Its feels nice today! (moderate), with "
	} else if weather.Main.Temp >= 273 && weather.Main.Temp < 290 {
		temperatureOutput = "Okay its a little chilly (cold), with "
	} else if weather.Main.Temp >= 0 && weather.Main.Temp < 273 {
		temperatureOutput = "Its REALLY cold, with "
	} else if weather.Main.Temp == 0 {
		temperatureOutput = "Every atom is frozen solid with, "
	} else {
		temperatureOutput = "Couldn't decide what to output *shrug*"
	}

	// now append the weather condition and temp
	// have to convert into byte array for response payload
	returnBytes := []byte(temperatureOutput)
	if weather.LocationData != nil && len(weather.
		LocationData) > 0 {
		returnBytes = append(returnBytes, []byte(weather.LocationData[0].Description)...)
	} else {
		returnBytes = append(returnBytes, []byte("No additional data available")...)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(returnBytes)
	if err != nil {
		fmt.Print("error writing bytes to header: ", err)
		http.Error(w, "Error writing bytes to response header", 422)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
}
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	//http.HandleFunc("/", getRoot)

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/getWeatherReport/:lat/:lon/:apiKey", getWeatherReport)
	log.Fatal(http.ListenAndServe(":6666", router))

}
