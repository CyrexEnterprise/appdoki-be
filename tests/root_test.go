package tests

import (
	"appdoki-be/app"
	"appdoki-be/tests/seeder"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var apiURL string
var seeds *seeder.Seeder

func TestMain(m *testing.M) {
	apiURL = os.Getenv("API_URL")
	seeds = seeder.NewSeeder(os.Getenv("DB_URI"))
	seeds.TruncateAll()
	m.Run()
}

func TestAPI_Root(t *testing.T) {
	t.Run("/", testRoot)
	t.Run("/users", testUsers)
	//t.Run("/beers", testBeers)
}

func testRoot(t *testing.T) {
	t.Run("expect GET / to respond with 200 and API info", func (t *testing.T) {
		res, err := http.Get(apiURL + "/")
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

		var response app.HomeResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal("expected to be able to parse response body")
		}
	})
}