package main

import (
	"github.com/google/logger"
	"modoowiki/model/config"
	"os"
)

func main() {

	lg := logger.Init("wiki", true, false, os.Stdout)
	lg.SetLevel(logger.Level(1))

	conf, err := config.Get()
	if err != nil {
		logger.Fatal(err)
	}

	w, err := Init(conf, lg)
	if err != nil {
		logger.Fatal(err)
	}

	err = w.Start()
	if err != nil {
		logger.Info(err)
	}

	defer w.Close()

}
