package upload

import (
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

// ImageToCloudinary takes a cards image link from the hearthstone API
// and creates a copy hosted on my cloudinary account.
func ImageToCloudinary(url, alias string) error {
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
		alias, unixTimestamp, privKeyFile)
	shaVal := sha1.Sum([]byte(shaStr))
	apiStr := fmt.Sprintf("file=%s&public_id=%s&timestamp=%d&api_key=%s&signature=%x",
		url, alias, unixTimestamp, pubKeyFile, shaVal)

	client := &http.Client{}
	req, err := http.NewRequest("POST", cloudinaryURL, bytes.NewBufferString(apiStr))
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	return err
}
