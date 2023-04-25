package tibia

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"tibia-stats/slices"
)

type OverviewWorld struct {
	Name string `json:"name"`
}

type worlds struct {
	RegularWorlds []OverviewWorld `json:"regular_worlds"`
}

type worldsResponse struct {
	Worlds worlds `json:"worlds"`
}

type OnlinePlayers struct {
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Vocation string `json:"vocation"`
}

type World struct {
	Name          string          `json:"name"`
	OnlinePlayers []OnlinePlayers `json:"online_players"`
}

type WorldResponse struct {
	World World `json:"world"`
}

type WR struct {
	Worlds WorldResponse `json:"worlds"`
}

func (ac *ApiClient) FetchWorlds() ([]OverviewWorld, error) {
	url := fmt.Sprintf("%s/v3/worlds", ac.baseUrl)

	body, err := ac.get(url)
	if err != nil {
		return nil, err
	}

	data := worldsResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return slices.MapSlice(data.Worlds.RegularWorlds, func(world OverviewWorld) OverviewWorld {
		return OverviewWorld{Name: world.Name}
	}), nil
}

func (ac *ApiClient) FetchOnlinePlayers(world string) ([]OnlinePlayers, error) {
	url := fmt.Sprintf("%s/v3/world/%s", ac.baseUrl, world)

	body, err := ac.get(url)
	var data WR
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Worlds.World.OnlinePlayers, nil
}

func (ac *ApiClient) get(url string) ([]byte, error) {
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
	return io.ReadAll(resp.Body)
}
