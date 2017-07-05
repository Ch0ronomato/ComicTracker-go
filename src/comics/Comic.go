package comics

import (
	"golang.org/x/net/html"
	"strings"
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

func ComicFromHTML(top *html.Node) *Comic {
	// uses BFS to search through tree to find the right parent container, then DFS to get the data
	var bfs func(nodes[]*html.Node) (name *html.Node, date *html.Node)
	bfs = func(nodes []*html.Node) (name *html.Node, date *html.Node) {
		var frontier []*html.Node
		for _, node := range nodes {
			if node == nil {
				continue
			} else if is_content(node) {
				// depth first search
				name = node.FirstChild
				date = node.LastChild
				return
			} else {
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Data == "div" {
						frontier = append(frontier, c)
					}
				}
			}
		}

		if len(frontier) > 0 {
			return bfs(frontier)
		} else {
			return nil, nil
		}
	}

	name_node, date_node := bfs([] *html.Node{top})
	if name_node == nil || name_node.FirstChild == nil || name_node.FirstChild.LastChild == nil {
		fmt.Print("Name node has null elements!\n")
		return nil
	}
	if date_node == nil || date_node.LastChild == nil {
		fmt.Print("Date node has null elements!\n")
		return nil
	}
	name := name_node.FirstChild.LastChild.Data
	date := date_node.LastChild.Data
	return NewComic(name, []string{"author"}, date, 2.32)
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
	// check the parent for the right container class.
	if n.Parent == nil {
		return false
	}

	parent := n.Parent
	is_parent := false
	if parent.Type == html.ElementNode && parent.Data == "div" {
		for _, a := range parent.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "upcoming_releases") {
				is_parent = true
			}
		}
	}

	// check the current node for the grid class
	is_child := false
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "grid") {
				is_child = true
			}
		}
	}
	return is_parent && is_child
}

func is_content(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "book__content") {
				return true
			}
		}
	}
	return false
}