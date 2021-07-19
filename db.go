package main

import "modoowiki/model"

type DB interface {
	Close()
	GetPage(title string) (model.Page, error)
	InitPage(title string) error
	DeletePage(title string) error
	SetContent(page model.Page) error
}
