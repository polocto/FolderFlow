/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/polocto/FolderFlow/cmd"
	"github.com/polocto/FolderFlow/internal/logger"
)

func main() {
	logger.Log.Info("Application started")
	cmd.Execute()
}
