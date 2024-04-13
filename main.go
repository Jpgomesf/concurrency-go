package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const apiKey = ""

func fetchWeather(city string, wg *sync.WaitGroup, ch chan<- string) interface{} {
	var data struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}

	defer wg.Done()

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Error fetching weather for %s: %s\n", city, err)
		return data
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Printf("Error decoding weather data for %s: %s\n", city, err)
		return data
	}

	ch <- fmt.Sprintf("Current city %s", city)

	return data

}

func main() {
	startNow := time.Now()

	cities := []string{"Toronto", "London", "Paris", "Tokyo"}

	ch := make(chan string)
	var wg sync.WaitGroup

	for _, city := range cities {
		wg.Add(1)
		go fetchWeather(city, &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		fmt.Println(result)
	}
	fmt.Println("This operation took", time.Since(startNow))
}
