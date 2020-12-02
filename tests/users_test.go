package tests

import (
	"appdoki-be/app/repositories"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAPI_ListUsers(t *testing.T) {
	res, err := http.Get(apiURL + "/users")
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

	var response []repositories.User
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatal("expected to be able to parse response body")
	}

	if len(response) != 5 {
		t.Fatal("expected to be get 5 users")
	}
}
