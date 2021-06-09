package config

import "embed"

type Config struct {
	DB
	Web
}

//go:embed file
var instance embed.FS

func Get() (config Config, err error) {

	config.DB, err = getDB()
	config.Web, err = getWeb()

	return

}
