package tibia

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
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
	Value int
}

type highscore struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type highscores struct {
	HighscoreList []highscore `json:"highscore_list"`
}

type highscoresResponse struct {
	Highscores highscores `json:"highscores"`
}

type HighscoreClient struct {
	httpClient *http.Client
	baseUrl    string
}

func newHttpClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 10
	transport.MaxConnsPerHost = 10
	transport.MaxIdleConnsPerHost = 10

	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
}

func NewHighscoreClient() *HighscoreClient {
	return &HighscoreClient{
		httpClient: newHttpClient(),
		baseUrl:    "https://api.tibiadata.com",
	}
}

func (hc *HighscoreClient) FetchHighscore(world string, profession Profession, highscoreType HighscoreType) ([]HighscoreResponse, error) {
	url := hc.baseUrl + "/v3/highscores/" + world + "/" + string(highscoreType) + "/" + string(profession)
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

	highscoreData := highscoresResponse{}
	err = json.Unmarshal(body, &highscoreData)
	if err != nil {
		return nil, err
	}

	return mapSlice(highscoreData.Highscores.HighscoreList, mapHighscore), nil
}

type highscoreResult struct {
	profession Profession
	response   []HighscoreResponse
	err        error
}

func (hc *HighscoreClient) FetchAllHighscores(world string, highscoreType HighscoreType) ([]HighscoreResponse, error) {
	retChannel := make(chan highscoreResult, 4)
	for _, profession := range AllProfessions {
		go func(profession2 Profession) {
			res, err := hc.FetchHighscore(world, profession2, highscoreType)
			retChannel <- highscoreResult{profession2, res, err}
		}(profession)
	}

	var err error
	var result []HighscoreResponse

	for i := 0; i < len(AllProfessions); i++ {
		highscoreResult := <-retChannel
		if highscoreResult.err != nil {
			log.Printf("Failed to fetch profession %s", highscoreResult.profession)
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

func mapSlice[IN any, OUT any](input []IN, mapper func(IN) OUT) []OUT {
	res := make([]OUT, 0)

	for _, value := range input {
		res = append(res, mapper(value))
	}

	return res
}
