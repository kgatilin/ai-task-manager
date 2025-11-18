package presenters

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/domain/repositories"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/components"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/queries"
	"github.com/kgatilin/darwinflow-pub/pkg/plugins/task_manager/presentation/tui/viewmodels"
	"github.com/muesli/reflow/wordwrap"
)

// DocumentViewerKeyMap defines keybindings for document viewer
type DocumentViewerKeyMap struct {
	Up          key.Binding
	Down        key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	JumpToStart key.Binding // shift+left - jump to beginning
	JumpToEnd   key.Binding // shift+right - jump to end
	Approve     key.Binding // 'a' - approve document
	Disapprove  key.Binding // 'd' - disapprove document
	Refresh     key.Binding // 'r' - refresh document
	Quit        key.Binding
	Back        key.Binding
	Help        key.Binding
}

// NewDocumentViewerKeyMap creates default keybindings for document viewer
func NewDocumentViewerKeyMap() DocumentViewerKeyMap {
	return DocumentViewerKeyMap{
		Up:   components.NewUpKey(),
		Down: components.NewDownKey(),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b", "shift+up"),
			key.WithHelp("pgup/shift+↑", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdn", "shift+down"),
			key.WithHelp("pgdn/shift+↓", "page down"),
		),
		JumpToStart: key.NewBinding(
			key.WithKeys("shift+left"),
			key.WithHelp("shift+←", "jump to start"),
		),
		JumpToEnd: key.NewBinding(
			key.WithKeys("shift+right"),
			key.WithHelp("shift+→", "jump to end"),
		),
		Approve: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "approve"),
		),
		Disapprove: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "disapprove"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Quit: components.NewQuitKey(),
		Back: components.NewBackKey(),
		Help: components.NewHelpKey(),
	}
}

// ShortHelp returns essential keybindings
func (k DocumentViewerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Approve, k.Disapprove, k.Refresh, k.Back, k.Quit}
}

// FullHelp returns all keybindings
func (k DocumentViewerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.PageUp, k.PageDown},
		{k.JumpToStart, k.JumpToEnd},
		{k.Approve, k.Disapprove, k.Refresh},
		{k.Back, k.Help, k.Quit},
	}
}

// DocumentViewerPresenter presents the document viewer view
type DocumentViewerPresenter struct {
	documentID     string
	viewModel      *viewmodels.DocumentViewModel
	help           components.Help
	keys           DocumentViewerKeyMap
	showFullHelp   bool
	width          int
	height         int
	repo           repositories.DocumentRepository
	ctx            context.Context
	scrollHelper   *components.ScrollHelper
	err            error
	isLoading      bool

	// Rendered markdown cache
	renderedContent string
}

// NewDocumentViewerPresenter creates a new document viewer presenter
func NewDocumentViewerPresenter(documentID string, repo repositories.DocumentRepository, ctx context.Context) *DocumentViewerPresenter {
	return &DocumentViewerPresenter{
		documentID:   documentID,
		help:         components.NewHelp(),
		keys:         NewDocumentViewerKeyMap(),
		showFullHelp: false,
		width:        80,  // Default width until WindowSizeMsg arrives
		height:       24,
		repo:         repo,
		ctx:          ctx,
		scrollHelper: components.NewScrollHelper(),
		isLoading:    true,
	}
}

// Init loads document data on initialization
func (p *DocumentViewerPresenter) Init() tea.Cmd {
	// Request window size immediately so we can properly size the viewport on load
	return tea.Batch(
		p.loadDocumentCmd(),
		tea.WindowSize(),
	)
}

