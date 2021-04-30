package services

import (
	"assistant/utils"
	"errors"
	"github.com/bwmarrin/discordgo"
)

func HandleRouteToWeather(subRoute string, flags map[string]string)(discordgo.MessageEmbed, error){
	var weatherEmbed = discordgo.MessageEmbed{}
	// Check if command is valid
	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0 {
			return weatherEmbed, nil
		} else {
			// Test for weather embed
			weatherEmbed.Title = "Weather forecast"
			weatherEmbed.Description = "Forecast for the day"

			// Create fields
			temperature := discordgo.MessageEmbedField{Name: "Temperature", Value: "16C"}
			humidity := discordgo.MessageEmbedField{Name: "Humidity", Value: "33%", Inline: true}
			pressure := discordgo.MessageEmbedField{Name: "Pressure", Value: "1024 pHa", Inline: true}
			wind := discordgo.MessageEmbedField{Name: "Wind", Value: "2 m/s", Inline: true}
			visibility := discordgo.MessageEmbedField{Name: "Visibility", Value: "24 km", Inline: true}
			fields := []*discordgo.MessageEmbedField{&temperature, &humidity, &pressure, &wind, &visibility}

			// Create footer
			footer := discordgo.MessageEmbedFooter{Text: "Data provided by datasource"}

			// Set footer and fields
			weatherEmbed.Fields = fields
			weatherEmbed.Footer = &footer

			return weatherEmbed, nil
		}
	default:
		return weatherEmbed, errors.New("sub route not recognized")
	}
}