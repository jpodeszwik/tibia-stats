package tibia

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestHighscoreClient_FetchAllHighscores(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !strings.HasSuffix(req.RequestURI, "/1") {
			rw.Write([]byte(`{"highscores":{"highscore_list":[]}}`))
		} else if strings.Contains(req.RequestURI, "knight") {
			rw.Write([]byte(`{"highscores":{"highscore_list":[{"name":"Knight Name","value":123456789}]}}`))
		} else if strings.Contains(req.RequestURI, "paladin") {
			rw.Write([]byte(`{"highscores":{"highscore_list":[{"name":"Paladin Name", "value":23456789}]}}`))
		} else if strings.Contains(req.RequestURI, "sorcerer") {
			rw.Write([]byte(`{"highscores":{"highscore_list":[{"name":"Sorcerer Name", "value": 3456789}]}}`))
		} else if strings.Contains(req.RequestURI, "druid") {
			rw.Write([]byte(`{"highscores":{"highscore_list":[{"name":"Druid Name","value":456789}]}}`))
		}
	}))
	defer server.Close()

	hc := &ApiClient{
		httpClient: server.Client(),
		baseUrl:    server.URL,
	}

	response, err := hc.FetchAllHighscores("Peloria", Exp)
	if err != nil {
		t.Errorf("Got error response %s", err)
	}

	expectedResponse := []HighscoreResponse{
		{Name: "Knight Name", Value: 123456789},
		{Name: "Paladin Name", Value: 23456789},
		{Name: "Sorcerer Name", Value: 3456789},
		{Name: "Druid Name", Value: 456789}}

	sortByExp(response)
	sortByExp(expectedResponse)

	if !reflect.DeepEqual(expectedResponse, response) {
		t.Errorf("Response: %v does not match expected: %v", response, expectedResponse)
	}
}

func sortByExp(hr []HighscoreResponse) {
	sort.Slice(hr, func(i, j int) bool {
		return hr[i].Value < hr[j].Value
	})
}

func TestHighscoreClient_FetchHighscore(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"highscores":{"highscore_list":[{"name":"Knight Name","value":123456789}]}}`))
	}))
	defer server.Close()

	hc := &ApiClient{
		httpClient: server.Client(),
		baseUrl:    server.URL,
	}

	response, err := hc.FetchHighscore("Peloria", Knight, Exp, 1)
	if err != nil {
		t.Errorf("Got error response %s", err)
	}

	expectedResponse := []HighscoreResponse{{Name: "Knight Name", Value: 123456789}}
	if !reflect.DeepEqual(expectedResponse, response) {
		t.Errorf("Response: %v does not match expected: %v", response, expectedResponse)
	}
}
