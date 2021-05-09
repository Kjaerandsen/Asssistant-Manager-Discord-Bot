package utils

type Fridge struct {
	Ingredients []string
}

//Recipe struct holds all data from the spoonacular API used to find recipes
type Recipe []struct {
	Name                 string           `json:"title"`
	Image                string           `json:"image"`
	MissedIngredients    []IngredientList `json:"missedIngredients"`   //Missing ingredients from
	UsedIngredientsCount int              `json:"usedIngredientCount"` //Amounts of Ingredients used
	UsedIngredients      []IngredientList `json:"usedIngredients"`
}

type IngredientList struct {
	IngredientName string `json:"name"`
}
