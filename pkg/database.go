package pkg

import (
	"YadroProject/types"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Storage interface {
	Save() error
	Add(types.Comics, int) error
}

type DB struct {
	cfg  Config
	data types.Data
}

func NewDB(cfg Config) (*DB, error) {

	flag, err := FileIsExist(cfg.Xkcd.DbFile)
	if err != nil {
		return nil, err
	}
	if !flag {
		return &DB{
			cfg: cfg,
			data: types.Data{
				Comics: make(map[int]types.Comics),
			},
		}, nil
	}

	file, err := os.Open(cfg.Xkcd.DbFile)
	defer file.Close()
	if err != nil {
		return nil, fmt.Errorf("error open file \"%s\": %w", cfg.Xkcd.DbFile, err)
	}
	jdb := &DB{
		cfg: cfg,
		data: types.Data{
			Comics: make(map[int]types.Comics),
		},
	}
	if err = json.NewDecoder(file).Decode(&jdb.data.Comics); err != nil {
		return nil, fmt.Errorf("error decode json from \"%s\": %w", cfg.Xkcd.DbFile, err)
	}

	return jdb, nil
}

func (db *DB) Add(comics types.Comics, id int) error {
	if _, ok := db.data.Comics[id]; ok {
		return errors.New(fmt.Sprintf("Comics with id %d already exitst", id))
	}
	db.data.Comics[id] = comics
	return nil
}

func (d *DB) Save() error {

	file, err := os.OpenFile(d.cfg.Xkcd.DbFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Errorf("Some kind of error in Openingfile %v", err)
		return err
	}

	defer file.Close()
	if err = json.NewEncoder(file).Encode(d.data.Comics); err != nil {
		return err
	}
	return nil
}

func FileIsExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
