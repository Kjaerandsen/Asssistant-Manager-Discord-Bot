package services

import (
	"assistant/utils"
	"errors"
	"github.com/bwmarrin/discordgo"
)

func HandleRouteToHelper(subRoute string, flags map[string]string)([]discordgo.MessageEmbed, error){
	var mainHelp []discordgo.MessageEmbed
	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		commands := discordgo.MessageEmbed{
		Title: "Available Services",
		Description: "Services are used for specifying what type resource or functionality you want to use. \n\n" +
			"Use arrows to flip between help pages or \n use the command <prefix> help <command> for specific help",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Prefix", Value: "<@!834015714200649758>"},
			{Name: "Supported Commands", Value: "1. Adding\n 2. Checking\n 3. Checking\n 4. Adding/Checking/Deleting\n 5. Adding/Checking/Deleting\n", Inline: true},
			{Name: "Available Services", Value: "1. Reminder\n 2. Weather\n 3. News\n 4. Bills\n 5. Meal Planner", Inline: true},
			},
		Footer: &discordgo.MessageEmbedFooter{Text:"page 1 of 2"},
		}
		subcommands := discordgo.MessageEmbed{
			Title: "Available Commands",
			Description: "Commands are used for doing specific things with a service, like adding to a resource or viewing information. \n\n" +
				"Use arrows to flip between help pages or \n use the command <prefix> help <command> for specific help",
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Checking", Value: "<get>, <check> or <view>"},
				{Name: "Adding", Value: "<add> or <set>"},
				{Name: "Deleting", Value: "<delete> or <remove>"},
			},
			Footer: &discordgo.MessageEmbedFooter{Text:"page 2 of 2"},
		}
		mainHelp = []discordgo.MessageEmbed{commands, subcommands}
		return mainHelp, nil
	default:
		return mainHelp, errors.New("sub route not recognized")
	}
}