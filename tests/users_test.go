package tests

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAPI_ListUsers(t *testing.T) {
	t.SkipNow()
	res, err := http.Get("http://localhost:4040/")
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

	//var statusInfo QueryResult
	//err = json.Unmarshal(jsonData, &statusInfo)
	//if err != nil {
	//	return nil, e.WithStack(err)
	//}

	t.Log(string(body))
	t.Log(body)
}