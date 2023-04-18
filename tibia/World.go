package tibia

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"tibia-stats/slices"
)

type WorldResponse struct {
	Name string
}

type world struct {
	Name string `json:"name"`
}

type worlds struct {
	RegularWorlds []world `json:"regular_worlds"`
}

type worldsResponse struct {
	Worlds worlds `json:"worlds"`
}

func (ac *ApiClient) FetchWorlds() ([]WorldResponse, error) {
	url := fmt.Sprintf("%s/v3/worlds", ac.baseUrl)

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

	data := worldsResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return slices.MapSlice(data.Worlds.RegularWorlds, func(world world) WorldResponse {
		return WorldResponse{Name: world.Name}
	}), nil
}
