package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/jfreeland/mpdq/mpdqlib"
	"github.com/spf13/cobra"
	"github.com/zencoder/go-dash/mpd"
)

// http://dash.akamaized.net/dash264/TestCasesHD/1b/qualcomm/1/MultiRate.mpd
// https://github.com/google/shaka-player/blob/master/docs/design/dash-manifests.md

var (
	r              io.Reader
	err            error
	manifest       *mpd.MPD
	representation string

	mpdBase, mpdFile, mpdURL string
)

var rootCmd = &cobra.Command{
	Use:   "mpqd",
	Short: "mpqd attempts to be a friendly way to parse DASH manifests",
	Long:  "mpqd attempts to be a friendly way to parse DASH manifests",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		getManifest()
	},
	Run: func(cmd *cobra.Command, args []string) {
		mpdqlib.PrintManifest(manifest)
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&mpdFile, "file", "f", "", "mpd filename")
	rootCmd.PersistentFlags().StringVarP(&mpdURL, "url", "u", "", "mpd url")
}

func getManifest() {
	if mpdFile == "" && mpdURL == "" {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			r = os.Stdin
		}
	} else if mpdFile != "" && mpdURL != "" {
		log.Fatalf("can only take a file or url, not both\n")
	} else if mpdFile != "" && mpdURL == "" {
		if r, err = os.Open(mpdFile); err != nil {
			log.Fatalf("could not open file %v err=%v\n", mpdFile, err)
		}
	} else if mpdFile == "" && mpdURL != "" {
		resp, err := http.Get(mpdURL)
		if err != nil {
			log.Fatalf("could not fetch manifest: %v\n", err)
		}
		r = resp.Body
		mpdBase, _ = path.Split(mpdURL)
	}

	manifestBody, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("could not read manifest body: %v\n", err)
	}

	parsed, err := mpd.ReadFromString(string(manifestBody))
	if err != nil {
		log.Fatalf("could not parse manifest: %v\n", err)
	}
	manifest = parsed
}
