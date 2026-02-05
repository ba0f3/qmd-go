package main

import (
	"fmt"
	"os"

	"github.com/ba0f3/qmd-go/internal/config"
	"github.com/ba0f3/qmd-go/internal/store"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qmd",
	Short: "Quick Markdown Search",
	Long: `An on-device search engine for everything you need to remember.
Index your markdown notes, meeting transcripts, documentation, and knowledge bases.`,
}

func getIndexName() string {
	name, _ := rootCmd.PersistentFlags().GetString("index")
	if name == "" {
		return "index"
	}
	return name
}

func getStorePath() (string, error) {
	return store.GetDefaultDbPath(getIndexName())
}

func openStore() (*store.Store, error) {
	path, err := getStorePath()
	if err != nil {
		return nil, err
	}
	return store.NewStore(path)
}

func initRoot() {
	config.CurrentIndexName = getIndexName()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("index", "", "Use named index (default: index)")
}
