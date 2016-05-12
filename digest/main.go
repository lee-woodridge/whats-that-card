package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	// "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Mechanic struct {
	Name string
}

type Card struct {
	Mechanics    []Mechanic // needs it's own type [{"name": "Charge"},{"name": "Divine Shield"}]
	Durability   int
	Locale       string
	Text         string
	HowToGet     string
	ImgGold      string
	Cost         int
	Flavor       string
	PlayerClass  string
	Img          string
	Attack       int
	Health       int
	Type         string
	Collectible  bool
	Faction      string
	InPlayText   string
	Elite        bool
	HowToGetGold string
	CardSet      string
	Name         string
	Artist       string
	Rarity       string
	Race         string
	CardId       string
}

const (
	hearthstoneAPI  = "https://omgvamp-hearthstone-v1.p.mashape.com/cards"
	cloudinaryURL   = "https://api.cloudinary.com/v1_1/elusive/image/upload"
	elasticURL      = "http://localhost:9200"
	mappingFilename = "mapping.json"
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

func uploadImageToCloudinary(card Card) error {
	if card.Img == "" {
		return nil
	}
	unixTimestamp := time.Now().Unix()
	privKeyFile, err := ioutil.ReadFile("./cloudinary.private.key")
	if err != nil {
		return err
	}
	// fmt.Printf("privKey: %s\n", privKeyFile)
	pubKeyFile, err := ioutil.ReadFile("./cloudinary.public.key")
	if err != nil {
		return err
	}
	// fmt.Printf("privKey: %s\n", pubKeyFile)
	shaStr := fmt.Sprintf("public_id=%s&timestamp=%d%s",
		card.CardId, unixTimestamp, privKeyFile)
	// fmt.Printf("shaStr: %s\n", shaStr)
	shaVal := sha1.Sum([]byte(shaStr))
	// fmt.Printf("shaVal: %x\n", shaVal)
	apiStr := fmt.Sprintf("file=%s&public_id=%s&timestamp=%d&api_key=%s&signature=%x",
		card.Img, card.CardId, unixTimestamp, pubKeyFile, shaVal)
	// fmt.Printf("apiStr: %s\n", apiStr)

	client := &http.Client{}
	req, err := http.NewRequest("POST", cloudinaryURL, bytes.NewBufferString(apiStr))
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	return err
}

func createCardsIndex() error {
	file, e := ioutil.ReadFile(mappingFilename)
	if e != nil {
		return e
	}
	// Setup http client.
	client := &http.Client{}
	buf := bytes.NewBuffer(file)
	req, err := http.NewRequest("PUT", elasticURL+"/hs", buf)
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	return err
}

func insertCardToElastic(card Card) error {
	// Setup http client.
	client := &http.Client{}
	// Marshal json and create byte buffer for PUT.
	b, err := json.Marshal(card)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)
	// Submit PUT to elastic using the cardId as unique key.
	req, err := http.NewRequest("PUT", elasticURL+"/hs/cards/"+card.CardId, buf)
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	// fmt.Printf("put %s with resp %s", elasticURL+"/hs/cards/external/"+card.CardId, res)
	return err
}

func main() {
	file, e := ioutil.ReadFile("./cards.json")
	// file, e := getCardsFromAPI()
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		return
	}
	var cards map[string][]Card
	json.Unmarshal(file, &cards)

	for _, val := range cards {
		for _, card := range val {
			err := insertCardToElastic(card)
			if err != nil {
				fmt.Errorf("%s\n", err)
				return
			}
		}
	}
}
