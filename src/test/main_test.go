package test

import (
	"testing"
	"golang.org/x/net/html"
	"fmt"
	"io/ioutil"
	"comics"
	"strings"
	"sync"
)
func TestComicMakingPipeline(t *testing.T) {
	names := [...]string{"The Black Monday Murders #6", "Crosswind #1", "Descender, Vol. 4: Orbital Mechanics TP", "Eclipse #8", "The Few #6 (Of 6)", "God Country #6", "Grrl Scouts: Magic Socks #2 (Of 6)", "Head Lopper #6", "Horizon #12", "I Hate Fairyland #13", "Invincible #137", "The Old Guard #5", "Plastic #3 (Of 5)", "Royal City #4", "September Mourning, Vol. 1", "Shirtless Bear-Fighter! #1 (Of 5)"}
	dates := [...]string{"June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017", "June 21, 2017"}
	var comic_channel chan *comics.Comic
	comic_channel = make(chan *comics.Comic)

	// defer close(comic_channel)
	raw, err := ioutil.ReadFile("../../data/All Series _ Image Comics.htm")
	doc, err := html.Parse(strings.NewReader(string(raw)))

	if err != nil {
		t.Error("Couldn't read html file")
	}

	server := comics.NewComicServer(comic_channel)
	var wg sync.WaitGroup
	var f func(n *html.Node, is_release bool, wg sync.WaitGroup) int
	f = func(n *html.Node, is_release bool, wg sync.WaitGroup) int {
		seen := 0
		if is_release {
			// we don't want to continue to dig through the tree from here
			if n.Type == html.ElementNode && n.Data == "div" {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					// fmt.Printf("Seeing comic book\n")
					seen += 1
					wg.Add(1)
					go func() {
						defer wg.Done()
						server.FoundComics <- n
					}()
				}
				fmt.Printf("Finishing processing comics, saw %d", seen)
			}
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				result := f(c, comics.IsNewRelease(n), wg)
				if seen < result {
					seen = result
				}
			}
		}
		return seen
	}
	seen := f(doc, false, wg)
	wg.Wait()

	// now validate all the comics parsed
	comic_count := 0
	var good_name bool = true
	var good_date bool = true
	for comic := range comic_channel {
		comic_count += 1
		for _, name := range names {
			good_name = good_name && bool(strings.Compare(name, comic.GetTitle()) == 0)
		}

		for _, date := range dates {
			good_date = good_date && bool(strings.Compare(date, comic.GetPublished()) == 0)
		}
		if comic_count == seen {
			break
		}
	}
	if !good_name || !good_date {
		t.Errorf("Comic did not have correct name or date")
	}

	fmt.Print("done waiting")
}