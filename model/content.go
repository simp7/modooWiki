package model

import "fmt"

type Block struct {
	Title    string
	Content  string
	Parent   *Block
	Children []Block
}

func (b *Block) FormatChild(depth int, prefix string) (result string) {

	for i, v := range b.Children {
		result += b.formatTitle(depth, prefix+string(i+1), i)
		result += b.Content
		result += v.FormatChild(depth+1, prefix+string(i+1))
	}

	return

}

func (b *Block) titlePrefix(parentPrefix string, idx int) string {
	return parentPrefix + string(idx+1) + "."
}

func (b *Block) formatTitle(depth int, prefix string, idx int) string {
	return wrapHeading(depth, prefix+string(idx+1)+".")
}

func wrapHeading(depth int, text string) string {
	return fmt.Sprintf("<%s>%s</%s>", heading(depth), text, heading(depth))
}

func heading(depth int) (heading string) {

	heading = "h"

	if depth > 5 {
		heading += "6"
	} else {
		heading += string(depth)
	}
	return

}

func (b *Block) Insert(child Block) {
	child.Parent = b
	b.Children = append(b.Children, child)
}
