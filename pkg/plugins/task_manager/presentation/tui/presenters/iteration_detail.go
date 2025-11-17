package presenters

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/entities"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/components"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
	"github.com/muesli/reflow/wordwrap"
)

// IterationDetailTab represents the active tab in iteration detail view
type IterationDetailTab int

const (
	IterationDetailTabTasks IterationDetailTab = iota
	IterationDetailTabACs
)

// IterationDetailKeyMap defines keybindings for iteration detail view
type IterationDetailKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Quit     key.Binding
	Back     key.Binding
	Help     key.Binding
	Tab      key.Binding
	Verify   key.Binding // Space - verify AC
	Skip     key.Binding // s - skip AC
	Fail     key.Binding // f - fail AC
	PageUp   key.Binding // pgup/b - page up
	PageDown key.Binding // pgdn - page down
	// Task state transitions
	InProgress key.Binding // i - todo → in-progress
	Review     key.Binding // r - in-progress → review
	Done       key.Binding // d - review → done (with AC verification)
	Reopen     key.Binding // o - done → todo
}

// NewIterationDetailKeyMap creates default keybindings for iteration detail
func NewIterationDetailKeyMap() IterationDetailKeyMap {
	return IterationDetailKeyMap{
		Up:    components.NewUpKey(),
		Down:  components.NewDownKey(),
		Enter: components.NewEnterKey(),
		Quit:  components.NewQuitKey(),
		Back:  components.NewBackKey(),
		Help:  components.NewHelpKey(),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch view"),
		),
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
		InProgress: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "in progress"),
		),
		Review: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "review"),
		),
		Done: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "done"),
		),
		Reopen: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "reopen"),
		),
	}
}

// ShortHelp returns keybindings based on active tab
func (k IterationDetailKeyMap) ShortHelp(activeTab IterationDetailTab) []key.Binding {
	if activeTab == IterationDetailTabTasks {
		return []key.Binding{k.Up, k.Down, k.Enter, k.InProgress, k.Review, k.Done, k.Tab, k.Back, k.Quit}
	}
	// ACs view
	return []key.Binding{k.Up, k.Down, k.Enter, k.Verify, k.Skip, k.Fail, k.Tab, k.Back, k.Quit}
}

// FullHelp returns all keybindings based on active tab
func (k IterationDetailKeyMap) FullHelp(activeTab IterationDetailTab) [][]key.Binding {
	if activeTab == IterationDetailTabTasks {
		return [][]key.Binding{
			{k.Up, k.Down, k.Enter},
			{k.PageUp, k.PageDown},
			{k.InProgress, k.Review, k.Done, k.Reopen},
			{k.Tab, k.Back, k.Help, k.Quit},
		}
	}
	// ACs view
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.PageUp, k.PageDown},
		{k.Verify, k.Skip, k.Fail},
		{k.Tab, k.Back, k.Help, k.Quit},
	}
}

// IterationDetailPresenter presents the iteration detail view
type IterationDetailPresenter struct {
	viewModel       *viewmodels.IterationDetailViewModel
	help            components.Help
	keys            IterationDetailKeyMap
	showFullHelp    bool
	activeTab       IterationDetailTab
	selectedIndex   int
	width           int
	height          int
	repo            domain.RoadmapRepository
	ctx             context.Context
	acListComponent *ACListComponent

	// Scrolling support
	scrollHelperTasks *components.ScrollHelper          // For tasks tab (single-line)
	scrollHelperACs   *components.ScrollHelperMultiline // For ACs tab (multi-line with expansion)
	terminalHeight    int
}

func NewIterationDetailPresenter(vm *viewmodels.IterationDetailViewModel, repo domain.RoadmapRepository, ctx context.Context) *IterationDetailPresenter {
	return NewIterationDetailPresenterWithTab(vm, repo, ctx, IterationDetailTabTasks)
}

