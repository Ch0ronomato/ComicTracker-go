package main
import (
	"net/http"
	"fmt"
	"io/ioutil"
	"golang.org/x/net/html"
	"strings"
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

func main () {
	res, err := http.Get("https://imagecomics.com/comics/series")
	if err != nil {
		fmt.Printf("I died: %s\n", err)
	}

	// read all tokens
	var f func(*html.Node, bool)
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Couldn't read html")
	}
	doc, err := html.Parse(strings.NewReader(string(raw)))

	if err != nil {
		fmt.Printf("Couldn't parse out the html")
	}

	f = func(n *html.Node, is_release bool) {
		if is_release {
			// we don't want to continue to dig through the tree from here
			if n.Type == html.ElementNode && n.Data == "div" {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					fmt.Printf("Seeing comic book")
				}
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, is_new_release(n))
		}
	}
	f(doc, false)
}
