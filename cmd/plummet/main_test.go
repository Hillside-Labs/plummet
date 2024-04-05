package main

import "github.com/stretchr/testify/assert"
import (
	"testing"
)

func TestCircularDependencyDetection(t *testing.T) {
	circularPlummetFile := PlummetFile{
		Targets: map[string]Target{
			"target1": {Deps: []string{"target2"}},
			"target2": {Deps: []string{"target1"}},
		},
	}

	visited := make(map[string]bool)
	db, _ := sql.Open("duckdb", ":memory:")
	err := executeTarget("target1", &circularPlummetFile, visited, db)
	assert.NotNil(t, err, "Expected a circular dependency error, but got none")
	assert.Contains(t, err.Error(), "circular dependency detected", "Expected a circular dependency error, got: %v", err)
}
