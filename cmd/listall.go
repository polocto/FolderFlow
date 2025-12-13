/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/polocto/FolderFlow/internal/core"
	"github.com/polocto/FolderFlow/internal/logger"
	"github.com/spf13/cobra"
)

// listallCmd represents the listall command
var listallCmd = &cobra.Command{
	Use:   "listall",
	Short: "List all extensions found in the specify folder",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var path string = "."
		if len(args) > 0 {
			path = args[0]
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			if cfg.Verbose {
				fmt.Printf("%s does not exist.", path)
			}
		}

		extensions, err := core.ListAllFilesExtensions(path, cfg.DryRun, cfg.Verbose)

		if err != nil {
			logger.Error(err.Error())
			return
		}

		if cfg.Verbose {
			fmt.Printf("Found %d unique extensions in %s\n", len(extensions), path)
		}

		for _, ext := range extensions {
			fmt.Println(ext)
		}
	},
}

func init() {
	rootCmd.AddCommand(listallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
