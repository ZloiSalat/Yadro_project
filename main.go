package main

import (
	"YadroProject/pkg"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cmd := pkg.NewFlagParser()
	stem := pkg.NewStemmer()
	cfg := pkg.NewConfig("config.yaml")
	store, err := pkg.NewDB(cfg)
	if err != nil {
		fmt.Errorf("Store init error", err)
	}

	a, err := pkg.NewAPIClient(store, *stem, cfg, cmd)
	if err != nil {
		log.Panicln(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		// Сохраняем данные в базу данных перед завершением
		if err := store.Save(); err != nil {
			log.Printf("Error saving data to database: %v", err)
		}
		os.Exit(1)
	}()

	a.Run()

}
