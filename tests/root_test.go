package tests

import (
	"appdoki-be/app"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var apiURL string

func TestMain(m *testing.M) {
	apiURL = os.Getenv("API_URL")
	m.Run()
}

func TestAPI_Root(t *testing.T) {
	res, err := http.Get(apiURL + "/")
 	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatal("expected 200 status code")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("expected to be able to read response body")
	}

	var response app.HomeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatal("expected to be able to parse response body")
	}
}