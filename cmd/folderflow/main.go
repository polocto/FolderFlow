/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/joho/godotenv"
	"github.com/polocto/FolderFlow/cmd"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}
