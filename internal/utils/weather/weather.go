package weather

import (
	"fmt"

	"github.com/Avimitin/go-bot/internal/net/browser"
)

// GetWeatherSingleLine will return weather in one line
func GetWeatherSingleLine(city string) string {
	url := fmt.Sprintf("https://wttr.in/%s?format=4", city)
	resp, err := browser.Browse(url)
	if err != nil {
		return "Error fetching weather"
	}
	return resp
}

// GetWeatherPic will return picture about given city's weather
func GetWeatherPic(city string) string {
	return fmt.Sprintf("https://wttr.in/%s_0.png", city)
}
