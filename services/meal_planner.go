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
			if ingredient, ok := flags[utils.Ingredient]; ok {
				singular := checkSingleIngredient(ingredient)
				if !singular {
					return mealEmbed, errors.New("Wrong use of flags, -Ingredient can only be used for singular ingredients")
				}
				ingredient = strings.ReplaceAll(ingredient, " ", "")
				recipes, nil := getRecipeByIngredient(ingredient, "10")
				return createRecipeMessages(recipes), nil

			} else if ingredients, ok := flags[utils.Ingredients]; ok {
				singular := checkSingleIngredient(ingredients)
				if !singular {
					return mealEmbed, errors.New("Wrong use of flags or bad syntax, -Ingredients can only be used for multiple ingredients \n given input: " + ingredients + "\n Example of good request: get meal -ingredients potato, chicken, ham")
				}
				//Replace spaces with no space as the url cannot generate without
				ingredients = strings.ReplaceAll(ingredients, " ", "")
				recipes, err := getRecipeByIngredient(ingredients, "10")
				if err != nil {
					return mealEmbed, err
				} else {
					return createRecipeMessages(recipes), nil
				}
			}
		}
		//No Recognizable flags
		recipes, err := getRecipeFromFridge(uid, 0)
		if err != nil {
			return mealEmbed, err
		}
		return createRecipeMessages(recipes), nil

	case utils.Add, utils.Set:
		if len(flags) != 0 {
			if ingredient, ok := flags[utils.Ingredient]; ok {
				singular := checkSingleIngredient(ingredient)
				if !singular {
					return mealEmbed, errors.New("Wrong use of flags use, -Ingredient can only be used for singular ingredients")
				}
				err := addToFridge(ingredient, uid)
				if err != nil {
					return mealEmbed, err
				}
				var info = discordgo.MessageEmbed{}
				info.Title = "Added ingredient " + ingredient
				mealEmbed = append(mealEmbed, info)
				return mealEmbed, nil
				//Check for multiple ingredient flag
			} else if ingredients, ok := flags[utils.Ingredients]; ok {
				singular := checkSingleIngredient(ingredients)
				if !singular {
					return mealEmbed, errors.New("Wrong use of flags use, -Ingredient can only be used for singular ingredients")
				}
				list, err := doParseIngredients(ingredients)
				//If parser found only one ingredient, or wrong syntax
				if err != nil {
					return mealEmbed, err
				}
				message := addMultipleIngredients(list, uid)
				return message, nil
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
				//If multiple ingredient flags are passed
			} else if ingredients, ok := flags[utils.Ingredients]; ok {
				singular := checkSingleIngredient(ingredients)
				if !singular {
					return mealEmbed, errors.New("Wrong use of flags use, -Ingredient can only be used for singular ingredients")
				}
				list, err := doParseIngredients(ingredients)
				//If parser found only one ingredient, or wrong syntax
				if err != nil {
					return mealEmbed, err
				}
				message := removeMultipleFromFridge(list, uid)
				return message, nil
			}
		} else {
			return mealEmbed, errors.New("Ingredient or Ingredients flag is needed to delete from fridge, see: help meals for instructions")
		}
		return mealEmbed, errors.New("Ingredient(s) flag not found in message. See: help meals for instructions")

	case utils.Help:
		return utils.MealPlannerHelper(), nil
	default:
		return mealEmbed, errors.New("sub route not recognized")
	}
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
		return utils.Recipe{}, err
	}
	return recipe, nil
}
func getRecipeByIngredient(ingredients string, number string) (utils.Recipe, error) {
	//TODO - convert number to int and check if its a number
	var recipes utils.Recipe
	//TODO - Implement new structs and make a new function for random
	if ingredients == "" {
		ingredients = "chicken,beef,pork,onion,potato,carrot,rhubarb,honey,butter,milk"
	}
	url := "https://api.spoonacular.com/recipes/findByIngredients?ingredients=" + ingredients + "&number=" + number + "&apiKey=" + utils.MealKey
	requestError := dataRequests.GetAndDecodeURL(url, &recipes)
	//Check if there was any errors in fetching and decoding the url
	if requestError != nil {
		fmt.Println("Hello what happened")
		return utils.Recipe{}, requestError
	}
	return recipes, nil
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
	ingredient = strings.TrimSpace(ingredient)
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
		print(err.Error())
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
	//for _, value := range list {
	//	value = strings.TrimSpace(value)
	//}
	return list, nil
}

func removeMultipleFromFridge(ingredients []string, uid string) []discordgo.MessageEmbed {
	//Initialize description strings
	var missedIngredients, foundIngredients string
	//Iterate over the list and try removing the items
	for _, ingredient := range ingredients {
		ingredient = strings.TrimSpace(ingredient)
		err := removeFromFridge(ingredient, uid)
		if err != nil {
			missedIngredients += " " + ingredient
		} else {
			foundIngredients += " " + ingredient
		}
	}
	var messageArray = []discordgo.MessageEmbed{}
	var info = discordgo.MessageEmbed{}
	foundField := discordgo.MessageEmbedField{Name: "Ingredients that got deleted: ", Value: foundIngredients}
	remainingField := discordgo.MessageEmbedField{Name: "Was not found and not deleted: ", Value: missedIngredients}
	info.Fields = []*discordgo.MessageEmbedField{&foundField, &remainingField}
	messageArray = append(messageArray, info)
	return messageArray
}

func addMultipleIngredients(ingredients []string, uid string) []discordgo.MessageEmbed {
	//Initialize description strings
	var errorMessage string
	//Iterate over the list and try removing the items
	for _, ingredient := range ingredients {
		ingredient = strings.TrimSpace(ingredient)
		err := addToFridge(ingredient, uid)
		if err != nil {
			errorMessage += err.Error() + "\n"
		}
	}
	if errorMessage == "" {
		errorMessage = "No errors found, all ingredients was added to fridge! \n Do view meals to see your fridge!"
	}
	var messageArray = []discordgo.MessageEmbed{}
	var info = discordgo.MessageEmbed{}
	errorField := discordgo.MessageEmbedField{Name: "Errors: ", Value: errorMessage}
	info.Fields = []*discordgo.MessageEmbedField{&errorField}
	messageArray = append(messageArray, info)
	return messageArray
}
