package images

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/lee-woodridge/whats-that-card/upload"
)

type Images struct {
	Config Config `json:"config"`
	Cards  Status `json:"cards"`
}

type Config struct {
	Base, Version string
}

type Status struct {
	Prerelease map[string]string `json:"pre"`
	Released   map[string]string `json:"rel"`
}

const (
	imagesFile       = "images.json"
	imagesJSONUrl    = "https://raw.githubusercontent.com/schmich/hearthstone-card-images/master/images.json"
	imageTemplateURL = "https://raw.githubusercontent.com/schmich/hearthstone-card-images/master/rel/%s.png"
)

func GetImages() {
	// Download images.json from github to use as current.
	imagesJSONResp, err := http.Get(imagesJSONUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting images json: %#v\n", err)
		return
	}
	defer imagesJSONResp.Body.Close()
	imagesJSONBytes, err := ioutil.ReadAll(imagesJSONResp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading images json: %#v\n", err)
		return
	}
	var images Images
	err = json.Unmarshal(imagesJSONBytes, &images)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling images json: %#v\n", err)
		return
	}
	// Check if we should fetch new images, using a locally cached images.json.
	if _, err := os.Stat(imagesFile); !os.IsNotExist(err) {
		imagesFile, err := ioutil.ReadFile(imagesFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening local images json: %#v\n", err)
			return
		}
		var localImages Images
		err = json.Unmarshal(imagesFile, &localImages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling local images json: %#v\n", err)
			return
		}

		if localImages.Config.Version == images.Config.Version {
			fmt.Printf("Local images up to date. Local version: %s, remove version: %s\n",
				localImages.Config.Version, images.Config.Version)
			return
		}
		fmt.Printf("Local images not up to date. Local version: %s, remove version: %s\n",
			localImages.Config.Version, images.Config.Version)
	} else {
		fmt.Printf("No local cache found\n")
	}
	fmt.Printf("Uploading new images...\n")

	// Use the remote images.json to link new image files to cloudinary
	for id := range images.Cards.Released {
		upload.ImageToCloudinary(fmt.Sprintf(imageTemplateURL, id), id)
	}

	// Save remote images.json to use for caching in future.
	fmt.Printf("Writing version file for caching\n")
	err = ioutil.WriteFile(imagesFile, imagesJSONBytes, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing version file: %#v", err)
	}
}
