package urlhandler

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetHomePage(t *testing.T) {
	resp, err := http.Get("/urls/3")
	if err != nil {
		fmt.Println(resp)
	} else {
		t.Errorf("failed, expected false for empty 'absdd' but found true")
	}

}
