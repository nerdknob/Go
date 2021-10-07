package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

//Create a struct that will be used to unmarshal the json received from the API call
type Time struct {
	Time_24  string `json:"time_24"`
	Time_12  string `json:"time_12"`
	Date     string `json:"date"`
	Timezone string `json:"timezone"`
}

func main() {
	//Format input location parameter
	var location string
	var cityName string
	if len(os.Args) == 3 {
		location = strings.ToLower(os.Args[1] + "+" + os.Args[2])
		cityName = strings.ToLower(os.Args[1] + " " + os.Args[2])
	} else {
		location = strings.ToLower(os.Args[1])
		cityName = strings.ToLower(os.Args[1])
	}

	//Set timeout for http request
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	//Sign up for a free dev account at https://ipgeolocation.io/signup.html to get your API key
	apiKey := "{{YourApiKeyGoesHere}}"
	url := "https://api.ipgeolocation.io/timezone?apiKey=" + apiKey + "&location=" + location

	//Get timezone information for specified location
	response, err := client.Get(url)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	//Extract json from request body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	//Convert the returned json into a struct so that we can work with the data
	var time Time
	json.Unmarshal(body, &time)

	//Print location time data to console
	fmt.Println("Location: " + strings.Title(cityName))
	fmt.Println("12 Hour:  " + time.Time_12)
	fmt.Println("24 Hour:  " + time.Time_24)
	fmt.Println("Date:     " + time.Date)
	fmt.Println("Timezone: " + time.Timezone)
}
