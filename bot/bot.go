package main

import (
	"assistant/DB"
	"assistant/services"
	"assistant/utils"
	"errors"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
	clock "time"
)

/*
	Discord bot identification variables
*/
var (
	Token      string
	FlagPrefix = "-"
	BotPrefix  = "<@!834015714200649758>"
)

func init() {
	Token = os.Getenv("BOT_TOKEN")
	if Token == "" {
		flag.StringVar(&Token, "t", "", "Bot Token")
		flag.Parse()
	}
}

func main() {
	discord, err := discordgo.New("Bot " + Token)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Initiates the database connection
	DB.DatabaseInit()

	checkReminders(discord)

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(router)
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func router(s *discordgo.Session, m *discordgo.MessageCreate) {

	var reply discordgo.MessageEmbed
	var replies []discordgo.MessageEmbed
	var err error

	// Ignore all messages created by the bot itself,
	// and anything that doesn't start with the prefix
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, BotPrefix) {
		return
	}

	//s.ChannelMessageDelete(m.ChannelID, m.ID)

	// TODO remove this, for testing purposes only
	// m.Author.ID is a unique identifier for the user typing the message
	//DB.Test(m.Author.ID)

	_, subRoute, route, flags, err := parseContent(m.Content)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	switch route {
	case utils.Weather:
		replies, err = services.HandleRouteToWeather(subRoute, flags, m.Author.ID)
	case utils.News:
		replies, err = services.HandleRouteToNews(subRoute, flags)
	case utils.Bills:
		reply, err = services.HandleRouteToBills(subRoute, flags, m.Author.ID)
	case utils.MealPlan:
		replies, err = services.HandleRouteToMeals(subRoute, flags, m.Author.ID)
	case utils.Config:
		reply, err = services.HandleRouteToConfig(subRoute, flags)
	case utils.Diag:
		reply, err = services.HandleRouteToDiag(subRoute, flags)
	case utils.Settings:
		reply, err = services.HandleRouteToSettings(subRoute, &BotPrefix, &FlagPrefix, flags)
	case utils.Help:
		replies, err = services.HandleRouteToHelper(subRoute, flags)
	case utils.Reminders:
		replies, err = services.HandleRouteToReminder(subRoute, flags, s, m)
	default:
		s.ChannelMessageSend(m.ChannelID, "command not recognized")
	}

	// Send reply
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		/*go func(messageID string, s *discordgo.Session) {
			time.Sleep(1 * time.Second)
			err := s.ChannelMessageDelete(m.ChannelID, messageID)
			fmt.Print(err)
		}(message.ID, s)*/
	} else if len(replies) > 1 {
		message, _ := s.ChannelMessageSendEmbed(m.ChannelID, &replies[0])
		go spinReaction(message.ID, m.ChannelID, replies, s)
		s.MessageReactionsRemoveAll(m.ChannelID, message.ID)
	} else if len(replies) == 1 {
		s.ChannelMessageSendEmbed(m.ChannelID, &replies[0])
	} else {
		s.ChannelMessageSendEmbed(m.ChannelID, &reply)
	}
}

func spinReaction(messageID string, channelID string, replies []discordgo.MessageEmbed, s *discordgo.Session) {
	// Add reactions
	s.MessageReactionAdd(channelID, messageID, "◀")
	s.MessageReactionAdd(channelID, messageID, "▶")
	s.MessageReactionAdd(channelID, messageID, "❌")

	var index = 0
	for i := 0; i < 30000; {
		/* Timeout Counter
		go func(i *int){
			for *i < 30 {

			}
		}(&i)
		*/

		if users, _ := s.MessageReactions(channelID, messageID, "◀", 2, "", ""); len(users) > 1 { // "◀"
			index -= 1

			if index < 0 {
				index = len(replies) - 1
			}

			s.ChannelMessageEditEmbed(channelID, messageID, &replies[index])
			s.MessageReactionRemove(channelID, messageID, "◀", users[0].ID)

			// Reset counter
			i = 0
		} else if users, _ := s.MessageReactions(channelID, messageID, "▶", 2, "", ""); len(users) > 1 { // "▶"
			index += 1

			if index >= len(replies) {
				index = 0
			}

			s.ChannelMessageEditEmbed(channelID, messageID, &replies[index])
			s.MessageReactionRemove(channelID, messageID, "▶", users[0].ID)

			// Reset counter
			i = 0
		} else if users, _ := s.MessageReactions(channelID, messageID, "❌", 2, "", ""); len(users) > 1 { // "❌"
			s.MessageReactionsRemoveAll(channelID, messageID)
			break
		}

		clock.Sleep(1 * clock.Second)
		i += 1
	}
}

func parseContent(content string) (string, string, string, map[string]string, error) {
	// Variables
	var prefix string
	var subCommand string
	var command string
	var potentialFlags []string

	// Split content
	s := strings.Split(content, " ")
	if len(s) < 3 {
		return "", "", "", nil, errors.New("invalid command syntax")
	}
	prefix, subCommand, command = s[0], s[1], s[2]

	if len(s) > 3 {
		potentialFlags = s[3:]
	}

	// Process flags
	var flags = make(map[string]string)
	currentFlag := ""

	if len(potentialFlags) != 0 {
		for _, element := range potentialFlags {
			if strings.HasPrefix(element, FlagPrefix) {
				if _, ok := flags[currentFlag]; ok {
					flags[currentFlag] = strings.TrimSpace(flags[currentFlag])
				}

				currentFlag = strings.TrimPrefix(element, FlagPrefix)
			} else {
				flags[currentFlag] += element + " "
			}
		}
	}

	return prefix, subCommand, command, flags, nil
}

func checkReminders(discord *discordgo.Session){
	reminderList := DB.RetrieveAll("reminders")
	for _, guild := range reminderList{
		for _, reminder := range guild.(map[string]interface{}){
			channel, _ := discord.Channel(reminder.(map[string]interface{})["channel"].(map[string]interface{})["ID"].(string))
			message := reminder.(map[string]interface{})["message"].(string)

			var users []*discordgo.User
			usersList := reminder.(map[string]interface{})["users"]
			for _, userObject := range usersList.([]interface{}){
				ID := userObject.(map[string]interface{})["ID"].(string)
				user,_ := discord.User(ID)
				users = append(users, user)
			}

			guildID := reminder.(map[string]interface{})["channel"].(map[string]interface{})["GuildID"].(string)
			messageID := reminder.(map[string]interface{})["channel"].(map[string]interface{})["LastMessageID"].(string)

			time := reminder.(map[string]interface{})["alarmTime"]
			difference := clock.Since(time.(clock.Time))
			if difference > 0{
				difference = 0
			}

			go services.RunReminder(difference, channel, message, users, guildID, messageID, discord)

		}
	}
}