package concurrency

import(
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

// fanOut: In this pattern, the parent goroutine creates 10 child goroutines
// and waits for them to signal their results.

// I usually use gRPC for implementations so if this seems wonky thats why (its also pretty barebones)
func fanOut(w http.ResponseWriter, r *http.Request, ps httprouter.Param) {
	
	children := 10
	ch := make(chan string, children)
	lat, long, apiKey := ps.ByName("lat"), ps.ByName("lon"), ps.ByName("apiKey")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s", lat, long, apiKey)

	// now we're gonna hit this api ten times...but concurrently (rip my free api tier)

	fmt.Println("response: ", string(response))
	for c := 0; c < children; c++ {
		go func(child int) {
			
			resp, err := http.Get(url)
			if err{
				resp := "error"
			}
			ch <- resp
			fmt.Println("child : sent signal :", child)
		}(c)
	}
	for children > 0 {
		d := <-ch
		children--
		fmt.Println(d)
		fmt.Println("parent : recv'd signal :", children)
	}

	time.Sleep(time.Second)
}