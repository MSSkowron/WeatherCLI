package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=ef95ee45d3e94837b26195852230703&q=London&days=1&aqi=no&alerts=no")
	if err != nil {
		log.Fatalf("error while fetching weather: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Weather API not unavailable: %d", res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error while reading body: %s", err.Error())
	}

	fmt.Println(string(bodyBytes))
}
