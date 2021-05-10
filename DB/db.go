package DB

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/iterator"
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
	sa := option.WithCredentialsFile("../DB/service-account.json")
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

// RetrieveAll Returns a map of all documents in the parameter collection
func RetrieveAll(collection string) map[string]interface{}{
	data := make(map[string]interface{})
	iter := Client.Collection(collection).Documents(Ctx) // Loop through all entries in param collection
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		// Add the document to the data map with the documents id as the map key
		data[doc.Ref.ID] = doc.Data()
	}

	// Returns the data
	return data
}

// Test function for adding test data to the database & retrieving data
func Test(i string){
	data := make(map[string]interface{})
	data["Location"] = "Oslo"
	data["n"] = "s"
	// Add a value to the map
	data["v"] = "v"
	// Delete a value from the map
	delete(data,"v")

	AddToDatabase("Users", i+"dwa", data)

	// Retrieves and prints all users
	data = RetrieveAll("Users")
	fmt.Println(data)

	fmt.Println("SINGLE ENTRY")

	// Retrieves and prints a single user entry
	data = RetrieveFromDatabase("Users", "181069578170793984")
	fmt.Println(data)
}