package services

import (
	"assistant/DB"
	"assistant/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleRouteToWeather(subRoute string, flags map[string]string, uid string) (discordgo.MessageEmbed, error) {
	var weatherEmbed = discordgo.MessageEmbed{}

	unitsOfMeasurement := make(map[string][3]string)
	unitsOfMeasurement["metric"] = [3]string{"C", " m/s", " km"}
	unitsOfMeasurement["imperial"] = [3]string{"F", " mile per hour", " miles"}
	unitsOfMeasurement["si"] = [3]string{"K", " m/s", " km"}
	currentCity := utils.DefaultCity
	currentUnit := strings.ToLower(utils.DefaultUnit)
	currentUnitsOfMeasurement := unitsOfMeasurement[currentUnit]
	fmt.Println(unitsOfMeasurement)

	// Check if command is valid
	switch subRoute {
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0 {
			if city, ok := flags[utils.Location]; ok {
				currentCity = city
			}

			if unit, ok := flags[utils.Units]; ok {
				unit = strings.ToLower(unit)
				unit = strings.TrimSpace(unit)
				if unit == "imperial" || unit == "si" {
					currentUnit = unit
					currentUnitsOfMeasurement = unitsOfMeasurement[unit]
				}
			}

			currentCity = strings.Title(currentCity)
			currentWeather := getWeather(currentCity, currentUnit)

			// Fill in the weather embed
			weatherEmbed.Title = "Weather forecast for " + currentCity
			weatherEmbed.Description = currentWeather.Weather[0].Description

			// Create fields
			temperature := discordgo.MessageEmbedField{Name: "Temperature", Value: fmt.Sprint(currentWeather.Main.Temp) + " " + currentUnitsOfMeasurement[0]}
			humidity := discordgo.MessageEmbedField{Name: "Humidity", Value: strconv.Itoa(currentWeather.Main.Humidity) + "%", Inline: true}
			pressure := discordgo.MessageEmbedField{Name: "Pressure", Value: strconv.Itoa(currentWeather.Main.Pressure) + " pHa", Inline: true}
			wind := discordgo.MessageEmbedField{Name: "Wind", Value: strconv.Itoa(int(currentWeather.Wind.Speed)) + currentUnitsOfMeasurement[1], Inline: true}
			visibility := discordgo.MessageEmbedField{Name: "Visibility", Value: strconv.Itoa(currentWeather.Visibility) + " " + currentUnitsOfMeasurement[2], Inline: true}
			fields := []*discordgo.MessageEmbedField{&temperature, &humidity, &pressure, &wind, &visibility}

			// Create footer
			footer := discordgo.MessageEmbedFooter{Text: "Data provided by openweathermap.org"}

			// Set footer and fields
			weatherEmbed.Fields = fields
			weatherEmbed.Footer = &footer

			return weatherEmbed, nil

		} else {
			var currentWeather utils.WeatherStruct
			var location string

			//default response if no flag is given
			data, err := DB.RetrieveFromDatabase("weather", uid)
			if err != nil {
				return weatherEmbed, err
			}
			// Checks if the default location is set, if not default to the defaultcity of the program
			if data["location"] == nil {
				location = utils.DefaultCity
			} else {
				location = data["location"].(string)
			}

			currentWeather = getWeather(location, utils.DefaultUnit)
			weatherEmbed.Title = "Weather forecast for " + location
			weatherEmbed.Description = currentWeather.Weather[0].Description

			// Create fields
			temperature := discordgo.MessageEmbedField{Name: "Temperature", Value: fmt.Sprint(currentWeather.Main.Temp) + "C"}
			humidity := discordgo.MessageEmbedField{Name: "Humidity", Value: strconv.Itoa(currentWeather.Main.Humidity) + "%", Inline: true}
			pressure := discordgo.MessageEmbedField{Name: "Pressure", Value: strconv.Itoa(currentWeather.Main.Pressure) + " pHa", Inline: true}
			wind := discordgo.MessageEmbedField{Name: "Wind", Value: strconv.Itoa(int(currentWeather.Wind.Speed)) + " m/s", Inline: true}
			visibility := discordgo.MessageEmbedField{Name: "Visibility", Value: strconv.Itoa(currentWeather.Visibility) + " km", Inline: true}
			fields := []*discordgo.MessageEmbedField{&temperature, &humidity, &pressure, &wind, &visibility}

			// Create footer
			footer := discordgo.MessageEmbedFooter{Text: "Data provided by openweathermap.org"}

			// Set footer and fields
			weatherEmbed.Fields = fields
			weatherEmbed.Footer = &footer

			return weatherEmbed, nil
		}
	case utils.Help:
		weatherEmbed = utils.WeatherHelper()
		return weatherEmbed, nil
	case utils.Set:
		// Get the city from the command
		if len(flags) != 0 {
			if city, ok := flags[utils.Location]; ok {
				currentCity = city
			}

			// Make the data structure and include the city
			data := make(map[string]interface{})
			data["location"] = currentCity
			// Add it to the database
			DB.AddToDatabase("weather", uid, data)

			weatherEmbed.Title = "Weather location update"
			weatherEmbed.Description = "New default weather location set"
			return weatherEmbed, nil
		} else {
			weatherEmbed.Title = "Something went wrong"
			weatherEmbed.Description = "No location specified when setting default weather location"
			return weatherEmbed, nil
		}
	default:
		// Error embed passed
		weatherEmbed.Title = "Something went wrong"
		weatherEmbed.Description = "Uknown flag was passed, please use @bot help weather to see what flags are available"

		return weatherEmbed, nil

		//return weatherEmbed, errors.New("sub route not recognized")
	}
}

/*
	Sends get request to the weather api for the given country
	and units of measurement system the data should be presented with
	returns a WeatherStruct
*/
func getWeather(city string, unit string) utils.WeatherStruct {

	response, err := http.Get(utils.WeatherAPI + city + "&units=" + unit + "&appid=" + utils.APIKey)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	responseObject := utils.WeatherStruct{}
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	return responseObject
}
