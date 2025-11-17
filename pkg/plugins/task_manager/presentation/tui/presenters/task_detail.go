package presenters

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/components"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
)

// TaskDetailKeyMap defines keybindings for task detail view
type TaskDetailKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding // Expand/collapse AC testing instructions
	Quit     key.Binding
	Back     key.Binding
	Help     key.Binding
	Verify   key.Binding // Space - verify AC
	Skip     key.Binding // s - skip AC
	Fail     key.Binding // f - fail AC with feedback
	PageUp   key.Binding // pgup/b - page up
	PageDown key.Binding // pgdn - page down
}

// NewTaskDetailKeyMap creates default keybindings for task detail
func NewTaskDetailKeyMap() TaskDetailKeyMap {
	return TaskDetailKeyMap{
		Up:    components.NewUpKey(),
		Down:  components.NewDownKey(),
		Enter: components.NewEnterKey(), // Note: Also used for expand/collapse AC testing instructions
		Quit:  components.NewQuitKey(),
		Back:  components.NewBackKey(),
		Help:  components.NewHelpKey(),
		Verify: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "verify AC"),
		),
		Skip: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "skip AC"),
		),
		Fail: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "fail AC"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b"),
			key.WithHelp("pgup/b", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdn"),
			key.WithHelp("pgdn", "page down"),
		),
	}
}

// ShortHelp returns keybindings for short help view
func (k TaskDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Verify, k.Skip, k.Fail, k.Back, k.Quit}
}

// FullHelp returns all keybindings for full help view
func (k TaskDetailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.PageUp, k.PageDown},
		{k.Verify, k.Skip, k.Fail},
		{k.Back, k.Help, k.Quit},
	}
}

// TaskDetailPresenter presents the task detail view with expandable ACs
type TaskDetailPresenter struct {
	viewModel       *viewmodels.TaskDetailViewModel
	help            components.Help
	keys            TaskDetailKeyMap
	showFullHelp    bool
	selectedIndex   int
	width           int
	height          int
	repo            domain.RoadmapRepository
	ctx             context.Context
	acListComponent *ACListComponent

	// Scrolling support
	scrollHelperACs  *components.ScrollHelperMultiline // For ACs (multi-line with expansion)
	terminalHeight   int
}

// NewTaskDetailPresenter creates a new task detail presenter
func NewTaskDetailPresenter(vm *viewmodels.TaskDetailViewModel, repo domain.RoadmapRepository, ctx context.Context) *TaskDetailPresenter {
	return NewTaskDetailPresenterWithSelection(vm, repo, ctx, 0)
}

// NewTaskDetailPresenterWithSelection creates a new task detail presenter with a specific selected index
func NewTaskDetailPresenterWithSelection(vm *viewmodels.TaskDetailViewModel, repo domain.RoadmapRepository, ctx context.Context, selectedIndex int) *TaskDetailPresenter {
	return &TaskDetailPresenter{
		viewModel:       vm,
		help:            components.NewHelp(),
		keys:            NewTaskDetailKeyMap(),
		showFullHelp:    false,
		selectedIndex:   selectedIndex,
		repo:            repo,
		ctx:             ctx,
		acListComponent: NewACListComponent(repo, ctx, true), // enableExpand=true for task detail
		width:           80,                                  // Default width until WindowSizeMsg arrives
		height:          24,

		// Scrolling support
		scrollHelperACs:  components.NewScrollHelperMultiline(),
		terminalHeight:   24,
	}
}

func (p *TaskDetailPresenter) Init() tea.Cmd {
	// Request terminal size immediately to get actual dimensions
	return tea.WindowSize()
}

