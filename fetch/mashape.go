package fetch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	mashapeKeyEnv = "MASHAPE_KEY"
)

func makeMashapeCall(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Get mashape api key from env variable (as not to check in to source control).
	// Can push env variable to heroku with:
	//		heroku config:add MASHAPE_KEY="$MASHAPE_KEY"
	mashapeKey := os.Getenv(mashapeKeyEnv)
	if len(mashapeKey) == 0 {
		return nil, fmt.Errorf("%s must be set", mashapeKeyEnv)
	}
	req.Header.Set("X-Mashape-Key", mashapeKey)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
