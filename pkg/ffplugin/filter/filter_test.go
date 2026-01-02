// Copyright 2026 Paul Sade
// GPLv3 - See LICENSE for details.


// pkg/ffplugin/filter/filter_test.go
package filter

import (
	"fmt"
	"io/fs"
	"testing"
)

// MockFilter is a mock implementation of the Filter interface for testing.
type MockFilter struct{}

func (m *MockFilter) Match(path string, info fs.FileInfo) (bool, error) {
	return true, nil
}

func (m *MockFilter) Selector() string {
	return "mock"
}

func (m *MockFilter) LoadConfig(config map[string]interface{}) error {
	return nil
}

func TestRegisterFilter(t *testing.T) {
	// Clear the registry before testing
	filterRegistry = make(map[string]func() Filter)

	// Test: Successful registration
	RegisterFilter("mock", func() Filter { return &MockFilter{} })
	if _, exists := filterRegistry["mock"]; !exists {
		t.Error("Filter was not registered")
	}
}

func TestRegisterFilterEmptyName(t *testing.T) {
	// Clear the registry before testing
	filterRegistry = make(map[string]func() Filter)

	// Test: Register with empty name (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Error("RegisterFilter should panic for empty name")
		} else if r != "filter name cannot be empty" {
			t.Errorf("Unexpected panic message: %v", r)
		}
	}()
	RegisterFilter("", func() Filter { return &MockFilter{} })
}

func TestRegisterFilterDuplicate(t *testing.T) {
	// Clear the registry before testing
	filterRegistry = make(map[string]func() Filter)

	// Register a filter
	RegisterFilter("mock", func() Filter { return &MockFilter{} })

	// Test: Register duplicate filter (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Error("RegisterFilter should panic for duplicate filter")
		} else if r != fmt.Sprintf("filter '%s' is already registered", "mock") {
			t.Errorf("Unexpected panic message: %v", r)
		}
	}()
	RegisterFilter("mock", func() Filter { return &MockFilter{} })
}

// TestNewFilterSuccess tests the successful creation of a filter.
func TestNewFilterSuccess(t *testing.T) {
	// Clear the registry before testing
	filterRegistry = make(map[string]func() Filter)

	// Register a mock filter
	RegisterFilter("mock", func() Filter { return &MockFilter{} })

	// Test: Create the filter
	filter, err := NewFilter("mock")
	if err != nil {
		t.Fatalf("NewFilter failed: %v", err)
	}

	// Verify the filter is of the correct type
	if _, ok := filter.(*MockFilter); !ok {
		t.Error("NewFilter did not return a *MockFilter")
	}
}

// TestNewFilterUnknown tests the error case when the filter is unknown.
func TestNewFilterUnknown(t *testing.T) {
	// Clear the registry before testing
	filterRegistry = make(map[string]func() Filter)

	// Test: Try to create an unknown filter
	_, err := NewFilter("unknown")
	if err == nil {
		t.Fatal("NewFilter should return an error for unknown filter")
	}

	// Verify the error message
	expectedError := "unknown filter: unknown"
	if err.Error() != expectedError {
		t.Errorf("NewFilter returned error %q, expected %q", err.Error(), expectedError)
	}
}

// --- End of pkg/ffplugin/filter/filter_test.go ---
