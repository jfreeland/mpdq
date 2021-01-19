package mpdqlib

import (
	"log"
	"strconv"

	"github.com/zencoder/go-dash/mpd"
)

func filterBandwidth(manifest *mpd.MPD, query query) *mpd.MPD {
	value, err := strconv.Atoi(query.value)
	if err != nil {
		log.Fatalf("value is not an integer: %v", query.value)
	}
	for _, period := range manifest.Periods {
		var filteredAdaptationSets []*mpd.AdaptationSet
		for _, as := range period.AdaptationSets {
			var filteredRepresentations []*mpd.Representation
			for _, r := range as.Representations {
				repBandwidth := int(*r.Bandwidth)
				if repBandwidth == 0 {
					continue
				}

				switch query.op {
				case "=", "==":
					if repBandwidth == value {
						filteredRepresentations = append(filteredRepresentations, r)
					}
				case "!=":
					if repBandwidth != value {
						filteredRepresentations = append(filteredRepresentations, r)
					}
				case "<":
					if repBandwidth < value {
						filteredRepresentations = append(filteredRepresentations, r)
					}
				case "<=":
					if repBandwidth <= value {
						filteredRepresentations = append(filteredRepresentations, r)
					}
				case ">":
					if repBandwidth > value {
						filteredRepresentations = append(filteredRepresentations, r)
					}
				case ">=":
					if repBandwidth >= value {
						filteredRepresentations = append(filteredRepresentations, r)
					}
				}

				as.Representations = filteredRepresentations
				if len(as.Representations) != 0 {
					filteredAdaptationSets = append(filteredAdaptationSets, as)
				}
			}
		}

		period.AdaptationSets = filteredAdaptationSets
	}
	return manifest
}
