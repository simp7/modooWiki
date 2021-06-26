package main

import (
	"github.com/kataras/golog"
	"modoowiki/model/config"
)

func main() {

	level := golog.DebugLevel

	logger := golog.New()
	logger.Level = level

	conf, err := config.Get()
	if err != nil {
		logger.Fatal(err)
	}

	w, err := Init(conf, level)
	if err != nil {
		logger.Fatal(err)
	}

	err = w.Start()
	if err != nil {
		logger.Info(err)
	}

	defer w.Close()

}
