package services

import (
	"assistant/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Endpoint URI
const endpoint = "https://api.bing.microsoft.com/v7.0/news/search"  // For queried requests
const CategoryEndpoint = "https://api.bing.microsoft.com/v7.0/news"	// For category requests
const articleCount = 10                                             // Articles per page

// API key
var token = utils.NewsKey

// NewsRequest type of request for news
type NewsRequest int

const ( // Enum values
	Search NewsRequest = iota
	Trending
	Category
	Webhooks
)

func HandleRouteToNews(subRoute, userID, channelID string, flags map[string]string) ([]discordgo.MessageEmbed, error) {
	var newsEmbed []discordgo.MessageEmbed
	var page = 1
	var requestType NewsRequest
	var req *http.Request
	var err error

	// Determine user intention based on first flag
	for flag, _ := range flags {
		var intentionDetermined = false
		if !intentionDetermined {
			switch flags[flag] {
			case flags[utils.Query]:
				requestType = Search
			case flags[utils.Category]:
				requestType = Category
			case flags[utils.Hooks]:
				requestType = Webhooks
			default:
				requestType = Trending
			}
			intentionDetermined = true
		}
	}


	switch subRoute {
	case utils.Get, utils.View, utils.Check:
		// Check if request for webhook
		if requestType == Webhooks {
			// Get webhooks funksjon
			// newsEmbed, nil
		}
		// Declare a new GET request.
		if requestType == Category {	// Category URL
			req, err = http.NewRequest(http.MethodGet, CategoryEndpoint, nil)
		} else {	// Trending and Search URL
			req, err = http.NewRequest(http.MethodGet, endpoint, nil)
		}
		if err != nil {
			print(err) // Print HTTP request error serverside
			err = fmt.Errorf("news: error %v could not generate get request", http.StatusInternalServerError)
			return nil, err // Not found error on discord
		}

		if len(flags) != 0 {
			// Get news
			if requestType == Category {
				newsEmbed, err = getNewsByCategory(req, flags)
			} else {
				newsEmbed, err = getNewsByQuery(req, flags, page, requestType)
			}
			if err != nil {
				return nil, err
			}
			return newsEmbed, nil // Response to discord
		} else {
			// Get news
			newsEmbed, err = getNewsTrending(req, page)
			if err != nil {
				return nil, err
			}
			return newsEmbed, nil
		}
	case utils.Add, utils.Set:

		return newsEmbed, nil
	case utils.Delete, utils.Remove:
		if len(flags) != 0 {
			return newsEmbed, nil
		} else {
			return newsEmbed, errors.New("flags are needed")
		}
	default:
		return newsEmbed, errors.New("sub route not recognized")
	}
}

// getNews Uses an HTTP GET request to get search or trending type news
// returns an array of discord message embeds after processing
func getNewsByQuery(req *http.Request, flags map[string]string, page int, reqType NewsRequest) ([]discordgo.MessageEmbed, error) {
	req, page, err := generateNewsHTTPRequest(flags, req, reqType) // HTTP request
	if err != nil {
		print(err) // Print Url generation error serverside
		// Internal error on discord
		err = fmt.Errorf("news: error %v bad request, could not fetch resource, %s", http.StatusBadRequest, err)
		return nil, err
	}
	// Send request and read the results
	body, err := sendNewsHTTPRequest(req)
	if err != nil {
		print(err) // Print HTTP request error serverside
		// Internal error on discord
		err = fmt.Errorf("news: error %v could not read resource", http.StatusInternalServerError)
		return nil, err
	}
	// Decode and generate response message
	newsEmbed, err := generateNewsEmbedResponse(body, page, reqType)
	if err != nil { // In case of decoding error
		print(err) // Print decoding error serverside
		// Internal error on discord
		err = fmt.Errorf("news: error %v could not decode resource", http.StatusInternalServerError)
		return nil, err
	}
	return newsEmbed, nil // Array of discord embedded messages
}

// getNewsTrending Uses an HTTP GET request to get trending type news
// returns an array of discord message embeds after processing
func getNewsTrending(req *http.Request, page int) ([]discordgo.MessageEmbed, error) {
	var reqType = Trending

	req, page, err := generateNewsHTTPRequest(nil, req, reqType) // HTTP request
	if err != nil {
		print(err) // Print Url generation error serverside
		// Internal error on discord
		err = fmt.Errorf("news: error %v bad request, could not fetch resource, %s", http.StatusBadRequest, err)
		return nil, err
	}
	// Send request and read the results
	body, err := sendNewsHTTPRequest(req)
	if err != nil {
		print(err) // Print HTTP request error serverside
		// Internal error on discord
		err = fmt.Errorf("news: error %v could not read resource", http.StatusInternalServerError)
		return nil, err
	}
	// Decode and generate response message
	newsEmbed, err := generateNewsEmbedResponse(body, page, reqType)
	if err != nil { // In case of decoding error
		print(err) // Print decoding error serverside
		// Internal error on discord
		err = fmt.Errorf("news: error %v could not decode resource", http.StatusInternalServerError)
		return nil, err
	}
	return newsEmbed, nil // Array of discord embedded messages
}

// getNewsByCategory Uses an HTTP GET request to get category type news
// returns an array of discord message embeds after processing
func getNewsByCategory(req *http.Request, flags map[string]string) ([]discordgo.MessageEmbed, error) {
	var reqType = Category

	req, page, err := generateNewsHTTPRequest(flags, req, reqType) // HTTP request
	if err != nil {
		print(err) // Print Url generation error serverside
		// Internal error on discord
		err = fmt.Errorf("news: error %v bad request, could not fetch resource, %s", http.StatusBadRequest, err)
		return nil, err
	}
	// Send request and read the results
	body, err := sendNewsHTTPRequest(req)
	if err != nil {
		print(err) // Print HTTP request error serverside
		// Internal error on discord
		err = fmt.Errorf("news: error %v could not read resource", http.StatusInternalServerError)
		return nil, err
	}
	// Decode and generate response message
	newsEmbed, err := generateNewsEmbedResponse(body, page, reqType)
	if err != nil { // In case of decoding error
		print(err) // Print decoding error serverside
		// Internal error on discord
		err = fmt.Errorf("news: error %v could not decode resource", http.StatusInternalServerError)
		return nil, err
	}
	return newsEmbed, nil // Array of discord embedded messages
}

// registerNewsWebhook
func registerNewsWebhook(uID, cID string, flags map[string]string) (string, error) {
	var webhook utils.NewsWebhook
	var err error
	webhook.Address = cID

	if flags[utils.Timeout] != "" { 	// Page option handling
		webhook.Timeout, err = strconv.Atoi(flags[utils.Timeout])
		if err != nil { // Set page to 1 in case invalid character
			webhook.Timeout = 30
		}
	}

	if flags[utils.Category] != "" {	// Category priority
		flags[utils.Category], err = getNewsCategory(flags)
		if err != nil {

		}
		webhook.RequestType = "category"
	} else if flags[utils.Query] != "" {	// Search secondary
		webhook.RequestType = "search"
	} else {	// Trending default
		webhook.RequestType = "search"
	}

	return "category", nil // Array of discord embedded messages
}

// generateNewsHTTPRequest generates an HTTP request based on flags provided
// returns http request and requested page
func generateNewsHTTPRequest(flags map[string]string, req *http.Request, reqType NewsRequest) (*http.Request, int, error) {
	// Adding parameters for HTTPS request
	var page = 1
	var category string
	var err error
	// Base URI
	var param = req.URL.Query()

	if reqType == Trending {
		param.Add("q", "")
		param.Add("mkt", "en-ww")		// US/Global market
		param.Add("sortBy", "relevance")
		param.Add("freshness", "Day")
		param.Add("since", fmt.Sprintf("%v", time.Now().Unix()))
	}

	if reqType == Category {
		category, err = getNewsCategory(flags)
		if err != nil { // Returns failure to set category error
			return nil, 0, err
		}
		param.Add("category", category)

		if category == "Sports" {	// US sports are not as global (and lame)
			param.Add("mkt", "en-gb")		// British market
		} else {	// Otherwise US market because most categories, and they include Worldwide/Regional markets
			param.Add("mkt", "en-us")		// US/Global market
		}
	}

	if flags != nil { // Search / Trending / Category w/optional flags
		if flags[utils.SortBy] == "date" {
			param.Del("sortBy") // Remove existing key
			param.Add("sortBy", flags[utils.SortBy])
		} else { // Default is relevance
			param.Del("sortBy") // Remove existing key
			param.Add("sortBy", "relevance")
		}

		if reqType == Search || reqType == Trending {	// Trending and Search exclusive parameter
			// Page number
			if flags[utils.Page] != "" { 	// Page option handling
				page, err = strconv.Atoi(flags[utils.Page])
				if err != nil { // Set page to 1 in case invalid character
					page = 1
				}
			}
			// Number of search elements to skip based on page number
			param.Add("offset", fmt.Sprintf("%v", articleCount*(page-1)))

			// Image source
			param.Add("originalImg", "true") 	// Include original image
		}

		if reqType == Search {	// Search parameters
			// Optional parameters
			param.Add("mkt", "en-ww")		// Global market
			// Freshness flag (Articles by Day/Week/Month)
			if strings.ToLower(flags[utils.Freshness]) == "day" {
				param.Add("freshness", "Day")
			} else if strings.ToLower(flags[utils.Freshness]) == "week" {
				param.Add("freshness", "Week")
			} else if strings.ToLower(flags[utils.Freshness]) == "month" {
				param.Add("freshness", "Month")
			} else {
				param.Add("freshness", "Week")
			}
			// Since flag
			if flags[utils.Since] != "" {
				if flags[utils.SortBy] == "date" {
					layout := "2006-01-02"
					since, err := time.Parse(layout, flags[utils.Since])
					if err != nil {
						// Error on failed parsing of since date
						return nil, 0, fmt.Errorf("since date flag is invalid, %s should be formatted like YYYY-MM-DD", flags[utils.Since])
					}
					param.Add("since", fmt.Sprintf("%v", since.Unix())) // Query flag required for searching
				}
			}
			// Fixed parameters
			param.Add("q", flags[utils.Query]) // Query flag required for searching
		}
	}

	// Encode HTTP parameters
	req.URL.RawQuery = param.Encode()
	// Insert the subscription-key header.
	req.Header.Add("Ocp-Apim-Subscription-Key", token)
	// Generated request and page number
	return req, page, err
}

// sendNewsHTTPRequest sends an HTTP request
// returns a json body
func sendNewsHTTPRequest(req *http.Request) ([]byte, error) {
	// Instantiate a client.
	client := new(http.Client)
	// Send the request to Bing API.
	resp, err := client.Do(req)
	if err != nil {
		// Print HTTP request error serverside
		print(err)
		// Not found error on discord
		err = fmt.Errorf("news: error %v could not find resource", http.StatusNotFound)
		return nil, err
	}
	// Close the connection.
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			print("News: Warning request connection not closing")
		}
	}(resp.Body)
	// Read the results
	body, err := ioutil.ReadAll(resp.Body)
	return body, err // JSON body
}

