package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"jack-henry-http-server/resolvers"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to NWS Weather API!\n")
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/getWeatherReport/:lat/:lon/", resolvers.GetWeatherReport)

	log.Println("Server starting on :9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}
