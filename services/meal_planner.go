package services

import (
	dataRequests "assistant/DataRequests"
	"assistant/utils"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func HandleRouteToMeals(subRoute string, flags map[string]string) (discordgo.MessageEmbed, error) {
	var mealEmbed = discordgo.MessageEmbed{}

	switch subRoute {
	case utils.Get, utils.View, utils.Check:
		if len(flags) != 0 {
			return mealEmbed, nil
		} else { //Get from fridge
			recipes, err := getRecipeFromFridge()
			fmt.Println(len(recipes))
			if err != nil {
				return mealEmbed, err
			}
			fmt.Println(len(recipes))

			mealEmbed.Title = recipes[0].Name
			mealEmbed.Footer.IconURL = recipes[0].Image
			mealEmbed.Description = "nice need to parse missing ingredients, hope the picture worked"
			fmt.Println("we made it here bro")

			return mealEmbed, nil
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

func getRecipeFromFridge() ([]utils.Recipe, error) {
	//Use a test fridge until we have an implementation of UserData
	fridge := createTestFridge()
	//Check if fridge is empty
	if len(fridge.Ingredients) == 0 {
		return []utils.Recipe{}, errors.New("fridge is empty")
	}
	//Fridge is not empty
	var ingredientString string
	for _, ingredient := range fridge.Ingredients {
		ingredientString += ingredient + ","
	}
	//Create url and recipe struct for holding data
	url := "https://api.spoonacular.com/recipes/findByIngredients?ingredients=" + "milk" + "&number=2&apiKey=c7939a239ecd43c49c1654aff9d387d6"
	var recipe []utils.Recipe
	//Use GetAndDecode function to decode it into recipe struct
	err := dataRequests.GetAndDecodeURL(url, recipe)
	//Check if there was any errors in fetching and decoding the url
	if err != nil {
		fmt.Println("NO")
		return []utils.Recipe{}, err
	}
	return recipe, nil
}

//createTestFridge returns a fridge with some ingredients
func createTestFridge() utils.Fridge {
	var fridge utils.Fridge
	fridge.Ingredients = append(fridge.Ingredients, "Apple", "Milk", "Chicken", "Butter")
	return fridge
}