// generateNewsEmbedResponse decodes a JSON response for Bing news API
// generates and returns an array of discord embedded messages of the news articles
func generateNewsEmbedResponse(body []byte, page int, reqType NewsRequest) ([]discordgo.MessageEmbed, error) {
	var newsResponse = utils.Newspaper{}
	var embed []discordgo.MessageEmbed
	var footer string
	var err error
	// Decode into struct
	err = json.Unmarshal(body, &newsResponse)
	if err != nil {
		return nil, err
	}
	// Iterate over search results and embed result values.
	for i, result := range newsResponse.Value {
		article := discordgo.MessageEmbed{}
		// Title an description
		article.Title = result.Name
		article.Description = result.Description
		// URL to article
		article.URL = result.URL
		// Publishing date
		article.Timestamp = result.DatePublished.Format("2006-01-02 15:04:05")
		// Provider
		article.Provider = &discordgo.MessageEmbedProvider{
			URL:  "",
			Name: result.Provider[0].Name,
		}
		// Image
		// Get original image only on search requests (API requirement) otherwise use thumbnail
		if result.Image.ContentURL != "" {
			if reqType == Search {
				article.Image = &discordgo.MessageEmbedImage{
					URL: result.Image.ContentURL,
				}
			} else { // For adding newer bing endpoints, this should handle their images
				article.Image = &discordgo.MessageEmbedImage{
					URL: result.Image.Thumbnail.ContentURL,
				}
			}
		}
		// Result count, page number and total search results
		if reqType == Search {
			footer = fmt.Sprintf("Result %v of %v\tPage %v\tTotal results: %v", i+1, len(newsResponse.Value), page, newsResponse.TotalEstimatedMatches)
		} else if reqType == Trending { // No estimated total results when using trending requests
			footer = fmt.Sprintf("Result %v of %v\tPage %v", i+1, len(newsResponse.Value), page)
		} else if reqType == Category {
			footer = fmt.Sprintf("Result %v of %v", i+1, len(newsResponse.Value))
		}

		article.Footer = &discordgo.MessageEmbedFooter{
			Text: footer,
		}
		// Insert article
		embed = append(embed, article)
	}
	return embed, nil // Array of discord messages
}

