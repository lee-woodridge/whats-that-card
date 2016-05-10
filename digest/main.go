package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	hearthstoneAPI = "https://omgvamp-hearthstone-v1.p.mashape.com/cards"
	elasticURL     = "http://localhost:9200"
)

func getCardsFromAPI() ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", hearthstoneAPI, nil)
	if err != nil {
		return nil, err
	}
	// Get mashape api key from file (as not to check in to source control).
	file, err := ioutil.ReadFile("./mashape.key")
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Mashape-Key", string(file))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	return bytes, err
}

func main() {
	file, e := ioutil.ReadFile("./cards.json")
	// file, e := getCardsFromAPI()
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		return
	}
	var cards map[string][]interface{}
	json.Unmarshal(file, &cards)

	for key, _ := range cards {
		fmt.Printf("%s\n", key)
	}
}
