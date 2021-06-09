package main

import (
	"github.com/kataras/golog"
	"log"
	"simpleWiki/model/config"
)

func main() {

	conf, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	w, err := Init(conf)
	if err != nil {
		log.Fatal(err)
	}

	w.Logger().Level = golog.DebugLevel
	err = w.Start()
	if err != nil {
		w.Logger().Info(err)
	}

	defer w.End()

}