// Update handles messages
func (p *DocumentViewerPresenter) Update(msg tea.Msg) (Presenter, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		p.help.SetWidth(msg.Width)

		// Calculate available viewport height for scrolling
		// Must match View() rendering logic
		if !p.isLoading && p.viewModel != nil {
			// Count actual header lines (varies based on track/iteration attachment)
			headerLines := 4 // title + type + status + blank line (minimum)
			if p.viewModel.TrackID != nil && *p.viewModel.TrackID != "" {
				headerLines++
			}
			if p.viewModel.IterationNumber != nil && *p.viewModel.IterationNumber > 0 {
				headerLines++
			}

			// Count help lines (short help is ~2 lines, full help is ~4-6 lines)
			helpLines := 2 // Default for short help
			if p.showFullHelp {
				helpLines = 6 // Estimate for full help
			}

			// Account for: scroll indicators (max 2) + separator before help (1)
			reservedLines := headerLines + 2 + 1 + helpLines

			availableHeight := msg.Height - reservedLines
			if availableHeight < 5 {
				availableHeight = 5 // Minimum viewport
			}

			p.scrollHelper.SetViewportHeight(availableHeight)

			// Re-render content with new width
			p.renderMarkdown()
		} else {
			// Default viewport size during loading
			p.scrollHelper.SetViewportHeight(msg.Height - 10)
		}

		return p, nil

	case tea.KeyMsg:
		if p.isLoading {
			return p, nil // Ignore key input while loading
		}

		switch {
		case key.Matches(msg, p.keys.Quit):
			return p, tea.Quit
		case key.Matches(msg, p.keys.Back):
			return p, func() tea.Msg { return BackMsgNew{} }
		case key.Matches(msg, p.keys.Help):
			p.showFullHelp = !p.showFullHelp
			// Recalculate viewport height when help changes size
			if p.viewModel != nil {
				headerLines := 4 // title + type + status + blank line (minimum)
				if p.viewModel.TrackID != nil && *p.viewModel.TrackID != "" {
					headerLines++
				}
				if p.viewModel.IterationNumber != nil && *p.viewModel.IterationNumber > 0 {
					headerLines++
				}

				helpLines := 2 // Short help
				if p.showFullHelp {
					helpLines = 6 // Full help estimate
				}

				reservedLines := headerLines + 2 + 1 + helpLines
				availableHeight := p.height - reservedLines
				if availableHeight < 5 {
					availableHeight = 5
				}
				p.scrollHelper.SetViewportHeight(availableHeight)
			}
		case key.Matches(msg, p.keys.Up):
			if p.scrollHelper.ViewportOffset() > 0 {
				p.scrollOffsetUp()
			}
		case key.Matches(msg, p.keys.Down):
			contentLines := strings.Split(p.renderedContent, "\n")
			p.scrollOffsetDown(len(contentLines))
		case key.Matches(msg, p.keys.PageUp):
			contentLines := strings.Split(p.renderedContent, "\n")
			p.scrollHelper.ScrollPageUp(len(contentLines))
		case key.Matches(msg, p.keys.PageDown):
			contentLines := strings.Split(p.renderedContent, "\n")
			p.scrollHelper.ScrollPageDown(len(contentLines))
		case key.Matches(msg, p.keys.JumpToStart):
			p.scrollHelper.ScrollToStart()
		case key.Matches(msg, p.keys.JumpToEnd):
			contentLines := strings.Split(p.renderedContent, "\n")
			p.scrollHelper.ScrollToEnd(len(contentLines))
		case key.Matches(msg, p.keys.Approve):
			if p.viewModel != nil {
				return p, p.approveDocumentCmd()
			}
		case key.Matches(msg, p.keys.Disapprove):
			if p.viewModel != nil {
				return p, p.disapproveDocumentCmd()
			}
		case key.Matches(msg, p.keys.Refresh):
			if p.viewModel != nil {
				return p, p.loadDocumentCmd()
			}
		}

	case DocumentLoadedMsg:
		if msg.Error != nil {
			p.err = msg.Error
		} else {
			p.viewModel = msg.ViewModel
			p.isLoading = false
			p.renderMarkdown()
		}
		return p, nil

	case DocumentActionCompletedMsg:
		// Reload document to get updated status
		return p, p.loadDocumentCmd()

	case ErrorMsg:
		p.err = msg.Err
		return p, nil
	}

	return p, nil
}

