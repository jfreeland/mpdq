package cmd

import (
	"log"
	"strings"

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
		q := strings.TrimSpace(args[0])
		pieces := strings.Split(q, " ")
		if len(pieces) != 3 {
			log.Fatalf("expected '{name} {op} {value}' got %v", q)
		}
		mpdqlib.Query(manifest, pieces[0], pieces[1], pieces[2])
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
