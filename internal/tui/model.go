package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lando/cami/internal/agent"
	"github.com/lando/cami/internal/config"
	"github.com/lando/cami/internal/deploy"
	"github.com/lando/cami/internal/discovery"
)

// ViewState represents the current view
type ViewState int

const (
	ViewAgentSelection ViewState = iota
	ViewLocationManagement
	ViewDeployment
	ViewResults
	ViewDiscovery
)

// Model represents the TUI state
type Model struct {
	state           ViewState
	agents          []*agent.Agent
	selectedAgents  map[int]bool
	config          *config.Config
	deployLocation  *config.DeployLocation
	deployResults   []*deploy.Result
	width           int
	height          int
	cursor          int
	viewportOffset  int // For scrolling the agent list
	message         string
	err             error

	// Location management
	locationCursor  int
	addingLocation  bool
	newLocationName string
	newLocationPath string
	inputField      int // 0 = name, 1 = path

	// Discovery
	discoveryResult    *discovery.DiscoveryResult
	discoveryLocationIdx int
	discoveryAgentIdx   int
	discoveryLoading    bool
}

// scanCompleteMsg is sent when discovery scan completes
type scanCompleteMsg struct {
	result *discovery.DiscoveryResult
	err    error
}

// keyMap defines keyboard shortcuts
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Select   key.Binding
	Deploy   key.Binding
	Locations key.Binding
	Quit     key.Binding
	Back     key.Binding
	Help     key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys(" ", "x"),
		key.WithHelp("space/x", "select"),
	),
	Deploy: key.NewBinding(
		key.WithKeys("enter", "d"),
		key.WithHelp("enter/d", "deploy"),
	),
	Locations: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "locations"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	checkboxStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	versionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	warningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)
)

// NewModel creates a new TUI model
func NewModel(agents []*agent.Agent, cfg *config.Config) Model {
	return Model{
		state:          ViewAgentSelection,
		agents:         agents,
		selectedAgents: make(map[int]bool),
		config:         cfg,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case scanCompleteMsg:
		m.discoveryLoading = false
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("Error scanning locations: %v", msg.err)
		} else {
			m.discoveryResult = msg.result
		}
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case ViewAgentSelection:
			return m.updateAgentSelection(msg)
		case ViewLocationManagement:
			return m.updateLocationManagement(msg)
		case ViewDeployment:
			return m.updateDeployment(msg)
		case ViewResults:
			return m.updateResults(msg)
		case ViewDiscovery:
			return m.updateDiscovery(msg)
		}
	}

	return m, nil
}

// View renders the current view
func (m Model) View() string {
	switch m.state {
	case ViewAgentSelection:
		return m.viewAgentSelection()
	case ViewLocationManagement:
		return m.viewLocationManagement()
	case ViewDeployment:
		return m.viewDeployment()
	case ViewResults:
		return m.viewResults()
	case ViewDiscovery:
		return m.viewDiscovery()
	default:
		return "Unknown view"
	}
}

func (m Model) updateAgentSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
			m.adjustViewport()
		}
	case key.Matches(msg, keys.Down):
		if m.cursor < len(m.agents)-1 {
			m.cursor++
			m.adjustViewport()
		}
	case key.Matches(msg, keys.Select):
		m.selectedAgents[m.cursor] = !m.selectedAgents[m.cursor]
	case key.Matches(msg, keys.Locations):
		m.state = ViewLocationManagement
		m.locationCursor = 0
	case key.Matches(msg, keys.Deploy):
		// Check if any agents are selected
		if len(m.selectedAgents) == 0 {
			m.message = "Please select at least one agent"
			return m, nil
		}
		m.state = ViewDeployment
		m.cursor = 0
	case msg.String() == "i":
		// Enter discovery view and trigger scan
		m.state = ViewDiscovery
		m.discoveryLocationIdx = 0
		m.discoveryAgentIdx = 0
		m.discoveryLoading = true
		return m, m.scanLocations()
	}
	return m, nil
}

