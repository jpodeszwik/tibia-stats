package tibia

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type GuildMember struct {
	Name string `json:"name"`
}

type Guild struct {
	Name    string        `json:"name"`
	Members []GuildMember `json:"members"`
}

type Guilds struct {
	Guild Guild `json:"guild"`
}

type guildResponse struct {
	Guilds Guilds `json:"guilds"`
}

type GuildClient struct {
	httpClient *http.Client
	baseUrl    string
}

func NewGuildClient() *GuildClient {
	return &GuildClient{
		httpClient: newHttpClient(),
		baseUrl:    "https://api.tibiadata.com",
	}
}

func (hc *GuildClient) FetchGuild(guild string) (*Guild, error) {
	url := fmt.Sprintf("%s/v3/guild/%s", hc.baseUrl, guild)
	log.Printf("Fetching: %s", url)

	resp, err := hc.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Printf("Failed to close body %s", err)
			}
		}()
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	guildResponse := guildResponse{}
	err = json.Unmarshal(body, &guildResponse)
	if err != nil {
		return nil, err
	}

	return &guildResponse.Guilds.Guild, nil
}
