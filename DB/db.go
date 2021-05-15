package DB

import (
	"assistant/utils"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
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
func DatabaseInit() {
	// Connects to the firestore database
	Ctx = context.Background()
	sa := option.WithCredentialsFile("DB/service-account.json")
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
func AddToDatabase(collection string, entry string, values map[string]interface{}) {
	fmt.Println(collection, entry)

	_, err := Client.Collection(collection).Doc(entry).Set(Ctx,
		values)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func RemoveFromDatabase(collection string, entry string) error{
	_, err := Client.Collection(collection).Doc(entry).Delete(Ctx)
	if err != nil {
		return errors.New("could not delete entry from database")
	}
	return nil
}

// RetrieveFromDatabase Function that retrieves the map data of a document in a collection in firestore
func RetrieveFromDatabase(collection string, entry string) (map[string]interface{}, error) {
	data, err := Client.Collection(collection).Doc(entry).Get(Ctx)
	if err != nil {
		fmt.Println(err)
		// Creates an empty output
		output := make(map[string]interface{})
		return output, errors.New("Could not find your storage, are you sure you have initialized a storage for this service?")
	}

	// Returns the data
	return data.Data(), nil
}

// RetrieveAll Returns a map of all documents in the parameter collection
func RetrieveAll(collection string) map[string]interface{} {
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
func Test(i string) {
	data := make(map[string]interface{})
	data["Location"] = "Oslo"
	data["n"] = "s"
	// Add a value to the map
	data["v"] = "v"
	// Delete a value from the map
	delete(data, "v")

	AddToDatabase("Users", i+"dwa", data)

	// Retrieves and prints all users
	data = RetrieveAll("Users")
	fmt.Println(data)

	fmt.Println("SINGLE ENTRY")

	// Retrieves and prints a single user entry
	data, err := RetrieveFromDatabase("Users", "181069578170793984")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

// PostNewsWebHook Adds a Webhook to the news collection
func PostNewsWebHook(userID, hookID string, values utils.NewsWebhook) (string, error) {
	_, err := Client.Collection("news").
		Doc(userID).
		Collection("webhooks").
		Doc(hookID).
		Set(Ctx, values)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return hookID, nil	// Returning the hookID serves as a confirmation of registration
}

// GetNewsWebHooksByUser Gets all webhooks registered by a discord user using their discord id
func GetNewsWebHooksByUser(filter, userID string) (utils.NewsWebhooks, error) {
	//var data []utils.NewsWebhook
	iter := Client.Collection("news").
		Doc(userID).
		Collection("webhooks").
		Documents(Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		reqType, err := doc.DataAt("requestType")
		if err != nil {
			break
		}
		switch filter {
		case "trending":
			if reqType.(string) == "trending" {
				// Add the document to the data map with the documents id as the map key
				//data = append(data, doc.DataTo(&data))
			}
		case "category":
			if reqType.(string) == "category" {
				// Add the document to the data map with the documents id as the map key
				//data[doc.Ref.ID] = doc.Data()
			}
		case "search":
			if reqType.(string) == "category" {
				// Add the document to the data map with the documents id as the map key
				//data[doc.Ref.ID] = doc.Data()
			}
		default:	// Gets all request types webhooks
			// Add the document to the data map with the documents id as the map key
			//data[doc.Ref.ID] = doc.Data()
		}
	}

	return nil, nil
}

// GetNewsWebHooksByID Gets a webhook registered by a discord user using the hook id
func GetNewsWebHooksByID(userID, hookID string) (utils.NewsWebhook, error) {
	 var data utils.NewsWebhook
	doc, err := Client.Collection("news").
		Doc(userID).
		Collection("webhooks").
		Doc(hookID).
		Get(Ctx)
	if err != nil {
		fmt.Println(err)
		return data, err
	}
	err = doc.DataTo(&data)
	if err != nil {
		fmt.Println(err)
		return data, err
	}
	return data, nil
}

// DeleteNewsWebHookByID Deletes a Webhook registered by a discord user using the hook id
// Returns String with confirmation of webhook deletion
func DeleteNewsWebHookByID(userID, hookID string) (string, error) {
	_, err := Client.Collection("news").
		Doc(userID).
		Collection("webhooks").
		Doc(hookID).
		Delete(Ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return fmt.Sprintf("Webhook %s has been deleted", hookID), nil
}

// DeleteNewsWebHooksByUser Deletes all webhooks registered by a discord user using their discord id
func DeleteNewsWebHooksByUser(userID string) (string, error) {
	col := Client.Collection("news").Doc(userID).Collection("webhooks")
	// Delete all documents for webhooks collection
	err := deleteCollection(Ctx, Client, col, 500) // 500 batch size is the limit
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return fmt.Sprintf("All Webhooks for %s have been deleted", userID), nil
}

// deleteCollection Retrieved from firestore documentation
// Takes values for context, client, a reference to the collection and batchsize
func deleteCollection(ctx context.Context, client *firestore.Client,
	ref *firestore.CollectionRef, batchSize int) error {
	for {
		// Get a batch of documents
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to a
		// WriteBatch.
		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}
