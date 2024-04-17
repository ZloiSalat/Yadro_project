package pkg

import (
	"YadroProject/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type APIClient struct {
	store Storage
	stem  Stemmer
	cfg   Config
	cmd   *Parser
	types types.Data
}

func NewAPIClient(store Storage, stem Stemmer, cfg Config, cmd *Parser) (*APIClient, error) {
	return &APIClient{
		store: store,
		stem:  stem,
		cfg:   cfg,
		cmd:   cmd,
		types: types.Data{
			Comics: make(map[int]types.Comics),
		},
	}, nil
}

func (a *APIClient) GetComicsByID(id int) (types.Comics, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d/info.0.json", a.cfg.Xkcd.Source, id))
	if err != nil {
		return types.Comics{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Comics{}, fmt.Errorf("status code not OK: %d", resp.StatusCode)
	}

	var d types.Comics
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return types.Comics{}, err
	}

	return d, nil
}

func (a *APIClient) Run() error {

	c, _ := a.cmd.ParseFlag()
	if !c {
		fmt.Printf("Config: DBSize %v, Source %v, DBFile %v, End_comics %v ", a.cfg.Xkcd.DbSize, a.cfg.Xkcd.Source, a.cfg.Xkcd.DbFile, a.cfg.Xkcd.End_comics)
		return nil
	}
	for i := 100; i < 100+a.cfg.Xkcd.End_comics; i++ {
		data := make(map[int]interface{})

		if err := json.Unmarshal([]byte(a.cfg.Xkcd.DbFile), &data); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return err
		}

		if _, ok := data[i]; !ok {
			continue
		}
		comics, err := a.GetComicsByID(i)
		if err != nil {
			return fmt.Errorf("error get comics from database: %w", err)
		}
		fmt.Println(comics)

		words := a.stem.stemmedWords(comics)

		for i, word := range words {
			words[i] = ` ` + word + ` `
		}

		result := strings.Join(words, ",")

		d := types.Comics{
			Keywords: result,
			URL:      comics.URL,
		}

		a.types.Comics[i] = d
		a.store.Save(a.types.Comics)
	}
	return nil
}
