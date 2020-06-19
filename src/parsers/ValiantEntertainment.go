package parsers

import (
	"ch0ronomato/comictracker/comics"
	"golang.org/x/net/html"
	"strings"
	"errors"
	"time"
	"fmt"
	"net/http"
	"io/ioutil"
)

const VALIANT_ENTERTAINMENT_SRC = "Valiant Entertainment"

type ValiantEntertainmentParser struct {
	url string
	yearId string
}

func (vep *ValiantEntertainmentParser) Name() string {
	return "valiant_entertainment"
}

func (vep *ValiantEntertainmentParser) ComicFromHTML(top *html.Node) (*comics.Comic, error) {
	_, month, _ := time.Now().Date()
	month_str := month.String()
	is_new := false
	if top.Data != "li" {
		return nil, errors.New("Incorrect element handed to ComicFromHTML")
	}

	for _,a := range top.Attr {
		if a.Key == "class" && strings.Contains(a.Val, month_str) {
			is_new = true
		}
	}

	if !is_new {
		return nil, errors.New("Comic wasn't new")
	}

	date_node := top.FirstChild.FirstChild
	name_node := top.LastChild.FirstChild.FirstChild.FirstChild // cause that's not a bad idea....
	if date_node == nil {
		return nil, errors.New("Date node was nil from parser or not where it was expected")
	}

	if name_node == nil {
		return nil, errors.New("Name node was nil from parser or not where it was expected")
	}
	imgPath, err := getImagePath(top)
	if err != nil {
		fmt.Printf("No image for %s", name_node.Data)
	}
	return comics.NewComic(name_node.Data, []string{"Author"}, strings.Trim(date_node.Data, " "), imgPath,2.32, VALIANT_ENTERTAINMENT_SRC), nil
}

func (vep *ValiantEntertainmentParser) IsNewRelease(src *html.Node) bool {
	// check the parent node
	idValue := "comic-" + vep.yearId
	parentDistance := 3
	start := src
	hasParent := false
	isULChild := src.Data == "ul"
	for parentDistance > 0 {
		if start == nil || start.Parent == nil {
			break
		}
		start = start.Parent
		parentDistance--
	}

	if parentDistance > 0 || start.Data != "div"{
		return false
	}

	for _, a := range start.Attr {
		if a.Key == "id" && strings.Contains(a.Val, idValue) {
			hasParent = true
			break
		}
	}

	return hasParent && isULChild
}

func getImagePath(top *html.Node) (string, error) {
	if top.LastChild == nil || top.LastChild.FirstChild == nil || top.LastChild.FirstChild.FirstChild == nil {
		return "", errors.New("Unexpected HTML format")
	}
	link := ""
	for _, a := range top.LastChild.FirstChild.LastChild.Attr {
		if a.Key == "href" && strings.Contains(a.Val, "http") {
			link = a.Val
		}
	}

	if link == "" {
		return "",errors.New("No image link ")
	}

	// load this page, and get img.
	pageRaw, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer pageRaw.Body.Close()

	rawPageHtml, err := ioutil.ReadAll(pageRaw.Body)
	if err != nil {
		return "", err
	}

	stripped_html := strings.Replace(strings.Replace(string(rawPageHtml), "\n", "", -1), "\t", "", -1)
	doc, err := html.Parse(strings.NewReader(stripped_html))

	if err != nil {
		return "", err
	}

	var bfs func(n *html.Node) (string);
	bfs = func(n *html.Node) (string){
		if n.Data == "picture" {
			for _, a := range n.LastChild.Attr {
				if a.Key == "srcset" {
					return a.Val
				}
			}
			return ""
		} else {
			src := ""
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				src = bfs(child)
				if src != "" {
					break
				}
			}
			return src
		}
	}
	return bfs(doc), nil
}