package fetch

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSelectiveUnmarshalling(t *testing.T) {
	s := `{"patch":"1","something_else":0}`
	var card apiInfo
	err := json.Unmarshal([]byte(s), &card)
	if err != nil {
		fmt.Printf("%#v", err)
		t.Fail()
	}
}
