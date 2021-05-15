package services

import (
	"assistant/DB"
	"assistant/utils"
	"errors"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	clock "time"
)

func HandleRouteToReminder(subRoute string, flags map[string]string, s *discordgo.Session, m *discordgo.MessageCreate)([]discordgo.MessageEmbed, error){
	var reminderEmbed = []discordgo.MessageEmbed{}
	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		reminderList, err := DB.RetrieveFromDatabase("reminders", m.GuildID)
		if err != nil {
			return nil, errors.New("could not find any reminders on this guild")
		}
		for _, remindersInGuild := range reminderList{
			author := remindersInGuild.(map[string]interface{})["author"]
			if author.(map[string]interface{})["ID"] == m.Author.ID{
				var embed discordgo.MessageEmbed
				embed.Title = "📌 Reminder #" + remindersInGuild.(map[string]interface{})["channel"].(map[string]interface{})["LastMessageID"].(string)
				embed.Description = "I will remind about \"" +
					remindersInGuild.(map[string]interface{})["message"].(string) + "\"."
				embed.Footer = &discordgo.MessageEmbedFooter{Text: "For " + remindersInGuild.(map[string]interface{})["alarmTime"].(clock.Time).String()}
				reminderEmbed = append(reminderEmbed, embed)
			}
		}
		return reminderEmbed, nil

	case utils.Add, utils.Set:
		if len(flags) != 0{
			if _,ok := flags["time"]; ok{
				// Parsing input to an amount of time
				split := strings.Split(flags["time"], " ")
				count, type_ := split[0], strings.ToLower(split[1])

				var amount int
				amount, err := strconv.Atoi(count)
				if err != nil {
					return nil, errors.New("specified time is not an integer")
				}

				// Parsing int time into time time
				var time clock.Duration
				if type_ == "day" || type_ == "days" {
					time = clock.Duration(amount) * clock.Hour * 24
				} else if type_ == "hour" || type_ == "hours" {
					time = clock.Duration(amount) * clock.Hour
				} else if type_ == "minute" || type_ == "minutes" {
					time = clock.Duration(amount) * clock.Minute
				} else {
					time = clock.Duration(amount) * clock.Second
				}

				alarmTime := clock.Now().Add(time)

				// Check for max value
				if time >= 2592000*clock.Second {
					return nil, errors.New("time exceeds maximum of 30 days")
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
				reminder := map[string]interface{}{
					"alarmTime": alarmTime,
					"author": m.Author,
					"message": flags["message"],
					"users": users,
					"channel": channel,
				}

				reminderList, err := DB.RetrieveFromDatabase("reminders", m.GuildID)
				if err != nil{
					DB.AddToDatabase("reminders", m.GuildID, nil)
				}
				reminderList[m.ID] = reminder
				DB.AddToDatabase("reminders", m.GuildID, reminderList)

				// Creating the initial reply
				var embed discordgo.MessageEmbed
				embed.Title = "📌 Reminder"
				embed.Description = "I will remind "
				if _, ok := flags["users"]; !ok {
					embed.Description += "you "
				} else {
					embed.Description += "mentioned users "
				}
				embed.Description += "about \"" + flags["message"] + "\" in " + flags["time"] + "."

				go runReminder(time, channel, users, flags, s, m)

				return []discordgo.MessageEmbed{embed}, nil
			} else {
				return nil, errors.New("time flag is needed")
			}
		} else {
			return nil, errors.New("flags are needed")
		}

	case utils.Delete, utils.Remove:
		if len(flags) != 0{
			reminderList, err := DB.RetrieveFromDatabase("reminders", m.GuildID)
			if err != nil {
				return nil, errors.New("could not find any reminders on this guild")
			}
			id := strings.TrimSpace(flags["id"])
			if _, ok := reminderList[id]; !ok{
				return nil, errors.New("could not find the reminder")
			}
			delete(reminderList, id)
			DB.AddToDatabase("reminders", m.GuildID, reminderList)
			return []discordgo.MessageEmbed{{Title: "📌 Reminder has been successfully deleted"}}, nil
		} else {
			return nil, errors.New("need to specify the -id tag")
		}
	case utils.Help:
		return utils.ReminderHelper(), nil
	default:
		return nil, errors.New("sub route not recognized")
	}
}

func runReminder(time clock.Duration, channel *discordgo.Channel, users []*discordgo.User, flags map[string]string, s *discordgo.Session, m *discordgo.MessageCreate){
	clock.Sleep(time)

	// Check if the reminder is still in the database
	reminderList, _ := DB.RetrieveFromDatabase("reminders", m.GuildID)
	if _, ok := reminderList[m.ID]; !ok{
		return
	}

	// Create reply
	var reply discordgo.MessageEmbed
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
}