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

// NewsAnswer This struct formats the answer provided by the Bing News Search API.
type NewsAnswer struct {
	ReadLink       string `json:"readLink"`
	QueryContext   struct {
		OriginalQuery   string   `json:"originalQuery"`
		AdultIntent     bool        `json:"adultIntent"`
	} `json:"queryContext"`
	TotalEstimatedMatches   int  `json:"totalEstimatedMatches"`
	Sort  []struct {
		Name    string  `json:"name"`
		ID       string    `json:"id"`
		IsSelected       bool  `json:"isSelected"`
		URL      string   `json:"url"`
	}  `json:"sort"`
	Value   []struct   {
		Name     string   `json:"name"`
		URL   string    `json:"url"`
		Image   struct   {
			Thumbnail   struct  {
				ContentUrl  string  `json:"thumbnail"`
				Width   int  `json:"width"`
				Height  int   `json:"height"`
			} `json:"thumbnail"`
		} `json:"image"`
		Description  string  `json:"description"`
		Provider  []struct   {
			Type   string    `json:"_type"`
			Name  string     `json:"name"`
		} `json:"provider"`
		DatePublished   string   `json:"datePublished"`
	} `json:"value"`
}
