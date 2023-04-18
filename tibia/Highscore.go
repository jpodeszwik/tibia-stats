package tibia

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"tibia-stats/slices"
)

type Profession string

const (
	Knight   Profession = "knight"
	Paladin             = "paladin"
	Druid               = "druid"
	Sorcerer            = "sorcerer"
)

var AllProfessions = [4]Profession{Knight, Paladin, Druid, Sorcerer}

type HighscoreType string

const (
	Exp HighscoreType = "exp"
)

type HighscoreResponse struct {
	Name  string
	Value int64
}

type highscore struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type highscores struct {
	HighscoreList []highscore `json:"highscore_list"`
}

type highscoresResponse struct {
	Highscores highscores `json:"highscores"`
}

func (ac *ApiClient) FetchHighscore(world string, profession Profession, highscoreType HighscoreType, page int) ([]HighscoreResponse, error) {
	url := fmt.Sprintf("%s/v3/highscores/%s/%s/%s/%d", ac.baseUrl, world, highscoreType, profession, page)

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

	highscoreData := highscoresResponse{}
	err = json.Unmarshal(body, &highscoreData)
	if err != nil {
		return nil, err
	}

	return slices.MapSlice(highscoreData.Highscores.HighscoreList, mapHighscore), nil
}

type highscoreResult struct {
	profession Profession
	response   []HighscoreResponse
	err        error
}

func (ac *ApiClient) FetchAllHighscores(world string, highscoreType HighscoreType) ([]HighscoreResponse, error) {
	retChannel := make(chan highscoreResult, 4*20)
	for _, profession := range AllProfessions {
		for page := 1; page <= 20; page++ {
			go func(profession2 Profession, page2 int) {
				res, err := ac.FetchHighscore(world, profession2, highscoreType, page2)
				retChannel <- highscoreResult{profession2, res, err}
			}(profession, page)
		}
	}

	var err error
	var result []HighscoreResponse

	for i := 0; i < 4*20; i++ {
		highscoreResult := <-retChannel
		if highscoreResult.err != nil {
			log.Printf("Failed to fetch profession %s, %v", highscoreResult.profession, highscoreResult.err)
			err = highscoreResult.err
		}
		result = append(result, highscoreResult.response...)
	}

	return result, err
}

func mapHighscore(h highscore) HighscoreResponse {
	return HighscoreResponse{
		Name:  h.Name,
		Value: h.Value,
	}
}
