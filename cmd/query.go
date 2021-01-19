package cmd

import (
	"log"

	"github.com/jfreeland/mpdq/mpdqlib"
	"github.com/spf13/cobra"
)

// listPeriodsCmd represents the playback command
var queryCmd = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q"},
	Short:   "Run a query against a manifest",
	Long:    "Run a query against a manifest.  Only return matching components.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("did not receive a query")
		}
		if len(args) > 1 || args[0] == "" {
			log.Fatalf("query must quoted and in the form of '{param} {operand} {value}'")
		}
		mpdqlib.Query(manifest, args[0])
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
