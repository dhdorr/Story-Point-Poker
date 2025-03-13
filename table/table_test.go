package table

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_GenerateTableSession(t *testing.T) {
	test, err := GenerateTableSession(url.Values{
		"tableID":        []string{"test"},
		"passcode":       []string{"test"},
		"cardLayout":     []string{"test"},
		"username":       []string{"test"},
		"numCards":       []string{"1"},
		"numRounds":      []string{"2"},
		"roundTimeLimit": []string{"3"},
		"playerMax":      []string{"4"},
	})
	if err != nil {
		t.Errorf("Invalid form keys or values: %s", err)
	}
	fmt.Printf("TEST: %v \n", test)
}
