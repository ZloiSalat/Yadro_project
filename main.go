package main

import (
	"YadroProject/pkg"
	"log"
)

func main() {

	cmd := pkg.NewFlagParser()
	stem := pkg.NewStemmer()
	cfg := pkg.NewConfig()
	store := pkg.NewData(cfg.Xkcd.Source)

	a, err := pkg.NewAPIClient(store, *stem, *cfg, cmd)
	if err != nil {
		log.Panicln(err)
	}
	a.Run()

}
