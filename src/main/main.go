package main
import (
	"net/http"
	"fmt"
	"io/ioutil"
	"golang.org/x/net/html"
	"strings"
	"comics"
)

func is_new_release(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "upcoming_releases") {
				return true
			}
		}
	}
	return false
}

func found_comic(server *comics.ComicServer, n *html.Node) {
	server.FoundComics <- n
}

func main () {
	res, err := http.Get("https://imagecomics.com/comics/series")
	if err != nil {
		fmt.Printf("I died: %s\n", err)
	}

	// read all tokens
	var f func(*html.Node, bool)
	raw, err := ioutil.ReadAll(res.Body)
	var comic_channel chan *comics.Comic
	comic_channel = make(chan *comics.Comic)

	defer res.Body.Close()
	defer close(comic_channel)
	
	if err != nil {
		fmt.Printf("Couldn't read html")
	}
	doc, err := html.Parse(strings.NewReader(string(raw)))

	if err != nil {
		fmt.Printf("Couldn't parse out the html")
	}

	server := comics.NewComicServer(comic_channel)
	f = func(n *html.Node, is_release bool) {
		if is_release {
			// we don't want to continue to dig through the tree from here
			if n.Type == html.ElementNode && n.Data == "div" {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					fmt.Printf("Seeing comic book\n")
					go found_comic(server, n)
					comic := <- comic_channel
					fmt.Printf("Finished here with %s\n", comic.ToString())
				}
				fmt.Printf("Finishing processing comic")
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, is_new_release(n))
		}
	}
	f(doc, false)
}
