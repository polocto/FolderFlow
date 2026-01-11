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

// pkg/ffplugin/filter/regex_test.go
package filter

import (
	"io/fs"
	"regexp"
	"testing"
	"time"
)

// MockFileInfo is a mock implementation of fs.FileInfo for testing.
type MockFileInfo struct{}

func (m *MockFileInfo) Name() string       { return "test.txt" }
func (m *MockFileInfo) Size() int64        { return 1024 }
func (m *MockFileInfo) Mode() fs.FileMode  { return 0644 }
func (m *MockFileInfo) ModTime() time.Time { return time.Now() }
func (m *MockFileInfo) IsDir() bool        { return false }
func (m *MockFileInfo) Sys() interface{}   { return nil }

func TestRegexFilterMatch(t *testing.T) {
	// Create a RegexFilter with a pattern that matches "test.txt"
	filter := &RegexFilter{
		Patterns:   []string{"^test.t.*|.*txt$"},
		compiledRe: []*regexp.Regexp{regexp.MustCompile(`^test.t.*|.*txt$`)},
	}

	// Test: Match a file that matches the pattern
	match, err := filter.Match("test.txt", &MockFileInfo{})
	if err != nil {
		t.Fatalf("Match returned error: %v", err)
	}
	if !match {
		t.Error("Match should return true for 'test.txt'")
	}

	// Test: Match a file that does not match the pattern
	match, err = filter.Match("test.md", &MockFileInfo{})
	if err != nil {
		t.Fatalf("Match returned error: %v", err)
	}
	if match {
		t.Error("Match should return false for 'test.md'")
	}

	// Test: Match an empty filename
	match, err = filter.Match("", &MockFileInfo{})
	if err != nil {
		t.Fatalf("Match returned error: %v", err)
	}
	if match {
		t.Error("Match should return false for empty filename")
	}

	// Test: Match a filename with no patterns
	filterEmpty := &RegexFilter{
		Patterns:   []string{},
		compiledRe: []*regexp.Regexp{},
	}
	match, err = filterEmpty.Match("anyfile.txt", &MockFileInfo{})
	if err != nil {
		t.Fatalf("Match returned error: %v", err)
	}
	if match {
		t.Error("Match should return false when no patterns are defined")
	}

	// Test: Match a filename with multiple patterns
	filterMultiple := &RegexFilter{
		Patterns:   []string{`^test\.txt$`, `^example\.md$`},
		compiledRe: []*regexp.Regexp{regexp.MustCompile(`^test\.txt$`), regexp.MustCompile(`^example\.md$`)},
	}
	match, err = filterMultiple.Match("example.md", &MockFileInfo{})
	if err != nil {
		t.Fatalf("Match returned error: %v", err)
	}
	if !match {
		t.Error("Match should return true for 'example.md'")
	}
	match, err = filterMultiple.Match("other.doc", &MockFileInfo{})
	if err != nil {
		t.Fatalf("Match returned error: %v", err)
	}
	if match {
		t.Error("Match should return false for 'other.doc'")
	}
}

func TestRegexFilterSelector(t *testing.T) {
	filter := &RegexFilter{}
	selector := filter.Selector()
	if selector != "regex" {
		t.Errorf("Selector() returned %s, expected 'regex'", selector)
	}
}

func TestRegexFilterLoadConfig(t *testing.T) {
	filter := &RegexFilter{}

	// Test: Load valid patterns
	err := filter.LoadConfig(map[string]interface{}{
		"patterns": []string{`test.txt`, `example.md`},
	})
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if len(filter.Patterns) != 2 {
		t.Error("LoadConfig did not load patterns correctly")
	}
	if len(filter.Patterns) != 2 {
		t.Error("LoadConfig did not compile patterns correctly")
	}

	// Test: Load invalid patterns (invalid regex)
	err = filter.LoadConfig(map[string]interface{}{
		"patterns": []string{"test.txt", "[invalid"},
	})
	if err == nil {
		t.Error("LoadConfig should return error for invalid regex")
	}

	// Test: Load missing patterns
	err = filter.LoadConfig(map[string]interface{}{"patterns": nil})
	if err == nil {
		t.Error("LoadConfig should return error for missing patterns")
	}
}

func TestRegexFilterLoadConfigInvalidConfig(t *testing.T) {
	filter := &RegexFilter{}

	// Test: Load config with invalid patterns type
	err := filter.LoadConfig(map[string]interface{}{
		"patterns": "not a slice",
	})
	if err == nil {
		t.Error("LoadConfig should return error for invalid patterns type")
	}
}
