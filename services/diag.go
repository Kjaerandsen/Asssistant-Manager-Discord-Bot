package services

import (
	"assistant/utils"
	"errors"
	"github.com/bwmarrin/discordgo"
)

func HandleRouteToDiag(subRoute string, flags map[string]string)(discordgo.MessageEmbed, error){
	var diagEmbed = discordgo.MessageEmbed{}
	switch subRoute{
	case utils.View:
		if len(flags) != 0{
			return diagEmbed, nil
		} else {
			return diagEmbed, nil
		}
	case utils.Set:
		if len(flags) != 0{
			return diagEmbed, nil
		} else {
			return diagEmbed, errors.New("flags are needed")
		}
	case utils.Remove:
		if len(flags) != 0{
			return diagEmbed, nil
		} else {
			return diagEmbed, errors.New("flags are needed")
		}
	default:
		return diagEmbed, errors.New("sub route not recognized")
	}
}