// NewIterationDetailPresenterWithTab creates a new iteration detail presenter with a specific active tab
func NewIterationDetailPresenterWithTab(vm *viewmodels.IterationDetailViewModel, repo domain.RoadmapRepository, ctx context.Context, activeTab IterationDetailTab) *IterationDetailPresenter {
	return NewIterationDetailPresenterWithSelection(vm, repo, ctx, activeTab, 0)
}

// NewIterationDetailPresenterWithSelection creates a new iteration detail presenter with a specific active tab and selected index
func NewIterationDetailPresenterWithSelection(vm *viewmodels.IterationDetailViewModel, repo domain.RoadmapRepository, ctx context.Context, activeTab IterationDetailTab, selectedIndex int) *IterationDetailPresenter {
	return &IterationDetailPresenter{
		viewModel:       vm,
		help:            components.NewHelp(),
		keys:            NewIterationDetailKeyMap(),
		showFullHelp:    false,
		activeTab:       activeTab,
		selectedIndex:   selectedIndex,
		repo:            repo,
		ctx:             ctx,
		acListComponent: NewACListComponent(repo, ctx, true), // enableExpand=true (same behavior as task detail)
		width:           80,                                    // Default width until WindowSizeMsg arrives
		height:          24,

		// Initialize scroll helpers
		scrollHelperTasks: components.NewScrollHelper(),
		scrollHelperACs:   components.NewScrollHelperMultiline(),
		terminalHeight:    24,
	}
}



func (p *IterationDetailPresenter) Init() tea.Cmd {
	// Request terminal size immediately to get actual dimensions
	return tea.WindowSize()
}

