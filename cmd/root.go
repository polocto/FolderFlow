/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/polocto/FolderFlow/internal/logger"
	"github.com/spf13/cobra"
)

type AppConfig struct {
	Verbose bool
	DryRun  bool
}

var cfg AppConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "FolderFlow",
	Short: "Manage and organise your folder as you want",
	Long:  "FolderFlow is a command-line tool that helps you filter and move image and video files from a source directory to specified destination folders (e.g., 'images/', 'videos/'), while keeping the original folder structure. It also supports creating a 'regroup' folder with symlinks or hard links to all moved files for easy access.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return logger.Init(cfg.Verbose)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVar(&cfg.Verbose, "verbose", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&cfg.DryRun, "dry-run", false, "dry run (no changes)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
