package comics

import (
	"golang.org/x/net/html"
	"strings"
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
	// uses BFS to search through tree to find the right parent container, then DFS to get the data
	bfs := []*html.Node{node}
	for _, node := range bfs {
		if (is_content(node)) {

		}
	}

	return NewComic("title", []string{"author"}, "10/10/10", 2.32)
}

func (c *Comic) Save() bool {
	// @todo: Implement
	return false
}

func (c *Comic) GetTitle() string {
	return c.title
}

func (c *Comic) GetWriters() []string {
	return c.writers
}

func (c *Comic) GetPublished() string {
	return c.published
}

func (c *Comic) GetPrice() float32 {
	return c.price
}

func IsNewRelease(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "upcoming_releases") {
				return true
			}
		}
	}
	return false
}

func is_content(n *html.Node) bool {
	return false
}