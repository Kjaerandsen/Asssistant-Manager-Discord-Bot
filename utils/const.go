package utils

/*
	Route constants
*/
const (
	Reminders = "reminder"
	Weather   = "weather"
	Config    = "config"
	Bills     = "bills"
	News      = "news"
	MealPlan  = "meals"
	Diag      = "diag"
	Settings  = "settings"
)

/*
	Subroute constants
*/
const (
	Get    = "get"
	View   = "view"
	Check  = "check"
	Add    = "add"
	Set    = "set"
	Delete = "delete"
	Remove = "remove"
	Help   = "help"
)

/*
	Flag constants
*/
const (
	Language    = "lang"
	Location    = "location"
	Units       = "units"
	Ingredient  = "ingredient"
	Ingredients = "ingredients"
	Name        = "name"
	Value       = "value"
	FlagPrefix  = "flagPref"
	BotPrefix   = "botPref"
	Query       = "q"
	Category    = "category"
	Page        = "p"
	Freshness   = "fresh"
	Since       = "since"
	SortBy      = "sort"
	Hooks      = "hooks"
	Timeout	   = "timeout"
)

/*
	Weather_forecast constants
*/
const (
	WeatherAPI      = "http://api.openweathermap.org/data/2.5/weather?q="
	WeatherKey      = "94aad1fbb7ae86f5de4cf9aafc51e2e2"
	DefaultCity     = "Gjøvik"
	DefaultUnit     = "Metric"
	DefaultLanguage = "en"
)

/*
	news_tracker constants
*/
const (
	/*
		Categories activated/deactivated
	*/
	BusinessCategory             = true
	EntertainmentCategory        = true
	HealthCategory               = true
	PoliticsCategory             = true
	ProductsCategory             = true
	ScienceAndTechnologyCategory = true
	SportsCategory               = true
	USCategory                   = true
	WorldCategory                = true
	AfricaCategory               = true
	AmericasCategory             = true
	AsiaCategory                 = true
	EuropeCategory               = true
	MiddleEastCategory           = true
)

/*
	API Key constants
*/
const (
	NewsKey = "d24e5c8886074fd18f3d305532291ae8"
)

/*
	MealPlanner constants
*/

const (
	MealPlannerAPI = "https://api.spoonacular.com"
	MealKey        = "c7939a239ecd43c49c1654aff9d387d6"
)
