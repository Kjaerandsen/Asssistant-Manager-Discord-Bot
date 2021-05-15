package services

import (
	"assistant/DB"
	dataRequests "assistant/DataRequests"
	"assistant/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleRouteToMeals(subRoute string, flags map[string]string, uid string) ([]discordgo.MessageEmbed, error) {
	var mealEmbed = []discordgo.MessageEmbed{}

	switch subRoute {
	case utils.View, utils.Check:
		return createViewMessage(uid)
	case utils.Get:
		if len(flags) != 0 {
			return mealEmbed, nil
		} else {
			//Get from fridge
			recipes, err := getRecipeFromFridge(uid, 0)
			if err != nil {
				return mealEmbed, err
			}
			return createRecipeMessages(recipes), nil
		}
	case utils.Add, utils.Set:
		if len(flags) != 0 {
			if ingredient, ok := flags[utils.Ingredient]; ok {
				err := addToFridge(ingredient, uid)
				if err != nil {
					return mealEmbed, err
				}
				var info = discordgo.MessageEmbed{}
				info.Title = "Added ingredient " + ingredient
				mealEmbed = append(mealEmbed, info)
				return mealEmbed, nil
			}
			return mealEmbed, errors.New("ingredient flag is required")
		} else {
			return mealEmbed, errors.New("flags are needed")
		}
	case utils.Delete, utils.Remove:
		if len(flags) != 0 {
			if ingredient, ok := flags[utils.Ingredient]; ok {
				//Check if its a singular ingredient
				singular := checkSingleIngredient(ingredient)
				if !singular {
					return mealEmbed, errors.New("Wrong use of flags use, -Ingredient can only be used for singular ingredients")
				}
				err := removeFromFridge(ingredient, uid)
				if err != nil {
					return mealEmbed, err
				}
				var info = discordgo.MessageEmbed{}
				info.Title = "Removed ingredient: " + ingredient
				mealEmbed = append(mealEmbed, info)
				return mealEmbed, nil
			} else if ingredients, ok := flags[utils.Ingredients]; ok {
				list, err := doParseIngredients(ingredients)

				if err != nil {
					return mealEmbed, err
				}
				for index, ing := range list {
					err = removeFromFridge(ing, uid)

					if err != nil {
						//Create error string
						deletedElements := strings.Join(list[0:index], ",")
						remainingElements := strings.Join(list[index:len(list)-1], ",")
						return mealEmbed, errors.New("An element was not found in fridge:" + list[index] + "\n" + "deleted elemments from fridge: " + deletedElements + "\n remaining elements: " + remainingElements)
					}
				}

			}
			return mealEmbed, errors.New("Ingredient(s) flag not found in message. See: help meals for instructions")
		} else {
			return mealEmbed, errors.New("Ingredient or Ingredients flag is needed to delete from fridge, see: help meals for instructions")
		}
	case utils.Help:
		message, err := createHelpMessage()
		return message, err
	default:
		return mealEmbed, errors.New("sub route not recognized")
	}
}

func createHelpMessage() ([]discordgo.MessageEmbed, error) {
	var messageList []discordgo.MessageEmbed
	message := utils.MealPlannerHelper()
	messageList = append(messageList, message)
	return messageList, nil
}

func getRecipeFromFridge(uid string, number int) (utils.Recipe, error) {
	//Use a test fridge until we have an implementation of UserData
	fridge, err := retrieveFridgeIngredients(uid)
	if err != nil {
		fmt.Println(err)
	}
	//Check if fridge is empty
	if len(fridge.Ingredients) == 0 {
		return utils.Recipe{}, errors.New("fridge is empty")
	}
	//Fridge is not empty
	var ingredientString string
	for _, ingredient := range fridge.Ingredients {
		ingredientString += ingredient + ","
	}
	ingredients := strings.ReplaceAll(ingredientString, " ", "")
	fmt.Println(ingredientString)
	//Set number for url
	var numberString string
	if number < 0 {
		numberString = "10"
	}
	//Create url and recipe struct for holding data
	url := "https://api.spoonacular.com/recipes/findByIngredients?ingredients=" + ingredients + "&number=" + numberString + "&apiKey=" + utils.MealKey
	var recipe utils.Recipe
	//Use GetAndDecode function to decode it into recipe struct
	requestError := dataRequests.GetAndDecodeURL(url, &recipe)
	//Check if there was any errors in fetching and decoding the url
	if requestError != nil {
		fmt.Println("Hello what happened")
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
		footer := discordgo.MessageEmbedFooter{Text: "Data provided by"}

		// Set footer and fields
		recipeMessage.Fields = fields
		recipeMessage.Footer = &footer
		messageArray = append(messageArray, recipeMessage)
	}
	return messageArray
}

