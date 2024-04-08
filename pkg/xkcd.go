package pkg

import (
	"YadroProject/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type APIClient struct {
	store Storage
	stem  Stemmer
	cfg   Config
	cmd   *Parser
}

func NewAPIClient(store Storage, stem Stemmer, cfg Config, cmd *Parser) (*APIClient, error) {
	return &APIClient{
		store: store,
		stem:  stem,
		cfg:   cfg,
		cmd:   cmd,
	}, nil
}

func (a *APIClient) GetComicsByID(id int) (types.Num, error) {
	d := types.Num{}
	url := fmt.Sprintf(a.cfg.Xkcd.Source+"/"+strconv.Itoa(id)+"/info.0.json", nil)
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	resp, err := c.Get(url)
	if resp.StatusCode == http.StatusNotFound {
		return types.Num{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return types.Num{}, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return types.Num{}, err
	}

	return d, nil
}

func (a *APIClient) Run() error {
	o, n, err := a.cmd.ParseFlagOandN(a.cfg.Xkcd.DbSize)
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		comics, err := a.GetComicsByID(n)
		if err != nil {
			return fmt.Errorf("error get comics from database: %w", err)
		}
		if err = a.stem.stemmedWords(comics); err != nil {

		}
	}

}
