package services

import (
	"assistant/utils"
	"errors"
	"github.com/bwmarrin/discordgo"
)

func HandleRouteToReminder(subRoute string, flags map[string]string)(discordgo.MessageEmbed, error){
	var reminderEmbed discordgo.MessageEmbed
	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0{
			return reminderEmbed, nil
		} else {
			return reminderEmbed, nil
		}
	case utils.Add, utils.Set:
		if len(flags) != 0{
			return reminderEmbed, nil
		} else {
			return reminderEmbed, errors.New("flags are needed")
		}
	case utils.Delete, utils.Remove:
		if len(flags) != 0{
			return reminderEmbed, nil
		} else {
			return reminderEmbed, errors.New("flags are needed")
		}
	default:
		return reminderEmbed, errors.New("sub route not recognized")
	}
}