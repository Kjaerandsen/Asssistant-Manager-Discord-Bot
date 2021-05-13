package services

import (
	dataRequests "assistant/DataRequests"
	"assistant/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleRouteToMeals(subRoute string, flags map[string]string) ([]discordgo.MessageEmbed, error) {
	var mealEmbed = []discordgo.MessageEmbed{}

	switch subRoute {
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0 {
			return mealEmbed, nil
		} else { //Get from fridge
			recipes, err := getRecipeFromFridge()
			if err != nil {
				return mealEmbed, err
			}
			return createRecipeMessages(recipes), nil
		}
	case utils.Add, utils.Set:
		if len(flags) != 0 {
			return mealEmbed, nil
		} else {
			return mealEmbed, errors.New("flags are needed")
		}
	case utils.Delete, utils.Remove:
		if len(flags) != 0 {
			return mealEmbed, nil
		} else {
			return mealEmbed, errors.New("flags are needed")
		}
	default:
		return mealEmbed, errors.New("sub route not recognized")
	}
}

func getRecipeFromFridge() (utils.Recipe, error) {
	//Use a test fridge until we have an implementation of UserData
	fridge := createTestFridge()
	//Check if fridge is empty
	if len(fridge.Ingredients) == 0 {
		return utils.Recipe{}, errors.New("fridge is empty")
	}
	//Fridge is not empty
	var ingredientString string
	for _, ingredient := range fridge.Ingredients {
		ingredientString += ingredient + ","
	}
	//Create url and recipe struct for holding data
	url := "https://api.spoonacular.com/recipes/findByIngredients?ingredients=" + "chicken, pork, beef, apple, pineapple" + "&number=5&apiKey=c7939a239ecd43c49c1654aff9d387d6"
	var recipe utils.Recipe
	//Use GetAndDecode function to decode it into recipe struct
	err := dataRequests.GetAndDecodeURL(url, &recipe)

	//Check if there was any errors in fetching and decoding the url
	if err != nil {
		fmt.Println("Hei fra getrecipe")
		return utils.Recipe{}, err
	}
	return recipe, nil
}

func createRecipeMessages(recipes utils.Recipe) []discordgo.MessageEmbed {
	var messageArray []discordgo.MessageEmbed
	for _, recipe := range recipes {
		recipeMessage := discordgo.MessageEmbed{}
		recipeMessage.Title = recipe.Name
		recipeMessage.Image = &discordgo.MessageEmbedImage{URL: recipe.Image}

		var missedIngredients, usedIngredients string
		for _, ingredients := range recipe.MissedIngredients {
			missedIngredients += strings.Title(ingredients.IngredientName) + "\n"
		}
		//Embed missed ingredients
		fieldMissed := discordgo.MessageEmbedField{Name: "Missed ingredients: ", Value: missedIngredients}
		fields := []*discordgo.MessageEmbedField{&fieldMissed}
		if recipe.UsedIngredientsCount > 0 {
			for _, ingredients := range recipe.UsedIngredients {
				usedIngredients += strings.Title(ingredients.IngredientName)
			}
			fieldUsed := discordgo.MessageEmbedField{Name: "Used Ingredients: ", Value: usedIngredients}
			fields = append(fields, &fieldUsed)
		}

		// Create footer
		footer := discordgo.MessageEmbedFooter{Text: "Data provided by https://api.spoonacular.com"}

		// Set footer and fields
		recipeMessage.Fields = fields
		recipeMessage.Footer = &footer
		messageArray = append(messageArray, recipeMessage)
	}
	return messageArray
}

//createTestFridge returns a fridge with some ingredients
func createTestFridge() utils.Fridge {
	var fridge utils.Fridge
	fridge.Ingredients = append(fridge.Ingredients, "Apple", "Milk", "Chicken", "Butter")
	return fridge
}
