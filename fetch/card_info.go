package fetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type apiInfo struct {
	Version string `json:"patch"`
}

const (
	versionFile = "info.json"
	cardsFile   = "cards.json"
	versionAPI  = "https://omgvamp-hearthstone-v1.p.mashape.com/info"
	cardsAPI    = "https://omgvamp-hearthstone-v1.p.mashape.com/cards"
)

// Fetch fetches new cards from the Mashape API, if there is a new version than
// the one we have locally.
func Fetch() {
	versionBytes, err := makeMashapeCall(versionAPI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting remote version: %#v", err)
		return
	}
	// Check if we should fetch a new version, using a locally cached version file.
	if _, err := os.Stat(versionFile); !os.IsNotExist(err) {
		// Get version from the cached version file.
		file, err := ioutil.ReadFile(versionFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading info: %#v", err)
			return
		}
		var localInfo apiInfo
		err = json.Unmarshal(file, &localInfo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling info: %#v", err)
			return
		}

		// Get the current version from the API.
		var remoteInfo apiInfo
		err = json.Unmarshal(versionBytes, &remoteInfo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling remote version: %#v", err)
			return
		}
		if remoteInfo.Version == localInfo.Version {
			fmt.Printf("Local cards up to date. Local version: %s, remove version: %s\n",
				remoteInfo.Version, localInfo.Version)
			return
		}
		fmt.Printf("Local cards not up to date. Local version: %s, remove version: %s\n",
			remoteInfo.Version, localInfo.Version)
	} else {
		fmt.Printf("No local cache found\n")
	}
	fmt.Printf("Fetching new cards...\n")
	cardBytes, err := makeMashapeCall(cardsAPI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting remote cards: %#v", err)
		return
	}
	err = ioutil.WriteFile(cardsFile, cardBytes, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing cards file: %#v", err)
	}
	fmt.Printf("Writing version file for caching\n")
	err = ioutil.WriteFile(versionFile, versionBytes, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing version file: %#v", err)
	}
	fmt.Print("Success!")
}
