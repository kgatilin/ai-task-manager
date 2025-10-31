package task_manager

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kgatilin/darwinflow-pub/pkg/pluginsdk"
)

// ViewMode represents the current TUI screen
type ViewMode int

// View constants
const (
	ViewRoadmapList ViewMode = iota
	ViewTrackDetail
	ViewError
	ViewLoading
)

// AppModel is the main TUI application model
type AppModel struct {
	ctx        context.Context
	repository RoadmapRepository
	logger     pluginsdk.Logger

	// State
	currentView ViewMode
	error       error

	// Data
	roadmap     *RoadmapEntity
	tracks      []*TrackEntity
	currentTrack *TrackEntity
	tasks       []*TaskEntity

	// UI state
	selectedTrackIdx int
	selectedTaskIdx  int
	width            int
	height           int

	// Timestamps for debouncing
	lastUpdate time.Time
}

// Message types for Bubble Tea

type RoadmapLoadedMsg struct {
	Roadmap *RoadmapEntity
	Tracks  []*TrackEntity
	Error   error
}

type TrackDetailLoadedMsg struct {
	Track *TrackEntity
	Tasks []*TaskEntity
	Error error
}

type ErrorMsg struct {
	Error error
}

type BackMsg struct{}

// NewAppModel creates a new TUI app model
func NewAppModel(
	ctx context.Context,
	repository RoadmapRepository,
	logger pluginsdk.Logger,
) *AppModel {
	return &AppModel{
		ctx:             ctx,
		repository:      repository,
		logger:          logger,
		currentView:     ViewLoading,
		selectedTrackIdx: 0,
		selectedTaskIdx: 0,
		lastUpdate:      time.Now(),
	}
}

// Init initializes the application
func (m *AppModel) Init() tea.Cmd {
	return m.loadRoadmap
}

// loadRoadmap is a tea.Cmd that loads the active roadmap
func (m *AppModel) loadRoadmap() tea.Msg {
	roadmap, err := m.repository.GetActiveRoadmap(m.ctx)
	if err != nil {
		return RoadmapLoadedMsg{Error: err}
	}

	tracks, err := m.repository.ListTracks(m.ctx, roadmap.ID, TrackFilters{})
	if err != nil {
		return RoadmapLoadedMsg{Error: err}
	}

	return RoadmapLoadedMsg{
		Roadmap: roadmap,
		Tracks:  tracks,
	}
}

// loadTrackDetail is a tea.Cmd that loads track details and tasks
func (m *AppModel) loadTrackDetail(trackID string) tea.Cmd {
	return func() tea.Msg {
		track, err := m.repository.GetTrack(m.ctx, trackID)
		if err != nil {
			return TrackDetailLoadedMsg{Error: err}
		}

		tasks, err := m.repository.ListTasks(m.ctx, TaskFilters{TrackID: trackID})
		if err != nil {
			return TrackDetailLoadedMsg{Error: err}
		}

		return TrackDetailLoadedMsg{
			Track: track,
			Tasks: tasks,
		}
	}
}

// Update processes messages and updates state
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.currentView == ViewTrackDetail {
				m.currentView = ViewRoadmapList
				m.selectedTaskIdx = 0
				return m, nil
			}
		}

		// View-specific key handling
		switch m.currentView {
		case ViewRoadmapList:
			return m.handleRoadmapListKeys(msg)
		case ViewTrackDetail:
			return m.handleTrackDetailKeys(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case RoadmapLoadedMsg:
		if msg.Error != nil {
			m.currentView = ViewError
			m.error = msg.Error
			return m, nil
		}
		m.roadmap = msg.Roadmap
		m.tracks = msg.Tracks
		m.currentView = ViewRoadmapList
		m.selectedTrackIdx = 0
		m.lastUpdate = time.Now()

	case TrackDetailLoadedMsg:
		if msg.Error != nil {
			m.currentView = ViewError
			m.error = msg.Error
			return m, nil
		}
		m.currentTrack = msg.Track
		m.tasks = msg.Tasks
		m.currentView = ViewTrackDetail
		m.selectedTaskIdx = 0
		m.lastUpdate = time.Now()

	case ErrorMsg:
		m.currentView = ViewError
		m.error = msg.Error
	}

	return m, nil
}

// handleRoadmapListKeys processes key presses on roadmap list view
func (m *AppModel) handleRoadmapListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.selectedTrackIdx < len(m.tracks)-1 {
			m.selectedTrackIdx++
		}
	case "k", "up":
		if m.selectedTrackIdx > 0 {
			m.selectedTrackIdx--
		}
	case "enter":
		if m.selectedTrackIdx < len(m.tracks) {
			trackID := m.tracks[m.selectedTrackIdx].ID
			m.currentView = ViewLoading
			return m, m.loadTrackDetail(trackID)
		}
	case "r":
		m.currentView = ViewLoading
		return m, m.loadRoadmap
	}
	return m, nil
}

// handleTrackDetailKeys processes key presses on track detail view
func (m *AppModel) handleTrackDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.selectedTaskIdx < len(m.tasks)-1 {
			m.selectedTaskIdx++
		}
	case "k", "up":
		if m.selectedTaskIdx > 0 {
			m.selectedTaskIdx--
		}
	}
	return m, nil
}

