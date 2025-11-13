package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lando/cami/internal/discovery"
)

func (m Model) viewAgentSelection() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("CAMI - Claude Agent Management Interface"))
	b.WriteString("\n\n")

	// Agent list
	b.WriteString("Select agents to deploy:\n\n")

	// Calculate viewport parameters
	overhead := 9
	if m.message != "" {
		overhead += 2
	}

	// Group agents by category for display
	type displayItem struct {
		isCategory  bool
		categoryName string
		agentIdx    int
	}

	var displayItems []displayItem
	categoryOrder := []string{"core", "specialized", "infrastructure", "integration", "design", "meta", "uncategorized"}

	for _, category := range categoryOrder {
		var categoryAgents []int
		for i, ag := range m.agents {
			agentCategory := ag.Category
			if agentCategory == "" {
				agentCategory = "uncategorized"
			}
			if agentCategory == category {
				categoryAgents = append(categoryAgents, i)
			}
		}

		if len(categoryAgents) > 0 {
			// Add category header
			displayCategory := category
			if category != "uncategorized" {
				displayCategory = strings.ToUpper(category[:1]) + category[1:]
			} else {
				displayCategory = "Uncategorized"
			}
			displayItems = append(displayItems, displayItem{
				isCategory:  true,
				categoryName: displayCategory,
			})

			// Add agents in this category
			for _, idx := range categoryAgents {
				displayItems = append(displayItems, displayItem{
					isCategory: false,
					agentIdx:   idx,
				})
			}
		}
	}

	// Calculate lines needed per item (category headers are 1 line, agents can be 2 lines with description)
	maxVisibleLines := m.height - overhead
	if maxVisibleLines < 1 {
		maxVisibleLines = 1
	}

	// Show scroll indicator at top if not at beginning
	if m.viewportOffset > 0 {
		b.WriteString(versionStyle.Render(fmt.Sprintf("  ↑ scroll up...\n")))
	}

	// Render items starting from viewportOffset
	currentLine := 0
	currentDisplayIdx := m.viewportOffset

	for currentDisplayIdx < len(displayItems) && currentLine < maxVisibleLines {
		item := displayItems[currentDisplayIdx]

		if item.isCategory {
			// Render category header
			b.WriteString(titleStyle.Render(fmt.Sprintf("── %s ──", item.categoryName)))
			b.WriteString("\n")
			currentLine++
		} else {
			// Render agent
			i := item.agentIdx
			ag := m.agents[i]
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			checkbox := "[ ]"
			if m.selectedAgents[i] {
				checkbox = checkboxStyle.Render("[✓]")
			}

			name := ag.Name
			if m.cursor == i {
				name = selectedStyle.Render(name)
			}

			version := versionStyle.Render(fmt.Sprintf("v%s", ag.Version))

			line := fmt.Sprintf("%s %s %s %s", cursor, checkbox, name, version)
			b.WriteString(line)
			b.WriteString("\n")
			currentLine++

			// Show description for selected item
			if m.cursor == i && currentLine < maxVisibleLines {
				desc := "  " + ag.Description
				if len(desc) > 80 {
					desc = desc[:77] + "..."
				}
				b.WriteString(versionStyle.Render(desc))
				b.WriteString("\n")
				currentLine++
			}
		}

		currentDisplayIdx++
	}

	// Show scroll indicator at bottom if more items below
	if currentDisplayIdx < len(displayItems) {
		b.WriteString(versionStyle.Render(fmt.Sprintf("  ↓ scroll down...\n")))
	}

	// Message
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(warningStyle.Render(m.message))
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("space: select  •  enter: deploy  •  l: locations  •  i: discovery  •  q: quit"))

	return b.String()
}

