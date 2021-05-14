package services

import (
	"assistant/DB"
	"assistant/utils"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func HandleRouteToBills(subRoute string, flags map[string]string, uid string)(discordgo.MessageEmbed, error){
	var billsEmbed = discordgo.MessageEmbed{}
	switch subRoute {
	case utils.Get, utils.View, utils.Check:
		// Get the bills
		bills := DB.RetrieveFromDatabase("bills", uid)
		// Set the title
		billsEmbed.Title = "Bills:"
		// Add the bills to the embed
		for bill, value := range bills {
			billsEmbed.Description = billsEmbed.Description + "Bill: " + fmt.Sprintf("%s\n", bill) +
				" value: " + fmt.Sprintf("%v\n\n", value) + ""
		}
		// Return the embed
		return billsEmbed, nil
	case utils.Add, utils.Set:
		if len(flags) != 0 {
			if name, ok := flags[utils.Name]; ok {
				if value, ok := flags[utils.Value]; ok {
					billAdd(uid,strings.TrimSpace(name),value)
					billsEmbed.Title = "Bill added"
					billsEmbed.Description = "Name: " + name + ", value: " + value
					return billsEmbed, nil
				} else {
					return billsEmbed, errors.New("flags are needed")
				}
			}
			return billsEmbed, errors.New("flags are needed")
		} else {
			return billsEmbed, errors.New("flags are needed")
		}
	case utils.Delete, utils.Remove:
		if len(flags) != 0 {
			if name, ok := flags[utils.Name]; ok {
				billRemove(uid,strings.TrimSpace(name))
				billsEmbed.Title = "Bill Removed"
				billsEmbed.Description = "Name: " + name
				return billsEmbed, nil
			}
			return billsEmbed, errors.New("flags are needed")
		} else {
			return billsEmbed, errors.New("flags are needed")
		}
	default:
		return billsEmbed, errors.New("sub route not recognized")
	}
}

// billAdd Add a bill to the database
func billAdd(uid string, bill string, value string) {
	// Retrieve the already existing bills
	bills := DB.RetrieveFromDatabase("bills", uid)
	// Add the new bill
	bills[bill] = value
	// Send the update to the database
	DB.AddToDatabase("bills", uid, bills)

	return
}

// billRemove Remove a bill from the database if it exists
func billRemove(uid string, name string) {
	// Retrieve the bills from the database
	bills := DB.RetrieveFromDatabase("bills", uid)
	// Remove the bill
	delete(bills, name)
	// Update the database entry
	DB.AddToDatabase("bills", uid, bills)

	return
}