// View renders the document viewer
func (p *DocumentViewerPresenter) View() string {
	if p.isLoading {
		return fmt.Sprintf("%s Loading document...", components.Styles.LoadingStyle.Render("●"))
	}

	if p.err != nil {
		return components.Styles.ErrorMessageStyle.Render(fmt.Sprintf("Error: %v", p.err))
	}

	if p.viewModel == nil {
		return components.Styles.MetadataStyle.Render("Document not found")
	}

	// Render header first to count actual lines
	var header strings.Builder

	// Title with ID and type badge
	header.WriteString(components.Styles.TitleStyle.Render(p.viewModel.Title))
	header.WriteString(" ")
	header.WriteString(components.Styles.MetadataStyle.Render(fmt.Sprintf("[%s]", p.viewModel.ID)))
	header.WriteString("\n")

	// Type badge
	typeBadge := fmt.Sprintf("Type: %s", p.viewModel.TypeLabel)
	header.WriteString(components.Styles.MetadataStyle.Render(typeBadge))
	header.WriteString("\n")

	// Status badge
	statusColor := p.viewModel.StatusColor
	statusStyled := getStatusColor(statusColor)(p.viewModel.StatusLabel)
	statusBadge := fmt.Sprintf("Status: %s %s", p.viewModel.Icon, statusStyled)
	header.WriteString(components.Styles.MetadataStyle.Render(statusBadge))

	// Position indicator (show scroll position)
	totalLines := len(strings.Split(p.renderedContent, "\n"))
	position := p.scrollHelper.ScrollPosition(totalLines)
	header.WriteString("  ")
	header.WriteString(components.Styles.MetadataStyle.Render(fmt.Sprintf("[%s]", position)))
	header.WriteString("\n")

	// Attachment info if attached
	if p.viewModel.TrackID != nil && *p.viewModel.TrackID != "" {
		header.WriteString(components.Styles.MetadataStyle.Render(fmt.Sprintf("Track: %s", *p.viewModel.TrackID)))
		header.WriteString("\n")
	}
	if p.viewModel.IterationNumber != nil && *p.viewModel.IterationNumber > 0 {
		header.WriteString(components.Styles.MetadataStyle.Render(fmt.Sprintf("Iteration: #%d", *p.viewModel.IterationNumber)))
		header.WriteString("\n")
	}

	header.WriteString("\n")

	// Count actual header lines
	headerText := header.String()
	headerLines := strings.Count(headerText, "\n")

	// Render help to count help lines
	var helpText string
	if p.showFullHelp {
		helpText = p.help.FullHelpView(p.keys.FullHelp())
	} else {
		helpText = p.help.ShortHelpView(p.keys.ShortHelp())
	}
	helpLines := strings.Count(helpText, "\n") + 1 // +1 for the line itself

	// Calculate content viewport: contentLines + scroll indicators (up to 2) + separator before help (1)
	contentLines := strings.Split(p.renderedContent, "\n")
	start, end := p.scrollHelper.VisibleRange(len(contentLines))

	// Build scroll indicators
	var scrollIndicators strings.Builder
	if start > 0 {
		scrollIndicators.WriteString(components.Styles.MetadataStyle.Render("  ↑ More content above"))
		scrollIndicators.WriteString("\n")
	}

	scrollIndicatorTop := start > 0 // Track if we have top indicator

	var scrollIndicatorBottom bool
	if end < len(contentLines) {
		scrollIndicatorBottom = true
	}

	// Calculate total lines we'll render
	scrollIndicatorLines := 0
	if scrollIndicatorTop {
		scrollIndicatorLines++
	}
	if scrollIndicatorBottom {
		scrollIndicatorLines++
	}

	// Build content
	var content strings.Builder
	for i := start; i < end && i < len(contentLines); i++ {
		content.WriteString(contentLines[i])
		content.WriteString("\n")
	}

	if scrollIndicatorBottom {
		scrollIndicators.WriteString(components.Styles.MetadataStyle.Render("  ↓ More content below"))
		scrollIndicators.WriteString("\n")
	}

	// Assemble final output
	var b strings.Builder
	b.WriteString(headerText)
	b.WriteString(scrollIndicators.String())
	b.WriteString(content.String())
	b.WriteString("\n") // Separator before help
	b.WriteString(helpText)

	// Count total rendered lines
	totalRenderedLines := headerLines + scrollIndicatorLines + (end - start) + 1 + helpLines

	// Pad to terminal height if needed
	if totalRenderedLines < p.height {
		padding := p.height - totalRenderedLines
		for i := 0; i < padding; i++ {
			b.WriteString("\n")
		}
	}

	output := b.String()

	// Enforce terminal height limit (truncate if over)
	outputLines := strings.Split(output, "\n")
	if len(outputLines) > p.height {
		output = strings.Join(outputLines[:p.height], "\n")
	}

	return output
}

