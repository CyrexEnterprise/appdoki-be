package tests

import (
	"appdoki-be/app/repositories"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAPI_ListBeerTransfers(t *testing.T) {
	t.Run("expect GET /beers to have default pagination", func(t *testing.T) {
		res, err := http.Get(apiURL + "/beers")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 status code, got %d", res.StatusCode)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal("expected to be able to read response body")
		}

		var response []repositories.BeerTransferFeedItem
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal("expected to be able to parse response body")
		}

		if len(response) != 20 {
			t.Fatalf("expected to be get 20 transfers, got %d instead", len(response))
		}
	})

	t.Run("expect GET /beers to accept partial pagination params", func(t *testing.T) {
		res, err := http.Get(apiURL + "/beers?limit=1")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 status code, got %d", res.StatusCode)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal("expected to be able to read response body")
		}

		var response []repositories.BeerTransferFeedItem
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal("expected to be able to parse response body")
		}

		if len(response) != 1 {
			t.Fatalf("expected to be get 1 transfer, got %d instead", len(response))
		}
	})

	t.Run("expect GET /beers to accept pagination params", func(t *testing.T) {
		res, err := http.Get(apiURL + "/beers?limit=1&givenAt=2020-01-01&op=gt")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 status code, got %d", res.StatusCode)
		}
	})
}
