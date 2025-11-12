package docs

import (
	"strings"
	"testing"

	"github.com/lando/cami/internal/agent"
)

func TestGenerateAgentSectionWithTimestamp(t *testing.T) {
	// Create a test agent
	agents := []*agent.Agent{
		{
			Name:        "test-agent",
			Version:     "1.0.0",
			Description: "A test agent for validation",
		},
	}

	// Generate the section
	section := generateAgentSection("Test Agents", agents)

	// Verify the timestamp is present in the start marker
	if !strings.Contains(section, "Last Updated:") {
		t.Error("Expected timestamp in start marker, but not found")
	}

	// Verify the marker format matches the pattern
	if !sectionMarkerPattern.MatchString(section) {
		t.Error("Generated marker does not match expected pattern")
	}

	// Verify the end marker is still present
	if !strings.Contains(section, sectionMarkerEnd) {
		t.Error("Expected end marker not found")
	}

	t.Logf("Generated section:\n%s", section)
}

func TestMergeContentBackwardCompatibility(t *testing.T) {
	// Test with old marker format (no timestamp)
	oldContent := `# Project Documentation

<!-- CAMI-MANAGED: DEPLOYED-AGENTS -->
## Deployed Agents

Old content here

<!-- /CAMI-MANAGED: DEPLOYED-AGENTS -->

Some other content
`

	newSection := `<!-- CAMI-MANAGED: DEPLOYED-AGENTS | Last Updated: 2025-10-09T14:30:00-05:00 -->
## Deployed Agents

New content here

<!-- /CAMI-MANAGED: DEPLOYED-AGENTS -->
`

	result := mergeContent(oldContent, newSection)

	// Verify old marker was replaced
	if strings.Contains(result, "Old content here") {
		t.Error("Old content was not replaced")
	}

	// Verify new content is present
	if !strings.Contains(result, "New content here") {
		t.Error("New content was not inserted")
	}

	// Verify timestamp is present
	if !strings.Contains(result, "Last Updated:") {
		t.Error("Timestamp not found in merged content")
	}

	// Verify other content was preserved
	if !strings.Contains(result, "Some other content") {
		t.Error("Non-managed content was not preserved")
	}

	t.Logf("Merged content:\n%s", result)
}

func TestMergeContentWithNewMarkerFormat(t *testing.T) {
	// Test with new marker format (with timestamp)
	existingContent := `# Project Documentation

<!-- CAMI-MANAGED: DEPLOYED-AGENTS | Last Updated: 2025-10-08T10:00:00-05:00 -->
## Deployed Agents

Previous content

<!-- /CAMI-MANAGED: DEPLOYED-AGENTS -->

Some other content
`

	newSection := `<!-- CAMI-MANAGED: DEPLOYED-AGENTS | Last Updated: 2025-10-09T14:30:00-05:00 -->
## Deployed Agents

Updated content

<!-- /CAMI-MANAGED: DEPLOYED-AGENTS -->
`

	result := mergeContent(existingContent, newSection)

	// Verify old timestamp was replaced
	if strings.Contains(result, "2025-10-08T10:00:00-05:00") {
		t.Error("Old timestamp was not replaced")
	}

	// Verify new timestamp is present
	if !strings.Contains(result, "2025-10-09T14:30:00-05:00") {
		t.Error("New timestamp not found")
	}

	// Verify content was updated
	if !strings.Contains(result, "Updated content") {
		t.Error("Content was not updated")
	}

	t.Logf("Merged content:\n%s", result)
}
