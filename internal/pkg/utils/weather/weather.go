package weather

import (
	"fmt"

	"github.com/Avimitin/go-bot/internal/pkg/browser"
)

// GetWeatherSingleLine will return weather in one line
func GetWeatherSingleLine(city string) string {
	url := "https://wttr.in/" + city + "?format=%l的天气:+%c+温度:%t+湿度:%h+降雨量:%p"
	resp, err := browser.Browse(url)
	if err != nil {
		return "Error fetching weather"
	}
	return string(resp)
}

// GetWeatherPic will return picture about given city's weather
func GetWeatherPic(city string) string {
	return fmt.Sprintf("https://wttr.in/%s.png", city)
}
