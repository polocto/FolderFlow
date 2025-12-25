/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/polocto/FolderFlow/internal/logger"
	"github.com/spf13/cobra"
)

type AppConfig struct {
	DryRun bool
}

var verbose bool
var debug bool

var cfg AppConfig

var closeLogger func() error

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "FolderFlow",
	Short: "Manage and organise your folder as you want",
	Long:  "FolderFlow is a command-line tool that helps you filter and move image and video files from a source directory to specified destination folders (e.g., 'images/', 'videos/'), while keeping the original folder structure. It also supports creating a 'regroup' folder with symlinks or hard links to all moved files for easy access.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		closeLogger, err = logger.Init(verbose, debug || os.Getenv("DEBUG") != "")
		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	defer func() {
		if closeLogger != nil {
			_ = closeLogger()
		}
	}()
	if err := rootCmd.Execute(); err != nil {
		slog.Error("command execution failed", "error", err)
		return err
	}
	return nil
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&cfg.DryRun, "dry-run", false, "dry run (no changes)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debugging logs")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
