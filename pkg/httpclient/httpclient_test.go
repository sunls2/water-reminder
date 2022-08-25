package httpclient

import (
	"testing"
)

const (
	checkURL = "https://sunls.me/check"
)

func TestGet(t *testing.T) {
	resp, err := Get(checkURL)
	if err != nil {
		t.Fatal(err)
	}
	var plain string
	if plain = resp.PlainText(); len(plain) == 0 {
		t.Fatal("response body is empty")
	}
	t.Log(plain)
}

func TestPost(t *testing.T) {
	resp, err := Post(checkURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	var plain string
	if plain = resp.PlainText(); len(plain) == 0 {
		t.Fatal("response body is empty")
	}
	t.Log(plain)
}
