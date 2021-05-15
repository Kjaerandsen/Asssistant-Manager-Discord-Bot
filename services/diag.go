func HandleRouteToDiag(subRoute string, flags map[string]string)(discordgo.MessageEmbed, error){
	var diagEmbed = discordgo.MessageEmbed{}
	switch subRoute{
	case utils.View, utils.Get, utils.Check:
		diagEmbed.Title = "Diagnostics"
		diagEmbed.Fields = []*discordgo.MessageEmbedField{
			{
				Name: "NewsAPI",
				Value: getStatusCode("___"),
			},
			{Name: "OpenWeatherMap", Value: getStatusCode("http://pro.openweathermap.org/data/2.5/forecast/daily?q=gj%C3%B8vik&cnt=3&appid=94aad1fbb7ae86f5de4cf9aafc51e2e2")},
			{Name: "Big News Search API", Value: "___"},
			{Name: "Spoonacular", Value: getStatusCode("https://api.spoonacular.com/recipes/716429/information?includeNutrition=false&apiKey=" + utils.MealKey)},

		}
		return diagEmbed, nil
	default:
		return diagEmbed, errors.New("sub route not recognized")
	}
}

func getStatusCode(url string)string{
	resp, _ := http.Get(url)
	return resp.Status
}
