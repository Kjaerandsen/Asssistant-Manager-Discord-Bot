package services

import (
	"assistant/utils"
	"errors"
)

func HandleRouteToWeather(subRoute string, flags map[string]string)(string, error){
	// Check if command is valid
	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0 {
			return "Getting weather with flags...", nil
		} else {
			return "Getting default weather...", nil
		}
	default:
		return "", errors.New("sub route not recognized")
	}
}