package services

import (
	"assistant/utils"
	"errors"
)

func HandleRouteToConfig(subRoute string, flags map[string]string)(string, error){
	switch subRoute{
	case utils.View:
		if len(flags) != 0{
			return "...", nil
		} else {
			return "...", nil
		}
	case utils.Set:
		if len(flags) != 0{
			return "...", nil
		} else {
			return "", errors.New("flags are needed")
		}
	case utils.Remove:
		if len(flags) != 0{
			return "...", nil
		} else {
			return "", errors.New("flags are needed")
		}
	default:
		return "", errors.New("sub route not recognized")
	}
}