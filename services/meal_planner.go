package services

import (
	"assistant/utils"
	"errors"
	"github.com/bwmarrin/discordgo"
)

func HandleRouteToMeals(subRoute string, flags map[string]string)(discordgo.MessageEmbed, error){
	var mealEmbed = discordgo.MessageEmbed{}
	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0{
			return mealEmbed, nil
		} else {
			return mealEmbed, nil
		}
	case utils.Add, utils.Set:
		if len(flags) != 0{
			return mealEmbed, nil
		} else {
			return mealEmbed, errors.New("flags are needed")
		}
	case utils.Delete, utils.Remove:
		if len(flags) != 0{
			return mealEmbed, nil
		} else {
			return mealEmbed, errors.New("flags are needed")
		}
	default:
		return mealEmbed, errors.New("sub route not recognized")
	}
}
