module github.com/ch0ronomato/comictracker 

go 1.14

require golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
require ch0ronomato/comictracker/parsers v0.0.0
replace ch0ronomato/comictracker/parsers v0.0.0 => /go/src/parsers
require ch0ronomato/comictracker/comics v0.0.0
replace ch0ronomato/comictracker/comics v0.0.0 => /go/src/comics
