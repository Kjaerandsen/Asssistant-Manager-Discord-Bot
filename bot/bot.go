package main

import (
	"assistant/services"
	"assistant/utils"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	clock "time"

	"github.com/bwmarrin/discordgo"
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
	//DB.DatabaseInit()

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
	case utils.Reminders:
		replies, err = services.HandleRouteToReminder(subRoute, flags)
		if err != nil { // Error handling
			break
		}

		// Prase time to seconds
		var time clock.Duration
		split := strings.Split(flags["time"], " ")
		count, type_ := split[0], split[1]
		var i int
		i, err = strconv.Atoi(count)
		if err != nil {
			break
		}

		if type_ == "day" || type_ == "days" {
			time = clock.Duration(i*24*60*60) * clock.Second
		} else if type_ == "hour" || type_ == "hours" {
			time = clock.Duration(i*60*60) * clock.Second
		} else if type_ == "minute" || type_ == "minutes" {
			time = clock.Duration(i*60) * clock.Second
		} else {
			time = clock.Duration(i) * clock.Second
		}

		// Check for max value
		if time >= 2592000*clock.Second {
			err = errors.New("time exceeds maximum of 30 days")
			break
		}

		// Get users
		var users []*discordgo.User
		if _, ok := flags["users"]; !ok {
			users = append(users, m.Author)
		} else {
			users = m.Mentions[1:]
		}

		// Get channel
		var channel *discordgo.Channel
		if _, ok := flags["channel"]; !ok {
			channel, _ = s.Channel(m.ChannelID)
		} else {
			channel = func(str string) *discordgo.Channel {
				str = str[1:]
				str = str[1:]
				str = str[:len(str)-1]
				channel, _ := s.Channel(str)
				return channel
			}(strings.Split(flags["channel"], " ")[0])
		}

		// Add reminder to database

		// Create coroutine and make it wait
		go func(time clock.Duration, channel *discordgo.Channel, users []*discordgo.User) {
			clock.Sleep(time)
			// Check if the reminder is still in the database

			reply.Title = "📌 Reminder"
			footer := ""
			for _, user :=  range users{
				u, _ := s.GuildMember(m.GuildID, user.ID)
				footer += u.Nick + " "
			}

			reply.Footer = &discordgo.MessageEmbedFooter{Text: footer}

			reply.Description = "Message: " + flags["message"]

			var mentions string
			for _, user := range users {
				mentions += " " + user.Mention()
			}

			if _, ok := flags["channel"]; ok {
				s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
					Content: mentions,
					Embed: &reply,
				})
			} else if users[0] == m.Author {
				// Send in DM channel
				dmchannel, _ := s.UserChannelCreate(users[0].ID)
				s.ChannelMessageSendEmbed(dmchannel.ID, &reply)
			} else {
				// Send in default channel
				s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
					Content: mentions,
					Embed: &reply,
				})
			}
		}(time, channel, users)
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
	default:
		s.ChannelMessageSend(m.ChannelID, "command not recognized")
	}

	// Send reply
	if err != nil {
		message, _ := s.ChannelMessageSend(m.ChannelID, err.Error())
		go func(messageID string, s *discordgo.Session) {
			time.Sleep(1 * time.Second)
			err := s.ChannelMessageDelete(m.ChannelID, messageID)
			fmt.Print(err)
		}(message.ID, s)
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

		time.Sleep(1 * time.Second)
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