// loadDocumentCmd loads the document from repository
func (p *DocumentViewerPresenter) loadDocumentCmd() tea.Cmd {
	return func() tea.Msg {
		vm, err := queries.LoadDocumentData(p.ctx, p.repo, p.documentID)
		if err != nil {
			return DocumentLoadedMsg{ViewModel: nil, Error: err}
		}
		return DocumentLoadedMsg{ViewModel: vm, Error: nil}
	}
}

// approveDocumentCmd updates document status to published
func (p *DocumentViewerPresenter) approveDocumentCmd() tea.Cmd {
	return func() tea.Msg {
		doc, err := p.repo.FindDocumentByID(p.ctx, p.documentID)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to load document: %w", err)}
		}

		// Update status to published
		err = doc.UpdateStatus("published")
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to update document status: %w", err)}
		}

		// Save updated document
		err = p.repo.UpdateDocument(p.ctx, doc)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to save document: %w", err)}
		}

		return DocumentActionCompletedMsg{}
	}
}

// disapproveDocumentCmd updates document status to draft
func (p *DocumentViewerPresenter) disapproveDocumentCmd() tea.Cmd {
	return func() tea.Msg {
		doc, err := p.repo.FindDocumentByID(p.ctx, p.documentID)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to load document: %w", err)}
		}

		// Update status to draft
		err = doc.UpdateStatus("draft")
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to update document status: %w", err)}
		}

		// Save updated document
		err = p.repo.UpdateDocument(p.ctx, doc)
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to save document: %w", err)}
		}

		return DocumentActionCompletedMsg{}
	}
}

// renderMarkdown renders the markdown content using glamour
func (p *DocumentViewerPresenter) renderMarkdown() {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(p.width - 4), // Account for padding
	)

	if err == nil {
		rendered, err := renderer.Render(p.viewModel.Content)
		if err == nil {
			p.renderedContent = rendered
			return
		}
	}

	// Fallback to simple text wrapping if glamour fails
	wrapped := wordwrap.String(p.viewModel.Content, p.width-4)
	p.renderedContent = wrapped
}

// getStatusColor returns a function that styles text with the given color code
func getStatusColor(colorCode string) func(string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(colorCode))
	return func(text string) string {
		return style.Render(text)
	}
}

// scrollOffsetUp decreases scroll offset by 1 (up)
// Scrolls up by one line using direct offset manipulation
func (p *DocumentViewerPresenter) scrollOffsetUp() {
	contentLines := strings.Split(p.renderedContent, "\n")
	p.scrollHelper.ScrollLineUp(len(contentLines))
}

// scrollOffsetDown increases scroll offset by 1 (down)
// Scrolls down by one line using direct offset manipulation
func (p *DocumentViewerPresenter) scrollOffsetDown(totalLines int) {
	p.scrollHelper.ScrollLineDown(totalLines)
}
