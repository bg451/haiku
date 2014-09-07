package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestUrlValidate(t *testing.T) {
	url := "https://www.youtube.com/watch?v=MPmObvuOMYA"
	ret, err := validateUrl(url)
	if err != nil {
		t.Error(err.Error())
	}
	if ret != "//www.youtube.com/embed/MPmObvuOMYA" {
		t.Error(fmt.Sprintf("urlValidate: Expected //youtube.com/embed/MPmObvuOMYA, got %s", ret))
	}

	//Check for failure
	url = "https://www.bing.com"
	_, err = validateUrl(url)
	if !strings.Contains(err.Error(), "Host not youtube") {
		t.Error("urlValidate: Allowed non youtube host")
	}
}
