package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qmd",
	Short: "Quick Markdown Search",
	Long: `An on-device search engine for everything you need to remember.
Index your markdown notes, meeting transcripts, documentation, and knowledge bases.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
