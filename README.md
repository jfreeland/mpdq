# mpdq

## NOTE

**This is very much a work in progress.**

This was primarily intended to be used by me as an easier way to interact with "dynamic" MPEG-DASH manifests, as used with LIVE streaming video. It scratches my itch for now but I'll probably continue to clean it up and add functionality over time.

For the time being I've removed all support for static manifests as I don't deal with video on demand much. Perhaps I'll add that back in later.

## Installation

```sh
go install github.com/jfreeland/mpdq
```

## Usage

At the moment this:

- has only been tested with 'dynamic' (live) DASH manifests
  - and more specifically works best with manifests that contain a SegmentTimeline as there's some additional troubleshooting I need to do for the no-SegmentTimeline case.
- only works for video representations

You can use this by passing in a DASH manifest via stdin, specifying a file, or specifying a URL. If you do not specify a command, this will print a highlighted version of the DASH manifest XML.

You can list representations (`lr`) or list segments for a presentation (`ls`). If you do not specify a representation (either with the number or the name), this will choose the highest bandwidth representation when listing segments.

You can query a manifest (`q`) but there's enormous amounts of of room for improvement. TODO's are in the relevant files.

When you're listing segments, this will alternate colors to highlight different periods. It will also highlight segment durations that do not match the previous segment duration in red. If it looks like there's a gap between periods it will call that out.

### Add Color

```sh
curl https://www.website.com/path/to/some/master.mpd | mpdq
mpdq testdata/dynamic.mpd
mpdq -u https://www.website.com/path/to/some/master.mpd
```

### List representations

```sh
curl https://www.website.com/path/to/some/master.mpd | mpdq lr
mpdq lr -f testdata/dynamic.mpd
mpdq lr -u https://www.website.com/path/to/some/master.mpd
```

### List segments of highest bandwidth representation

```sh
curl https://www.website.com/path/to/some/master.mpd | mpdq ls
mpdq ls -f testdata/dynamic.mpd
mpdq ls -u https://www.website.com/path/to/some/master.mpd | more -R
```

### List segments for a specific representation

```sh
curl https://www.website.com/path/to/some/master.mpd | mpdq ls -r 540p-30fps-2436kbps
mpdq ls -r 540p-30fps-2436kbps -f testdata/dynamic.mpd
mpdq ls -r 540p-30fps-2436kbps -u https://www.website.com/path/to/some/master.mpd | more -R
```

### Watch a manifest

A URL must be provided to continually watch the manifest.

```sh
mpdq ls -w -u https://www.website.com/path/to/some/master.mpd
mpdq ls -w -r 540p-30fps-2436kbps -u https://www.website.com/path/to/some/master.mpd
```

### Query a manifest

```sh
curl https://www.website.com/path/to/some/master.mpd | mpdq q 'bandwidth >= 300000'
mpdq q 'bandwidth >= 300000' -f testdata/dynamic.mpd
mpdq q 'bandwidth >= 300000' -u https://www.website.com/path/to/some/master.mpd
```

## Development

```sh
curl https://www.website.com/path/to/some/master.mpd | go run . lr
go run . lr -f testdata/dynamic.mpd
go run . ls -u https://www.website.com/path/to/some/master.mpd
go run . q 'bandwidth >= 300000'
```

### Remove special characters and color from saved output

```sh
cat testdata/dynamic.mpd | sed -E "s/"$'\E'"\[([0-9]{1,3}((;[0-9]{1,3})*)?)?[m|K]//g" > d.mpd
```

## Manifests to test against

```sh
https://testassets.dashif.org/#testvector/list
https://livesim.dashif.org/livesim/testpic_2s/Manifest.mpd
https://livesim.dashif.org/livesim/start_1800/testpic_2s/Manifest.mpd
https://livesim.dashif.org/livesim/scte35_2/testpic_2s/Manifest.mpd
https://livesim.dashif.org/livesim/modulo_10/testpic_2s/Manifest.mpd
https://livesim.dashif.org/livesim/utc_direct-head/testpic_2s/Manifest.mpd (way off)
https://livesim.dashif.org/livesim/chunkdur_1/ato_7/testpic4_8s/Manifest300.mpd (not even close)
https://livesim.dashif.org/livesim/chunkdur_1/ato_7/testpic4_8s/Manifest.mpd (not even close)
https://livesim.dashif.org/livesim/testpic_2s/Manifest.mpd#t=posix:1465406946
https://livesim.dashif.org/livesim/testpic_2s/Manifest.mpd#t=posix:now
https://livesim.dashif.org/livesim/utc_direct/testpic_2s/Manifest.mpd
https://livesim.dashif.org/livesim/utc_head/testpic_2s/Manifest.mpd
https://livesim.dashif.org/livesim/utc_ntp/testpic_2s/Manifest.mpd
https://livesim.dashif.org/livesim/utc_sntp/testpic_2s/Manifest.mpd
```
