package tibia

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"tibia-stats/utils/slices"
)

type GuildMemberResponse struct {
	Name  string
	Level int
}

type GuildResponse struct {
	Name    string
	Members []GuildMemberResponse
}

type OverviewGuild struct {
	Name string
}

type guildMember struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type guild struct {
	Name    string        `json:"name"`
	Members []guildMember `json:"members"`
}

type Guilds struct {
	Guild guild `json:"guild"`
}

type guildResponse struct {
	Guilds Guilds `json:"guilds"`
}

type overviewGuild struct {
	Name string `json:"name"`
}

type overviewGuilds struct {
	Active    []overviewGuild `json:"active"`
	Formation []overviewGuild `json:"formation"`
}

type guildsOverviewResponse struct {
	Guilds overviewGuilds `json:"guilds"`
}

func (ac *ApiClient) FetchGuild(guildName string) (*GuildResponse, error) {
	url := fmt.Sprintf("%s/v3/guild/%s", ac.baseUrl, guildName)

	resp, err := ac.httpClient.Get(url)
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

	guild := guildResponse.Guilds.Guild

	return &GuildResponse{
		Name: guild.Name,
		Members: slices.MapSlice(guild.Members, func(in guildMember) GuildMemberResponse {
			return GuildMemberResponse{
				Name:  in.Name,
				Level: in.Level,
			}
		}),
	}, nil
}

func (ac *ApiClient) FetchGuilds(world string) ([]OverviewGuild, error) {
	url := fmt.Sprintf("%s/v3/guilds/%s", ac.baseUrl, world)

	resp, err := ac.httpClient.Get(url)
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

	overviewGuildsResponse := guildsOverviewResponse{}
	err = json.Unmarshal(body, &overviewGuildsResponse)
	if err != nil {
		return nil, err
	}

	return slices.MapSlice(overviewGuildsResponse.Guilds.Active, func(in overviewGuild) OverviewGuild {
		return OverviewGuild{Name: in.Name}
	}), nil
}
