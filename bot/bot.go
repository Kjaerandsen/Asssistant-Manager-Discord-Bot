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
	Prefix = "@bot "
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

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, Prefix){
		return
	}

	_, command, object, flags, err := parseContent(m.Content)
	if err != nil{
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	switch object{
	case "weather":
		if command == "get"{
			if len(flags) != 0 { // Needs a better check
				message := "Getting weather"
				if val, ok := flags["-location"]; ok {
					message += " in " + val
				}

				if val, ok := flags["-time"]; ok {
					message += " for " + val
				}

				s.ChannelMessageSend(m.ChannelID, message)
			} else {
				// Defaulting
				s.ChannelMessageSend(m.ChannelID, "Getting weather ..")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Command not recognized...")
		}
	case "news":
		// To be implemented
	case "reminders":
		// To be implemented
	case "bills":
		// To be implemented
	}

}

func parseContent(content string)(string, string, string, map[string]string, error){
	// Variables
	var prefix string
	var command string
	var object string
	var potentialFlags []string

	// Split content
	s := strings.Split(content, " ")
	if len(s) < 3{
		return "", "", "", nil, errors.New("invalid command syntax")
	}
	prefix, command, object = s[0], s[1], s[2]

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

	return prefix, command, object, flags, nil
}