func (p *IterationDetailPresenter) Update(msg tea.Msg) (Presenter, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		p.terminalHeight = msg.Height
		p.help.SetWidth(msg.Width)

		// Calculate available viewport height for scrolling
		// Account for: title (1) + metadata (4-5) + progress (1) + tab headers (2) + help (2)
		headerHeight := 11
		footerHeight := 2 // Help text
		availableHeight := msg.Height - headerHeight - footerHeight
		if availableHeight < 5 {
			availableHeight = 5 // Minimum height
		}

		p.scrollHelperTasks.SetViewportHeight(availableHeight)
		p.scrollHelperACs.SetViewportHeight(availableHeight)

		// Ensure current selection is visible with new viewport height
		if p.activeTab == IterationDetailTabTasks {
			totalTasks := len(p.viewModel.TODOTasks) + len(p.viewModel.InProgressTasks) + len(p.viewModel.ReviewTasks) + len(p.viewModel.DoneTasks)
			p.scrollHelperTasks.EnsureVisible(totalTasks, p.selectedIndex)
		} else {
			lineCounts := p.calculateACLineCounts()
			p.scrollHelperACs.EnsureVisibleMultiline(lineCounts, p.selectedIndex)
		}

	case tea.KeyMsg:
		// Component handles feedback input if active
		if handled, cmd := p.acListComponent.UpdateFeedback(msg); handled {
			// Check if Enter was pressed (submit)
			if msg.Type == tea.KeyEnter {
				acID, feedback := p.acListComponent.SubmitFeedback()
				return p, p.acListComponent.FailAC(acID, feedback, p.activeTab, p.selectedIndex)
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
		case key.Matches(msg, p.keys.Tab):
			// Switch between tasks and ACs
			if p.activeTab == IterationDetailTabTasks {
				p.activeTab = IterationDetailTabACs
			} else {
				p.activeTab = IterationDetailTabTasks
			}
			p.selectedIndex = 0
		case key.Matches(msg, p.keys.Up):
			if p.activeTab == IterationDetailTabTasks {
				totalTasks := len(p.viewModel.TODOTasks) + len(p.viewModel.InProgressTasks) + len(p.viewModel.ReviewTasks) + len(p.viewModel.DoneTasks)
				if p.selectedIndex > 0 {
					p.selectedIndex--
					p.scrollHelperTasks.EnsureVisible(totalTasks, p.selectedIndex)
				}
			} else {
				if p.selectedIndex > 0 {
					p.selectedIndex--
					lineCounts := p.calculateACLineCounts()
					p.scrollHelperACs.EnsureVisibleMultiline(lineCounts, p.selectedIndex)
				}
			}
		case key.Matches(msg, p.keys.Down):
			maxIndex := p.getMaxIndex()
			if p.selectedIndex < maxIndex {
				p.selectedIndex++
				if p.activeTab == IterationDetailTabTasks {
					totalTasks := len(p.viewModel.TODOTasks) + len(p.viewModel.InProgressTasks) + len(p.viewModel.ReviewTasks) + len(p.viewModel.DoneTasks)
					p.scrollHelperTasks.EnsureVisible(totalTasks, p.selectedIndex)
				} else {
					lineCounts := p.calculateACLineCounts()
					p.scrollHelperACs.EnsureVisibleMultiline(lineCounts, p.selectedIndex)
				}
			}
		case key.Matches(msg, p.keys.PageUp):
			if p.activeTab == IterationDetailTabTasks {
				totalTasks := len(p.viewModel.TODOTasks) + len(p.viewModel.InProgressTasks) + len(p.viewModel.ReviewTasks) + len(p.viewModel.DoneTasks)
				newIndex := p.scrollHelperTasks.PageUp(totalTasks)
				p.selectedIndex = newIndex
			}
		case key.Matches(msg, p.keys.PageDown):
			if p.activeTab == IterationDetailTabTasks {
				totalTasks := len(p.viewModel.TODOTasks) + len(p.viewModel.InProgressTasks) + len(p.viewModel.ReviewTasks) + len(p.viewModel.DoneTasks)
				newIndex := p.scrollHelperTasks.PageDown(totalTasks, p.selectedIndex)
				p.selectedIndex = newIndex
			}
		case key.Matches(msg, p.keys.Enter):
			if p.activeTab == IterationDetailTabTasks {
				// Navigate to task detail
				taskID := p.getSelectedTaskID()
				if taskID != "" {
					return p, func() tea.Msg {
						return TaskSelectedMsg{TaskID: taskID}
					}
				}
			} else if p.activeTab == IterationDetailTabACs {
				// Expand/collapse AC testing instructions (same as TaskDetail)
				acID := p.getSelectedACID()
				if acID != "" {
					// Find and toggle the AC in viewModel
					for _, group := range p.viewModel.TaskACs {
						for _, ac := range group.ACs {
							if ac.ID == acID {
								ac.IsExpanded = !ac.IsExpanded
								// Recalculate line counts with new expansion state
								lineCounts := p.calculateACLineCounts()
								p.scrollHelperACs.EnsureVisibleMultiline(lineCounts, p.selectedIndex)
								return p, nil
							}
						}
					}
				}
			}
		case key.Matches(msg, p.keys.Verify):
			if p.activeTab == IterationDetailTabACs {
				acID := p.getSelectedACID()
				if acID != "" {
					return p, p.acListComponent.VerifyAC(acID, p.activeTab, p.selectedIndex)
				}
			}
		case key.Matches(msg, p.keys.Skip):
			if p.activeTab == IterationDetailTabACs {
				acID := p.getSelectedACID()
				if acID != "" {
					return p, p.acListComponent.SkipAC(acID, p.activeTab, p.selectedIndex)
				}
			}
		case key.Matches(msg, p.keys.Fail):
			if p.activeTab == IterationDetailTabACs {
				acID := p.getSelectedACID()
				if acID != "" {
					return p, p.acListComponent.StartFeedback(acID)
				}
			}
		case key.Matches(msg, p.keys.InProgress):
			if p.activeTab == IterationDetailTabTasks {
				task := p.getSelectedTask()
				if task != nil && task.Status == "todo" {
					return p, p.transitionTaskStatus(task.ID, "in-progress", p.activeTab, p.selectedIndex)
				}
			}
		case key.Matches(msg, p.keys.Review):
			if p.activeTab == IterationDetailTabTasks {
				task := p.getSelectedTask()
				if task != nil && task.Status == "in-progress" {
					return p, p.transitionTaskStatus(task.ID, "review", p.activeTab, p.selectedIndex)
				}
			}
		case key.Matches(msg, p.keys.Done):
			if p.activeTab == IterationDetailTabTasks {
				task := p.getSelectedTask()
				if task != nil && task.Status == "review" {
					return p, p.transitionTaskToDone(task.ID, p.activeTab, p.selectedIndex)
				}
			}
		case key.Matches(msg, p.keys.Reopen):
			if p.activeTab == IterationDetailTabTasks {
				task := p.getSelectedTask()
				if task != nil && task.Status == "done" {
					return p, p.transitionTaskStatus(task.ID, "todo", p.activeTab, p.selectedIndex)
				}
			}
		}
	}

	return p, nil
}

