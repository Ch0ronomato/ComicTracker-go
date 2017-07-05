package parsers

import (
	"golang.org/x/net/html"
	"comics"
)

type ComicParser interface {
	ComicFromHTML(*html.Node) (*comics.Comic, error)
	IsNewRelease(*html.Node) bool
	Name() string
}