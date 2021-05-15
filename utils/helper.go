package utils

import "github.com/bwmarrin/discordgo"

func WeatherHelper() discordgo.MessageEmbed {
	var embed = discordgo.MessageEmbed{
		Title: "Weather help",
		Description: "Weather commands fetches the current weather report for a specified location. \n" +
			"Available flags",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "-location", Value: "Parameter: 'city name' \n" +
				"						Example: get weather -location Oslo"},
			{Name: "-units", Value: "Parameter: 'system of measurement' \n " +
				"					Available systems: Metric, Imperial and SI \n" +
				"					Example: get weather -units Imperial"},
			{Name: "-lang", Value: "Parameter: 'Language code' \n " +
				"					Available language codes: en, no, es, fr, de, zh_tw \n" +
				"					Example: get weather -lang no"},
			{Name: "Set default values:", Value: "Example: set weather -location london -units imperial -lang fr"},
		},
	}
	return embed
}

func ReminderHelper() discordgo.MessageEmbed {
	var embed = discordgo.MessageEmbed{
		Title: "News help",
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

func NewsHelper() discordgo.MessageEmbed {
	var embed = discordgo.MessageEmbed{
		Title: "News command fetches news articles with a specified topic or something",
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

func MealPlannerHelper() discordgo.MessageEmbed {
	var embed = discordgo.MessageEmbed{
		Title: "Meal Planner Help",
		Description: "The meal planner service provides recipe ideas, these can be based on your fridge or you can be randomized \n \n" +
			"Available flags:",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "-ingredient", Value: "Parameter: 'any ingredient' \n" +
				"						Example: get meals -ingredient potato"},
			{Name: "-ingredients", Value: "Parameter: 'multiple ingredients, separated with commas' \n \n" +
				"					Example: add meals -ingredients chicken, apple, potato \n \n " +
				"-Available commands"},
			{Name: "Get", Value: "Get is used for receiving recipes without any flags, this command will default to using your fridge\n " +
				"					his command has to be used in the combination of ingredient or ingredients flag to search for recipes with \n" +
				"                   given ingredients \n" +
				"                   Example with flag: get meals -ingredients chicken, flour, potato \n \n" +
				"					Example from fridge: get meals"},
			{Name: "-Add & Set", Value: "Add and Set are used for adding an ingredient to your virtual fridge \n " +
				"					This command has to be used in the combination of ingredient or ingredients flag \n" +
				"					Example: add meals -ingredients chicken, apple, potato"},
			{Name: "Check & View", Value: "Check and View are used for viewing your virtual fridge \n " +
				"					This command takes no flags and will ignore it them, if any \n \n" +
				"					Example: view meals"},
			{Name: "Remove & Delete", Value: "Remove and Delete are used for deleting ingredients from your virtual fridge \n " +
				"					his command has to be used in the combination of ingredient or ingredients flag \n \n" +
				"					Example: delete meals -ingredients chicken, apple, potato"},
		},
	}
	return embed
}
