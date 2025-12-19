/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"

	"github.com/polocto/FolderFlow/internal/classify"
	"github.com/polocto/FolderFlow/internal/config"
	"github.com/spf13/cobra"
)

var configFile string

// classifyCmd represents the classify command
var classifyCmd = &cobra.Command{
	Use:     "classify",
	Aliases: []string{"class"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.LoadConfig(configFile)

		if err != nil {
			slog.Error("An error occured while loading the config", "error", err)
			return
		}

		if err := classify.Classify(*conf, cfg.DryRun); err != nil {
			slog.Error("An error occured while classing the documents", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(classifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// classifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	classifyCmd.Flags().StringVarP(&configFile, "file", "f", "", "")
}