func (p *IterationDetailPresenter) View() string {
	var b strings.Builder

	// Title
	b.WriteString(components.Styles.TitleStyle.Render(fmt.Sprintf("Iteration #%d: %s", p.viewModel.Number, p.viewModel.Name)))
	b.WriteString("\n\n")

	// Metadata
	// Goal with text wrapping
	if p.viewModel.Goal != "" {
		availableWidth := p.width - 6 // Account for "Goal: " prefix
		if availableWidth < 20 {
			availableWidth = 20 // Minimum width
		}
		wrappedGoal := wordwrap.String(p.viewModel.Goal, availableWidth)
		b.WriteString(components.Styles.MetadataStyle.Render(fmt.Sprintf("Goal: %s", wrappedGoal)))
		b.WriteString("\n")
	}

	// Deliverable with text wrapping
	if p.viewModel.Deliverable != "" {
		availableWidth := p.width - 13 // Account for "Deliverable: " prefix
		if availableWidth < 20 {
			availableWidth = 20 // Minimum width
		}
		wrappedDeliverable := wordwrap.String(p.viewModel.Deliverable, availableWidth)
		b.WriteString(components.Styles.MetadataStyle.Render(fmt.Sprintf("Deliverable: %s", wrappedDeliverable)))
		b.WriteString("\n")
	}

	// Status with color
	statusText := getStatusStyle(p.viewModel.StatusColor).Render(p.viewModel.Status)
	b.WriteString(components.Styles.MetadataStyle.Render(fmt.Sprintf("Status: %s", statusText)))
	b.WriteString("\n")
	if p.viewModel.StartedAt != "" {
		b.WriteString(components.Styles.MetadataStyle.Render(fmt.Sprintf("Started: %s", p.viewModel.StartedAt)))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Progress bar
	progressText := fmt.Sprintf("Progress: %d/%d tasks (%.0f%%)",
		p.viewModel.Progress.Completed,
		p.viewModel.Progress.Total,
		p.viewModel.Progress.Percent*100)
	b.WriteString(components.Styles.ProgressStyle.Render(progressText))
	b.WriteString("\n\n")

	// Tab headers
	if p.activeTab == IterationDetailTabTasks {
		b.WriteString(components.Styles.ActiveTabStyle.Render("Tasks"))
		b.WriteString("  ")
		b.WriteString(components.Styles.TabStyle.Render("Acceptance Criteria"))
	} else {
		b.WriteString(components.Styles.TabStyle.Render("Tasks"))
		b.WriteString("  ")
		b.WriteString(components.Styles.ActiveTabStyle.Render("Acceptance Criteria"))
	}
	b.WriteString("\n\n")

	// Content based on active tab
	if p.activeTab == IterationDetailTabTasks {
		p.renderTasksView(&b)
	} else {
		p.renderACsView(&b)
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
		b.WriteString(p.help.FullHelpView(p.keys.FullHelp(p.activeTab)))
	} else {
		b.WriteString(p.help.ShortHelpView(p.keys.ShortHelp(p.activeTab)))
	}

	return b.String()
}

// calculateACLineCounts returns the line count for each AC based on expansion state
// Collapsed AC = 1 line, Expanded AC = header + testing instruction lines
func (p *IterationDetailPresenter) calculateACLineCounts() []int {
	lineCounts := make([]int, 0)
	for _, group := range p.viewModel.TaskACs {
		for _, ac := range group.ACs {
			if ac.IsExpanded && ac.TestingInstructions != "" {
				// Count lines in testing instructions + header + spacing
				lines := strings.Count(ac.TestingInstructions, "\n") + 3 // +3 for header, content, spacing
				lineCounts = append(lineCounts, lines)
			} else {
				// Collapsed AC is 1 line
				lineCounts = append(lineCounts, 1)
			}
		}
	}
	return lineCounts
}

func (p *IterationDetailPresenter) renderTasksView(b *strings.Builder) {
	// Build flat task list with section info
	type taskItem struct {
		task        *viewmodels.TaskRowViewModel
		section     string
		sectionName string
	}
	allTasks := make([]taskItem, 0)

	for _, task := range p.viewModel.TODOTasks {
		allTasks = append(allTasks, taskItem{task: task, section: "todo", sectionName: "TODO"})
	}
	for _, task := range p.viewModel.InProgressTasks {
		allTasks = append(allTasks, taskItem{task: task, section: "in-progress", sectionName: "IN PROGRESS"})
	}
	for _, task := range p.viewModel.ReviewTasks {
		allTasks = append(allTasks, taskItem{task: task, section: "review", sectionName: "REVIEW"})
	}
	for _, task := range p.viewModel.DoneTasks {
		allTasks = append(allTasks, taskItem{task: task, section: "done", sectionName: "DONE"})
	}

	if len(allTasks) == 0 {
		b.WriteString(components.Styles.MetadataStyle.Render("No tasks in this iteration"))
		return
	}

	// Get visible range from scroll helper
	start, end := p.scrollHelperTasks.VisibleRange(len(allTasks))

	// Scroll indicator (above)
	if start > 0 {
		b.WriteString(components.Styles.MetadataStyle.Render("  ↑ More tasks above"))
		b.WriteString("\n")
	}

	// Render visible tasks with section headers
	currentSection := ""
	for i := start; i < end; i++ {
		item := allTasks[i]

		// Render section header if new section
		if item.section != currentSection {
			currentSection = item.section
			b.WriteString(components.Styles.SectionStyle.Render(item.sectionName))
			b.WriteString("\n")
		}

		// Render task with colored status
		statusText := getStatusStyle(item.task.StatusColor).Render(item.task.Status)
		var output string
		if i == p.selectedIndex {
			output = components.Styles.SelectedStyle.Render(fmt.Sprintf("  %s: %s - %s", item.task.ID, item.task.Title, statusText))
		} else {
			output = fmt.Sprintf("  %s: %s - %s", item.task.ID, item.task.Title, statusText)
		}
		b.WriteString(output)
		b.WriteString("\n")
	}

	// Scroll indicator (below)
	if end < len(allTasks) {
		b.WriteString(components.Styles.MetadataStyle.Render("  ↓ More tasks below"))
		b.WriteString("\n")
	}
}

func (p *IterationDetailPresenter) renderACsView(b *strings.Builder) {
	if len(p.viewModel.TaskACs) == 0 {
		b.WriteString(components.Styles.MetadataStyle.Render("No acceptance criteria"))
		return
	}

	// Build flat AC list with task context
	type acItem struct {
		ac             *viewmodels.IterationACViewModel
		taskID         string
		taskTitle      string
		isFirstInGroup bool
	}
	allACs := make([]acItem, 0)

	for _, group := range p.viewModel.TaskACs {
		for i, ac := range group.ACs {
			allACs = append(allACs, acItem{
				ac:             ac,
				taskID:         group.Task.ID,
				taskTitle:      group.Task.Title,
				isFirstInGroup: i == 0,
			})
		}
	}

	if len(allACs) == 0 {
		b.WriteString(components.Styles.MetadataStyle.Render("No acceptance criteria"))
		return
	}

	// Get visible range from multiline scroll helper
	lineCounts := p.calculateACLineCounts()
	firstItem, lastItem, lineOffset := p.scrollHelperACs.VisibleRangeMultiline(lineCounts)

	// Scroll indicator (above)
	if firstItem > 0 {
		b.WriteString(components.Styles.MetadataStyle.Render("  ↑ More ACs above"))
		b.WriteString("\n")
	}

	// Render visible ACs: group by task and use ACListComponent for AC rendering
	currentTaskID := ""
	acIndex := 0 // Track index within the visible AC section (not global)

	for i := firstItem; i <= lastItem && i < len(allACs); i++ {
		item := allACs[i]

		// Render task header if new task group
		if item.taskID != currentTaskID {
			currentTaskID = item.taskID
			b.WriteString(components.Styles.SectionStyle.Render(fmt.Sprintf("Task: %s - %s", item.taskID, item.taskTitle)))
			b.WriteString("\n")
		}

		// Determine if this AC should show partial content (due to lineOffset)
		skipLines := 0
		if i == firstItem {
			skipLines = lineOffset
		}

		// Render AC header (unless skipped by lineOffset)
		if skipLines == 0 {
			ac := &IterationACViewModelWrapper{IterationACViewModel: item.ac}
			// Create a single-item list and render it
			singleACList := []ACViewModel{ac}
			// Determine if this AC is selected (selectedIndex tracks global index, not visible index)
			isSelected := i == p.selectedIndex
			selectedIndexInList := -1
			if isSelected {
				selectedIndexInList = 0
			}
			p.acListComponent.RenderACList(b, singleACList, selectedIndexInList, p.width)
		} else {
			// If lineOffset skips this AC's header, we still need to show partial content
			if item.ac.IsExpanded && item.ac.TestingInstructions != "" {
				instructionLines := strings.Split(item.ac.TestingInstructions, "\n")
				for j := skipLines; j < len(instructionLines); j++ {
					b.WriteString(fmt.Sprintf("    %s\n", instructionLines[j]))
				}
			}
		}

		acIndex++
	}

	// Scroll indicator (below)
	if lastItem < len(allACs)-1 {
		b.WriteString(components.Styles.MetadataStyle.Render("  ↓ More ACs below"))
		b.WriteString("\n")
	}
}

func (p *IterationDetailPresenter) getMaxIndex() int {
	if p.activeTab == IterationDetailTabTasks {
		return len(p.viewModel.TODOTasks) +
			len(p.viewModel.InProgressTasks) +
			len(p.viewModel.ReviewTasks) +
			len(p.viewModel.DoneTasks) - 1
	}
	// ACs view - count total ACs across all task groups
	totalACs := 0
	for _, group := range p.viewModel.TaskACs {
		totalACs += len(group.ACs)
	}
	return totalACs - 1
}

// getSelectedTaskID returns the task ID of the currently selected task
func (p *IterationDetailPresenter) getSelectedTaskID() string {
	if p.activeTab != IterationDetailTabTasks {
		return ""
	}

	index := p.selectedIndex
	todoLen := len(p.viewModel.TODOTasks)
	inProgressLen := len(p.viewModel.InProgressTasks)
	reviewLen := len(p.viewModel.ReviewTasks)

	if index < todoLen {
		return p.viewModel.TODOTasks[index].ID
	}
	index -= todoLen

	if index < inProgressLen {
		return p.viewModel.InProgressTasks[index].ID
	}
	index -= inProgressLen

	if index < reviewLen {
		return p.viewModel.ReviewTasks[index].ID
	}
	index -= reviewLen

	if index < len(p.viewModel.DoneTasks) {
		return p.viewModel.DoneTasks[index].ID
	}

	return ""
}

// getSelectedTask returns the task view model of the currently selected task
func (p *IterationDetailPresenter) getSelectedTask() *viewmodels.TaskRowViewModel {
	if p.activeTab != IterationDetailTabTasks {
		return nil
	}

	index := p.selectedIndex
	todoLen := len(p.viewModel.TODOTasks)
	inProgressLen := len(p.viewModel.InProgressTasks)
	reviewLen := len(p.viewModel.ReviewTasks)

	if index < todoLen {
		return p.viewModel.TODOTasks[index]
	}
	index -= todoLen

	if index < inProgressLen {
		return p.viewModel.InProgressTasks[index]
	}
	index -= inProgressLen

	if index < reviewLen {
		return p.viewModel.ReviewTasks[index]
	}
	index -= reviewLen

	if index < len(p.viewModel.DoneTasks) {
		return p.viewModel.DoneTasks[index]
	}

	return nil
}

// getSelectedACID returns the AC ID of the currently selected AC from grouped ACs
func (p *IterationDetailPresenter) getSelectedACID() string {
	if p.activeTab != IterationDetailTabACs {
		return ""
	}

	// Navigate through grouped ACs to find the AC at selectedIndex
	index := p.selectedIndex
	for _, group := range p.viewModel.TaskACs {
		if index < len(group.ACs) {
			return group.ACs[index].ID
		}
		index -= len(group.ACs)
	}

	return ""
}

// transitionTaskStatus transitions a task to a new status using repository
func (p *IterationDetailPresenter) transitionTaskStatus(taskID, newStatus string, activeTab IterationDetailTab, currentSelectedIndex int) tea.Cmd {
	return func() tea.Msg {
		// Fetch task
		task, err := p.repo.GetTask(p.ctx, taskID)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to get task: %w", err)}
		}

		// Update task status (newStatus is already a valid string)
		task.Status = newStatus
		task.UpdatedAt = time.Now()

		// Save
		err = p.repo.UpdateTask(p.ctx, task)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to update task: %w", err)}
		}

		return TaskTransitionCompletedMsg{ActiveTab: activeTab, SelectedIndex: currentSelectedIndex}
	}
}

