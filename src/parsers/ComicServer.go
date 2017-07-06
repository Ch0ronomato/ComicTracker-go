package parsers

import (
	"golang.org/x/net/html"
	"comics"
)

type ComicServer struct {
	FoundComics chan *html.Node
	Comics []comics.Comic
	ResponseChan chan *comics.Comic
}

func NewComicServer(ResponseChan chan *comics.Comic, parser ComicParser) *ComicServer {
	server := &ComicServer {
		FoundComics: make(chan *html.Node),
		Comics: make([]comics.Comic, 0, 50),
		ResponseChan: ResponseChan,
	}

	go server.loop(parser)
	return server
}

func (c *ComicServer) loop(parser ComicParser) {
	for comic_html := range c.FoundComics {
		if comic_html != nil {
			resp, err := parser.ComicFromHTML(comic_html)
			if err != nil {
				c.ResponseChan <- nil
			} else {
				c.ResponseChan <- resp
			}
		}
	}
}