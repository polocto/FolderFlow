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

package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/polocto/FolderFlow/cmd"
	_ "github.com/polocto/FolderFlow/internal/filter"
	_ "github.com/polocto/FolderFlow/internal/strategy"
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
