package services

import (
	"assistant/utils"
	"errors"
	"github.com/bwmarrin/discordgo"
)

func HandleRouteToConfig(subRoute string, flags map[string]string)(discordgo.MessageEmbed, error){
	var configEmbed = discordgo.MessageEmbed{}
	switch subRoute{
	case utils.View:
		if len(flags) != 0{
			return configEmbed, nil
		} else {
			return configEmbed, nil
		}
	case utils.Set:
		if len(flags) != 0{
			return configEmbed, nil
		} else {
			return configEmbed, errors.New("flags are needed")
		}
	case utils.Remove:
		if len(flags) != 0{
			return configEmbed, nil
		} else {
			return configEmbed, errors.New("flags are needed")
		}
	default:
		return configEmbed, errors.New("sub route not recognized")
	}
}