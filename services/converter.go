package services

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

func HandleRouteToConverter(subRoute string, flags map[string]string)(discordgo.MessageEmbed, error){
	var convertEmbed = discordgo.MessageEmbed{}
	var type_ string
	var to = flags["to"]
	var from = flags["from"]

	var to_value string
	var from_value string

	if len(flags) != 3 {
		return convertEmbed, errors.New("invalid amount of flags")
	}

	switch subRoute{
	case "time": // Only supports timezones as of yet
		type_ = "Time"

	case "unit":
		type_ = "Unit"
		to_value = flags["value"]


	case "currency":
		type_ = "Currency"

	default:
		return convertEmbed, errors.New("sub route not recognized")
	}

	// Either Unit, Currency or Time
	convertEmbed.Title = type_ + " Converter"

	// Create fields
	fromField := discordgo.MessageEmbedField{Name: to, Value: to_value, Inline: true}
	equals := discordgo.MessageEmbedField{Name: " ", Value: " = ", Inline: true}
	toField := discordgo.MessageEmbedField{Name: from, Value: from_value, Inline: true}
	fields := []*discordgo.MessageEmbedField{&toField, &equals, &fromField}

	// Create footer
	footer := discordgo.MessageEmbedFooter{Text: ""}

	// Set footer and fields
	convertEmbed.Fields = fields
	convertEmbed.Footer = &footer

	return convertEmbed, nil
}
