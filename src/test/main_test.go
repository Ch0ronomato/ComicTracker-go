package test

import (
	"testing"
	"golang.org/x/net/html"
	"fmt"
	"io/ioutil"
	"comics"
	"strings"
	"sync"
	"log"
	"parsers"
)

func TestComicFromHtml(t *testing.T) {
	raw_html := `<div class="grid-cell u-size1of2 u-md-size1of4 u-sm-size1of2"><div class="book"><div class="book__img" style="height: 247px;"><img src="./All Series _ Image Comics_files/Crosswind_01-1.png" alt="Crosswind #1" class="book"></div><div class="book__content"><p class="u-m0"><a href="https://parsers.com/comics/releases/crosswind-1">Crosswind #1</a></p><p class="u-m0">June 21, 2017</p></div></div></div>`
	comic_html, err := html.Parse(strings.NewReader(raw_html))
	if err != nil {
		log.Fatal("error!")
	}
	// need to find the right element....
	var f func(*html.Node) *html.Node
	f = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "div" {
			return n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			ret := f(c)
			if ret != nil {
				return ret
			}
		}
		return nil
	}
	comic_html = f(comic_html)
	_, parser, err := parsers.CreateParser(map[string]string{
		"PARSER": "imagecomicsparser",
	}, make(chan *comics.Comic))
	if err != nil {
		t.Errorf("Got an error %s", err)
	}
	comic, err := parser.ComicFromHTML(comic_html)
	if err != nil {
		t.Errorf("Got an error %s", err)
	}

	if comic.GetTitle() == "" {
		t.Error("Nope")
	}
}

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

	server, parser, err := parsers.CreateParser(map[string]string {
		"PARSER": "imagecomicsparser",
	}, comic_channel)
	var wg sync.WaitGroup
	var f func(n *html.Node, is_release bool, wg *sync.WaitGroup) int
	f = func(n *html.Node, is_release bool, wg *sync.WaitGroup) int {
		seen := 0
		if is_release {
			// we don't want to continue to dig through the tree from here
			if n.Type == html.ElementNode && n.Data == "div" {
				for comic_html := n.FirstChild; comic_html != nil; comic_html = comic_html.NextSibling {
					if comic_html.Type == html.ElementNode && strings.EqualFold(comic_html.Data, "div") {
						// fmt.Printf("Seeing comic book\n")
						seen += 1
						go func(to_parse *html.Node) {
							server.FoundComics <- to_parse
						}(comic_html)
					}
				}
				fmt.Printf("Finishing processing comics, saw %d\n", seen)
			}
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				result := f(c, parser.IsNewRelease(c), wg)
				if seen < result {
					seen = result
				}
			}
		}
		return seen
	}
	seen := f(doc, false, &wg)

	// now validate all the comics parsed
	comic_count := 0
	for comic := range comic_channel {
		fmt.Printf("Eating comic %d\n", comic_count)
		var good_name bool = false
		var good_date bool = false
		comic_count += 1
		if comic == nil {
			// t.Errorf("Comic %d was nil", comic_count - 1)
			continue
		}
		for _, name := range names {
			good_name = good_name || bool(strings.Compare(name, comic.GetTitle()) == 0)
		}

		if !good_name {
			t.Errorf("Comic %d did not have correct title: %s\n", comic_count - 1, comic.GetTitle())
		}


		for _, date := range dates {
			good_date = good_date || bool(strings.Compare(date, comic.GetPublished()) == 0)
		}

		if !good_date {
			t.Errorf("Comic %d did not have correct date: %s\n", comic_count - 1, comic.GetPublished())
		}

		if comic_count == seen {
			break
		}
	}
}