func (p *TaskDetailPresenter) Update(msg tea.Msg) (Presenter, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		p.terminalHeight = msg.Height
		p.help.SetWidth(msg.Width)

		// Calculate available viewport height
		headerHeight := 12  // Task header, description, track info, iteration membership
		footerHeight := 2   // Help text
		availableHeight := msg.Height - headerHeight - footerHeight
		if availableHeight < 1 {
			availableHeight = 1
		}

		p.scrollHelperACs.SetViewportHeight(availableHeight)
		return p, nil

	case tea.KeyMsg:
		// Component handles feedback input if active
		if handled, cmd := p.acListComponent.UpdateFeedback(msg); handled {
			// Check if Enter was pressed (submit)
			if msg.Type == tea.KeyEnter {
				acID, feedback := p.acListComponent.SubmitFeedback()
				return p, p.acListComponent.FailAC(acID, feedback, IterationDetailTabTasks, p.selectedIndex)
			}
			// Check if Quit was pressed while in feedback mode
			if key.Matches(msg, p.keys.Quit) {
				return p, tea.Quit
			}
			return p, cmd
		}

		// Normal key handling when feedback input is not active
		switch {
		case key.Matches(msg, p.keys.Quit):
			return p, tea.Quit
		case key.Matches(msg, p.keys.Back):
			return p, func() tea.Msg { return BackMsgNew{} }
		case key.Matches(msg, p.keys.Help):
			p.showFullHelp = !p.showFullHelp
		case key.Matches(msg, p.keys.Up):
			if p.selectedIndex > 0 {
				p.selectedIndex--
				lineCounts := p.calculateACLineCounts()
				p.scrollHelperACs.EnsureVisibleMultiline(lineCounts, p.selectedIndex)
			}
		case key.Matches(msg, p.keys.Down):
			maxIndex := len(p.viewModel.AcceptanceCriteria) - 1
			if p.selectedIndex < maxIndex {
				p.selectedIndex++
				lineCounts := p.calculateACLineCounts()
				p.scrollHelperACs.EnsureVisibleMultiline(lineCounts, p.selectedIndex)
			}
		case key.Matches(msg, p.keys.PageUp):
			if p.selectedIndex > 0 {
				// Jump up by viewport height
				viewportHeight := p.scrollHelperACs.ViewportOffset()
				if viewportHeight == 0 {
					viewportHeight = 10 // Default if not yet set
				}
				newIndex := p.selectedIndex - viewportHeight
				if newIndex < 0 {
					newIndex = 0
				}
				p.selectedIndex = newIndex
				lineCounts := p.calculateACLineCounts()
				p.scrollHelperACs.EnsureVisibleMultiline(lineCounts, p.selectedIndex)
			}
		case key.Matches(msg, p.keys.PageDown):
			totalACs := len(p.viewModel.AcceptanceCriteria)
			if p.selectedIndex < totalACs-1 {
				viewportHeight := p.scrollHelperACs.ViewportOffset()
				if viewportHeight == 0 {
					viewportHeight = 10 // Default if not yet set
				}
				newIndex := p.selectedIndex + viewportHeight
				if newIndex >= totalACs {
					newIndex = totalACs - 1
				}
				p.selectedIndex = newIndex
				lineCounts := p.calculateACLineCounts()
				p.scrollHelperACs.EnsureVisibleMultiline(lineCounts, p.selectedIndex)
			}
		case key.Matches(msg, p.keys.Enter):
			// Expand/collapse AC testing instructions
			if p.selectedIndex >= 0 && p.selectedIndex < len(p.viewModel.AcceptanceCriteria) {
				ac := p.viewModel.AcceptanceCriteria[p.selectedIndex]
				ac.IsExpanded = !ac.IsExpanded

				// Recalculate line counts with new expansion state
				lineCounts := p.calculateACLineCounts()

				// Ensure expanded AC is visible (scroll to show expanded content)
				p.scrollHelperACs.EnsureVisibleMultiline(lineCounts, p.selectedIndex)
			}
		case key.Matches(msg, p.keys.Verify):
			if p.selectedIndex >= 0 && p.selectedIndex < len(p.viewModel.AcceptanceCriteria) {
				acID := p.viewModel.AcceptanceCriteria[p.selectedIndex].ID
				return p, p.acListComponent.VerifyAC(acID, IterationDetailTabTasks, p.selectedIndex)
			}
		case key.Matches(msg, p.keys.Skip):
			if p.selectedIndex >= 0 && p.selectedIndex < len(p.viewModel.AcceptanceCriteria) {
				acID := p.viewModel.AcceptanceCriteria[p.selectedIndex].ID
				return p, p.acListComponent.SkipAC(acID, IterationDetailTabTasks, p.selectedIndex)
			}
		case key.Matches(msg, p.keys.Fail):
			if p.selectedIndex >= 0 && p.selectedIndex < len(p.viewModel.AcceptanceCriteria) {
				acID := p.viewModel.AcceptanceCriteria[p.selectedIndex].ID
				return p, p.acListComponent.StartFeedback(acID)
			}
		}
	}

	return p, nil
}