func (m Model) viewLocationManagement() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Deployment Locations"))
	b.WriteString("\n\n")

	if m.addingLocation {
		b.WriteString("Add New Location:\n\n")

		nameField := m.newLocationName
		pathField := m.newLocationPath

		if m.inputField == 0 {
			nameField = selectedStyle.Render(nameField + "_")
		}
		if m.inputField == 1 {
			pathField = selectedStyle.Render(pathField + "_")
		}

		b.WriteString(fmt.Sprintf("Name: %s\n", nameField))
		b.WriteString(fmt.Sprintf("Path: %s\n", pathField))

		b.WriteString("\n")
		b.WriteString(helpStyle.Render("tab: switch field  •  enter: save  •  esc: cancel"))
	} else {
		if len(m.config.Locations) == 0 {
			b.WriteString("No deployment locations configured.\n\n")
			b.WriteString("Press 'a' to add a location.\n")
		} else {
			// Calculate viewport parameters for locations
			overhead := 5
			if m.message != "" {
				overhead += 2
			}
			linesPerLocation := 1
			// +1 for "Add new location" option
			totalItems := len(m.config.Locations) + 1
			maxVisibleLocations := (m.height - overhead) / linesPerLocation
			if maxVisibleLocations < 1 {
				maxVisibleLocations = 1
			}

			// Calculate visible range
			startIdx := m.locationViewportOffset
			endIdx := m.locationViewportOffset + maxVisibleLocations
			if endIdx > totalItems {
				endIdx = totalItems
			}

			// Show scroll indicator at top if not at beginning
			if startIdx > 0 {
				b.WriteString(versionStyle.Render(fmt.Sprintf("  ↑ %d more above...\n", startIdx)))
			}

			// Render visible locations
			for i := startIdx; i < endIdx && i < len(m.config.Locations); i++ {
				loc := m.config.Locations[i]
				cursor := " "
				if m.locationCursor == i {
					cursor = ">"
				}

				name := loc.Name
				if m.locationCursor == i {
					name = selectedStyle.Render(name)
				}

				path := versionStyle.Render(loc.Path)
				b.WriteString(fmt.Sprintf("%s %s - %s\n", cursor, name, path))
			}

			// Add "New Location" option if visible in viewport
			addLocationIdx := len(m.config.Locations)
			if addLocationIdx >= startIdx && addLocationIdx < endIdx {
				cursor := " "
				if m.locationCursor == addLocationIdx {
					cursor = ">"
				}
				addText := "+ Add new location"
				if m.locationCursor == addLocationIdx {
					addText = selectedStyle.Render(addText)
				}
				b.WriteString(fmt.Sprintf("\n%s %s\n", cursor, addText))
			}

			// Show scroll indicator at bottom if more items below
			if endIdx < totalItems {
				b.WriteString(versionStyle.Render(fmt.Sprintf("  ↓ %d more below...\n", totalItems-endIdx)))
			}
		}

		if m.message != "" {
			b.WriteString("\n")
			b.WriteString(warningStyle.Render(m.message))
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(helpStyle.Render("a: add  •  d: delete  •  esc: back  •  q: quit"))
	}

	return b.String()
}

func (m Model) viewDeployment() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Select Deployment Location"))
	b.WriteString("\n\n")

	// Show selected agents
	var selectedNames []string
	for idx, selected := range m.selectedAgents {
		if selected {
			selectedNames = append(selectedNames, m.agents[idx].Name)
		}
	}
	b.WriteString(fmt.Sprintf("Deploying: %s\n\n", strings.Join(selectedNames, ", ")))

	if len(m.config.Locations) == 0 {
		b.WriteString(errorStyle.Render("No deployment locations configured."))
		b.WriteString("\n\n")
		b.WriteString("Press 'esc' to go back and configure locations.\n")
	} else {
		b.WriteString("Select destination:\n\n")

		for i, loc := range m.config.Locations {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			name := loc.Name
			if m.cursor == i {
				name = selectedStyle.Render(name)
			}

			path := versionStyle.Render(loc.Path)
			b.WriteString(fmt.Sprintf("%s %s - %s\n", cursor, name, path))
		}

		b.WriteString("\n")
		b.WriteString(helpStyle.Render("enter: deploy  •  esc: back  •  q: quit"))
	}

	return b.String()
}

func (m Model) viewResults() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Deployment Results"))
	b.WriteString("\n\n")

	successCount := 0
	conflictCount := 0
	errorCount := 0

	for _, result := range m.deployResults {
		var status string
		if result.Success {
			status = successStyle.Render("✓")
			successCount++
		} else if result.Conflict {
			status = warningStyle.Render("⚠")
			conflictCount++
		} else {
			status = errorStyle.Render("✗")
			errorCount++
		}

		b.WriteString(fmt.Sprintf("%s %s: %s\n", status, result.Agent.Name, result.Message))
	}

	// Summary
	b.WriteString("\n")
	b.WriteString("Summary:\n")
	if successCount > 0 {
		b.WriteString(successStyle.Render(fmt.Sprintf("  ✓ %d deployed successfully\n", successCount)))
	}
	if conflictCount > 0 {
		b.WriteString(warningStyle.Render(fmt.Sprintf("  ⚠ %d conflicts (files already exist)\n", conflictCount)))
	}
	if errorCount > 0 {
		b.WriteString(errorStyle.Render(fmt.Sprintf("  ✗ %d errors\n", errorCount)))
	}

	if conflictCount > 0 {
		b.WriteString("\n")
		b.WriteString(versionStyle.Render("Note: Files with conflicts were not overwritten."))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("enter: return  •  q: quit"))

	return b.String()
}

