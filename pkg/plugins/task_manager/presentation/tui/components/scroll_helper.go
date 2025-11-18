package components

import "fmt"

// ScrollHelper manages viewport scrolling for single-line item lists.
// It tracks viewport offset and ensures selected items remain visible.
// This helper is reusable across all presenters (Dashboard, IterationDetail, TaskDetail).
type ScrollHelper struct {
	viewportOffset int // Index of first visible item
	viewportHeight int // Number of visible items
}

// NewScrollHelper creates a new scroll helper with default viewport height.
func NewScrollHelper() *ScrollHelper {
	return &ScrollHelper{
		viewportOffset: 0,
		viewportHeight: 10, // Default, will be updated on WindowSizeMsg
	}
}

// SetViewportHeight updates the viewport height (call on WindowSizeMsg or terminal resize).
func (s *ScrollHelper) SetViewportHeight(height int) {
	if height < 1 {
		height = 1 // Minimum 1 line
	}
	s.viewportHeight = height
}

// EnsureVisible adjusts viewport offset to keep selectedIndex visible.
// totalItems: total number of items in list
// selectedIndex: currently selected item index
func (s *ScrollHelper) EnsureVisible(totalItems, selectedIndex int) {
	if totalItems == 0 {
		s.viewportOffset = 0
		return
	}

	// Clamp selected index to valid range
	if selectedIndex < 0 {
		selectedIndex = 0
	}
	if selectedIndex >= totalItems {
		selectedIndex = totalItems - 1
	}

	// If selected item is above visible area, scroll up
	if selectedIndex < s.viewportOffset {
		s.viewportOffset = selectedIndex
	}

	// If selected item is below visible area, scroll down
	if selectedIndex >= s.viewportOffset+s.viewportHeight {
		s.viewportOffset = selectedIndex - s.viewportHeight + 1
	}

	// Clamp offset to valid range
	maxOffset := totalItems - s.viewportHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if s.viewportOffset > maxOffset {
		s.viewportOffset = maxOffset
	}
	// Note: viewportOffset cannot be negative here due to prior logic:
	// - Line 46 sets offset = selectedIndex (>= 0 after clamping)
	// - Line 51 sets offset = selectedIndex - height + 1, but only executes when
	//   selectedIndex >= offset + height, ensuring result is non-negative
}

// VisibleRange returns the start and end indices of visible items.
// Returns (start, end) where items[start:end] should be rendered.
func (s *ScrollHelper) VisibleRange(totalItems int) (start, end int) {
	start = s.viewportOffset
	end = start + s.viewportHeight
	if end > totalItems {
		end = totalItems
	}
	return start, end
}

// PageUp moves selection up by one page height and adjusts viewport.
// Returns the new selected index.
func (s *ScrollHelper) PageUp(totalItems int) int {
	newSelectedIndex := s.viewportOffset - s.viewportHeight
	if newSelectedIndex < 0 {
		newSelectedIndex = 0
	}
	s.EnsureVisible(totalItems, newSelectedIndex)
	return newSelectedIndex
}

// PageDown moves selection down by one page height and adjusts viewport.
// currentSelectedIndex: current selected index
// Returns the new selected index.
func (s *ScrollHelper) PageDown(totalItems, currentSelectedIndex int) int {
	newSelectedIndex := currentSelectedIndex + s.viewportHeight
	if newSelectedIndex >= totalItems {
		newSelectedIndex = totalItems - 1
	}
	s.EnsureVisible(totalItems, newSelectedIndex)
	return newSelectedIndex
}

// ViewportOffset returns the current viewport offset (for debugging/inspection).
func (s *ScrollHelper) ViewportOffset() int {
	return s.viewportOffset
}

// ViewportHeight returns the current viewport height.
func (s *ScrollHelper) ViewportHeight() int {
	return s.viewportHeight
}

// ScrollLineUp scrolls up by one line (for documents, not item-based lists).
// Clamps offset to valid range [0, maxOffset].
func (s *ScrollHelper) ScrollLineUp(totalLines int) {
	if s.viewportOffset > 0 {
		s.viewportOffset--
	}
}

// ScrollLineDown scrolls down by one line (for documents, not item-based lists).
// Clamps offset to valid range [0, maxOffset].
func (s *ScrollHelper) ScrollLineDown(totalLines int) {
	maxOffset := totalLines - s.viewportHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if s.viewportOffset < maxOffset {
		s.viewportOffset++
	}
}

// ScrollPageUp scrolls up by one viewport height (for documents).
// Keeps 1 line overlap for context (standard pager behavior).
// Clamps offset to valid range [0, maxOffset].
func (s *ScrollHelper) ScrollPageUp(totalLines int) {
	// Scroll by (viewport - 1) to keep last line visible
	scrollAmount := s.viewportHeight - 1
	if scrollAmount < 1 {
		scrollAmount = 1
	}
	s.viewportOffset -= scrollAmount
	if s.viewportOffset < 0 {
		s.viewportOffset = 0
	}
}

// ScrollPageDown scrolls down by one viewport height (for documents).
// Keeps 1 line overlap for context (standard pager behavior).
// Clamps offset to valid range [0, maxOffset].
func (s *ScrollHelper) ScrollPageDown(totalLines int) {
	maxOffset := totalLines - s.viewportHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	// Scroll by (viewport - 1) to keep last line visible
	scrollAmount := s.viewportHeight - 1
	if scrollAmount < 1 {
		scrollAmount = 1
	}
	s.viewportOffset += scrollAmount
	if s.viewportOffset > maxOffset {
		s.viewportOffset = maxOffset
	}
}

