package mpdqlib

import (
	"log"
	"strings"

	"github.com/thoas/go-funk"
	"github.com/zencoder/go-dash/mpd"
)

// This takes a lot of inspiration from CBS Interactive (Thank you!)
// https://github.com/cbsinteractive/bakery/blob/master/filters/dash.go
// https://github.com/cbsinteractive/bakery/blob/master/filters/filter.go

// Query returns a parsed manifest with elements matching query string
func Query(manifest *mpd.MPD, query string) {
	queryMPD(manifest, query, true)
}

func queryMPD(manifest *mpd.MPD, query string, print bool) {
	validParams := []string{"b", "bw", "bandwidth", "c", "codec", "fr", "fps", "l", "lang", "language", "t", "type", "ts", "timescale", "w", "width"}
	validOps := []string{"=", "==", "!=", "<", "<=", ">", ">="}
	queryList := getQueries(query)
	for _, queryString := range queryList {
		query := parseQuery(queryString)
		if !funk.Contains(validParams, query.param) {
			log.Fatalf("invalid query parameter %v\nvalid params: %v\n", query.param, validParams)
		}
		if !funk.Contains(validOps, query.op) {
			log.Fatalf("invalid query operand %v\nvalid operands: %v\n", query.op, validOps)
		}

		switch query.param {
		case "b", "bw", "bandwidth":
			filterBandwidth(manifest, query)
		case "c", "codec":
			// filterCodecs == audio, video, text, image
			log.Fatal("not there yet")
		case "fr", "fps":
			// filterFrameRate == ...
			log.Fatal("not there yet")
		case "h", "height":
			// filterHeight == ...
			log.Fatal("not there yet")
		case "l", "lang", "language":
			// filterLanguage == audio, video, text, image
			log.Fatal("not there yet")
		case "t", "type":
			// filterContentType == ...
			log.Fatal("not there yet")
		case "ts", "timescale":
			// filterTimescale == ...
			log.Fatal("not there yet")
		case "w", "width":
			// filterTimescale == ...
			log.Fatal("not there yet")
		}

		if print {
			PrintManifest(manifest)
		}
	}
}

type query struct {
	param, op, value string
}

func getQueries(query string) []string {
	var queries []string
	if strings.Contains(query, ",") {
		queries = strings.Split(query, ",")
	}
	if strings.Contains(query, "&&") {
		queries = strings.Split(query, "&&")
	}
	if len(queries) == 0 {
		queries = append(queries, query)
	}
	return queries
}

func parseQuery(queryString string) query {
	parts := strings.Split(queryString, " ")
	if len(parts) != 3 {
		log.Fatalf("expected '{name} {op} {value}' got %v", queryString)
	}
	query := query{
		param: parts[0],
		op:    parts[1],
		value: parts[2],
	}
	return query
}
