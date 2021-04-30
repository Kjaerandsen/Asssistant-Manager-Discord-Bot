package services

import (
	"assistant/utils"
	"errors"
	"github.com/bwmarrin/discordgo"
)

func HandleRouteToNews(subRoute string, flags map[string]string)(discordgo.MessageEmbed, error){
	var newsEmbed = discordgo.MessageEmbed{}
	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0{
			return newsEmbed, nil
		} else {
			return newsEmbed, nil
		}
	case utils.Add, utils.Set:
		if len(flags) != 0{
			return newsEmbed, nil
		} else {
			return newsEmbed, errors.New("flags are needed")
		}
	case utils.Delete, utils.Remove:
		if len(flags) != 0{
			return newsEmbed, nil
		} else {
			return newsEmbed, errors.New("flags are needed")
		}
	default:
		return newsEmbed, errors.New("sub route not recognized")
	}
}