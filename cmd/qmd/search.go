package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tobi/qmd-go/internal/store"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Full-text search",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")
		limit, _ := cmd.Flags().GetInt("n")

		s, err := store.NewStore("")
		if err != nil {
			fmt.Printf("Error opening store: %v\n", err)
			os.Exit(1)
		}
		defer s.Close()

		results, err := s.SearchFTS(query, limit)
		if err != nil {
			fmt.Printf("Search failed: %v\n", err)
			os.Exit(1)
		}

		if len(results) == 0 {
			fmt.Println("No results found.")
			return
		}

		for _, r := range results {
			fmt.Printf("%s (Score: %.2f)\n", r.Filepath, r.Score)
			fmt.Printf("Title: %s\n", r.Title)
			fmt.Println()
		}
	},
}

func init() {
	searchCmd.Flags().IntP("n", "n", 5, "Number of results")
	rootCmd.AddCommand(searchCmd)
}
