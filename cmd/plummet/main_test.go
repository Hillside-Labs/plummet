package main

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
	err := executeTarget("target1", &circularPlummetFile, visited)
	if err == nil {
		t.Errorf("Expected a circular dependency error, but got none")
	} else if !strings.Contains(err.Error(), "circular dependency detected") {
		t.Errorf("Expected a circular dependency error, got: %v", err)
	}
}
