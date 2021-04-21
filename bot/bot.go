package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Variables used for command line parameters
var (
	Token string
	Prefix = "@bot"
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
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, Prefix){
		return
	}

	_, subCommand, command, flags, err := parseContent(m.Content)
	if err != nil{
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	switch command{
	case "weather":
		reply, err := handleCommandToWeather(subCommand, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Send reply
		s.ChannelMessageSend(m.ChannelID, reply)
	case "news":
		reply, err := handleCommandsToNews(subCommand, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Send reply
		s.ChannelMessageSend(m.ChannelID, reply)
	case "reminders":
		reply, err:= handleCommandsToReminder(subCommand, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Send reply
		s.ChannelMessageSend(m.ChannelID, reply)
	case "bills":
		reply, err:= handleCommandsToBill(subCommand, flags)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		// Send reply
		s.ChannelMessageSend(m.ChannelID, reply)
	case "config":
		var reply string

		if subCommand == "view"{
			// Get config file for user
			reply = "config-view"
		} else if subCommand == "set"{
			// Set values in config file for user
			reply = "config-set"
		} else {
			//
			reply = "subCommand not recognized"
		}

		s.ChannelMessageSend(m.ChannelID, reply)
	default:
		s.ChannelMessageSend(m.ChannelID, "command not recognized")
	}
}

func handleCommandToWeather(subCommand string, flags map[string]string)(string, error){
	// Check if command is valid
	switch subCommand{
	case "get", "view", "check":
		if len(flags) != 0 {
			return "Getting weather with flags...", nil
		} else {
			return "Getting default weather...", nil
		}
	default:
		return "", errors.New("sub command not recognized")
	}
}

func handleCommandsToNews(subCommand string, flags map[string]string)(string, error){
	switch subCommand{
	case "get", "view", "check":
		if len(flags) != 0{
			return "...", nil
		} else {
			return "...", nil
		}
	case "add", "set":
		if len(flags) != 0{
			return "...", nil
		} else {
			return "", errors.New("flags are needed")
		}
	case "delete", "remove":
		if len(flags) != 0{
			return "...", nil
		} else {
			return "", errors.New("flags are needed")
		}
	default:
		return "", errors.New(("sub command not recognized"))
	}
}

func handleCommandsToReminder(subCommand string, flags map[string]string)(string, error){
	switch subCommand{
	case "get", "view", "check":
		if len(flags) != 0{
			return "...", nil
		} else {
			return "...", nil
		}
	case "add", "set":
		if len(flags) != 0{
			return "...", nil
		} else {
			return "", errors.New("flags are needed")
		}
	case "delete", "remove":
		if len(flags) != 0{
			return "...", nil
		} else {
			return "", errors.New("flags are needed")
		}
	default:
		return "", errors.New(("sub command not recognized"))
	}
}

func handleCommandsToBill(subCommand string, flags map[string]string)(string, error){
	switch subCommand{
	case "get", "view", "check":
		if len(flags) != 0{
			return "...", nil
		} else {
			return "...", nil
		}
	case "add", "set":
		if len(flags) != 0{
			return "...", nil
		} else {
			return "", errors.New("flags are needed")
		}
	case "delete", "remove":
		if len(flags) != 0{
			return "...", nil
		} else {
			return "", errors.New("flags are needed")
		}
	default:
		return "", errors.New(("sub command not recognized"))
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
			if strings.HasPrefix(element, "-"){
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
