package urlhandler

import (
	"testing"
)

/* func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}*/

func TestEmptyURL(t *testing.T) {

	result1 := isValidURL("")
	result2 := isValidURL("http:///")

	if result1 == true || result2 == true {
		t.Errorf("failed, expected false for empty string but found true")
	}
}

func TestValidURL(t *testing.T) {

	result := isValidURL("http://google.com/")

	if result == false {
		t.Errorf("failed, expected true for 'http://google.com/' but found false")
	}
}

func TestInvalidURL(t *testing.T) {

	result := isValidURL("absdd")
	if result == true {
		t.Errorf("failed, expected false for empty 'absdd' but found true")
	}
	result = isValidURL("http://google/")
	if result == true {
		t.Errorf("failed, expected false for 'http://google/' but found true")
	}
}
