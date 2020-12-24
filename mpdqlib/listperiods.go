package mpdqlib

import (
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zencoder/go-dash/mpd"
)

// ListPeriods lists the periods for a given representation
// TODO: I don't think this is terribly interesting or relevant.  Probably going to delete this.
func ListPeriods(manifest *mpd.MPD, representation string) {
	r := getOneVideoRepresentation(manifest, representation)
	listPeriods(manifest, r)
}

func listPeriods(manifest *mpd.MPD, r ListRepresentation) {
	// var (
	// 	duration                uint64
	// 	durationColor, rowColor []int
	// )

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	// TODO: I don't know if I should be attempting to print wall clock time (which I think I am?) or presentation time or something else?
	table.SetHeader([]string{"period", "wall clock time?", "duration", "number", "path"})
	//defer table.Render()

	startTimeLayout := "2006-01-02T15:04:05.000Z"
	startTime, err := time.Parse(startTimeLayout, *manifest.AvailabilityStartTime)
	if err != nil {
		fmt.Printf("startTime parse error: %v\n", err)
	}
	fmt.Printf("availabilityStartTime: %v\n", startTime.String())

	for _, period := range manifest.Periods {
		pstartTime := period.Start.String()
		currentTime, err := mpd.ParseDuration(pstartTime)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		periodStartTime := startTime.Add(currentTime)
		fmt.Printf("%v\n", periodStartTime)
	}
}
