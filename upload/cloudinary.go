package upload

import (
	. "github.com/lee-woodridge/whats-that-card/card"

	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	cloudinaryURL = "https://api.cloudinary.com/v1_1/elusive/image/upload"
)

func UploadImageToCloudinary(card Card) error {
	if card.Img == "" {
		return nil
	}
	unixTimestamp := time.Now().Unix()
	privKeyFile, err := ioutil.ReadFile("./cloudinary.private.key")
	if err != nil {
		return err
	}
	pubKeyFile, err := ioutil.ReadFile("./cloudinary.public.key")
	if err != nil {
		return err
	}
	shaStr := fmt.Sprintf("public_id=%s&timestamp=%d%s",
		card.CardId, unixTimestamp, privKeyFile)
	shaVal := sha1.Sum([]byte(shaStr))
	apiStr := fmt.Sprintf("file=%s&public_id=%s&timestamp=%d&api_key=%s&signature=%x",
		card.Img, card.CardId, unixTimestamp, pubKeyFile, shaVal)

	client := &http.Client{}
	req, err := http.NewRequest("POST", cloudinaryURL, bytes.NewBufferString(apiStr))
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	return err
}
