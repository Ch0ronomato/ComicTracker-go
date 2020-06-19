package main
import (
	"net/http"
	"fmt"
	"io/ioutil"
	"golang.org/x/net/html"
	"ch0ronomato/comictracker/comics"
	"ch0ronomato/comictracker/parsers"
	"strings"
	"sync"
)

func found_comic(server *parsers.ComicServer, n *html.Node) {
	server.FoundComics <- n
}

func DownloadComicSource(url string, parser_name string, main_wg *sync.WaitGroup) {
	var f func(*html.Node, bool)
	comic_channel := make(chan *comics.Comic)
	defer close(comic_channel)


	res, err := http.Get(url)
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

	// create a comic parser and run through the html object, finding comics
	server, parser, err := parsers.CreateParser(map[string]string {
		"PARSER": parser_name,
	}, comic_channel)
	f = func(n *html.Node, is_release bool) {
		if is_release {
			// we don't want to continue to dig through the tree from here
			if n.Type == html.ElementNode {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					go found_comic(server, c)
					comic := <- comic_channel
					if comic == nil {
						// fmt.Printf("Comic parser %s has skipped comic\n", parser_name)
					} else {
						fmt.Printf("Finished here with %s, %s (picture %s) from %s\n", comic.GetTitle(), comic.GetPublished(), comic.GetImgPath(), comic.GetSource())
					}
				}
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, parser.IsNewRelease(c))
		}
	}
	f(doc, false)
	main_wg.Done()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go DownloadComicSource("https://www.imagecomics.com/comics/series","imagecomicsparser", &wg)
	go DownloadComicSource("http://valiantentertainment.com/events/", "valiantentertainmentparser", &wg)
	wg.Wait()
	fmt.Print("Here")
}