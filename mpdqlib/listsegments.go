package mpdqlib

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/zencoder/go-dash/mpd"
)

// ListSegments lists the segments for a given representation
func ListSegments(manifest *mpd.MPD, watch bool, mpdURL, representation, mpdBase string) {
	r := getOneVideoRepresentation(manifest, representation)
	if *manifest.Type == "dynamic" && !watch {
		listDynamicSegments(manifest, r, mpdBase, true)
	} else if *manifest.Type == "dynamic" && watch {
		watchDynamicSegments(manifest, r, mpdURL)
		//fmt.Println("forgot how to do this")
	} else if *manifest.Type == "static" {
		//listStaticSegments(manifest, r, mpdBase)
		fmt.Println("forgot how to do this")
	} else {
		fmt.Printf("i have no idea what i'm doing with this manifest type: %v\n", *manifest.Type)
	}
}

type segment struct {
	pidx      int
	durationC []int

	durationS, path, period, number string

	endTime, startTime time.Time
}

type templateOptions struct {
	pidx                                        int
	period                                      *mpd.Period
	rep                                         *mpd.Representation
	mpdBase                                     string
	now, availabilityStartTime, periodStartTime time.Time
	previousSegmentEndTime                      *time.Time
	sTemplate                                   *mpd.SegmentTemplate
}

func listDynamicSegments(manifest *mpd.MPD, r ListRepresentation, mpdBase string, print bool) []segment {
	var (
		previousSegmentEndTime      time.Time
		timeBetweenPeriods          time.Duration
		allSegments, periodSegments []segment
		sTemplate                   *mpd.SegmentTemplate
	)
	table := getTable()
	if print {
		printHeader(table)
		defer table.Render()
	}

	mpdBase = getMPDBase(manifest.BaseURL, mpdBase)
	now := time.Now().UTC()
	availabilityStartTime := getAvailabilityStartTime(*manifest.AvailabilityStartTime)

	for pidx, period := range manifest.Periods {
		rowColor := getRowColor(pidx)
		periodStart, err := mpd.ParseDuration(period.Start.String())
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		periodStartTime := availabilityStartTime.Add(periodStart)
		if pidx != 0 {
			timeBetweenPeriods = periodStartTime.Sub(previousSegmentEndTime)
			gapSegment, gap := checkForGap(table, pidx, period.ID, timeBetweenPeriods, previousSegmentEndTime)
			if gap {
				if print {
					printSegment(table, tablewriter.Color(tablewriter.FgRedColor), gapSegment)
				}
				allSegments = append(allSegments, gapSegment)
			}
		}
		for _, as := range period.AdaptationSets {
			if as.SegmentTemplate != nil {
				sTemplate = as.SegmentTemplate
			}
			for _, rep := range as.Representations {
				if r.ID != *rep.ID {
					continue
				}
				if rep.SegmentTemplate != nil {
					sTemplate = rep.SegmentTemplate
				}
				if sTemplate.SegmentTimeline == nil {
					periodSegments = parseSegmentTemplateNoTimeline(templateOptions{
						pidx:                   pidx,
						period:                 period,
						rep:                    rep,
						mpdBase:                mpdBase,
						now:                    now,
						availabilityStartTime:  availabilityStartTime,
						periodStartTime:        periodStartTime,
						previousSegmentEndTime: &previousSegmentEndTime,
						sTemplate:              sTemplate,
					})
				} else {
					periodSegments = parseSegmentTemplateWithTimeline(templateOptions{
						pidx:                   pidx,
						period:                 period,
						rep:                    rep,
						mpdBase:                mpdBase,
						now:                    now,
						availabilityStartTime:  availabilityStartTime,
						periodStartTime:        periodStartTime,
						previousSegmentEndTime: &previousSegmentEndTime,
						sTemplate:              sTemplate,
					})
				}
				for _, segment := range periodSegments {
					allSegments = append(allSegments, segment)
					if print {
						printSegment(table, rowColor, segment)
					}
				}
			}
		}
	}
	return allSegments
}

func parseSegmentTemplateWithTimeline(opts templateOptions) []segment {
	var (
		currentTime time.Time
		duration    uint64
		segNum      int
		segments    []segment
	)
	sTimeline := opts.sTemplate.SegmentTimeline
	segNum = int(*opts.sTemplate.StartNumber)
	for _, s := range sTimeline.Segments {
		repeat := getRepeatCount(s.RepeatCount)
		for idx := 0; idx <= repeat; idx++ {
			sTimescale := uint64(*opts.sTemplate.Timescale)
			expectedDuration := s.Duration / sTimescale * uint64(time.Second)
			durationColor := getDurationColor(duration, expectedDuration, opts.pidx)
			duration = s.Duration / sTimescale * uint64(time.Second)
			if idx == 0 {
				// TODO: I think this is right?  Needs to be validated.
				// https://github.com/google/shaka-player/blob/master/docs/design/dash-manifests.md#calculating-presentation-times
				currentTime = opts.periodStartTime.Add(time.Duration((*s.StartTime - *opts.sTemplate.PresentationTimeOffset) / sTimescale * uint64(time.Second)))
			} else {
				currentTime = currentTime.Add(time.Duration(duration))
			}
			durationS := strconv.FormatUint(s.Duration/uint64(*opts.sTemplate.Timescale), 10) + "s"
			path := cleanFilePath(opts.mpdBase, *opts.sTemplate.Media, *opts.rep.ID, int(segNum))
			endTime := currentTime.Add(time.Duration(duration))
			*opts.previousSegmentEndTime = endTime
			segments = append(segments, segment{
				durationC: durationColor,
				durationS: durationS,
				endTime:   endTime.Round(time.Second),
				number:    strconv.Itoa(int(segNum)),
				path:      path,
				pidx:      opts.pidx,
				period:    opts.period.ID,
				startTime: currentTime.Round(time.Second),
			})
			segNum++
		}
	}
	return segments
}

