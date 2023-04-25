package tibia

import (
	"encoding/json"
	"fmt"
)

type CharactersGuild struct {
	Name string `json:"name"`
}

type Character struct {
	Name  string          `json:"name"`
	Guild CharactersGuild `json:"guild"`
}

type Deaths struct {
	Time   string `json:"time"`
	Level  int    `json:"level"`
	Reason string `json:"reason"`
}

type Characters struct {
	Deaths    []Deaths  `json:"deaths"`
	Character Character `json:"character"`
}

type CharacterResponse struct {
	Characters Characters `json:"characters"`
}

func (ac *ApiClient) FetchCharacter(name string) (*Characters, error) {
	url := fmt.Sprintf("%s/v3/character/%s", ac.baseUrl, name)

	body, err := ac.get(url)
	var data CharacterResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data.Characters, nil
}
