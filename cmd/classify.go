// Copyright (c) 2026 Paul Sade.
//
// This file is part of the FolderFlow project.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3,
// as published by the Free Software Foundation (see the LICENSE file).
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU General Public License for more details.

package cmd

import (
	"log/slog"

	"github.com/polocto/FolderFlow/internal/classify"
	"github.com/polocto/FolderFlow/internal/config"
	"github.com/polocto/FolderFlow/internal/stats"
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
		var s stats.Stats
		if classifier, err := classify.NewClassifier(*conf, &s, cfg.DryRun); err != nil {
			slog.Error("An error occured while configuring classification")
		} else if err := classifier.Classify(); err != nil {
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
	classifyCmd.Flags().StringVarP(&configFile, "config", "c", "", "path of the YAML config file")
}
