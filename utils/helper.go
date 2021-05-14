package utils

import "github.com/bwmarrin/discordgo"

func WeatherHelper()discordgo.MessageEmbed{
	var embed = discordgo.MessageEmbed{
		Title:  "Weather help",
		Description: "Weather commands fetches the current weather report for a specified location. \n" +
			"Available flags",
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

func ReminderHelper()discordgo.MessageEmbed{
	var embed = discordgo.MessageEmbed{
		Title:  "News help",
		Description: "Reminder command sets a reminder for specified user(s) \n" +
			"Available flags",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "-time", Value: "Parameter: 'time in seconds/minutes/hours/days' \n" +
				"					Example: set reminder -time 3 minutes"},
			{Name: "-message", Value: "Parameter: 'reminder message' \n " +
				"					Example: set reminder -message Remember to deliver assignment"},
			{Name: "-channel", Value: "Parameter: 'tagged channel' \n " +
				"					Example: set reminder -channel #general"},
			{Name: "-users", Value: "Parameter: 'tagged users' \n " +
				"					Example: set reminder -users @everyone"},
		},
	}
	return embed
}

func NewsHelper()discordgo.MessageEmbed{
	var embed = discordgo.MessageEmbed{
		Title:  "News command fetches news articles with a specified topic or something",
		Description: "" +
			"Available flags",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "-q", Value: "Parameter: 'topic query' \n" +
				"					Example: get news -q Bitcoin"},
			{Name: "-p", Value: "Parameter: 'page number' \n " +
				"					Example: get news -p 10"},
		},
	}
	return embed
}

func MealPlannerHelper()discordgo.MessageEmbed{
	var embed = discordgo.MessageEmbed{
		Title:  "News help",
		Description: "Meal Planner command does something, idk \n" +
			"Available flags",
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