// retrieveFridge Retrieve fridge with its ingredients from the database
func retrieveFridge(uid string) (map[string]interface{}, error) {
	fmt.Println("Retrieve from fridge command")
	// Retrieve the fridge entry from the database
	fridge, err := DB.RetrieveFromDatabase("fridge", uid)
	if err != nil {
		fmt.Println(err)
		//return fridge, err
	}
	return fridge, nil
}

// retrieveFridgeIngredients Retrieves the ingredients in the format required for recipe searching
func retrieveFridgeIngredients(uid string) (utils.Fridge, error) {
	// Output variable
	var fridgeIngredients utils.Fridge
	// Retrieve the fridge from the database
	fridge, err := retrieveFridge(uid)
	if err != nil {
		return utils.Fridge{}, err
	}
	// Add the fridge ingredients to the fridgeIngredients
	for ingredient := range fridge {
		fridgeIngredients.Ingredients = append(fridgeIngredients.Ingredients, ingredient)
	}
	// Return the formatted ingredients
	return fridgeIngredients, nil
}

// addToFridge Adds an ingredient to the database
func addToFridge(ingredient string, uid string) error {
	fmt.Println("Add to fridge command")
	// Retrieve the fridge from the database
	fridge, err := retrieveFridge(uid)
	if err != nil {
		return err
	}
	// Add the ingredient to the fridge
	fridge[ingredient] = "1"
	// Send the updated fridge to the database
	DB.AddToDatabase("fridge", uid, fridge)

	return nil
}

// removeFromFridge Removes an ingredient from the database
func removeFromFridge(ingredient string, uid string) error {
	// Retrieve the fridge from the database
	fridge, err := retrieveFridge(uid)
	if err != nil {
		return err
	}

	if _, ok := fridge[ingredient]; ok {
		// Remove the ingredient
		delete(fridge, ingredient)
		// Send the updated fridge to the database
		DB.AddToDatabase("fridge", uid, fridge)
		return nil
	} else {
		return errors.New("Ingredient: " + ingredient + " is not in fridge")
	}
}

func createViewMessage(uid string) ([]discordgo.MessageEmbed, error) {
	fridge, err := retrieveFridgeIngredients(uid)
	if err != nil {
		return []discordgo.MessageEmbed{}, err
	}
	//fridge := createTestFridge()

	var messageList []discordgo.MessageEmbed
	var message discordgo.MessageEmbed

	//Create message
	message.Title = "Your Fridge"

	var ingredients string
	for _, ingredient := range fridge.Ingredients {
		ingredients += ingredient + "\n"
	}
	if len(fridge.Ingredients) < 1 {
		ingredients = "There are no ingredients stored in your fridge"
	}
	//Embed ingredients to message
	fridgeContent := discordgo.MessageEmbedField{Name: "Ingredients: ", Value: ingredients}
	fields := []*discordgo.MessageEmbedField{&fridgeContent}

	message.Fields = fields
	messageList = append(messageList, message)
	return messageList, nil
}

//checkSingleIngredient is a simple test for single ingredients values
func checkSingleIngredient(ingredient string) bool {
	list := strings.Split(ingredient, ",") //Split on ,
	if len(list) > 1 {
		return false
	}
	return true
}

//doParseIngredients parses the ingredients string to a slice
func doParseIngredients(ingredients string) ([]string, error) {
	list := strings.Split(ingredients, ",") //Split on ,

	if len(list) < 1 {
		return list, errors.New("Error during parsing, please seperate using commas like: -ingredients potato,chicken,sausage")
	}
	//Remove trailing and start spaces
	for _, value := range list {
		value = strings.TrimSpace(value)
	}
	return list, nil
}
