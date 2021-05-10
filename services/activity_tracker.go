package services

import (
	"assistant/utils"
	"errors"
	"github.com/bwmarrin/discordgo"
)

func HandleRouteToReminder(subRoute string, flags map[string]string)(discordgo.MessageEmbed, error){
	var reminderEmbed = discordgo.MessageEmbed{}
	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0{
			return reminderEmbed, errors.New("function not implemented")
		} else {
			return reminderEmbed, errors.New("function not implemented")
		}
	case utils.Add, utils.Set:
		if len(flags) != 0{
			reminderEmbed.Title = "📌 Reminder"
			reminderEmbed.Description = "I will remind "
			if _, ok := flags["users"]; !ok {
				reminderEmbed.Description += "you "
			} else {
				reminderEmbed.Description += "mentioned users "
			}

			reminderEmbed.Description += "about \"" + flags["message"] + "\" in " + flags["time"] + "."


			// Create footer
			footer := discordgo.MessageEmbedFooter{Text: "Data provided by _____"}

			// Set footer and fields
			reminderEmbed.Footer = &footer

			return reminderEmbed, nil
		} else {
			return reminderEmbed, errors.New("flags are needed")
		}
	case utils.Delete, utils.Remove:
		if len(flags) != 0{
			return reminderEmbed, errors.New("function not implemented")
		} else {
			return reminderEmbed, errors.New("flags are needed")
		}
	default:
		return reminderEmbed, errors.New("sub route not recognized")
	}
}