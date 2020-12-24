package cmd

import (
	"github.com/jfreeland/mpdq/mpdqlib"
	"github.com/spf13/cobra"
)

// listRepresentationsCmd represents the playback command
var listRepresentationsCmd = &cobra.Command{
	Use:     "listrepresentations",
	Aliases: []string{"lr", "r"},
	Short:   "Lists representations for a manifest",
	Long:    "List representations for a manifest",
	Run: func(cmd *cobra.Command, args []string) {
		mpdqlib.ListVideoRepresentations(manifest)
	},
}

func init() {
	rootCmd.AddCommand(listRepresentationsCmd)
}