// https://livesim.dashif.org/livesim/periods_20/testpic_2s/Manifest.mpd
func parseSegmentTemplateNoTimeline(opts templateOptions) []segment {
	var (
		segNum   int
		segments []segment
	)
	// TODO: We'll look back 60 seconds for now
	currentTime := opts.now.Add(time.Duration(-60 * time.Second))

	if int(*opts.sTemplate.StartNumber) != 0 {
		segNum = int(*opts.sTemplate.StartNumber)
	} else {
		segNum = int(opts.now.Sub(opts.periodStartTime).Seconds()) / int(*opts.sTemplate.Duration)
	}
	for currentTime.Before(opts.now) {
		path := cleanFilePath(opts.mpdBase, *opts.sTemplate.Media, *opts.rep.ID, int(segNum))
		endTime := currentTime.Add(time.Duration(*opts.sTemplate.Duration) * time.Second)
		*opts.previousSegmentEndTime = opts.periodStartTime.Add(time.Duration(*opts.sTemplate.Duration) * time.Second)
		durationColor := getDurationColor(0, 0, opts.pidx)
		durationS := fmt.Sprintf("%vs", *opts.sTemplate.Duration)
		segments = append(segments, segment{
			durationC: durationColor,
			durationS: durationS,
			endTime:   endTime.Round(time.Second),
			number:    strconv.Itoa(int(segNum)),
			path:      path,
			pidx:      opts.pidx,
			period:    opts.period.ID,
			startTime: currentTime.Round(time.Second),
		})
		segNum++
		currentTime = currentTime.Add(time.Duration(*opts.sTemplate.Duration) * time.Second)
	}
	return segments
}

// TODO: Collapse this into listDynamicSegments and just pass the watch flag
func watchDynamicSegments(manifest *mpd.MPD, r ListRepresentation, mpdURL string) {
	fmt.Println("watching last 5 minutes and moving forward")
	table := getTable()
	printHeader(table)

	if mpdURL == "" {
		panic("can't watch a manifest without a url")
	}
	var segments []segment
	mpdBase := getMPDBase(manifest.BaseURL, mpdURL)
	start := time.Now().UTC().Add(time.Duration(-300) * time.Second)
	seen := make(map[string]bool)
	updatePeriod := getUpdatePeriod(manifest.MinimumUpdatePeriod)

	segments = listDynamicSegments(manifest, r, mpdBase, false)
	for _, segment := range segments {
		if start.Before(segment.startTime) && !seen[segment.number] {
			printSegment(table, getRowColor(segment.pidx), segment)
		}
		seen[segment.number] = true
	}
	table.Render()
	table.ClearRows()

	for range time.Tick(updatePeriod) {
		manifest := getManifest(mpdURL)

		newSegments := 0
		segments = listDynamicSegments(manifest, r, mpdBase, false)
		for _, segment := range segments {
			if start.Before(segment.startTime) && !seen[segment.number] {
				newSegments++
				printSegment(table, getRowColor(segment.pidx), segment)
			}
			seen[segment.number] = true
		}
		if newSegments == 0 {
			color.New(color.FgRed).Println("did not see a new segment on this fetch: slow playlist update?")
		}
		table.Render()
		table.ClearRows()
	}
}

// TODO: I haven't touched this in a bit.  It probably needs a lot of work.
// func listStaticSegments(manifest *mpd.MPD, r ListRepresentation, mpdBase string) {
// 	table := tablewriter.NewWriter(os.Stdout)
// 	table.SetAutoWrapText(false)
// 	table.SetHeader([]string{"period", "presentation time offset", "duration", "number", "path"})
// 	defer table.Render()

// 	duration, err := mpd.ParseDuration(*manifest.MediaPresentationDuration)
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 	}
// 	durationS := int(duration.Seconds())
// 	var countedS int

// 	var segNum int
// 	for _, period := range manifest.Periods {
// 		var sTemplate *mpd.SegmentTemplate
// 		for _, as := range period.AdaptationSets {
// 			if as.SegmentTemplate != nil {
// 				sTemplate = as.SegmentTemplate
// 			}
// 			for _, rep := range as.Representations {
// 				if r.ID != *rep.ID {
// 					continue
// 				}
// 				if rep.SegmentTemplate != nil {
// 					sTemplate = rep.SegmentTemplate
// 				}
// 				if sTemplate.StartNumber != nil {
// 					segNum = int(*sTemplate.StartNumber)
// 				} else {
// 					segNum = 1
// 				}
// 				// TODO: find examples of and, figure out how to, handle multiperiod static manifests
// 				if len(manifest.Periods) == 1 {
// 					for countedS <= durationS {
// 						segTime := *sTemplate.Duration / *sTemplate.Timescale
// 						path := cleanFilePath(mpdBase, *sTemplate.Media, *rep.ID, segNum)
// 						rowColor := tablewriter.Colors{tablewriter.Normal}
// 						data := []string{period.ID, strconv.Itoa(countedS), strconv.Itoa(int(segTime)) + "s", strconv.Itoa(int(segNum)), path}
// 						table.Rich(data, []tablewriter.Colors{rowColor, rowColor, rowColor, rowColor, rowColor})
// 						countedS += int(segTime)
// 						segNum++
// 					}
// 				}
// 			}
// 		}
// 	}
// }
