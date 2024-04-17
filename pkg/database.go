package pkg

import (
	"YadroProject/types"
	"encoding/json"
	"os"
)

type Storage interface {
	Save(map[int]types.Comics) error
}

type DB struct {
	cfg   Config
	types types.Data
}

func NewDB(cfg Config) *DB {
	return &DB{
		cfg: cfg,
	}
}

func (d *DB) Save(data map[int]types.Comics) error {

	file, err := os.OpenFile(d.cfg.Xkcd.DbFile, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	if err = json.NewEncoder(file).Encode(data); err != nil {
		return err
	}
	return nil
}
