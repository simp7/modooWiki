package model

import "time"

type Page struct {
	Root         Block
	LastModified time.Time
}

func NewPage(key string) Page {
	block := Block{Title: key, Parent: nil}
	return Page{block, time.Now()}
}

func (p Page) Key() string {
	return p.Root.Title
}

func (p Page) String() (result string) {
	result = wrapHeading(0, p.Key())
	result += p.Root.formatChild(0, "")
	return
}
