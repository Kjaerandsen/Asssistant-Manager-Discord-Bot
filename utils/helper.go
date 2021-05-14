package utils

import "github.com/bwmarrin/discordgo"

func WeatherHelper()discordgo.MessageEmbed{
	var embed = discordgo.MessageEmbed{
		Title:  "Weather help",
		Description: "Weather commands fetches the current weather report for a specified location. \n" +
			"Available Commands",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "-location", Value: "Parameter: 'city name' \n" +
				"						Example: get weather -location Oslo"},
			{Name: "-units", Value: "Parameter: 'system of measurement' \n " +
				"					Available systems: Metric, Imperial and SI \n" +
				"					Example: get weather -units Imperial"},
		},
	}
	return embed
}

func NewsHelper()discordgo.MessageEmbed{
	var embed = discordgo.MessageEmbed{
		Title:  "News help",
		Description: "News commands fetches current news articles for a specified topic. \n" +
			"Available Commands",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "-location", Value: "Parameter: 'city name' \n" +
				"						Example: get weather -location Oslo"},
			{Name: "-units", Value: "Parameter: 'system of measurement' \n " +
				"					Available systems: Metric, Imperial and SI \n" +
				"					Example: get weather -units Imperial"},
		},
	}
	return embed
}