func (p *TaskDetailPresenter) View() string {
	var b strings.Builder

	// Calculate available width (leave some margin)
	availableWidth := p.width - 4
	if availableWidth < 40 {
		availableWidth = 40 // Minimum width
	}

	// Title
	b.WriteString(components.Styles.TitleStyle.Render(fmt.Sprintf("Task: %s", p.viewModel.ID)))
	b.WriteString("\n\n")

	// Metadata with width wrapping
	titleText := lipgloss.NewStyle().Width(availableWidth).Render(p.viewModel.Title)
	b.WriteString(components.Styles.TitleStyle.Render(titleText))
	b.WriteString("\n")

	statusText := lipgloss.NewStyle().Width(availableWidth).Render(fmt.Sprintf("Status: %s", p.viewModel.Status))
	b.WriteString(components.Styles.MetadataStyle.Render(statusText))
	b.WriteString("\n")

	if p.viewModel.Branch != "" {
		branchText := lipgloss.NewStyle().Width(availableWidth).Render(fmt.Sprintf("Branch: %s", p.viewModel.Branch))
		b.WriteString(components.Styles.MetadataStyle.Render(branchText))
		b.WriteString("\n")
	}

	createdText := lipgloss.NewStyle().Width(availableWidth).Render(fmt.Sprintf("Created: %s", p.viewModel.CreatedAt))
	b.WriteString(components.Styles.MetadataStyle.Render(createdText))
	b.WriteString("\n")

	updatedText := lipgloss.NewStyle().Width(availableWidth).Render(fmt.Sprintf("Updated: %s", p.viewModel.UpdatedAt))
	b.WriteString(components.Styles.MetadataStyle.Render(updatedText))
	b.WriteString("\n\n")

	// Track info with width wrapping
	if p.viewModel.TrackInfo != nil {
		b.WriteString(components.Styles.SectionStyle.Render("Track"))
		b.WriteString("\n")
		trackText := lipgloss.NewStyle().Width(availableWidth).Render(
			fmt.Sprintf("  %s: %s (%s)", p.viewModel.TrackInfo.ID, p.viewModel.TrackInfo.Title, p.viewModel.TrackInfo.Status))
		b.WriteString(trackText)
		b.WriteString("\n\n")
	}

	// Iteration membership with width wrapping
	if len(p.viewModel.Iterations) > 0 {
		b.WriteString(components.Styles.SectionStyle.Render("Iterations"))
		b.WriteString("\n")
		for _, iter := range p.viewModel.Iterations {
			iterText := lipgloss.NewStyle().Width(availableWidth).Render(
				fmt.Sprintf("  #%d %s (%s)", iter.Number, iter.Name, iter.Status))
			b.WriteString(iterText)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Description with width wrapping
	if p.viewModel.Description != "" {
		b.WriteString(components.Styles.SectionStyle.Render("Description"))
		b.WriteString("\n")
		descText := lipgloss.NewStyle().Width(availableWidth).Render(p.viewModel.Description)
		b.WriteString(descText)
		b.WriteString("\n\n")
	}

	// Acceptance Criteria
	b.WriteString(components.Styles.SectionStyle.Render("Acceptance Criteria"))
	b.WriteString("\n")

	if len(p.viewModel.AcceptanceCriteria) == 0 {
		b.WriteString(components.Styles.MetadataStyle.Render("  No acceptance criteria"))
		b.WriteString("\n")
	} else {
		p.renderACs(&b, availableWidth)
	}

	// Feedback input component renders inline at bottom if active
	feedbackView := p.acListComponent.ViewFeedback(p.width)
	if feedbackView != "" {
		b.WriteString(feedbackView)
		return b.String()
	}

	// Help view
	b.WriteString("\n")
	if p.showFullHelp {
		b.WriteString(p.help.FullHelpView(p.keys.FullHelp()))
	} else {
		b.WriteString(p.help.ShortHelpView(p.keys.ShortHelp()))
	}

	return b.String()
}

// calculateACLineCounts returns line counts for each AC (collapsed = 1, expanded = N)
func (p *TaskDetailPresenter) calculateACLineCounts() []int {
	lineCounts := make([]int, len(p.viewModel.AcceptanceCriteria))
	for i, ac := range p.viewModel.AcceptanceCriteria {
		if ac.IsExpanded && ac.TestingInstructions != "" {
			// Count lines in testing instructions + header
			lines := strings.Count(ac.TestingInstructions, "\n") + 2 // +2 for header and content
			lineCounts[i] = lines
		} else {
			// Collapsed AC is 1 line
			lineCounts[i] = 1
		}
	}
	return lineCounts
}

func (p *TaskDetailPresenter) renderACs(b *strings.Builder, availableWidth int) {
	acs := p.viewModel.AcceptanceCriteria
	if len(acs) == 0 {
		return
	}

	// Get visible range from multiline scroll helper
	lineCounts := p.calculateACLineCounts()
	firstItem, lastItem, lineOffset := p.scrollHelperACs.VisibleRangeMultiline(lineCounts)

	// Scroll indicator (above)
	if firstItem > 0 {
		b.WriteString("  ↑ More ACs above\n")
	}

	// Render visible ACs
	for i := firstItem; i <= lastItem && i < len(acs); i++ {
		ac := acs[i]

		// Determine if this AC should show partial content (due to lineOffset)
		skipLines := 0
		if i == firstItem {
			skipLines = lineOffset
		}

		// Highlight selected
		prefix := "  "
		if i == p.selectedIndex {
			prefix = "> "
		}

		// Render AC header (unless skipped by lineOffset)
		if skipLines == 0 {
			statusIcon := "○" // default
			switch ac.Status {
			case "verified":
				statusIcon = "✓"
			case "failed":
				statusIcon = "✗"
			case "pending-review":
				statusIcon = "⧗"
			case "skipped":
				statusIcon = "⊘"
			}

			// Show clipboard icon if AC has testing instructions
			description := ac.Description
			if ac.TestingInstructions != "" {
				description = "📋 " + description
			}

			b.WriteString(fmt.Sprintf("%s%s %s\n", prefix, statusIcon, description))
		}

		// If expanded, render testing instructions (respecting lineOffset)
		if ac.IsExpanded && ac.TestingInstructions != "" {
			instructionLines := strings.Split(ac.TestingInstructions, "\n")
			for j, line := range instructionLines {
				// Skip lines before lineOffset (only for first visible item)
				if i == firstItem && j+1 < skipLines {
					continue
				}
				b.WriteString(fmt.Sprintf("    %s\n", line))
			}
		}
	}

	// Scroll indicator (below)
	if lastItem < len(acs)-1 {
		b.WriteString("  ↓ More ACs below\n")
	}
}

