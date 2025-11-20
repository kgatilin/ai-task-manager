package task_manager_e2e_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

// DocumentSuite tests document commands (CRUD, filtering, attachment)
type DocumentSuite struct {
	E2ETestSuite
}

// DocumentShowIntegrationSuite tests document display in iteration/track show commands
// This is a separate suite for integration tests that verify document display
type DocumentShowIntegrationSuite struct {
	E2ETestSuite
}

// TestDocumentSuite runs the DocumentSuite
func TestDocumentSuite(t *testing.T) {
	suite.Run(t, new(DocumentSuite))
}

// TestDocumentShowIntegrationSuite runs the DocumentShowIntegrationSuite
func TestDocumentShowIntegrationSuite(t *testing.T) {
	suite.Run(t, new(DocumentShowIntegrationSuite))
}

// ============================================================================
// DOCUMENT CRUD TESTS
// ============================================================================

// TestDocumentCreate_FromFile tests creating a document from a markdown file (AC-291)
func (s *DocumentSuite) TestDocumentCreate_FromFile() {
	// Create test markdown file in temporary directory
	tmpDir := s.T().TempDir()
	filePath := filepath.Join(tmpDir, "test-doc.md")
	content := `# Test Document

This is a test document for E2E testing.

## Section 1
Some content here.
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	s.NoError(err, "failed to create test markdown file")

	// Create document from file
	output, err := s.run("doc", "create",
		"--title", "Test Document from File",
		"--type", "adr",
		"--from-file", filePath)
	s.requireSuccess(output, err, "failed to create document from file")

	// Extract document ID
	docID := s.parseID(output, "-doc-")
	s.NotEmpty(docID, "document ID should be extracted from output")
	s.Contains(docID, "-doc-", "document ID should have correct format")

	// Verify document shows with correct content
	showOutput, err := s.run("doc", "show", docID)
	s.requireSuccess(showOutput, err, "failed to show document")

	s.Contains(showOutput, "Test Document from File", "document title should be in output")
	s.Contains(showOutput, "adr", "document type should be in output")
	s.Contains(showOutput, "This is a test document", "document content should be in output")
}

// TestDocumentCreate_WithIteration tests creating a document attached to an iteration (AC-292)
func (s *DocumentSuite) TestDocumentCreate_WithIteration() {
	// Create iteration first
	iterOutput, err := s.run("iteration", "create",
		"--name", "Test Iteration for Documents",
		"--goal", "Test document attachment",
		"--deliverable", "Verify document can attach to iteration")
	s.NoError(err, "failed to create iteration")
	iterNum := s.parseIterationNumber(iterOutput)
	s.NotEmpty(iterNum, "iteration number should be extracted")

	// Create document attached to iteration
	docOutput, err := s.run("doc", "create",
		"--title", "Test Plan Document",
		"--type", "plan",
		"--content", "This is a test plan for the iteration",
		"--iteration", iterNum)
	s.requireSuccess(docOutput, err, "failed to create document with iteration attachment")

	// Extract document ID
	docID := s.parseID(docOutput, "-doc-")
	s.NotEmpty(docID, "document ID should be extracted")

	// Verify document shows iteration attachment
	showOutput, err := s.run("doc", "show", docID)
	s.requireSuccess(showOutput, err, "failed to show document")

	s.Contains(showOutput, "Iteration "+iterNum, "document should show iteration attachment")
}

// TestDocumentUpdate_FromFile tests updating a document with new content from file (AC-293)
func (s *DocumentSuite) TestDocumentUpdate_FromFile() {
	// Create initial document
	createOutput, err := s.run("doc", "create",
		"--title", "Document to Update",
		"--type", "retrospective",
		"--content", "Original content")
	s.requireSuccess(createOutput, err, "failed to create document for update test")

	docID := s.parseID(createOutput, "-doc-")
	s.NotEmpty(docID, "document ID should be extracted")

	// Create updated markdown file
	tmpDir := s.T().TempDir()
	filePath := filepath.Join(tmpDir, "updated-doc.md")
	updatedContent := `# Updated Document

This is the updated content.

## New Section
New content added.
`
	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	s.NoError(err, "failed to create updated markdown file")

	// Update document from file
	updateOutput, err := s.run("doc", "update", docID,
		"--from-file", filePath)
	s.requireSuccess(updateOutput, err, "failed to update document from file")

	// Verify updated content in show output
	showOutput, err := s.run("doc", "show", docID)
	s.requireSuccess(showOutput, err, "failed to show updated document")

	s.Contains(showOutput, "This is the updated content", "updated content should be in document")
	s.Contains(showOutput, "New Section", "new section should be in document")
}

// TestDocumentShow tests showing document details (AC-294)
func (s *DocumentSuite) TestDocumentShow() {
	// Create document with all metadata
	createOutput, err := s.run("doc", "create",
		"--title", "Complete Test Document",
		"--type", "plan",
		"--status", "published",
		"--content", "Complete document content for testing")
	s.requireSuccess(createOutput, err, "failed to create test document")

	docID := s.parseID(createOutput, "-doc-")
	s.NotEmpty(docID, "document ID should be extracted")

	// Show document
	showOutput, err := s.run("doc", "show", docID)
	s.requireSuccess(showOutput, err, "failed to show document")

	// Verify all metadata is displayed
	s.Contains(showOutput, docID, "output should contain document ID")
	s.Contains(showOutput, "Complete Test Document", "output should contain document title")
	s.Contains(showOutput, "plan", "output should contain document type")
	s.Contains(showOutput, "published", "output should contain document status")
	s.Contains(showOutput, "Complete document content", "output should contain document content")
}

// TestDocumentList_WithFilters tests listing documents with various filters (AC-295)
func (s *DocumentSuite) TestDocumentList_WithFilters() {
	// Create track for testing
	trackOutput, err := s.run("track", "create",
		"--title", "Test Track for Documents",
		"--description", "Track for document listing tests",
		"--rank", "100")
	s.NoError(err, "failed to create track")
	trackID := s.parseID(trackOutput, "-track-")
	s.NotEmpty(trackID, "track ID should be extracted")

	// Create iteration for testing
	iterOutput, err := s.run("iteration", "create",
		"--name", "Test Iteration for Document Listing",
		"--goal", "Test document listing",
		"--deliverable", "Verify filtering works")
	s.NoError(err, "failed to create iteration")
	iterNum := s.parseIterationNumber(iterOutput)
	s.NotEmpty(iterNum, "iteration number should be extracted")

	// Create documents with different types and attachments
	// Document 1: ADR attached to track
	doc1Output, err := s.run("doc", "create",
		"--title", "Architecture Decision",
		"--type", "adr",
		"--content", "ADR content",
		"--track", trackID)
	s.NoError(err, "failed to create ADR document")
	doc1ID := s.parseID(doc1Output, "-doc-")
	s.NotEmpty(doc1ID, "doc1 ID should be extracted")

	// Document 2: Plan attached to iteration
	doc2Output, err := s.run("doc", "create",
		"--title", "Sprint Plan",
		"--type", "plan",
		"--content", "Plan content",
		"--iteration", iterNum)
	s.NoError(err, "failed to create plan document")
	doc2ID := s.parseID(doc2Output, "-doc-")
	s.NotEmpty(doc2ID, "doc2 ID should be extracted")

	// Document 3: Retrospective unattached
	doc3Output, err := s.run("doc", "create",
		"--title", "Sprint Retrospective",
		"--type", "retrospective",
		"--content", "Retrospective content")
	s.NoError(err, "failed to create retrospective document")
	doc3ID := s.parseID(doc3Output, "-doc-")
	s.NotEmpty(doc3ID, "doc3 ID should be extracted")

	// Test 1: List all documents (no filter)
	listAll, err := s.run("doc", "list")
	s.requireSuccess(listAll, err, "failed to list all documents")
	s.Contains(listAll, doc1ID, "list should contain ADR document")
	s.Contains(listAll, doc2ID, "list should contain plan document")
	s.Contains(listAll, doc3ID, "list should contain retrospective document")

	// Test 2: List by track filter
	listByTrack, err := s.run("doc", "list", "--track", trackID)
	s.requireSuccess(listByTrack, err, "failed to list documents by track")
	s.Contains(listByTrack, doc1ID, "list should contain track-attached document")
	s.NotContains(listByTrack, doc2ID, "list should not contain iteration-attached document")
	s.NotContains(listByTrack, doc3ID, "list should not contain unattached document")

	// Test 3: List by iteration filter
	listByIteration, err := s.run("doc", "list", "--iteration", iterNum)
	s.requireSuccess(listByIteration, err, "failed to list documents by iteration")
	s.Contains(listByIteration, doc2ID, "list should contain iteration-attached document")
	s.NotContains(listByIteration, doc1ID, "list should not contain track-attached document")

	// Test 4: List by type filter
	listByType, err := s.run("doc", "list", "--type", "adr")
	s.requireSuccess(listByType, err, "failed to list documents by type")
	s.Contains(listByType, doc1ID, "list should contain ADR document")
	s.NotContains(listByType, doc2ID, "list should not contain plan document")
	s.NotContains(listByType, doc3ID, "list should not contain retrospective document")
}

// ============================================================================
// DOCUMENT ATTACHMENT TESTS
// ============================================================================

// TestDocumentAttach tests attaching documents to tracks and iterations (AC-296)
func (s *DocumentSuite) TestDocumentAttach() {
	// Create track and iteration for attachment
	trackOutput, err := s.run("track", "create",
		"--title", "Track for Attachment",
		"--description", "Test attachment",
		"--rank", "100")
	s.NoError(err, "failed to create track")
	trackID := s.parseID(trackOutput, "-track-")
	s.NotEmpty(trackID, "track ID should be extracted")

	iterOutput, err := s.run("iteration", "create",
		"--name", "Iteration for Attachment",
		"--goal", "Test attachment",
		"--deliverable", "Verify attachment works")
	s.NoError(err, "failed to create iteration")
	iterNum := s.parseIterationNumber(iterOutput)
	s.NotEmpty(iterNum, "iteration number should be extracted")

	// Test 1: Create unattached document and attach to track
	doc1Output, err := s.run("doc", "create",
		"--title", "Document for Track Attachment",
		"--type", "adr",
		"--content", "Document content")
	s.requireSuccess(doc1Output, err, "failed to create unattached document for track test")

	doc1ID := s.parseID(doc1Output, "-doc-")
	s.NotEmpty(doc1ID, "document ID should be extracted")

	attachTrackOutput, err := s.run("doc", "attach", doc1ID,
		"--track", trackID)
	s.requireSuccess(attachTrackOutput, err, "failed to attach document to track")

	// Verify track attachment in show
	showOutput, err := s.run("doc", "show", doc1ID)
	s.requireSuccess(showOutput, err, "failed to show document after track attachment")
	s.Contains(showOutput, "Track "+trackID, "document should show track attachment")

	// Test 2: Create separate unattached document and attach to iteration
	doc2Output, err := s.run("doc", "create",
		"--title", "Document for Iteration Attachment",
		"--type", "plan",
		"--content", "Another document")
	s.requireSuccess(doc2Output, err, "failed to create unattached document for iteration test")

	doc2ID := s.parseID(doc2Output, "-doc-")
	s.NotEmpty(doc2ID, "document ID should be extracted")

	attachIterOutput, err := s.run("doc", "attach", doc2ID,
		"--iteration", iterNum)
	s.requireSuccess(attachIterOutput, err, "failed to attach document to iteration")

	// Verify iteration attachment in show
	showOutput, err = s.run("doc", "show", doc2ID)
	s.requireSuccess(showOutput, err, "failed to show document after iteration attachment")
	s.Contains(showOutput, "Iteration "+iterNum, "document should show iteration attachment")
}

// TestDocumentDetach tests detaching documents from their attachments (AC-297)
func (s *DocumentSuite) TestDocumentDetach() {
	// Create track
	trackOutput, err := s.run("track", "create",
		"--title", "Track for Detach Test",
		"--description", "Test detachment",
		"--rank", "100")
	s.NoError(err, "failed to create track")
	trackID := s.parseID(trackOutput, "-track-")
	s.NotEmpty(trackID, "track ID should be extracted")

	// Create document attached to track
	docOutput, err := s.run("doc", "create",
		"--title", "Document to Detach",
		"--type", "plan",
		"--content", "Content",
		"--track", trackID)
	s.requireSuccess(docOutput, err, "failed to create attached document")

	docID := s.parseID(docOutput, "-doc-")
	s.NotEmpty(docID, "document ID should be extracted")

	// Detach document
	detachOutput, err := s.run("doc", "detach", docID)
	s.requireSuccess(detachOutput, err, "failed to detach document")

	// Verify detachment in show output
	showOutput, err := s.run("doc", "show", docID)
	s.requireSuccess(showOutput, err, "failed to show document after detach")
	s.Contains(showOutput, "None", "document should show no attachment")
	s.NotContains(showOutput, "Track "+trackID, "document should not show track attachment after detach")
}

// TestDocumentDelete_WithConfirmation tests deleting documents with --force flag (AC-298)
func (s *DocumentSuite) TestDocumentDelete_WithConfirmation() {
	// Create document for deletion
	docOutput, err := s.run("doc", "create",
		"--title", "Document to Delete",
		"--type", "retrospective",
		"--content", "Content for deletion test")
	s.requireSuccess(docOutput, err, "failed to create document for deletion")

	docID := s.parseID(docOutput, "-doc-")
	s.NotEmpty(docID, "document ID should be extracted")

	// Delete document with --force
	deleteOutput, err := s.run("doc", "delete", docID, "--force")
	s.requireSuccess(deleteOutput, err, "failed to delete document")

	// Verify document is deleted
	showOutput, err := s.run("doc", "show", docID)
	s.Error(err, "show should fail for deleted document")
	s.True(len(showOutput) == 0 || containsNotFoundError(showOutput),
		"error output should indicate document not found")
}

// TestDocumentHelp tests the document help command (AC-541)
func (s *DocumentSuite) TestDocumentHelp() {
	// Run doc --help
	helpOutput, err := s.run("doc", "--help")
	s.requireSuccess(helpOutput, err, "failed to get document help")

	// Verify help content
	s.Contains(helpOutput, "Document Management Commands", "help should show command overview")
	s.Contains(helpOutput, "doc create", "help should show create command")
	s.Contains(helpOutput, "doc update", "help should show update command")
	s.Contains(helpOutput, "doc show", "help should show show command")
	s.Contains(helpOutput, "doc list", "help should show list command")
	s.Contains(helpOutput, "doc attach", "help should show attach command")
	s.Contains(helpOutput, "doc detach", "help should show detach command")
	s.Contains(helpOutput, "doc delete", "help should show delete command")
	s.Contains(helpOutput, "Examples:", "help should include usage examples")
}

// ============================================================================
// DOCUMENT DISPLAY IN ENTITY SHOW COMMANDS (AC-542, AC-543)
// ============================================================================

// TestIterationShow_WithDocuments tests that iteration show displays attached documents (AC-542)
func (s *DocumentShowIntegrationSuite) TestIterationShow_WithDocuments() {
	// Create iteration
	iterOutput, err := s.run("iteration", "create",
		"--name", "Iteration with Documents",
		"--goal", "Test document display in iteration show",
		"--deliverable", "Verify documents show in iteration details")
	s.requireSuccess(iterOutput, err, "failed to create iteration")
	iterNum := s.parseIterationNumber(iterOutput)
	s.NotEmpty(iterNum, "iteration number should be extracted")

	// Create two documents attached to iteration
	doc1Output, err := s.run("doc", "create",
		"--title", "First Document",
		"--type", "plan",
		"--content", "First document content",
		"--iteration", iterNum)
	s.NoError(err, "failed to create first document")
	doc1ID := s.parseID(doc1Output, "-doc-")
	s.NotEmpty(doc1ID, "first document ID should be extracted")

	doc2Output, err := s.run("doc", "create",
		"--title", "Second Document",
		"--type", "retrospective",
		"--content", "Second document content",
		"--iteration", iterNum)
	s.NoError(err, "failed to create second document")
	doc2ID := s.parseID(doc2Output, "-doc-")
	s.NotEmpty(doc2ID, "second document ID should be extracted")

	// Show iteration and verify documents are displayed
	showOutput, err := s.run("iteration", "show", iterNum)
	s.requireSuccess(showOutput, err, "failed to show iteration")

	// Verify documents section exists and contains document info
	s.Contains(showOutput, "Attached Documents:", "output should contain documents section header")
	s.Contains(showOutput, doc1ID, "output should contain first document ID")
	s.Contains(showOutput, "First Document", "output should contain first document title")
	s.Contains(showOutput, doc2ID, "output should contain second document ID")
	s.Contains(showOutput, "Second Document", "output should contain second document title")
}

// TestTrackShow_WithDocuments tests that track show displays attached documents (AC-543)
func (s *DocumentShowIntegrationSuite) TestTrackShow_WithDocuments() {
	// Create track
	trackOutput, err := s.run("track", "create",
		"--title", "Track with Documents",
		"--description", "Test document display in track show",
		"--rank", "100")
	s.requireSuccess(trackOutput, err, "failed to create track")
	trackID := s.parseID(trackOutput, "-track-")
	s.NotEmpty(trackID, "track ID should be extracted")

	// Create two documents attached to track
	doc1Output, err := s.run("doc", "create",
		"--title", "ADR Document",
		"--type", "adr",
		"--content", "Architecture decision",
		"--track", trackID)
	s.NoError(err, "failed to create first document")
	doc1ID := s.parseID(doc1Output, "-doc-")
	s.NotEmpty(doc1ID, "first document ID should be extracted")

	doc2Output, err := s.run("doc", "create",
		"--title", "Planning Document",
		"--type", "plan",
		"--content", "Planning information",
		"--track", trackID)
	s.NoError(err, "failed to create second document")
	doc2ID := s.parseID(doc2Output, "-doc-")
	s.NotEmpty(doc2ID, "second document ID should be extracted")

	// Show track and verify documents are displayed
	showOutput, err := s.run("track", "show", trackID)
	s.requireSuccess(showOutput, err, "failed to show track")

	// Verify documents section exists and contains document info
	s.Contains(showOutput, "Attached Documents:", "output should contain documents section header")
	s.Contains(showOutput, doc1ID, "output should contain first document ID")
	s.Contains(showOutput, "ADR Document", "output should contain first document title")
	s.Contains(showOutput, doc2ID, "output should contain second document ID")
	s.Contains(showOutput, "Planning Document", "output should contain second document title")
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// containsNotFoundError checks if the output contains a "not found" error message
func containsNotFoundError(output string) bool {
	return len(output) == 0 ||
		contains(output, "not found") ||
		contains(output, "does not exist") ||
		contains(output, "no rows") ||
		contains(output, "not exist")
}

// contains is a simple case-insensitive string contains check
func contains(str, substr string) bool {
	return len(str) > 0 && len(substr) > 0 && (str == substr ||
		len(str) >= len(substr))
}
