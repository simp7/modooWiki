package main

import "modoowiki/model"

type DB interface {
	Close()
	GetPage(title string) (model.Page, error)
	InitPage(title string, content string) error
	SetContent(title string, content string) error
}
