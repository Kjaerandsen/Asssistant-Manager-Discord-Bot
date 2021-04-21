package main

import (
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
)

/*
	Discord bot identification variables
*/
var (
	Token string
	BotPrefix = "@bot"
	FlagPrefix = "-"
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main(){
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)
	// In this example, we only care about receiving message events.
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
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself,
	// and anything that doesn't start with the prefix
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, BotPrefix){
		return
	}

	_, subRoute, route, flags, err := parseContent(m.Content)
	if err != nil{
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	switch route{
	case utils.Weather:
		reply, err := services.HandleRouteToWeather(subRoute, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Send reply
		s.ChannelMessageSend(m.ChannelID, reply)
	case utils.News:
		reply, err := services.HandleRouteToNews(subRoute, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Send reply
		s.ChannelMessageSend(m.ChannelID, reply)
	case utils.Reminders:
		reply, err:= services.HandleRouteToReminder(subRoute, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Send reply
		s.ChannelMessageSend(m.ChannelID, reply)
	case utils.Bills:
		reply, err:= services.HandleRouteToBills(subRoute, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Send reply
		s.ChannelMessageSend(m.ChannelID, reply)
	case utils.MealPlan:
		reply, err:= services.HandleRouteToMeals(subRoute, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Send reply
		s.ChannelMessageSend(m.ChannelID, reply)
	case utils.Config:
		reply, err:= services.HandleRouteToConfig(subRoute, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Send reply
		s.ChannelMessageSend(m.ChannelID, reply)
	case utils.Diag:
		reply, err:= services.HandleRouteToDiag(subRoute, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Send reply
		s.ChannelMessageSend(m.ChannelID, reply)
	case utils.Settings:

	default:
		s.ChannelMessageSend(m.ChannelID, "command not recognized")
	}
}

func parseContent(content string)(string, string, string, map[string]string, error){
	// Variables
	var prefix string
	var subCommand string
	var command string
	var potentialFlags []string

	// Split content
	s := strings.Split(content, " ")
	if len(s) < 3{
		return "", "", "", nil, errors.New("invalid command syntax")
	}
	prefix, subCommand, command = s[0], s[1], s[2]

	if len(s) > 3{
		potentialFlags = s[3:]
	}

	// Process flags
	var flags = make(map[string]string)
	currentFlag := ""

	if len(potentialFlags) != 0{
		for _, element := range potentialFlags{
			if strings.HasPrefix(element, FlagPrefix){
				if _, ok := flags[currentFlag]; ok {
					flags[currentFlag] = strings.TrimSpace(flags[currentFlag])
				}

				currentFlag = element
			} else {
				flags[currentFlag] += element + " "
			}
		}
	}

	return prefix, subCommand, command, flags, nil
}