func (m Model) viewDiscovery() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Agent Discovery & Updates"))
	b.WriteString("\n\n")

	if m.discoveryLoading {
		b.WriteString("Scanning deployment locations...\n")
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("esc: back  •  q: quit"))
		return b.String()
	}

	if m.discoveryResult == nil || len(m.discoveryResult.LocationStatuses) == 0 {
		b.WriteString(errorStyle.Render("No deployment locations configured."))
		b.WriteString("\n\n")
		b.WriteString("Press 'l' to configure locations.\n")
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("esc: back  •  q: quit"))
		return b.String()
	}

	// Location tabs
	b.WriteString("Locations: ")
	for i, locStatus := range m.discoveryResult.LocationStatuses {
		name := locStatus.Location.Name
		if i == m.discoveryLocationIdx {
			name = selectedStyle.Render("[" + name + "]")
		} else {
			name = versionStyle.Render(" " + name + " ")
		}
		b.WriteString(name)
		if i < len(m.discoveryResult.LocationStatuses)-1 {
			b.WriteString(" ")
		}
	}
	b.WriteString("\n\n")

	// Current location status
	if m.discoveryLocationIdx < len(m.discoveryResult.LocationStatuses) {
		locStatus := m.discoveryResult.LocationStatuses[m.discoveryLocationIdx]

		// Count statuses
		upToDate := 0
		updateAvailable := 0
		notDeployed := 0
		for _, agentStatus := range locStatus.AgentStatuses {
			switch agentStatus.Status {
			case discovery.StatusUpToDate:
				upToDate++
			case discovery.StatusUpdateAvailable:
				updateAvailable++
			case discovery.StatusNotDeployed:
				notDeployed++
			}
		}

		// Summary
		b.WriteString(fmt.Sprintf("Status: %s up-to-date, %s updates available, %s not deployed\n\n",
			successStyle.Render(fmt.Sprintf("%d", upToDate)),
			warningStyle.Render(fmt.Sprintf("%d", updateAvailable)),
			versionStyle.Render(fmt.Sprintf("%d", notDeployed))))

		// Agent list with scrolling
		b.WriteString("Agents:\n")

		// Calculate viewport parameters for discovery
		overhead := 12
		if m.message != "" {
			overhead += 2
		}
		linesPerAgent := 1
		maxVisibleAgents := (m.height - overhead) / linesPerAgent
		if maxVisibleAgents < 1 {
			maxVisibleAgents = 1
		}

		// Calculate visible range
		startIdx := m.discoveryViewportOffset
		endIdx := m.discoveryViewportOffset + maxVisibleAgents
		if endIdx > len(locStatus.AgentStatuses) {
			endIdx = len(locStatus.AgentStatuses)
		}

		// Show scroll indicator at top if not at beginning
		if startIdx > 0 {
			b.WriteString(versionStyle.Render(fmt.Sprintf("  ↑ %d more above...\n", startIdx)))
		}

		// Render visible agents
		for i := startIdx; i < endIdx; i++ {
			agentStatus := locStatus.AgentStatuses[i]
			cursor := " "
			if i == m.discoveryAgentIdx {
				cursor = ">"
			}

			symbol := discovery.GetStatusSymbol(agentStatus.Status)
			var statusStyle lipgloss.Style
			switch agentStatus.Status {
			case discovery.StatusUpToDate:
				statusStyle = successStyle
			case discovery.StatusUpdateAvailable:
				statusStyle = warningStyle
			case discovery.StatusNotDeployed:
				statusStyle = versionStyle
			}

			name := agentStatus.Agent.Name
			if i == m.discoveryAgentIdx {
				name = selectedStyle.Render(name)
			}

			versionInfo := ""
			if agentStatus.DeployedVersion != "" {
				if agentStatus.Status == discovery.StatusUpdateAvailable {
					versionInfo = fmt.Sprintf("v%s → v%s", agentStatus.DeployedVersion, agentStatus.AvailableVersion)
				} else {
					versionInfo = fmt.Sprintf("v%s", agentStatus.DeployedVersion)
				}
			} else {
				versionInfo = fmt.Sprintf("v%s available", agentStatus.AvailableVersion)
			}

			line := fmt.Sprintf("%s %s %-20s %s", cursor, statusStyle.Render(symbol), name, versionStyle.Render(versionInfo))
			b.WriteString(line)
			b.WriteString("\n")
		}

		// Show scroll indicator at bottom if more items below
		if endIdx < len(locStatus.AgentStatuses) {
			b.WriteString(versionStyle.Render(fmt.Sprintf("  ↓ %d more below...\n", len(locStatus.AgentStatuses)-endIdx)))
		}
	}

	// Message
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(successStyle.Render(m.message))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓/j/k: navigate  •  ←/→/h/l: switch location  •  u: update  •  U: update all  •  r: refresh  •  esc: back"))

	return b.String()
}

