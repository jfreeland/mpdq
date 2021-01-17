package cmd

import (
	"github.com/jfreeland/mpdq/mpdqlib"
	"github.com/spf13/cobra"
)

var (
	watch    bool
	lastTime string
)

// listSegmentsCmd represents the playback command
var listSegmentsCmd = &cobra.Command{
	Use:     "listsegments",
	Aliases: []string{"ls", "s"},
	Short:   "Lists segments for a manifest",
	Long:    "Lists segments for a manifest",
	Run: func(cmd *cobra.Command, args []string) {
		mpdqlib.ListSegments(manifest, watch, lastTime, mpdURL, representation, mpdBase)
	},
}

func init() {
	rootCmd.AddCommand(listSegmentsCmd)

	listSegmentsCmd.Flags().StringVarP(&representation, "representation", "r", "max", "the representation you want to list segments for")
	listSegmentsCmd.Flags().BoolVarP(&watch, "watch", "w", false, "continuously watch the manifest being updated")
	listSegmentsCmd.Flags().StringVarP(&lastTime, "last", "l", "300s", "how far back to list segments while watching the manifest")
}
