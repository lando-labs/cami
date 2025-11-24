package agent

import (
	"testing"
)

func TestGetPhaseWeightsByClass(t *testing.T) {
	tests := []struct {
		name     string
		class    string
		expected PhaseWeights
	}{
		{
			name:  "workflow-specialist",
			class: "workflow-specialist",
			expected: PhaseWeights{
				Research: 15,
				Execute:  70,
				Validate: 15,
			},
		},
		{
			name:  "technology-implementer",
			class: "technology-implementer",
			expected: PhaseWeights{
				Research: 30,
				Execute:  55,
				Validate: 15,
			},
		},
		{
			name:  "strategic-planner",
			class: "strategic-planner",
			expected: PhaseWeights{
				Research: 45,
				Execute:  30,
				Validate: 25,
			},
		},
		{
			name:  "unknown class returns default",
			class: "unknown",
			expected: PhaseWeights{
				Research: 30,
				Execute:  50,
				Validate: 20,
			},
		},
		{
			name:  "empty class returns default",
			class: "",
			expected: PhaseWeights{
				Research: 30,
				Execute:  50,
				Validate: 20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPhaseWeightsByClass(tt.class)
			if result != tt.expected {
				t.Errorf("GetPhaseWeightsByClass(%q) = %+v, want %+v", tt.class, result, tt.expected)
			}
		})
	}
}

func TestGetUserFriendlyClassName(t *testing.T) {
	tests := []struct {
		name     string
		class    string
		expected string
	}{
		{
			name:     "workflow-specialist",
			class:    "workflow-specialist",
			expected: "Task Automator",
		},
		{
			name:     "technology-implementer",
			class:    "technology-implementer",
			expected: "Feature Builder",
		},
		{
			name:     "strategic-planner",
			class:    "strategic-planner",
			expected: "System Architect",
		},
		{
			name:     "unknown class returns default",
			class:    "unknown",
			expected: "Agent",
		},
		{
			name:     "empty class returns default",
			class:    "",
			expected: "Agent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUserFriendlyClassName(tt.class)
			if result != tt.expected {
				t.Errorf("GetUserFriendlyClassName(%q) = %q, want %q", tt.class, result, tt.expected)
			}
		})
	}
}

func TestAgent_GetPhaseWeights(t *testing.T) {
	tests := []struct {
		name     string
		agent    *Agent
		expected PhaseWeights
	}{
		{
			name: "workflow-specialist agent",
			agent: &Agent{
				Name:  "k8s-pod-checker",
				Class: "workflow-specialist",
			},
			expected: PhaseWeights{
				Research: 15,
				Execute:  70,
				Validate: 15,
			},
		},
		{
			name: "technology-implementer agent",
			agent: &Agent{
				Name:  "express-api",
				Class: "technology-implementer",
			},
			expected: PhaseWeights{
				Research: 30,
				Execute:  55,
				Validate: 15,
			},
		},
		{
			name: "strategic-planner agent",
			agent: &Agent{
				Name:  "enterprise-architect",
				Class: "strategic-planner",
			},
			expected: PhaseWeights{
				Research: 45,
				Execute:  30,
				Validate: 25,
			},
		},
		{
			name: "agent without class gets default weights",
			agent: &Agent{
				Name: "legacy-agent",
			},
			expected: PhaseWeights{
				Research: 30,
				Execute:  50,
				Validate: 20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.agent.GetPhaseWeights()
			if result != tt.expected {
				t.Errorf("%s.GetPhaseWeights() = %+v, want %+v", tt.agent.Name, result, tt.expected)
			}
		})
	}
}
