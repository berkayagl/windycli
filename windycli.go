package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Windy struct {
	Location struct {
		Name      string `json:"name"`
		Country   string `json:"country"`
		Thrombin  string `json:"tz_id"`
		Localtime string `json:"localtime"`
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
				Cloud        float64 `json:"cloud"`
				ChanceOfSnow float64 `json:"chance_of_snow"`
				SnowCm       float64 `json:"snow_cm"`
				GustMph      float64 `json:"gust_mph"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	var api string
	var city string

	bannerText := `

##   ##    ####   ###  ##  ### ##   ##  ##    ## ##   ####       ####   
##   ##     ##      ## ##   ##  ##  ##  ##   ##   ##   ##         ##    
##   ##     ##     # ## #   ##  ##  ##  ##   ##        ##         ##    
## # ##     ##     ## ##    ##  ##   ## ##   ##        ##         ##    
# ### #     ##     ##  ##   ##  ##    ##     ##        ##         ##    
 ## ##      ##     ##  ##   ##  ##    ##     ##   ##   ##  ##     ##    
##   ##    ####   ###  ##  ### ##     ##      ## ##   ### ###    ####   
                                                                        	  
	
	`
	fmt.Println(bannerText)

	fmt.Print("Enter your Weather API : ")
	fmt.Scan(&api)

	fmt.Println("")

	fmt.Print("Enter your city : ")
	fmt.Scan(&city)

	fmt.Println("")

	res, errors := http.Get("http://api.weatherapi.com/v1/forecast.json?key=" + api + "&q=" + city + "&days=1&aqi=no&alerts=no")

	if errors != nil {
		panic(errors)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available!")
	}

	body, errors := io.ReadAll(res.Body)

	if errors != nil {
		panic(errors)
	}

	var windy Windy

	errors = json.Unmarshal(body, &windy)

	if errors != nil {
		panic(errors)
	}

	location, current, hours := windy.Location, windy.Current, windy.Forecast.Forecastday[0].Hour

	fmt.Printf(
		"City : %s \nCountry : %s \nCurrent temperature: %.0fC° \nWeather forecast : %s\nTime zone : %s\nLocal time : %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
		location.Thrombin,
		location.Localtime,
	)

	fmt.Println("")

	fmt.Println("--- Hourly weather forecast ---")

	fmt.Println("")

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		temperatureText := fmt.Sprintf("%.0fC°", hour.TempC)
		rainChanceText := fmt.Sprintf("%.0f", hour.ChanceOfRain)
		cloudText := fmt.Sprintf("%.0f", hour.Cloud)
		chanceOfSnowText := fmt.Sprintf("%.0f", hour.ChanceOfSnow)
		snowCmText := fmt.Sprintf("%.0f", hour.SnowCm)
		gustMphText := fmt.Sprintf("%.0f", hour.GustMph)

		output := fmt.Sprintf(
			"Hours : %s\nTemperature : %s\nRain chance : %s\nWeather forecast : %s\nCloud percent : %s\nSnow chance : %s\nSnow CM : %s\nGust mph : %s\n",
			date.Format("15:04"),
			temperatureText,
			rainChanceText,
			hour.Condition.Text,
			cloudText,
			chanceOfSnowText,
			snowCmText,
			gustMphText,
		)
		fmt.Println(output)
	}

}
