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

func (p Page) String() (result string) {
	result = wrapHeading(0, p.Root.Title)
	return p.Root.FormatChild(0, "")
}
