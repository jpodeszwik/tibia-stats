package tibia

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGuildClient_FetchGuild(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"guilds":{"guild":{"name":"Guild Name","members":[{"name":"Member Name"}]}}}`))
	}))
	defer server.Close()

	gc := &ApiClient{
		httpClient: server.Client(),
		baseUrl:    server.URL,
	}

	response, err := gc.FetchGuild("Guild Name")
	if err != nil {
		t.Errorf("Got error response %s", err)
	}

	expectedResponse := &GuildResponse{Name: "Guild Name", Members: []GuildMemberResponse{{Name: "Member Name"}}}
	if !reflect.DeepEqual(expectedResponse, response) {
		t.Errorf("Response: %v does not match expected: %v", response, expectedResponse)
	}
}

func TestApiClient_FetchGuilds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"guilds":{"active":[{"name":"Guild Name"}]}}`))
	}))
	defer server.Close()

	gc := &ApiClient{
		httpClient: server.Client(),
		baseUrl:    server.URL,
	}

	response, err := gc.FetchGuilds("World Name")
	if err != nil {
		t.Errorf("Got error response %s", err)
	}

	expectedResponse := []OverviewGuild{{Name: "Guild Name"}}
	if !reflect.DeepEqual(expectedResponse, response) {
		t.Errorf("Response: %v does not match expected: %v", response, expectedResponse)
	}
}
