package services

import (
	"assistant/DB"
	"assistant/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleRouteToWeather(subRoute string, flags map[string]string, uid string) (discordgo.MessageEmbed, error) {
	var weatherEmbed = discordgo.MessageEmbed{}

	//creates a map containing the available units
	location, unit, language := getDefaultValues(uid)
	currentCity := location
	currentUnit := strings.ToLower(unit)
	currentUnitsOfMeasurement := getCurrentUnitsOfMeasurement(currentUnit)
	currentLang := language

	// Check if command is valid
	switch subRoute {
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0 {
			//Check for -location flag
			if city, ok := flags[utils.Location]; ok {
				currentCity = city
			}
			//Check for -units flag
			if unit, ok := flags[utils.Units]; ok {
				unit = strings.ToLower(unit)
				unit = strings.TrimSpace(unit)
				if unit == "imperial" || unit == "si" || unit == "metric" {
					currentUnit = unit
					currentUnitsOfMeasurement = getCurrentUnitsOfMeasurement(currentUnit)
				}
			}
			//Check for -lang flag
			if lang, ok := flags[utils.Language]; ok {
				lang = strings.ToLower(lang)
				lang = strings.TrimSpace(lang)
				if lang == "es" || lang == "en" || lang == "no" || lang == "fr" || lang == "de" || lang == "zh_tw" {
					currentLang = lang
				}
			}

			currentCity = strings.Title(currentCity)
			//retrieve weather from API using given parameters
			currentWeather, err := getWeather(currentCity, currentUnit, currentLang)
			if err != nil {
				weatherEmbed = utils.WeatherHelper()
				return weatherEmbed, nil
			}

			// Fill in the weather embed
			weatherEmbed.Title = "Weather forecast for " + currentCity
			weatherEmbed.Description = strings.Title(currentWeather.Weather[0].Description)

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
			//default response if no flag is given
			currentWeather, err := getWeather(location, unit, language)
			if err != nil {
				weatherEmbed = utils.WeatherHelper()
				return weatherEmbed, nil
			}
			weatherEmbed.Title = "Weather forecast for " + location
			weatherEmbed.Description = strings.Title(currentWeather.Weather[0].Description)

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
		}
	case utils.Help:
		weatherEmbed = utils.WeatherHelper()
		return weatherEmbed, nil
	case utils.Set:

		if len(flags) != 0 {

			descrption := ""

			// Get the city from the command
			if city, ok := flags[utils.Location]; ok {
				currentCity = city
				descrption += "New default weather location set to " + currentCity + "\n"
			}

			// Get the units from the command
			if unit, ok := flags[utils.Units]; ok {
				unit = strings.ToLower(unit)
				unit = strings.TrimSpace(unit)
				if unit == "imperial" || unit == "si" || unit == "metric" {
					currentUnit = unit
					descrption += "New default units set to " + currentUnit + "\n"
				}
			}
			// Get the language from the command
			if lang, ok := flags[utils.Language]; ok {
				lang = strings.ToLower(lang)
				lang = strings.TrimSpace(lang)
				if lang == "es" || lang == "en" || lang == "no" || lang == "fr" || lang == "de" || lang == "zh_tw" {
					currentLang = lang
					descrption += "New default language set to " + currentLang + "\n"
				}
			}

			// Make the data structure and include the new values
			data := make(map[string]interface{})
			data["location"] = currentCity
			data["units"] = currentUnit
			data["language"] = currentLang
			// Add it to the database
			DB.AddToDatabase("weather", uid, data)

			weatherEmbed.Title = "Default values updated"
			weatherEmbed.Description = descrption
			return weatherEmbed, nil
		} else {
			weatherEmbed.Title = "Something went wrong"
			weatherEmbed.Description = "No parameter specified when setting default values"
			return weatherEmbed, nil
		}
	default:
		// Error embed passed
		weatherEmbed.Title = "Something went wrong"
		weatherEmbed.Description = "Sub route not recognized, please use '@bot help' to see which options are available"

		return weatherEmbed, nil
	}
}

/*
	Sends get request to the weather api for the given country
	and units of measurement system the data should be presented with
	returns a WeatherStruct
*/
func getWeather(city string, unit string, lang string) (utils.WeatherStruct, error) {

	weatherStruct := utils.WeatherStruct{}

	response, err := http.Get(utils.WeatherAPI + city + "&units=" + unit + "&lang=" + lang + "&appid=" + utils.WeatherKey)
	if err != nil {
		return weatherStruct, errors.New("possible HTTP error, or too many redirects")
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return weatherStruct, errors.New("possible HTTP error, or too many redirects")
	}

	responseObject := utils.WeatherStruct{}
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		return weatherStruct, errors.New("possible HTTP error, or too many redirects")
	}

	return responseObject, nil
}

func getDefaultValues(uid string) (string, string, string) {
	var location string
	var unit string
	var language string

	data, err := DB.RetrieveFromDatabase("weather", uid)
	if err != nil {
		//if weather collection doesn't exist create it with program default values as defaults
		data := make(map[string]interface{})
		data["location"] = utils.DefaultCity
		data["units"] = utils.DefaultUnit
		data["language"] = utils.DefaultLanguage
		// Add it to the database
		DB.AddToDatabase("weather", uid, data)
	}
	// Checks if the values are set, if not default to the default values of the program
	if data["location"] == nil {
		location = utils.DefaultCity
	} else {
		location = data["location"].(string)
	}

	if data["units"] == nil {
		unit = utils.DefaultUnit
	} else {
		unit = data["units"].(string)
	}

	if data["language"] == nil {
		language = utils.DefaultLanguage
	} else {
		language = data["language"].(string)
	}

	return location, unit, language
}

func getCurrentUnitsOfMeasurement(unit string) [3]string {
	unitsOfMeasurement := make(map[string][3]string)
	unitsOfMeasurement["metric"] = [3]string{"C", " m/s", " km"}
	unitsOfMeasurement["imperial"] = [3]string{"F", " mile per hour", " miles"}
	unitsOfMeasurement["si"] = [3]string{"K", " m/s", " km"}

	return unitsOfMeasurement[unit]
}
