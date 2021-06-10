package main

import (
	"github.com/kataras/golog"
	"modoowiki/model/config"
)

func main() {

	conf, err := config.Get()
	if err != nil {
		golog.Fatal(err)
	}

	w, err := Init(conf)
	if err != nil {
		golog.Fatal(err)
	}

	err = w.Start()
	if err != nil {
		w.Logger().Info(err)
	}

	defer w.End()

}