// adjustViewport ensures the cursor is visible within the viewport
func (m *Model) adjustViewport() {
	// Calculate available height for agent list
	// Title (3 lines) + "Select agents" (2 lines) + message (2 lines if present) + help (2 lines) = ~9 lines overhead
	overhead := 9
	if m.message != "" {
		overhead += 2
	}

	// Each agent takes 2 lines (main line + description when selected)
	linesPerAgent := 2

	maxVisibleAgents := (m.height - overhead) / linesPerAgent
	if maxVisibleAgents < 1 {
		maxVisibleAgents = 1
	}

	// Adjust viewport to keep cursor visible
	if m.cursor < m.viewportOffset {
		// Cursor moved above viewport
		m.viewportOffset = m.cursor
	} else if m.cursor >= m.viewportOffset+maxVisibleAgents {
		// Cursor moved below viewport
		m.viewportOffset = m.cursor - maxVisibleAgents + 1
	}
}

func (m Model) updateLocationManagement(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle text input mode first - when adding a location
	if m.addingLocation {
		switch msg.String() {
		case "esc":
			// Allow ESC to cancel adding location
			m.addingLocation = false
			m.newLocationName = ""
			m.newLocationPath = ""
		case "tab":
			// Tab switches between name and path fields
			m.inputField = (m.inputField + 1) % 2
		case "enter":
			// Enter saves the new location
			if m.newLocationName != "" && m.newLocationPath != "" {
				if err := m.config.AddDeployLocation(m.newLocationName, m.newLocationPath); err != nil {
					m.message = fmt.Sprintf("Error: %v", err)
				} else {
					m.config.Save()
					m.addingLocation = false
					m.message = "Location added successfully"
				}
			}
		case "backspace":
			// Handle backspace for text deletion
			if m.inputField == 0 && len(m.newLocationName) > 0 {
				m.newLocationName = m.newLocationName[:len(m.newLocationName)-1]
			} else if m.inputField == 1 && len(m.newLocationPath) > 0 {
				m.newLocationPath = m.newLocationPath[:len(m.newLocationPath)-1]
			}
		default:
			// Handle regular text input (including 'a', 'q', 'd', etc.)
			if len(msg.String()) == 1 {
				if m.inputField == 0 {
					m.newLocationName += msg.String()
				} else {
					m.newLocationPath += msg.String()
				}
			}
		}
		return m, nil
	}

	// Only process global shortcuts when NOT in text input mode
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, keys.Back):
		m.state = ViewAgentSelection
	case key.Matches(msg, keys.Up):
		if m.locationCursor > 0 {
			m.locationCursor--
		}
	case key.Matches(msg, keys.Down):
		if m.locationCursor < len(m.config.DeployLocations) {
			m.locationCursor++
		}
	case msg.String() == "a":
		// Start adding a new location
		m.addingLocation = true
		m.newLocationName = ""
		m.newLocationPath = ""
		m.inputField = 0
	case msg.String() == "d":
		// Delete the selected location
		if m.locationCursor < len(m.config.DeployLocations) {
			m.config.RemoveDeployLocation(m.locationCursor)
			m.config.Save()
			if m.locationCursor > 0 {
				m.locationCursor--
			}
		}
	}
	return m, nil
}

func (m Model) updateDeployment(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, keys.Back):
		m.state = ViewAgentSelection
	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, keys.Down):
		if m.cursor < len(m.config.DeployLocations)-1 {
			m.cursor++
		}
	case key.Matches(msg, keys.Deploy), msg.String() == "enter":
		if m.cursor < len(m.config.DeployLocations) {
			m.deployLocation = &m.config.DeployLocations[m.cursor]

			// Get selected agents
			var selectedAgents []*agent.Agent
			for idx, selected := range m.selectedAgents {
				if selected {
					selectedAgents = append(selectedAgents, m.agents[idx])
				}
			}

			// Deploy agents
			results, err := deploy.DeployAgents(selectedAgents, m.deployLocation.Path, false)
			if err != nil {
				m.err = err
				m.state = ViewAgentSelection
				return m, nil
			}

			m.deployResults = results
			m.state = ViewResults
		}
	}
	return m, nil
}

