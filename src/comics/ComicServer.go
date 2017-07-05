package comics

import (
	"golang.org/x/net/html"
	"fmt"
)

type ComicServer struct {
	FoundComics chan *html.Node
	Comics []Comic
	ResponseChan chan *Comic
}

func NewComicServer(ResponseChan chan *Comic) *ComicServer {
	server := &ComicServer {
		FoundComics: make(chan *html.Node),
		Comics: make([]Comic, 0, 50),
		ResponseChan: ResponseChan,
	}

	go server.loop()
	return server
}

func (c *ComicServer) loop() {
	for comic_html := range c.FoundComics {
		fmt.Printf("Queing comic\n")
		if comic_html != nil {
			c.ResponseChan <- ComicFromHTML(comic_html)
		}
	}
}