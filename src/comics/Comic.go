package comics

import (
	"golang.org/x/net/html"
	"fmt"
)

type Comic struct {
	title string
	writers []string
	published string
	price float32
}

func NewComic(title string, writers []string, published string, price float32) *Comic {
	return &Comic {
		title: title,
		writers: writers,
		published: published,
		price: price,
	}
}

func ComicFromHTML(node *html.Node) *Comic {
	fmt.Printf("I got called yay\n")
	return NewComic("title", []string{"author"}, "10/10/10", 2.32)
}

func (c *Comic) Save() bool {
	// @todo: Implement
	return false
}

func (c *Comic) ToString() string {
	return "title: " + c.title + ", published: " + c.published
}