// transitionTaskToDone transitions a task to done status with AC verification check
func (p *IterationDetailPresenter) transitionTaskToDone(taskID string, activeTab IterationDetailTab, currentSelectedIndex int) tea.Cmd {
	return func() tea.Msg {
		// Check for unverified ACs
		acs, err := p.repo.ListAC(p.ctx, taskID)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to check acceptance criteria: %w", err)}
		}

		// Filter for unverified ACs (status != verified and status != skipped)
		var unverifiedACs []string
		for _, ac := range acs {
			if ac.Status != entities.ACStatusVerified && ac.Status != entities.ACStatusSkipped {
				unverifiedACs = append(unverifiedACs, ac.ID)
			}
		}

		// Block transition if unverified ACs exist
		if len(unverifiedACs) > 0 {
			return ErrorMsg{
				Err: fmt.Errorf("cannot mark task as done: %d unverified acceptance criteria (%s). Please verify or skip ACs first",
					len(unverifiedACs),
					strings.Join(unverifiedACs, ", ")),
			}
		}

		// All ACs verified/skipped - proceed with transition
		// Fetch task
		task, err := p.repo.GetTask(p.ctx, taskID)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to get task: %w", err)}
		}

		// Update task status to done
		task.Status = "done"
		task.UpdatedAt = time.Now()

		// Save
		err = p.repo.UpdateTask(p.ctx, task)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to update task: %w", err)}
		}

		return TaskTransitionCompletedMsg{ActiveTab: activeTab, SelectedIndex: currentSelectedIndex}
	}
}
