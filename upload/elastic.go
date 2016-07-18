package upload

import (
	. "github.com/lee-woodridge/whats-that-card/card"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	elasticURL      = "http://search-whats-that-card-fja3htmk3m3ibexgzpdbd4kw7e.us-west-1.es.amazonaws.com/hs"
	mappingFilename = "./mapping.json"
)

// createCardsIndex sets up the correct ElasticSearch indexes to quickly
// perform the text searches I want to do on the card data.
func createCardsIndex() error {
	// Setup http client.
	client := &http.Client{}

	// Delete if already exists.
	req, err := http.NewRequest("DELETE", elasticURL, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	file, e := ioutil.ReadFile(mappingFilename)
	if e != nil {
		return e
	}
	buf := bytes.NewBuffer(file)
	req, err = http.NewRequest("PUT", elasticURL, buf)
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	return err
}

func insertCardToElastic(client *http.Client, card Card) error {
	// Marshal json and create byte buffer for PUT.
	b, err := json.Marshal(card)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)
	// Submit PUT to elastic using the cardId as unique key.
	req, err := http.NewRequest("PUT", elasticURL+"/cards/"+card.CardId, buf)
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	return err
}

// InsertCardsToElastic uploads the cards from the API into my ElasticSearch instance.
func InsertCardsToElastic(cards map[string][]Card) error {
	// Setup indexes to insert the cards first.
	if err := createCardsIndex(); err != nil {
		return err
	}

	// Setup http client.
	client := &http.Client{}

	for _, val := range cards { // card set -> list of cards
		for _, card := range val { // each card
			err := insertCardToElastic(client, card)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
