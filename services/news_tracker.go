package services

import (
	"assistant/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
)

// Endpoint URI
const endpoint = "https://api.bing.microsoft.com/v7.0/news/search"
var token = utils.NewsKey

func HandleRouteToNews(subRoute string, flags map[string]string)([]discordgo.MessageEmbed, error){
	var newsEmbed []discordgo.MessageEmbed
	var newsResponse = utils.NewsAnswer{}

	switch subRoute{
	case utils.Get, utils.View, utils.Check:
		// Declare a new GET request.
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			// Print HTTP request error serverside
			print(err)
			// Not found error on discord
			err = fmt.Errorf("news: error %s could not get resource", http.StatusNotFound)
			return newsEmbed, err
		}
		if len(flags) != 0 {
			// SEPARATE INTO URL SETUP FUNCTION
			///////////////////
			// Adding parameters for HTTPS request
			// Base URI
			param := req.URL.Query()
			// Search term
			param.Add("q", flags["q"])
			// Optional parameters
			// Encode HTTP parameters
			req.URL.RawQuery = param.Encode()
			// Insert the subscription-key header.
			req.Header.Add("Ocp-Apim-Subscription-Key", token)
			///////////////////
			// SEPARATE INTO A JSON BODY REQUESTING FUNCTION

			// Instantiate a client.
			client := new(http.Client)

			// Send the request to Bing API.
			resp, err := client.Do(req)
			if err != nil {
				// Print HTTP request error serverside
				print(err)
				// Not found error on discord
				err = fmt.Errorf("news: error %s could not resolve resource", http.StatusNotFound)
				return newsEmbed, err
			}

			// Close the connection.
			defer resp.Body.Close()

			// Read the results
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// Print HTTP request error serverside
				print(err)
				// Not found error on discord
				err = fmt.Errorf("news: error %s could not decode resource", http.StatusInternalServerError)
				return newsEmbed, err
			}
			////////////////
			// SEPARATE INTO JSON BODY PROCESSING FUNCTION
			// Decode into struct
			err = json.Unmarshal(body, &newsResponse)
			if err != nil {
				fmt.Println(err)
			}

			// Iterate over search results and embed result values.
			for _, result := range newsResponse.Value {
				article := discordgo.MessageEmbed{}
				// Title an description
				article.Title = result.Name
				article.Description = result.Description
				// Image
				//article.Image.URL = result.Image.Thumbnail.ContentUrl
				//article.Image.Width = result.Image.Thumbnail.Width
				//article.Image.Height = result.Image.Thumbnail.Height
				// URL to article
				article.URL = result.URL
				// Provider
				//article.Provider.Name = result.Provider[0].Name
				// Publishing date
				article.Timestamp = result.DatePublished
				newsEmbed = append(newsEmbed, article)
			}

			return newsEmbed, nil
		} else {
			// Reserved for no flags, return a message here for missing required 'q' param
			return newsEmbed, nil
		}
	case utils.Add, utils.Set:
		if len(flags) != 0{
			return newsEmbed, nil
		} else {
			return newsEmbed, errors.New("flags are needed")
		}
	case utils.Delete, utils.Remove:
		if len(flags) != 0{
			return newsEmbed, nil
		} else {
			return newsEmbed, errors.New("flags are needed")
		}
	default:
		return newsEmbed, errors.New("sub route not recognized")
	}
}