package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"simpleWiki/model"
)

type DB struct {
	Address              string `json:"address"`
	Port                 int    `json:"port"`
	Name                 string `json:"name"`
	model.PageCollection `json:"pageCollection"`
}

func getDB() (DB, error) {

	db := DB{}

	data, err := instance.ReadFile("file/db.json")
	if err != nil && err != io.EOF {
		return db, err
	}

	err = json.Unmarshal(data, &db)
	return db, err

}

func (d DB) Path() string {
	return fmt.Sprintf("mongodb://%s:%d", d.Address, d.Port)
}
