package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

//getAndDecodeURL recieves an url that you want to retrieve information from
//and decodes it into the struct that is passed along
func GetAndDecodeURL(url string, decodedJSON *interface{}) error {
	//Send a GET request to the url
	response, err := http.Get(url)
	if err != nil {
		return errors.New("possible HTTP error, or too many redirects")
	}
	//Check if the statuscode of the request is OK
	if response.StatusCode != http.StatusOK {
		return errors.New("something unexpected happened." + url + " do !diag to check the status of our services")
	}
	//Request was successful, decode it into struct
	err = json.NewDecoder(response.Body).Decode(&decodedJSON)
	if err != nil {
		return err
	}
	//No errors, return nil
	return nil
}