func (m Model) updateResults(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, keys.Back), msg.String() == "enter":
		m.state = ViewAgentSelection
		m.selectedAgents = make(map[int]bool)
		m.deployResults = nil
	}
	return m, nil
}

func (m Model) updateDiscovery(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.discoveryLoading || m.discoveryResult == nil {
		// Allow quitting or going back while loading
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Back):
			m.state = ViewAgentSelection
		}
		return m, nil
	}

	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, keys.Back):
		m.state = ViewAgentSelection
	case key.Matches(msg, keys.Up), msg.String() == "k":
		if m.discoveryAgentIdx > 0 {
			m.discoveryAgentIdx--
		}
	case key.Matches(msg, keys.Down), msg.String() == "j":
		if len(m.discoveryResult.AvailableAgents) > 0 && m.discoveryAgentIdx < len(m.discoveryResult.AvailableAgents)-1 {
			m.discoveryAgentIdx++
		}
	case msg.String() == "h", msg.String() == "left":
		if m.discoveryLocationIdx > 0 {
			m.discoveryLocationIdx--
		}
	case msg.String() == "l", msg.String() == "right":
		if len(m.discoveryResult.LocationStatuses) > 0 && m.discoveryLocationIdx < len(m.discoveryResult.LocationStatuses)-1 {
			m.discoveryLocationIdx++
		}
	case msg.String() == "r":
		// Refresh/rescan
		m.discoveryLoading = true
		return m, m.scanLocations()
	case msg.String() == "u":
		// Update single agent
		if len(m.discoveryResult.LocationStatuses) > 0 && m.discoveryLocationIdx < len(m.discoveryResult.LocationStatuses) {
			locStatus := m.discoveryResult.LocationStatuses[m.discoveryLocationIdx]
			if m.discoveryAgentIdx < len(locStatus.AgentStatuses) {
				agentStatus := locStatus.AgentStatuses[m.discoveryAgentIdx]
				if agentStatus.Status == discovery.StatusUpdateAvailable || agentStatus.Status == discovery.StatusNotDeployed {
					// Deploy single agent
					results, err := deploy.DeployAgents([]*agent.Agent{agentStatus.Agent}, locStatus.Location.Path, false)
					if err == nil && len(results) > 0 && results[0].Success {
						m.message = fmt.Sprintf("Updated %s to v%s", agentStatus.Agent.Name, agentStatus.Agent.Version)
						// Rescan to update status
						m.discoveryLoading = true
						return m, m.scanLocations()
					}
				}
			}
		}
	case msg.String() == "U":
		// Update all agents at current location
		if len(m.discoveryResult.LocationStatuses) > 0 && m.discoveryLocationIdx < len(m.discoveryResult.LocationStatuses) {
			locStatus := m.discoveryResult.LocationStatuses[m.discoveryLocationIdx]
			var agentsToUpdate []*agent.Agent
			for _, agentStatus := range locStatus.AgentStatuses {
				if agentStatus.Status == discovery.StatusUpdateAvailable || agentStatus.Status == discovery.StatusNotDeployed {
					agentsToUpdate = append(agentsToUpdate, agentStatus.Agent)
				}
			}
			if len(agentsToUpdate) > 0 {
				_, err := deploy.DeployAgents(agentsToUpdate, locStatus.Location.Path, false)
				if err == nil {
					m.message = fmt.Sprintf("Updated %d agents", len(agentsToUpdate))
					// Rescan to update status
					m.discoveryLoading = true
					return m, m.scanLocations()
				}
			}
		}
	}
	return m, nil
}

// scanLocations performs an asynchronous scan of all locations
func (m Model) scanLocations() tea.Cmd {
	return func() tea.Msg {
		result, err := discovery.ScanAllLocations(m.config.DeployLocations, m.agents)
		return scanCompleteMsg{
			result: result,
			err:    err,
		}
	}
}
