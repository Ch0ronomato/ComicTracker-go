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

func TestImageComicFromHtml(t *testing.T) {
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

func TestValiantComicFromHtml(t *testing.T) {
	raw_html := `<li class="white float-wrap  July"><div class="date float-left col-12 standard-40-bold no-float-mobile">Jul 26</div><div class="float-left col-85 no-float-mobile"><h3 class="standard-40-bold"><a class="white" href="http://valiantentertainment.com/comics faith-and-the-future-force/faith-and-the-future-force-1-of-4/">FAITH AND THE FUTURE FORCE #1</a></h3></div></li>`
	comic_html, err := html.Parse(strings.NewReader(raw_html))
	if err != nil {
		log.Fatal("error!")
	}
	// need to find the right element....
	var f func(*html.Node) *html.Node
	f = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "li" {
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
		"PARSER": "valiantentertainmentparser",
	}, make(chan *comics.Comic))
	if err != nil {
		t.Errorf("Got an error %s", err)
	}
	comic, err := parser.ComicFromHTML(comic_html)
	if err != nil {
		t.Errorf("Got an error %s", err)
	}

	if comic.GetTitle() == "" {
		t.Error("No comic name parsed")
	}

	if comic.GetTitle() != "FAITH AND THE FUTURE FORCE #1" {
		t.Errorf("Parsed comic name incorrectly")
	}
}

