package mpdqlib

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/antchfx/xmlquery"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/thoas/go-funk"
	"github.com/zencoder/go-dash/mpd"
)

// TODO: There must be a better way.  I was trying to find a 'simple' way to
// find the nodes that match a query and return all parents with ONLY the
// matching node, removing siblings that don't match, or formatting only the
// nodes that match the query with some other formatting.  IANASWE and so far
// it's taken more time than I'm prepared to invest at the moment so I will
// revisit later and leave this half solution here for now.

// TODO: Review this approach, https://github.com/cbsinteractive/bakery/blob/master/filters/dash.go.

// Query returns nodes and parents that have an attribute matching the query
// parameters
func Query(manifest *mpd.MPD, name, op, value string) {
	query(manifest, name, op, value)
}

func query(manifest *mpd.MPD, name, op, value string) {
	validOps := []string{"=", "!=", "<", "<=", ">", ">="}
	if !funk.Contains(validOps, op) {
		log.Fatalf("invalid operation: %v\n", op)
	}
	manifestString, err := manifest.WriteToString()
	if err != nil {
		log.Fatalf("unable to convert manifest to string: %v\n", err)
	}
	xml, err := xmlquery.Parse(strings.NewReader(manifestString))
	if err != nil {
		log.Fatalf("unable to parse manifest: %v\n", err)
	}
	queryString := strings.Join([]string{"//*[@", name, op, value, "]"}, "")
	finder := xmlquery.Find(xml, queryString)
	if len(finder) == 0 {
		fmt.Println("no matching values")
		return
	}
	for _, n := range finder {
		pretty := xmlfmt.FormatXML(n.OutputXML(true), "", "  ")
		err = quick.Highlight(os.Stdout, pretty, "xml", "terminal16m", "pygments")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("")
	}
}
