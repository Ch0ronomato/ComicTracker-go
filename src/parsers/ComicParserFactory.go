package parsers

import (
	"log"
	"errors"
	"fmt"
	"strings"
	"ch0ronomato/comictracker/comics"
	"time"
	"strconv"
)

type ComicParserFactory func (conf map[string]string) (ComicParser, error)
var parserFactories = make(map[string]ComicParserFactory)


func NewImageComicParser(_ map[string]string) (ComicParser, error) {
	// no conf needed
	return &ImageComicParser{

	}, nil
}

func NewValiantEntertainmentParser(_ map[string]string) (ComicParser, error) {
	year := time.Now().Year()
	strYear := strconv.Itoa(year)
	return &ValiantEntertainmentParser {
		 url: "http://valiantentertainment.com/events/",
		 yearId: strYear,
	}, nil
}

func Register(name string, factory ComicParserFactory) {
	if factory == nil {
		log.Panicf("Parser factory %s doesn't exist man!", name)
	}

	_, registered := parserFactories[name]
	if registered {
		log.Printf("%s is already registered, ignoring.", name)
	}
	parserFactories[name] = factory
}

func CreateParser(conf map[string]string, responseChan chan *comics.Comic) (*ComicServer, ComicParser, error) {
	parserName := conf["PARSER"]

	parserFactory, ok := parserFactories[parserName]
	if !ok {
		availableParsers := make([]string, len(parserFactories))
		for k, _ := range parserFactories {
			availableParsers = append(availableParsers, k)
		}
		return nil, nil, errors.New(fmt.Sprintf("Invalid parser. Must be one of %s",
			strings.Join(availableParsers, ",")))
	}
	parser, err := parserFactory(conf)
	return NewComicServer(responseChan, parser), parser, err
}

func init() {
	Register("imagecomicsparser", NewImageComicParser)
	Register("valiantentertainmentparser", NewValiantEntertainmentParser)
}