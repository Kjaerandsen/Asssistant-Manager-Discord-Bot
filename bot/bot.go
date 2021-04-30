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
	BotPrefix = "@news"
	FlagPrefix = "-"
)

func init() {
	Token = os.Getenv("BOT_TOKEN")
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
	discord.AddHandler(router)
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
func router(s *discordgo.Session, m *discordgo.MessageCreate) {
	var reply = discordgo.MessageEmbed{}
	var err error

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
		reply, err = services.HandleRouteToWeather(subRoute, flags)
	case utils.News:
		reply, err = services.HandleRouteToNews(subRoute, flags)
	case utils.Reminders:
		reply, err = services.HandleRouteToReminder(subRoute, flags)
	case utils.Bills:
		reply, err = services.HandleRouteToBills(subRoute, flags)
	case utils.MealPlan:
		reply, err = services.HandleRouteToMeals(subRoute, flags)
	case utils.Config:
		reply, err = services.HandleRouteToConfig(subRoute, flags)
	case utils.Diag:
		reply, err = services.HandleRouteToDiag(subRoute, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
	case utils.Settings:
	default:
		s.ChannelMessageSend(m.ChannelID, "command not recognized")
	}

	// Send reply
	if err != nil {			// Error handling
		s.ChannelMessageSend(m.ChannelID, err.Error())
	} else {				// Result
		s.ChannelMessageSendEmbed(m.ChannelID, &reply)
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