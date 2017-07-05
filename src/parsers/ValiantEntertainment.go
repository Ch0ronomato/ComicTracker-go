package parsers

import (
	"comics"
	"golang.org/x/net/html"
	"strings"
	"errors"
)

const VALIANT_ENTERTAINMENT_SRC = "Valiant Entertainment"

type ValiantEntertainmentParser struct {
	url string
	yearId string
}

func (vep *ValiantEntertainmentParser) Name() string {
	return "valiant_entertainment"
}

func (vep *ValiantEntertainmentParser) ComicFromHTML(top *html.Node) (*comics.Comic, error) {
	if top.Data != "li" {
		return nil, errors.New("Incorrect element handed to ComicFromHTML")
	}
	date_node := top.FirstChild.FirstChild
	name_node := top.LastChild.FirstChild.FirstChild.FirstChild // cause that's not a bad idea....
	if date_node == nil {
		return nil, errors.New("Date node was nil from parser or not where it was expected")
	}

	if name_node == nil {
		return nil, errors.New("Name node was nil from parser or not where it was expected")
	}
	return comics.NewComic(name_node.Data, []string{"Author"}, strings.Trim(date_node.Data, " "), 2.32, VALIANT_ENTERTAINMENT_SRC), nil
}

func (vep *ValiantEntertainmentParser) IsNewRelease(src *html.Node) bool {
	// check the parent node
	idValue := "comic-" + vep.yearId
	parentDistance := 3
	start := src
	hasParent := false
	isULChild := src.Data == "ul"
	for parentDistance > 0 {
		if start == nil || start.Parent == nil {
			break
		}
		start = start.Parent
		parentDistance--
	}

	if parentDistance > 0 || start.Data != "div"{
		return false
	}

	for _, a := range start.Attr {
		if a.Key == "id" && strings.Contains(a.Val, idValue) {
			hasParent = true
			break
		}
	}

	return hasParent && isULChild
}