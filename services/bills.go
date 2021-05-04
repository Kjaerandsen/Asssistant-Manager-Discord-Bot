package services

import (
	"assistant/utils"
	"errors"
	"github.com/bwmarrin/discordgo"
)

func HandleRouteToBills(subRoute string, flags map[string]string)(discordgo.MessageEmbed, error){
	var billsEmbed = discordgo.MessageEmbed{}
	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0{
			return billsEmbed, nil
		} else {
			return billsEmbed, nil
		}
	case utils.Add, utils.Set:
		if len(flags) != 0{
			return billsEmbed, nil
		} else {
			return billsEmbed, errors.New("flags are needed")
		}
	case utils.Delete, utils.Remove:
		if len(flags) != 0{
			return billsEmbed, nil
		} else {
			return billsEmbed, errors.New("flags are needed")
		}
	default:
		return billsEmbed, errors.New("sub route not recognized")
	}
}