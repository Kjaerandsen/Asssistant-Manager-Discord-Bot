package DB

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
)

// Ctx Firebase context and client used by Firestore functions throughout the program.
var Ctx context.Context
var Client *firestore.Client
// Collection name in Firestore

// DatabaseInit Initiates the database connection
func DatabaseInit(){
	// Connects to the firestore database
	Ctx = context.Background()
	sa := option.WithCredentialsFile("./DB/service-account.json")
	app, err := firebase.NewApp(Ctx, nil, sa)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	Client, err = app.Firestore(Ctx)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
}

// AddToDatabase Function that adds an document with an interface to a collection in firestore
func AddToDatabase(collection string, entry string, values map[string]interface{}){
	fmt.Println(collection,entry)

	_, err := Client.Collection(collection).Doc(entry).Set(Ctx,
		values)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// RetrieveFromDatabase Function that retrieves the map data of a document in a collection in firestore
func RetrieveFromDatabase(collection string, entry string) map[string]interface{}{
	data, err := Client.Collection(collection).Doc(entry).Get(Ctx)
	if err != nil{
		fmt.Println(err)
		// Creates an empty output
		output := make(map[string]interface{})
		return output
	}

	// Returns the data
	return data.Data()
}

// Test function for testing adding data to the database
func Test(i string){
	data := make(map[string]interface{})
	data["Location"] = "Oslo"
	data["n"] = "s"
	data["v"] = "v"

	AddToDatabase("Users", i, data)
}