func TestImageComicMakingPipeline(t *testing.T) {
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

func TestValiantEntertainmentMakingPipeline(t *testing.T) {
	names := [...]string{"RAPTURE #3","BLOODSHOT\xe2\x80\x99S DAY OFF #1","BRITANNIA: WE WHO ARE ABOUT TO DIE #4","GENERATION ZERO VOL. 2: HEROSCAPE TPB","HARBINGER RENEGADE #5 (ALL-NEW ARC! ALL-NEW JUMPING-ON POINT! \xe2\x80\x9cMASSACRE\xe2\x80\x9d \xe2\x80\x93 PART ONE)","NINJAK VOL. 6: THE SEVEN BLADES OF MASTER DARQUE TPB","X-O MANOWAR #4 (NEW ARC! \xe2\x80\x9cGENERAL\xe2\x80\x9d \xe2\x80\x93 PART 1)","SECRET WEAPONS #1","BRITANNIA: WE WHO ARE ABOUT TO DIE #2","RAI: THE HISTORY OF THE VALIANT UNIVERSE #1","DIVINITY III: HEROES OF THE GLORIOUS STALINVERSE TPB","X-O MANOWAR (2017) VOL. 1: SOLDIER TPB","FAITH #12","FAITH #11","X-O MANOWAR #2","DIVINITY III: STALINVERSE TPB","NINJAK #27","ETERNAL WARRIOR: AWAKENING #1","X-O MANOWAR DELUXE EDITION BOOK 5 HC","X-O MANOWAR #5","RAPTURE #2","X-O MANOWAR #3","A&A: THE ADVENTURES OF ARCHER & ARMSTRONG VOL. 3 \xe2\x80\x93 ANDROMEDA ESTRANGED TPB","FAITH #9","DIVINITY III: ESCAPE FROM GULAG 396 #1","WRATH OF THE ETERNAL WARRIOR VOL. 3: A DEAL WITH A DEVIL TPB","GENERATION ZERO #8","FAITH: HOLLYWOOD & VINE DELUXE EDITION HC","4001 A.D. DELUXE EDITION HC","FAITH VOL. 3: SUPERSTAR TPB","BLOODSHOT REBORN #0","GENERATION ZERO #9","BRITANNIA: WE WHO ARE ABOUT TO DIE #1","BLOODSHOT U.S.A. TPB","NINJAK #26","BRITANNIA: WE WHO ARE ABOUT TO DIE #3","FAITH #10 (ALL-NEW ARC! \xe2\x80\x9cTHE FAITHLESS\xe2\x80\x9d)","NINJAK VOL. 5: THE FIST & THE STEEL TPB","IMMORTAL BROTHERS: THE TALE OF THE GREEN KNIGHT #1","HARBINGER RENEGADE #4","X-O MANOWAR DELUXE EDITION BOOK 4 HC","HARBINGER RENEGADE VOL. 1: THE JUDGMENT OF SOLOMON TPB","SECRET WEAPONS #2","GENERATION ZERO #7","DIVINITY III: STALINVERSE #4","SAVAGE #4","DIVINITY III: SHADOWMAN & THE BATTLE OF NEW STALINGRAD #1","NINJAK #24","BRITANNIA TPB","SAVAGE #3","GENERATION ZERO VOL. 1: WE ARE THE FUTURE TPB","FAITH #8","BLOODSHOT U.S.A. #4","A&A: THE ADVENTURES OF ARCHER AND ARMSTRONG #12","DIVINITY III: ARIC, SON OF THE REVOLUTION #1","DIVINITY III: STALINVERSE #2","X-O MANOWAR #1","NINJAK #25","RAPTURE #1","HARBINGER RENEGADE #3","GENERATION ZERO #6 (ALL-NEW ARC! \xe2\x80\x9cHEROSCAPE\xe2\x80\x9d)","NINJAK #23 (NEW ARC! \xe2\x80\x9cTHE SEVEN BLADES OF MASTER DARQUE\xe2\x80\x9d)","X-O MANOWAR VOL. 13: SUCCESSION & OTHER TALES TPB","FAITH #7","DIVINITY III: STALINVERSE #3","FAITH AND THE FUTURE FORCE #1","A&A: THE ADVENTURES OF ARCHER & ARMSTRONG #11"}
	dates := [...]string {"Jul 19","Jul 5","Jul 19","Jul 19","Jul 12","Jul 12","Jun 28","Jun 28","May 17","Jun 14","Jun 14","Jun 28","Jun 7","May 3","Apr 26","May 3","May 17","May 10","Jul 26","Jul 26","Jun 21","May 24","May 24","Mar 1","Mar 15","Mar 15","Mar 29","Mar 1","Mar 29","May 10","Mar 22","Apr 19","Apr 26","Apr 26","Apr 19","Jun 21","Apr 5","Feb 8","Apr 12","Feb 22","Feb 22","Apr 19","Jul 19","Feb 15","Mar 29","Feb 15","Feb 8","Feb 8","Feb 15","Jan 25","Jan 25","Feb 1","Jan 25","Feb 1","Jan 18","Jan 25","Mar 22","Mar 29","May 24","Jan 18","Jan 18","Jan 11","Jan 11","Jan 4","Feb 22","Jul 26","Jan 4"}
	var comic_channel chan *comics.Comic
	comic_channel = make(chan *comics.Comic)

	// defer close(comic_channel)
	raw, err := ioutil.ReadFile("../../data/Events _ Valiant Entertainment.htm")
	if err != nil {
		t.Errorf("The file %s wasn't found!\n", "../../data/Events _ Valiant Entertainment.htm")
	}
	cleaned_html := strings.Replace(strings.Replace(string(raw), "\n", "", -1), "\t", "", -1)
	doc, err := html.Parse(strings.NewReader(cleaned_html))

	if err != nil {
		t.Error("Couldn't read html file")
	}

	server, parser, err := parsers.CreateParser(map[string]string {
		"PARSER": "valiantentertainmentparser",
	}, comic_channel)
	var wg sync.WaitGroup
	var f func(n *html.Node, is_release bool, wg *sync.WaitGroup) int
	f = func(n *html.Node, is_release bool, wg *sync.WaitGroup) int {
		seen := 0
		if is_release {
			// we don't want to continue to dig through the tree from here
			if n.Type == html.ElementNode && n.Data == "ul" {
				for comic_html := n.FirstChild; comic_html != nil; comic_html = comic_html.NextSibling {
					if comic_html.Type == html.ElementNode {
						// fmt.Printf("Seeing comic book\n")
						seen += 1
						go func(to_parse *html.Node) {
							server.FoundComics <- to_parse
						} (comic_html)
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