package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ba0f3/qmd-go/internal/config"
	"github.com/spf13/cobra"
)

var collectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "Manage collections",
}

var collectionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all collections",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if len(cfg.Collections) == 0 {
			fmt.Println("No collections found.")
			return
		}

		fmt.Println("Collections:")
		for name, col := range cfg.Collections {
			fmt.Printf("- %s (%s) [%s]\n", name, col.Path, col.Pattern)
		}
	},
}

var collectionAddCmd = &cobra.Command{
	Use:   "add [path]",
	Short: "Add a collection",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		absPath, _ := filepath.Abs(path)

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = filepath.Base(absPath)
		}
		pattern, _ := cmd.Flags().GetString("mask")

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		if _, exists := cfg.Collections[name]; exists {
			fmt.Printf("Collection '%s' already exists.\n", name)
			os.Exit(1)
		}

		cfg.Collections[name] = config.Collection{
			Path:    absPath,
			Pattern: pattern,
		}

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Collection '%s' added.\n", name)
	},
}

var collectionRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a collection",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		if _, exists := cfg.Collections[name]; !exists {
			fmt.Printf("Collection '%s' not found.\n", name)
			os.Exit(1)
		}

		delete(cfg.Collections, name)

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Collection '%s' removed.\n", name)
	},
}

func init() {
	collectionAddCmd.Flags().String("name", "", "Collection name")
	collectionAddCmd.Flags().String("mask", "**/*.md", "File pattern mask")

	collectionCmd.AddCommand(collectionListCmd)
	collectionCmd.AddCommand(collectionAddCmd)
	collectionCmd.AddCommand(collectionRemoveCmd)
	rootCmd.AddCommand(collectionCmd)
}
