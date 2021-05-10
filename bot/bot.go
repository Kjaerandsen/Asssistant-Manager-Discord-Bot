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
	"time"
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
	var reply discordgo.MessageEmbed
	var replies []discordgo.MessageEmbed
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
		replies, err = services.HandleRouteToNews(subRoute, flags)
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
	case utils.Settings:
		reply, err = services.HandleRouteToSettings(subRoute, &BotPrefix, &FlagPrefix, flags)
	default:
		s.ChannelMessageSend(m.ChannelID, "command not recognized")
	}

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
	} else if len(replies) > 0 {
		message, _ := s.ChannelMessageSendEmbed(m.ChannelID, &replies[0])
		go spinReaction(message.ID, m.ChannelID, replies, s)
	} else {
		s.ChannelMessageSendEmbed(m.ChannelID, &reply)
	}
/*
	switch reply. {
		case discordgo.MessageEmbed:
			message := reply.(discordgo.MessageEmbed)
			s.ChannelMessageSendEmbed(m.ChannelID, &message)
		case []discordgo.MessageEmbed:

				replies := []ds.MessageEmbed{
					{Title: "Mutli-page Embed Example", Description: "First Page", Fields: []*ds.MessageEmbedField{{Name: "Temperature", Value: "16C"},{Name: "Article", Value: "India reports highest daily coronavirus deaths"}, {Name: "Other", Value: "9382774"}}, Footer: &ds.MessageEmbedFooter{Text: "Data provided by datasource"}},
					{Title: "Mutli-page Embed Example", Description: "Second Page", Fields: []*ds.MessageEmbedField{{Name: "Temperature", Value: "25C"},{Name: "Article", Value: "Night-time splashdown for four ISS astronauts"}, {Name: "Other", Value: "32905934"}}, Footer: &ds.MessageEmbedFooter{Text: "Data provided by datasource"}},
					{Title: "Mutli-page Embed Example", Description: "Third Page", Fields: []*ds.MessageEmbedField{{Name: "Temperature", Value: "-12C"},{Name: "Article", Value: "The dogs that stayed in Chernobyl"}, {Name: "Other", Value: "0340043"}}, Footer: &ds.MessageEmbedFooter{Text: "Data provided by datasource"}},
				}

			//messages := reply.([]discordgo.MessageEmbed)
			message, _ := s.ChannelMessageSendEmbed(m.ChannelID, &reply.([]discordgo.MessageEmbed)[0])
			go spinReaction(message.ID, m.ChannelID, reply.([]discordgo.MessageEmbed), s)
	default:
			s.ChannelMessageSend(m.ChannelID, err.Error())
	}
	*/
}

func spinReaction(messageID string, channelID string, replies []discordgo.MessageEmbed, s *discordgo.Session,){
	// Add reactions
	s.MessageReactionAdd(channelID, messageID, "◀")
	s.MessageReactionAdd(channelID, messageID, "▶")
	s.MessageReactionAdd(channelID, messageID, "❌")

	var index = 0
	for i := 0; i < 30000;{
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