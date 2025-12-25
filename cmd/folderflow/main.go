/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/polocto/FolderFlow/cmd"
)

func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
