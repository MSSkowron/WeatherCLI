package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

const (
	apiBaseURL = "http://api.weatherapi.com/v1/forecast.json"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "Cracow"

	apiKeyPath, err := getAPIKeyPath()
	if err != nil {
		log.Fatalf("error getting user's API key: %s", err.Error())
	}

	weatherAPIKey, err := readAPIKeyFromFile(apiKeyPath)
	if err != nil {
		log.Fatalf("error reading API key: %s", err.Error())
	}

	if len(os.Args) > 2 {
		log.Fatalln("Invalid number of arguments. One argument is required that specifies the city.")
	}

	if len(os.Args) == 2 {
		q = os.Args[1]
	}

	weather, err := fetchWeather(q, weatherAPIKey)
	if err != nil {
		log.Fatalf("error while fetching weather: %s", err.Error())
	}

	printCurrentWeather(weather)
	printHourlyForecast(weather)
}

func getAPIKeyPath() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user's home directory: %w", err)
	}
	return homeDirectory + "/weatherCLI/apikey", nil
}

func readAPIKeyFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", fmt.Errorf("no API key found in the file")
}

func fetchWeather(city, apiKey string) (*Weather, error) {
	url := fmt.Sprintf("%s?key=%s&q=%s&days=1&aqi=no&alerts=no", apiBaseURL, apiKey, city)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Weather API returned status code: %d", res.StatusCode)
	}

	var weather Weather
	if err := json.NewDecoder(res.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &weather, nil
}

func printCurrentWeather(weather *Weather) {
	location, current := weather.Location, weather.Current
	msg := fmt.Sprintf("%s, %s: %.0f°C, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)
	color.Green(msg)
}

func printHourlyForecast(weather *Weather) {
	hours := weather.Forecast.Forecastday[0].Hour

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}

		msg := fmt.Sprintf("%s - %.0f°C, %.0f%%, %s\n", date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)
		if hour.ChanceOfRain < 40 {
			fmt.Print(msg)
		} else {
			color.Blue(msg)
		}
	}
}
