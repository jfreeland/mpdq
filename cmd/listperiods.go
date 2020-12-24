package cmd

import (
	"github.com/jfreeland/mpdq/mpdqlib"
	"github.com/spf13/cobra"
)

// listPeriodsCmd represents the playback command
var listPeriodsCmd = &cobra.Command{
	Use:     "listperiods",
	Aliases: []string{"lp", "p"},
	Short:   "Lists periods for a manifest",
	Long:    "Lists periods for a manifest",
	Run: func(cmd *cobra.Command, args []string) {
		mpdqlib.ListPeriods(manifest, representation)
	},
}

func init() {
	rootCmd.AddCommand(listPeriodsCmd)

	listPeriodsCmd.Flags().StringVarP(&representation, "representation", "r", "max", "the representation you want to list segments for")
}