// View renders the current view
func (m *AppModel) View() string {
	switch m.currentView {
	case ViewLoading:
		return m.renderLoading()
	case ViewError:
		return m.renderError()
	case ViewRoadmapList:
		return m.renderRoadmapList()
	case ViewTrackDetail:
		return m.renderTrackDetail()
	default:
		return "Unknown view"
	}
}

// renderLoading renders the loading screen
func (m *AppModel) renderLoading() string {
	return "Loading...\n\nPress q to quit"
}

// renderError renders the error screen
func (m *AppModel) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("196"))

	return errorStyle.Render(fmt.Sprintf("Error: %v", m.error)) + "\n\nPress esc to go back"
}

// renderRoadmapList renders the roadmap overview screen
func (m *AppModel) renderRoadmapList() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	trackItemStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		MarginBottom(0)

	selectedTrackStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		MarginBottom(0).
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("229"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Italic(true).
		MarginTop(1)

	var s string

	// Header
	if m.roadmap != nil {
		s += titleStyle.Render(fmt.Sprintf("Roadmap: %s", m.roadmap.ID)) + "\n"
		s += fmt.Sprintf("Vision: %s\n", m.roadmap.Vision)
		s += fmt.Sprintf("Success Criteria: %s\n", m.roadmap.SuccessCriteria)
		s += "\n"
	}

	// Tracks
	s += "Tracks:\n"
	if len(m.tracks) == 0 {
		s += "  No tracks yet\n"
	} else {
		for i, track := range m.tracks {
			statusIcon := getStatusIcon(track.Status)
			priorityIcon := getPriorityIcon(track.Priority)

			line := fmt.Sprintf("%s %s %s - %s", statusIcon, priorityIcon, track.ID, track.Title)

			if i == m.selectedTrackIdx {
				s += selectedTrackStyle.Render(line) + "\n"
			} else {
				s += trackItemStyle.Render(line) + "\n"
			}
		}
	}

	// Help text
	s += "\n"
	s += helpStyle.Render("Navigation: j/k or ↑/↓ | Enter: View track | r: Refresh | q: Quit")

	return s
}

// renderTrackDetail renders the track detail view
func (m *AppModel) renderTrackDetail() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Italic(true).
		MarginBottom(1)

	taskItemStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		MarginBottom(0)

	selectedTaskStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		MarginBottom(0).
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("229"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Italic(true).
		MarginTop(1)

	var s string

	if m.currentTrack == nil {
		return "Loading track details..."
	}

	// Header
	s += titleStyle.Render(fmt.Sprintf("Track: %s", m.currentTrack.ID)) + "\n"
	s += fmt.Sprintf("Title: %s\n", m.currentTrack.Title)
	s += fmt.Sprintf("Description: %s\n", m.currentTrack.Description)
	s += fmt.Sprintf("Status: %s | Priority: %s\n", m.currentTrack.Status, m.currentTrack.Priority)

	// Dependencies
	if len(m.currentTrack.Dependencies) > 0 {
		s += "\nDependencies:\n"
		for _, dep := range m.currentTrack.Dependencies {
			s += fmt.Sprintf("  → %s\n", dep)
		}
	}

	// Tasks
	s += fmt.Sprintf("\nTasks (%d):\n", len(m.tasks))
	if len(m.tasks) == 0 {
		s += subtitleStyle.Render("No tasks yet")
	} else {
		for i, task := range m.tasks {
			statusIcon := getStatusIcon(task.Status)
			priorityIcon := getPriorityIcon(task.Priority)

			line := fmt.Sprintf("%s %s %s - %s", statusIcon, priorityIcon, task.ID, task.Title)

			if i == m.selectedTaskIdx {
				s += selectedTaskStyle.Render(line) + "\n"
			} else {
				s += taskItemStyle.Render(line) + "\n"
			}
		}
	}

	// Help
	s += "\n"
	s += helpStyle.Render("Navigation: j/k or ↑/↓ | esc: Back | q: Quit")

	return s
}

// Helper functions for rendering

func getStatusIcon(status string) string {
	switch status {
	case "done", "complete":
		return "✓"
	case "in-progress":
		return "→"
	case "blocked":
		return "✗"
	case "waiting":
		return "⏸"
	default:
		return "○"
	}
}

func getPriorityIcon(priority string) string {
	switch priority {
	case "critical":
		return "🔴"
	case "high":
		return "🟠"
	case "medium":
		return "🟡"
	case "low":
		return "🟢"
	default:
		return "⚪"
	}
}

// Test helper methods - exported for testing

// SetRoadmap sets the roadmap for testing
func (m *AppModel) SetRoadmap(roadmap *RoadmapEntity) {
	m.roadmap = roadmap
}

// SetTracks sets the tracks for testing
func (m *AppModel) SetTracks(tracks []*TrackEntity) {
	m.tracks = tracks
}

// SetTasks sets the tasks for testing
func (m *AppModel) SetTasks(tasks []*TaskEntity) {
	m.tasks = tasks
}

// SetCurrentTrack sets the current track for testing
func (m *AppModel) SetCurrentTrack(track *TrackEntity) {
	m.currentTrack = track
}

// SetCurrentView sets the current view for testing
func (m *AppModel) SetCurrentView(view ViewMode) {
	m.currentView = view
}

// SetError sets the error for testing
func (m *AppModel) SetError(err error) {
	m.error = err
}

// SetDimensions sets the width and height for testing
func (m *AppModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}
