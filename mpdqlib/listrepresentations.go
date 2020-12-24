package mpdqlib

import (
	"os"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/zencoder/go-dash/mpd"
)

// ListRepresentation is used to store a list of representation ID's and bandwidth
type ListRepresentation struct {
	ID        string
	Bandwidth int
}

func representationExists(a ListRepresentation, list []ListRepresentation) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetVideoRepresentations returns an ordered list of representation ID's and bandwidth
func GetVideoRepresentations(manifest *mpd.MPD) []ListRepresentation {
	representations := make([]ListRepresentation, 0)
	for _, period := range manifest.Periods {
		for _, as := range period.AdaptationSets {
			for _, rep := range as.Representations {
				if as.MimeType != nil && *as.MimeType != "video/mp4" {
					continue
				}
				if rep.MimeType != nil && *rep.MimeType != "video/mp4" {
					continue
				}
				r := ListRepresentation{
					ID:        *rep.ID,
					Bandwidth: int(*rep.Bandwidth),
				}
				if !representationExists(r, representations) {
					representations = append(representations, r)
				}
			}
		}
	}

	sort.SliceStable(representations, func(i, j int) bool {
		return representations[i].Bandwidth > representations[j].Bandwidth
	})
	return representations
}

// GetMaxVideoRepresentation returns the highest bandwidth representation
func GetMaxVideoRepresentation(manifest *mpd.MPD) ListRepresentation {
	r := GetVideoRepresentations(manifest)
	return r[0]
}

// ListVideoRepresentations lists the representations in the manifest
func ListVideoRepresentations(manifest *mpd.MPD) {
	// if manifest.ID != "" {
	// 	fmt.Printf("received %v\n", manifest.ID)
	// } else {
	// 	for _, pi := range manifest.ProgramInformation {
	// 		if pi.Title != "" {
	// 			fmt.Printf("received %v\n", pi.Title)
	// 		}
	// 	}
	// }

	representations := GetVideoRepresentations(manifest)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "bandwidth"})
	for _, v := range representations {
		table.Append([]string{v.ID, strconv.Itoa(v.Bandwidth)})
	}
	table.Render()
}
