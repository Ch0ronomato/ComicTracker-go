package parsers

import (
	"golang.org/x/net/html"
	"ch0ronomato/comictracker/comics"
	"errors"
	"fmt"
	"strings"
)

const IMAGE_COMICS_SRC = "Image Comics"

type ImageComicParser struct {

}

func (icp *ImageComicParser) Name() string {
	return "image_comics"
}

func (icp *ImageComicParser) ComicFromHTML(top *html.Node) (*comics.Comic, error) {
	// uses BFS to search through tree to find the right parent container, then DFS to get the data
	var bfs func(nodes[]*html.Node) (name *html.Node, date *html.Node, image *html.Node)
	bfs = func(nodes []*html.Node) (name *html.Node, date *html.Node, image *html.Node) {
		var frontier []*html.Node
		for _, node := range nodes {
			if node == nil {
				continue
			} else if is_content(node) {
				// depth first search
				name = node.FirstChild
				date = node.LastChild
				if image != nil {
					return
				}
			} else if is_image(node) {
				image = node
				if name != nil && date != nil {
					return
				}
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
			return nil, nil, nil
		}
	}

	name_node, date_node, image_node := bfs([] *html.Node{top})
	if name_node == nil || name_node.FirstChild == nil || name_node.FirstChild.LastChild == nil {
		fmt.Print("Name node has null elements!\n")
		return nil, errors.New("Name node was not parsed. The website structure has probably changed")
	}
	if date_node == nil || date_node.LastChild == nil {
		fmt.Print("Date node has null elements!\n")
		return nil, errors.New("Date node was not parsed. The website structure has probably changed")
	}
	name := name_node.FirstChild.LastChild.Data
	date := date_node.LastChild.Data
	imgPath := ""
	for _, attr := range image_node.FirstChild.Attr {
		if attr.Key == "src" {
			imgPath = attr.Val
		}
	}
	return comics.NewComic(name, []string{"author"}, date, imgPath, 2.32, IMAGE_COMICS_SRC), nil
}

func (icp *ImageComicParser) IsNewRelease(n *html.Node) bool {
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

func is_image(n *html.Node) bool {	if n.Type == html.ElementNode && n.Data == "div" {
	for _, a := range n.Attr {
		if a.Key == "class" && strings.Contains(a.Val, "book__img") {
			return true
		}
	}
}
	return false
}