package services

import (
	"assistant/utils"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"unicode/utf8"
)

// HandleRouteToSettings experimental editing of bot and flag prefixes while bot is running.
// Works fine the first change, but no changes after second edit
func HandleRouteToSettings(subRoute string, botP, flagP *string, flags map[string]string)(discordgo.MessageEmbed, error){
	var settingsEmbed discordgo.MessageEmbed
	var botPrefixEmbed = discordgo.MessageEmbedField{Name: "Bot command prefix unchanged", Value: *botP}		// Fields
	var flagPrefixEmbed = discordgo.MessageEmbedField{Name: "Bot flag command prefix unchanged", Value: *flagP}
	// Set return message Title and description
	settingsEmbed.Title = "Settings"
	settingsEmbed.Description = "Current bot settings"

	// Check if command is valid
	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		// Create fields
		botPrefixEmbed = discordgo.MessageEmbedField{Name: "Bot call command", Value: *botP}
		flagPrefixEmbed = discordgo.MessageEmbedField{Name: "Bot flag command", Value: *flagP}
		fields := []*discordgo.MessageEmbedField{&botPrefixEmbed, &flagPrefixEmbed}

		// Embed fields
		settingsEmbed.Fields = fields
		return settingsEmbed, nil
	case utils.Set:
		// Create fields
		if len(flags) != 0 {
			// Check for existence of value for bot prefix
			if newBotPref, ok := flags["botPref"]; ok {
				if *botP != newBotPref {		// Check that new flag is not same as old flag
					// Set new flag prefix
					*botP = newBotPref
					// Set changed flag message
					flagPrefixEmbed = discordgo.MessageEmbedField{Name: "New Bot command", Value: *flagP}
				}
			}
			// Check for existence of value for flag prefix
			if newFlagPref, ok := flags["flagPref"]; ok {
				// Flag prefix can only be one character
				if utf8.RuneCountInString(newFlagPref) - 1 != 1 {
					err := fmt.Errorf("settings: new flag cannot be more than one character long")
					return settingsEmbed, err
				} else {
					if *flagP != newFlagPref {		// Check that new flag is not same as old flag
						// Set new flag prefix
						*flagP = newFlagPref
						// Set changed flag message
						flagPrefixEmbed = discordgo.MessageEmbedField{Name: "New Bot flag command", Value: *flagP}
					}
				}
			}

			// Create fields
			fields := []*discordgo.MessageEmbedField{&botPrefixEmbed, &flagPrefixEmbed}

			// Embed fields
			settingsEmbed.Fields = fields
			return settingsEmbed, nil
		} else {
			err := fmt.Errorf("settings: missing required flag, use %s help Settings", *botP)
			return settingsEmbed, err
		}
	case utils.Help:
		if len(flags) != 0 {
			return settingsEmbed, nil
		} else {
			return settingsEmbed, nil
		}
	default:
		return settingsEmbed, errors.New("sub route not recognized")
	}
}
