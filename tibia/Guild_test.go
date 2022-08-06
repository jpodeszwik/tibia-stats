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

	gc := &GuildClient{
		httpClient: server.Client(),
		baseUrl:    server.URL,
	}

	response, err := gc.FetchGuild("Guild Name")
	if err != nil {
		t.Errorf("Got error response %s", err)
	}

	expectedResponse := &Guild{Name: "Guild Name", Members: []GuildMember{{Name: "Member Name"}}}
	if !reflect.DeepEqual(expectedResponse, response) {
		t.Errorf("Response: %v does not match expected: %v", response, expectedResponse)
	}
}
