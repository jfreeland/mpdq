package mpdqlib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zencoder/go-dash/mpd"
)

// TODO: Return an err if we can't find the requested representation.
func getOneVideoRepresentation(manifest *mpd.MPD, representation string) ListRepresentation {
	var r ListRepresentation
	if representation == "max" {
		r = GetMaxVideoRepresentation(manifest)
	} else {
		reps := GetVideoRepresentations(manifest)
		for _, v := range reps {
			if representation == v.ID || representation == strconv.Itoa(v.Bandwidth) {
				r = v
				fmt.Printf("returned %v\n", r)
			}
		}
	}
	return r
}

func cleanFilePath(mpdBase, filePath, rep string, segNum int) string {
	if strings.Contains(filePath, "$RepresentationID$") && rep != "" && rep != "max" {
		re := regexp.MustCompile(`\$RepresentationID\$`)
		filePath = re.ReplaceAllString(filePath, rep)
	}
	if strings.Contains(filePath, "$Number$") && segNum != 0 {
		re := regexp.MustCompile(`\$Number\$`)
		filePath = re.ReplaceAllString(filePath, strconv.Itoa(segNum))
	} else if strings.Contains(filePath, "$Number") && segNum != 0 {
		// TODO: lazy
		re := regexp.MustCompile(`\$Number\%06d\$`)
		filePath = re.ReplaceAllString(filePath, fmt.Sprintf("%06d", segNum))
	}
	if mpdBase != "" {
		return mpdBase + filePath
	}
	return filePath
}

func getRowColor(pidx int) []int {
	if pidx%2 != 0 {
		return tablewriter.Colors{tablewriter.FgCyanColor}
	}
	return tablewriter.Colors{tablewriter.Normal}
}

func getDurationColor(duration uint64, expectedDuration uint64, pidx int) []int {
	c := getRowColor(pidx)
	if duration != expectedDuration {
		return tablewriter.Colors{tablewriter.FgRedColor}
	}
	return c
}

func getAvailabilityStartTime(manifestAvailabilityStartTime string) time.Time {
	// TODO: There's got to be a cleaner way to do this.
	startTimeLayout := "2006-01-02T15:04:05.000Z"
	availabilityStartTime, err := time.Parse(startTimeLayout, manifestAvailabilityStartTime)
	if err != nil {
		alternateStartTimeLayout := "2006-01-02T15:04:05Z"
		availabilityStartTime, err = time.Parse(alternateStartTimeLayout, manifestAvailabilityStartTime)
		if err != nil {
			fmt.Printf("error parsing availabilityStartTime: %v\n", err)
		}
	}
	//fmt.Printf("availabilityStartTime: %v\n", availabilityStartTime.String())
	return availabilityStartTime
}

func getManifest(mpdURL string) *mpd.MPD {
	resp, err := http.Get(mpdURL)
	if err != nil {
		log.Fatalf("could not fetch manifest: %v\n", err)
	}
	manifestBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("could not read manifest: %v\n", err)
	}
	manifest, err := mpd.ReadFromString(string(manifestBody))
	if err != nil {
		log.Fatalf("could not parse manifest: %v\n", err)
	}
	return manifest
}

func getMPDBase(fromManifest, fromURL string) string {
	var mpdBase string
	if fromManifest != "" {
		mpdBase = fromManifest
	} else {
		mpdBase, _ = path.Split(fromURL)
	}
	return mpdBase
}

func getRepeatCount(r *int) int {
	if r != nil {
		return *r
	}
	return 0
}

func getTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(false)
	table.SetColumnSeparator("  ")
	table.SetHeaderLine(false)
	return table
}

func getUpdatePeriod(manifestUpdatePeriod *string) time.Duration {
	re := regexp.MustCompile(`PT(\d+)S`)
	minUpdatePeriodMatch := re.FindStringSubmatch(*manifestUpdatePeriod)
	minUpdatePeriod, err := strconv.Atoi(minUpdatePeriodMatch[1])
	if err != nil {
		log.Fatalf("could not parse minimum update period: %v\n", err)
	}
	return time.Duration(minUpdatePeriod) * time.Second
}

func checkForGap(t *tablewriter.Table, pidx int, periodID string, timeBetweenPeriods time.Duration, previousSegmentEndTime time.Time) (segment, bool) {
	// TODO: I arbirarily chose 500ms as the cutoff for what might be a gap?
	if timeBetweenPeriods > time.Duration(500*time.Millisecond) {
		gapSegment := segment{
			durationC: tablewriter.Color(tablewriter.FgRedColor),
			durationS: timeBetweenPeriods.String(),
			endTime:   previousSegmentEndTime.Add(time.Duration(timeBetweenPeriods)).Round(time.Second),
			number:    fmt.Sprintf("GAP%v", periodID),
			path:      "possible gap detected",
			pidx:      pidx,
			period:    periodID,
			startTime: previousSegmentEndTime.Round(time.Second),
		}
		return gapSegment, true
	}
	return segment{}, false
}

func checkPeriodsSameDuration(manifest *mpd.MPD) (int, bool) {
	var (
		durations  []int
		startTimes []time.Time
	)
	availabilityStartTime := getAvailabilityStartTime(*manifest.AvailabilityStartTime)
	for _, period := range manifest.Periods {
		periodStartTime := availabilityStartTime.Add(time.Duration(*period.Start))
		startTimes = append(startTimes, periodStartTime)
	}
	for i := 0; i < len(startTimes)-1; i++ {
		duration := int(startTimes[i].Sub(startTimes[i+1]).Seconds()) * -1
		durations = append(durations, duration)
	}
	totalSeconds := 0
	for _, s := range durations {
		totalSeconds += s
	}
	if totalSeconds/len(durations) == durations[0] {
		return durations[0], true
	}
	return 0, false
}

func printHeader(t *tablewriter.Table) {
	c := tablewriter.Color(tablewriter.Normal)
	data := ([]string{"PERIOD ID", "SEG ID", "SEGMENT START TIME", "DURATION", "SEGMENT END TIME", "PATH"})
	if t != nil {
		t.Rich(data, []tablewriter.Colors{c, c, c, c, c, c})
	}
}

func printSegment(t *tablewriter.Table, color []int, segment segment) {
	data := []string{segment.period, segment.number, segment.startTime.String(), segment.durationS, segment.endTime.String(), segment.path}
	if t != nil {
		t.Rich(data, []tablewriter.Colors{color, color, color, segment.durationC, color, color})
	}
}
