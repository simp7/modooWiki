package model

import (
	"errors"
	"fmt"
	"strings"
)

type Block struct {
	Title    string
	Content  string
	Parent   *Block
	Children []Block
}

func (b *Block) formatChild(depth int, parentPrefix string) (result string) {

	for i, v := range b.Children {
		prefix := titlePrefix(parentPrefix, i)
		result += v.formatTitle(depth, prefix)
		result += v.formatContent()
		result += v.formatChild(depth+1, prefix)
	}

	return

}

func titlePrefix(parentPrefix string, num int) string {
	return fmt.Sprintf("%s%d.", parentPrefix, num)
}

func (b *Block) formatTitle(depth int, prefix string) string {
	return wrapHeading(depth, prefix+" "+b.Title)
}

func trim(text string) string {
	result := text
	result = strings.TrimPrefix(result, "\"")
	result = strings.TrimSuffix(result, "\"")
	return result
}

func wrapHeading(depth int, text string) string {
	return fmt.Sprintf("<%s>%s</%s>", heading(depth), trim(text), heading(depth))
}

func heading(depth int) (heading string) {

	if depth > 5 {
		heading = "h6"
	} else {
		heading = fmt.Sprintf("h%d", depth+1)
	}

	return

}

func (b *Block) formatContent() string {
	return "<body>" + b.Content + "</body>"
}

func (b *Block) Append(child Block) {
	child.Parent = b
	b.Children = append(b.Children, child)
}

func (b *Block) Insert(child Block, idx int) error {

	child.Parent = b

	if idx < len(b.Children) {
		b.Children = append(b.Children[:idx+1], b.Children[idx:]...)
		b.Children[idx] = child
		return nil
	}

	if idx == len(b.Children) {
		b.Append(child)
		return nil
	}

	return errors.New("index out of range when inserting to children")

}
