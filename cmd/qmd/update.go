package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tobi/qmd-go/internal/config"
	"github.com/tobi/qmd-go/internal/indexer"
	"github.com/tobi/qmd-go/internal/store"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update index for all collections",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		s, err := store.NewStore("")
		if err != nil {
			fmt.Printf("Error opening store: %v\n", err)
			return
		}
		defer s.Close()

		for name, col := range cfg.Collections {
			fmt.Printf("Updating collection '%s'...\n", name)
			if err := indexer.IndexFiles(s, name, col.Path, col.Pattern); err != nil {
				fmt.Printf("Error indexing collection '%s': %v\n", name, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