// ScrollToStart jumps to the beginning of the document.
func (s *ScrollHelper) ScrollToStart() {
	s.viewportOffset = 0
}

// ScrollToEnd jumps to the end of the document.
// totalLines: total number of lines in the document.
func (s *ScrollHelper) ScrollToEnd(totalLines int) {
	maxOffset := totalLines - s.viewportHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	s.viewportOffset = maxOffset
}

// ScrollPosition returns a string indicating the current scroll position.
// Returns "All" if all content fits, "Top" at beginning, "Bottom" at end,
// or percentage (e.g., "25%", "50%", "75%") for middle positions.
func (s *ScrollHelper) ScrollPosition(totalLines int) string {
	// All content visible (fits in viewport)
	if totalLines <= s.viewportHeight {
		return "All"
	}

	maxOffset := totalLines - s.viewportHeight
	if maxOffset <= 0 {
		return "All"
	}

	// At top
	if s.viewportOffset == 0 {
		return "Top"
	}

	// At bottom
	if s.viewportOffset >= maxOffset {
		return "Bot"
	}

	// Middle - calculate percentage
	percentage := int((float64(s.viewportOffset) / float64(maxOffset)) * 100)
	if percentage < 1 {
		percentage = 1
	}
	if percentage > 99 {
		percentage = 99
	}
	return fmt.Sprintf("%d%%", percentage)
}

// ScrollHelperMultiline manages viewport scrolling for multi-line item lists.
// Used when items can expand/collapse (e.g., ACs with testing instructions).
// Items are selectable by index, but scrolling is done by line count.
type ScrollHelperMultiline struct {
	viewportOffset int // Line offset (not item offset)
	viewportHeight int // Lines visible
}

// NewScrollHelperMultiline creates a new multiline scroll helper.
func NewScrollHelperMultiline() *ScrollHelperMultiline {
	return &ScrollHelperMultiline{
		viewportOffset: 0,
		viewportHeight: 10, // Default, will be updated on WindowSizeMsg
	}
}

// SetViewportHeight updates the viewport height (call on WindowSizeMsg).
func (s *ScrollHelperMultiline) SetViewportHeight(height int) {
	if height < 1 {
		height = 1 // Minimum 1 line
	}
	s.viewportHeight = height
}

// EnsureVisibleMultiline adjusts offset to keep selected item visible.
// itemLineCounts: number of lines for each item (e.g., [1, 5, 1] for collapsed/expanded items)
// selectedIndex: which item is selected
func (s *ScrollHelperMultiline) EnsureVisibleMultiline(itemLineCounts []int, selectedIndex int) {
	if len(itemLineCounts) == 0 {
		s.viewportOffset = 0
		return
	}

	// Clamp selected index to valid range
	if selectedIndex < 0 {
		selectedIndex = 0
	}
	if selectedIndex >= len(itemLineCounts) {
		selectedIndex = len(itemLineCounts) - 1
	}

	// Calculate which line the selected item starts at
	selectedLine := 0
	for i := 0; i < selectedIndex; i++ {
		selectedLine += itemLineCounts[i]
	}

	selectedItemLines := itemLineCounts[selectedIndex]

	// If selected item starts above visible area, scroll up
	if selectedLine < s.viewportOffset {
		s.viewportOffset = selectedLine
	}

	// If selected item ends below visible area, scroll down
	if selectedLine+selectedItemLines > s.viewportOffset+s.viewportHeight {
		s.viewportOffset = selectedLine + selectedItemLines - s.viewportHeight
	}

	// Clamp offset to valid range
	totalLines := 0
	for _, lines := range itemLineCounts {
		totalLines += lines
	}

	maxOffset := totalLines - s.viewportHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if s.viewportOffset > maxOffset {
		s.viewportOffset = maxOffset
	}
	// Note: viewportOffset cannot be negative here (same logic as ScrollHelper)
}

// VisibleRangeMultiline returns (firstVisibleItem, lastVisibleItem, lineOffset).
// firstVisibleItem: index of first visible item
// lastVisibleItem: index of last visible item (inclusive)
// lineOffset: number of lines to skip from firstVisibleItem
func (s *ScrollHelperMultiline) VisibleRangeMultiline(itemLineCounts []int) (firstItem, lastItem, lineOffset int) {
	if len(itemLineCounts) == 0 {
		return 0, -1, 0
	}

	// Find first visible item (line where viewport starts)
	currentLine := 0
	firstItem = 0
	for i := 0; i < len(itemLineCounts); i++ {
		if currentLine+itemLineCounts[i] > s.viewportOffset {
			firstItem = i
			break
		}
		currentLine += itemLineCounts[i]
	}

	// Calculate line offset within first item
	lineOffset = s.viewportOffset - currentLine
	// Note: lineOffset should not be negative as currentLine is the start of firstItem,
	// and viewportOffset should be >= currentLine by construction of the loop above

	// Find last visible item (line where viewport ends)
	currentLine = 0
	lastItem = len(itemLineCounts) - 1
	visibleEndLine := s.viewportOffset + s.viewportHeight
	for i := 0; i < len(itemLineCounts); i++ {
		if currentLine >= visibleEndLine {
			lastItem = i - 1
			// Note: lastItem cannot be negative as loop starts at i=0,
			// and currentLine >= visibleEndLine can only be true when i > 0
			// (since currentLine starts at 0 and visibleEndLine > 0)
			break
		}
		currentLine += itemLineCounts[i]
	}

	return firstItem, lastItem, lineOffset
}

// ViewportOffset returns the current line offset (for debugging/inspection).
func (s *ScrollHelperMultiline) ViewportOffset() int {
	return s.viewportOffset
}
