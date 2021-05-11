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
	"strconv"
	"strings"
	"syscall"
	"time"
	clock "time"
)

/*
	Discord bot identification variables
*/
var (
	Token string
	FlagPrefix = "-"
	BotPrefix = "@news"
)

func init() {
	Token = os.Getenv("BOT_TOKEN")
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	Token = os.Getenv("BOT_TOKEN")
}

func main() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Initiates the database connection
	DB.DatabaseInit()

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
	var reply discordgo.MessageEmbed
	var replies []discordgo.MessageEmbed
	var err error

	// Ignore all messages created by the bot itself,
	// and anything that doesn't start with the prefix
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, BotPrefix) {
		return
	}

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
		reply, err = services.HandleRouteToWeather(subRoute, flags)
	case utils.News:
		replies, err = services.HandleRouteToNews(subRoute, flags)
	case utils.Reminders:
		reply, err = services.HandleRouteToReminder(subRoute, flags)
		if err != nil {			// Error handling
			break
		}
		// Unsure where to put this, but I need to make a goroutine and get user info from inital message.
		// Prase time to seconds
		var time clock.Duration
		split := strings.Split(flags["-time"], " ")
		count, type_ := split[0], split[1]
		var i int
		i, err = strconv.Atoi(count)
		if err != nil{
			break
		}

		if type_ == "day" || type_ == "days" {
			time = clock.Duration(i * 24 * 60 * 60) * clock.Second
		} else if type_ == "hour" || type_ == "hours" {
			time = clock.Duration(i * 60 * 60) * clock.Second
		} else if type_ == "minute" || type_ == "minutes" {
			time = clock.Duration(i * 60) * clock.Second
		} else {
			time = clock.Duration(i) * clock.Second
		}

		// Check for max value
		if time >= 2592000 * clock.Second{
			err = errors.New("time exceeds maximum of 30 days")
			break
		}

		// Get users
		var users []*discordgo.User
		if _, ok := flags["-users"]; !ok {
			users = append(users, m.Author)
		} else {
			users = m.Mentions
		}

		// Get channel
		var channel *discordgo.Channel
		if _, ok := flags["-channel"]; !ok {
			channel, _ = s.Channel(m.ChannelID)
		} else {
			channel = func(str string)*discordgo.Channel{
				str = strings.TrimPrefix(str, "<")
				str = strings.TrimPrefix(str, "#")
				str = strings.TrimSuffix(str, ">")
				channel, _ := s.Channel(str)
				return channel
			}(flags["-channel"])
		}

		// Create coroutine and make it wait
		go func(time clock.Duration, channel *discordgo.Channel, users []*discordgo.User) {
			clock.Sleep(time)
			reply.Description = flags["-message"]

			var mentions string
			for _, user := range users{
				mentions += " " + user.Mention()
			}

		if channel.ID != m.ChannelID{
			// Send in specified channel
			s.ChannelMessageSend(channel.ID, mentions)
			s.ChannelMessageSendEmbed(channel.ID, &reply)
		}else if users[0] == m.Author{
			// Send in DM channel
			dmchannel, _ := s.UserChannelCreate(users[0].ID)
			s.ChannelMessageSendEmbed(dmchannel.ID, &reply)
		} else {
			// Send in default channel
			s.ChannelMessageSend(channel.ID, mentions)
			s.ChannelMessageSendEmbed(channel.ID, &reply)
		}
		}(time, channel, users)

	case utils.Bills:
		reply, err = services.HandleRouteToBills(subRoute, flags)
	case utils.MealPlan:
		reply, err = services.HandleRouteToMeals(subRoute, flags)
	case utils.Config:
		reply, err = services.HandleRouteToConfig(subRoute, flags)
	case utils.Diag:
		reply, err = services.HandleRouteToDiag(subRoute, flags)
	case utils.Settings:
		reply, err = services.HandleRouteToSettings(subRoute, &BotPrefix, &FlagPrefix, flags)
	default:
		s.ChannelMessageSend(m.ChannelID, "command not recognized")
	}

	// Send reply
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
	} else if len(replies) > 0 {
		message, _ := s.ChannelMessageSendEmbed(m.ChannelID, &replies[0])
		go spinReaction(message.ID, m.ChannelID, replies, s)
	} else {
		s.ChannelMessageSendEmbed(m.ChannelID, &reply)
	}
}

func spinReaction(messageID string, channelID string, replies []discordgo.MessageEmbed, s *discordgo.Session,){
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

		if users, _ := s.MessageReactions(channelID, messageID, "◀", 2, "", ""); len(users) > 1{ // "◀"
			index -= 1

			if index < 0{
				index = len(replies) - 1
			}

			s.ChannelMessageEditEmbed(channelID, messageID, &replies[index])
			s.MessageReactionRemove(channelID, messageID, "◀", users[0].ID)

			// Reset counter
			i = 0
		} else if users, _ := s.MessageReactions(channelID, messageID, "▶", 2, "", ""); len(users) > 1{ // "▶"
			index += 1

			if index >= len(replies){
				index = 0
			}

			s.ChannelMessageEditEmbed(channelID, messageID, &replies[index])
			s.MessageReactionRemove(channelID, messageID, "▶", users[0].ID)

			// Reset counter
			i = 0
		} else if users, _ := s.MessageReactions(channelID, messageID, "❌", 2, "", ""); len(users) > 1{ // "❌"
			s.MessageReactionsRemoveAll(channelID, messageID)
			break
		}

		time.Sleep(1 * time.Second)
		i += 1
	}

	// Done waiting, remove all reactions
	s.MessageReactionsRemoveAll(channelID, messageID)
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

	if len(potentialFlags) != 0 {
		for _, element := range potentialFlags{
			if strings.HasPrefix(element, FlagPrefix){
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


