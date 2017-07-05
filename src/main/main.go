package main
import (
	"net/http"
	"fmt"
	"io/ioutil"
	"golang.org/x/net/html"
	"strings"
	"comics"
)

func found_comic(server *comics.ComicServer, n *html.Node) {
	server.FoundComics <- n
}

func main () {
	var f func(*html.Node, bool)
	comic_channel := make(chan *comics.Comic)
	defer close(comic_channel)


	res, err := http.Get("https://imagecomics.com/comics/series")
	defer res.Body.Close()

	if err != nil {
		fmt.Printf("I died: %s\n", err)
	}

	// read body
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Couldn't read html")
	}
	stripped_html := strings.Replace(strings.Replace(string(raw), "\n", "", -1), "\t", "", -1)
	doc, err := html.Parse(strings.NewReader(stripped_html))

	if err != nil {
		fmt.Printf("Couldn't parse out the html\n")
	}

	// run through the html object, finding comics
	server := comics.NewComicServer(comic_channel)
	f = func(n *html.Node, is_release bool) {
		if is_release {
			// we don't want to continue to dig through the tree from here
			if n.Type == html.ElementNode && n.Data == "div" {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					fmt.Printf("Seeing comic book\n")
					go found_comic(server, n)
					comic := <- comic_channel
					fmt.Printf("Finished here with %s, %s\n", comic.GetTitle(), comic.GetPublished())
				}
				fmt.Printf("Finishing processing comic\n")
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, comics.IsNewRelease(n))
		}
	}
	f(doc, false)
}
