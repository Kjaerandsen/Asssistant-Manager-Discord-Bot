package services

import (
	"assistant/utils"
	"errors"
)

func HandleRouteToNews(subRoute string, flags map[string]string)(string, error){
	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0{
			return "...", nil
		} else {
			return "...", nil
		}
	case utils.Add, utils.Set:
		if len(flags) != 0{
			return "...", nil
		} else {
			return "", errors.New("flags are needed")
		}
	case utils.Delete, utils.Remove:
		if len(flags) != 0{
			return "...", nil
		} else {
			return "", errors.New("flags are needed")
		}
	default:
		return "", errors.New("sub route not recognized")
	}
}