// getNewsCategory Checks category input if valid
// returns a string and error value
func getNewsCategory(flags map[string]string) (string, error) {
	var category string

	var categories = map[string]bool{ // valid categories
		"Business":               utils.BusinessCategory,
		"Entertainment":          utils.EntertainmentCategory,
		"Health":                 utils.HealthCategory,
		"Politics":               utils.PoliticsCategory,
		"Products":               utils.ProductsCategory,
		"Science And Technology": utils.ScienceAndTechnologyCategory,
		"Sports":                 utils.SportsCategory,
		"US":                     utils.USCategory,
		"World":                  utils.WorldCategory,
		"Africa":                 utils.AfricaCategory,
		"World_Americas":         utils.AmericasCategory,
		"World_Asia":             utils.AsiaCategory,
		"World_Europe":           utils.EuropeCategory,
		"World_MiddleEast":       utils.MiddleEastCategory,
	}

	if flags[utils.Category] != "" { // Check if category set or not
		category = strings.Title(strings.ToLower(flags[utils.Category])) // lowercase, then uppercase first letter
		if val, ok := categories[category]; ok {                         // if category specified is found
			if !val { // check if category activated
				// Return message specifying deactivated category
				return "", fmt.Errorf("the %s category has been deactivated", flags[utils.Category])
			} else {
				// Check for Regional/whitespace exceptions
				switch category {
				case "Africa":
					category = "World_Africa"
				case "Americas":
					category = "World_Americas"
				case "Asia":
					category = "World_Asia"
				case "Europe":
					category = "World_Europe"
				case "Middle East":
					category = "World_MiddleEast"
				case "Science And Technology":
					category = "ScienceAndTechnology"
				}
			}
		} else {
			// Return message specifying wrong category
			return "", fmt.Errorf("%s is not a valid category", flags[utils.Category])
		}

	} else { // Default category
		// Default categories are always active even if deactivated in constants
		return "World", nil
	}

	return category, nil // Array of discord embedded messages
}