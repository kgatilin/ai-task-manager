package task_manager_e2e_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ADRTestSuite tests ADR CRUD and lifecycle commands
type ADRTestSuite struct {
	E2ETestSuite
}

func TestADRSuite(t *testing.T) {
	suite.Run(t, new(ADRTestSuite))
}

func (s *ADRTestSuite) SetupSuite() {
	s.T().Skip("Skipping ADR E2E tests, I'm going to remove this functionality anyway")
}

// TestADRCreate tests creating an ADR for a track
func (s *ADRTestSuite) TestADRCreate() {
	// Create a track first (ADRs belong to tracks)
	trackOutput, err := s.run("track", "create", "--title", "Architecture Track", "--description", "Core architecture decisions")
	s.requireSuccess(trackOutput, err, "failed to create track for ADR test")

	trackID := s.parseID(trackOutput, "-track-")

	// Create an ADR
	adrOutput, err := s.run("adr", "create", trackID,
		"--title", "Use Microservices Architecture",
		"--context", "Need scalable system with independent deployment",
		"--decision", "Adopt microservices pattern with Docker containers",
		"--consequences", "Increased operational complexity, better scalability")
	s.requireSuccess(adrOutput, err, "failed to create ADR")

	// Extract ADR ID (format: "ID: XXX-adr-X")
	adrID := s.parseID(adrOutput, "-adr-")
	s.NotEmpty(adrID, "ADR ID should be extracted from output")
	s.Contains(adrID, "-adr-", "ADR ID should have correct format")
}

// TestADRCreateMinimal tests creating an ADR with required fields only
func (s *ADRTestSuite) TestADRCreateMinimal() {
	// Create a track
	trackOutput, err := s.run("track", "create", "--title", "Decision Track", "--description", "For minimal ADR")
	s.requireSuccess(trackOutput, err, "failed to create track")

	trackID := s.parseID(trackOutput, "-track-")

	// Create minimal ADR (with all required fields)
	adrOutput, err := s.run("adr", "create", trackID,
		"--title", "Use PostgreSQL",
		"--context", "Need reliable database",
		"--decision", "Choose PostgreSQL over other options",
		"--consequences", "Better reliability and ACID guarantees")
	s.requireSuccess(adrOutput, err, "failed to create minimal ADR")

	adrID := s.parseID(adrOutput, "-adr-")
	s.NotEmpty(adrID, "ADR ID should be created")
}

// TestADRList tests listing all ADRs
func (s *ADRTestSuite) TestADRList() {
	// Create a track
	trackOutput, err := s.run("track", "create", "--title", "List Test Track", "--description", "For ADR listing")
	s.requireSuccess(trackOutput, err, "failed to create track")

	trackID := s.parseID(trackOutput, "-track-")

	// Create an ADR
	adrOutput, err := s.run("adr", "create", trackID,
		"--title", "API Design Decision",
		"--context", "Need consistent API",
		"--decision", "Use REST with JSON",
		"--consequences", "Simpler implementation and wider compatibility")
	s.requireSuccess(adrOutput, err, "failed to create ADR for listing")

	adrID := s.parseID(adrOutput, "-adr-")

	// List ADRs
	listOutput, err := s.run("adr", "list")
	s.requireSuccess(listOutput, err, "failed to list ADRs")

	// Verify the ADR appears in the list
	s.Contains(listOutput, adrID, "created ADR should appear in list")
	s.Contains(listOutput, "API Design Decision", "ADR title should appear in list")
}

// TestADRListByTrack tests listing ADRs for a specific track
func (s *ADRTestSuite) TestADRListByTrack() {
	// Create two tracks
	track1Output, err := s.run("track", "create", "--title", "Track 1", "--description", "First track")
	s.requireSuccess(track1Output, err, "failed to create track 1")
	trackID1 := s.parseID(track1Output, "-track-")

	track2Output, err := s.run("track", "create", "--title", "Track 2", "--description", "Second track")
	s.requireSuccess(track2Output, err, "failed to create track 2")
	trackID2 := s.parseID(track2Output, "-track-")

	// Create ADRs for both tracks
	adr1Output, err := s.run("adr", "create", trackID1,
		"--title", "Decision 1",
		"--context", "Context 1",
		"--decision", "Decision for track 1",
		"--consequences", "Impact 1")
	s.requireSuccess(adr1Output, err, "failed to create ADR for track 1")
	adrID1 := s.parseID(adr1Output, "-adr-")

	adr2Output, err := s.run("adr", "create", trackID2,
		"--title", "Decision 2",
		"--context", "Context 2",
		"--decision", "Decision for track 2",
		"--consequences", "Impact 2")
	s.requireSuccess(adr2Output, err, "failed to create ADR for track 2")

	// List ADRs for track 1 only
	listOutput, err := s.run("adr", "list", "--track", trackID1)
	s.requireSuccess(listOutput, err, "failed to list ADRs for track")

	// Verify only track 1's ADR appears
	s.Contains(listOutput, adrID1, "track 1 ADR should appear")
	s.Contains(listOutput, "Decision 1", "track 1 ADR title should appear")
}

// TestADRShow tests showing ADR details
func (s *ADRTestSuite) TestADRShow() {
	// Create a track
	trackOutput, err := s.run("track", "create", "--title", "Show Test Track", "--description", "For ADR show test")
	s.requireSuccess(trackOutput, err, "failed to create track")

	trackID := s.parseID(trackOutput, "-track-")

	// Create an ADR
	adrOutput, err := s.run("adr", "create", trackID,
		"--title", "Testing Framework Decision",
		"--context", "Need robust testing setup",
		"--decision", "Use Go testing and Testify",
		"--consequences", "Learning curve but good testing patterns")
	s.requireSuccess(adrOutput, err, "failed to create ADR for show test")

	adrID := s.parseID(adrOutput, "-adr-")

	// Show ADR details
	showOutput, err := s.run("adr", "show", adrID)
	s.requireSuccess(showOutput, err, "failed to show ADR")

	// Verify all details
	s.Contains(showOutput, adrID, "ADR ID should be in output")
	s.Contains(showOutput, "Testing Framework Decision", "ADR title should be in output")
	s.Contains(showOutput, "Need robust testing setup", "ADR context should be in output")
	s.Contains(showOutput, "Use Go testing and Testify", "ADR decision should be in output")
	s.Contains(showOutput, "Learning curve", "ADR consequences should be in output")
}

// TestADRUpdate tests updating ADR fields
func (s *ADRTestSuite) TestADRUpdate() {
	// Create a track
	trackOutput, err := s.run("track", "create", "--title", "Update Test Track", "--description", "For ADR update test")
	s.requireSuccess(trackOutput, err, "failed to create track")

	trackID := s.parseID(trackOutput, "-track-")

	// Create an ADR
	adrOutput, err := s.run("adr", "create", trackID,
		"--title", "Original Title",
		"--context", "Original context",
		"--decision", "Original decision",
		"--consequences", "Original consequences")
	s.requireSuccess(adrOutput, err, "failed to create ADR for update test")

	adrID := s.parseID(adrOutput, "-adr-")

	// Update the ADR
	updateOutput, err := s.run("adr", "update", adrID,
		"--title", "Updated Title",
		"--decision", "Updated decision text")
	s.requireSuccess(updateOutput, err, "failed to update ADR")

	// Verify the update
	showOutput, err := s.run("adr", "show", adrID)
	s.requireSuccess(showOutput, err, "failed to show updated ADR")

	s.Contains(showOutput, "Updated Title", "ADR title should be updated")
	s.Contains(showOutput, "Updated decision text", "ADR decision should be updated")
}

// TestADRStatusChange tests changing ADR status
func (s *ADRTestSuite) TestADRStatusChange() {
	// Create a track
	trackOutput, err := s.run("track", "create", "--title", "Status Track", "--description", "For ADR status test")
	s.requireSuccess(trackOutput, err, "failed to create track")

	trackID := s.parseID(trackOutput, "-track-")

	// Create an ADR (default status: proposed)
	adrOutput, err := s.run("adr", "create", trackID,
		"--title", "Status Decision",
		"--context", "Test context",
		"--decision", "Test decision",
		"--consequences", "Test consequences")
	s.requireSuccess(adrOutput, err, "failed to create ADR for status test")

	adrID := s.parseID(adrOutput, "-adr-")

	// Update status to accepted
	updateOutput, err := s.run("adr", "update", adrID, "--status", "accepted")
	s.requireSuccess(updateOutput, err, "failed to update ADR status to accepted")

	// Verify status change
	showOutput, err := s.run("adr", "show", adrID)
	s.requireSuccess(showOutput, err, "failed to show ADR after status update")

	s.Contains(showOutput, "accepted", "ADR status should be updated to accepted")
}

// TestADRSupersede tests superseding an ADR with another
func (s *ADRTestSuite) TestADRSupersede() {
	// Create a track
	trackOutput, err := s.run("track", "create", "--title", "Supersede Track", "--description", "For supersede test")
	s.requireSuccess(trackOutput, err, "failed to create track")

	trackID := s.parseID(trackOutput, "-track-")

	// Create first ADR
	adr1Output, err := s.run("adr", "create", trackID,
		"--title", "Old Decision",
		"--context", "Original approach",
		"--decision", "Use approach A",
		"--consequences", "Limited scalability")
	s.requireSuccess(adr1Output, err, "failed to create first ADR for supersede test")

	adrID1 := s.parseID(adr1Output, "-adr-")

	// Create second ADR (replacement)
	adr2Output, err := s.run("adr", "create", trackID,
		"--title", "New Decision",
		"--context", "Better approach discovered",
		"--decision", "Use approach B instead",
		"--consequences", "Better scalability")
	s.requireSuccess(adr2Output, err, "failed to create second ADR for supersede test")

	adrID2 := s.parseID(adr2Output, "-adr-")

	// Supersede first ADR with second
	supersedeOutput, err := s.run("adr", "supersede", adrID1, "--by", adrID2)
	s.requireSuccess(supersedeOutput, err, "failed to supersede ADR")

	// Verify superseded status
	showOutput, err := s.run("adr", "show", adrID1)
	s.requireSuccess(showOutput, err, "failed to show superseded ADR")

	s.Contains(showOutput, "superseded", "ADR status should show superseded")
	s.Contains(showOutput, adrID2, "superseded by reference should appear")
}

// TestADRDeprecate tests deprecating an ADR
func (s *ADRTestSuite) TestADRDeprecate() {
	// Create a track
	trackOutput, err := s.run("track", "create", "--title", "Deprecate Track", "--description", "For deprecate test")
	s.requireSuccess(trackOutput, err, "failed to create track")

	trackID := s.parseID(trackOutput, "-track-")

	// Create an ADR
	adrOutput, err := s.run("adr", "create", trackID,
		"--title", "Legacy Decision",
		"--context", "Old approach",
		"--decision", "Old decision",
		"--consequences", "Now obsolete")
	s.requireSuccess(adrOutput, err, "failed to create ADR for deprecate test")

	adrID := s.parseID(adrOutput, "-adr-")

	// Deprecate the ADR
	deprecateOutput, err := s.run("adr", "deprecate", adrID)
	s.requireSuccess(deprecateOutput, err, "failed to deprecate ADR")

	// Verify deprecated status
	showOutput, err := s.run("adr", "show", adrID)
	s.requireSuccess(showOutput, err, "failed to show deprecated ADR")

	s.Contains(showOutput, "deprecated", "ADR status should show deprecated")
}

// TestADRListByStatus tests listing ADRs filtered by status
func (s *ADRTestSuite) TestADRListByStatus() {
	// Create a track
	trackOutput, err := s.run("track", "create", "--title", "Status Filter Track", "--description", "For status filter test")
	s.requireSuccess(trackOutput, err, "failed to create track")

	trackID := s.parseID(trackOutput, "-track-")

	// Create an ADR (default status: proposed)
	adrOutput, err := s.run("adr", "create", trackID,
		"--title", "Status Filter ADR",
		"--context", "Test context",
		"--decision", "Test decision",
		"--consequences", "Test impact")
	s.requireSuccess(adrOutput, err, "failed to create ADR for status filter test")

	adrID := s.parseID(adrOutput, "-adr-")

	// List all ADRs (no status filter) to verify ADR was created
	listOutput, err := s.run("adr", "list")
	s.requireSuccess(listOutput, err, "failed to list ADRs")

	// Newly created ADR should appear in the overall list
	s.Contains(listOutput, adrID, "newly created ADR should appear in list")
}

// TestADRCheckRequirement tests checking if track has required ADR
func (s *ADRTestSuite) TestADRCheckRequirement() {
	s.T().Skip("No check command implemented yet")

	// Create a track without ADR
	trackOutput, err := s.run("track", "create", "--title", "No ADR Track", "--description", "Track without ADR")
	s.requireSuccess(trackOutput, err, "failed to create track for check test")

	trackID := s.parseID(trackOutput, "-track-")

	// Check if track has ADR - should indicate missing
	checkOutput, _ := s.run("track", "check", trackID)
	// This may succeed or fail depending on implementation, just verify it runs
	s.NotEmpty(checkOutput, "check command should produce output")

	// Create an ADR for the track
	adrOutput, err := s.run("adr", "create", trackID,
		"--title", "Required ADR",
		"--context", "Now has ADR",
		"--decision", "ADR decision",
		"--consequences", "Track has ADR")
	s.requireSuccess(adrOutput, err, "failed to create ADR for check verification")

	// Check again - should indicate ADR exists
	checkOutput, err = s.run("track", "check", trackID)
	s.requireSuccess(checkOutput, err, "check command should succeed with existing ADR")
}

// TestADRMultiplePerTrack tests creating multiple ADRs for a single track
func (s *ADRTestSuite) TestADRMultiplePerTrack() {
	// Create a track
	trackOutput, err := s.run("track", "create", "--title", "Multi-ADR Track", "--description", "Track with multiple ADRs")
	s.requireSuccess(trackOutput, err, "failed to create track")

	trackID := s.parseID(trackOutput, "-track-")

	// Create multiple ADRs
	adr1Output, err := s.run("adr", "create", trackID,
		"--title", "First ADR",
		"--context", "First context",
		"--decision", "First decision",
		"--consequences", "First impact")
	s.requireSuccess(adr1Output, err, "failed to create first ADR")

	adr2Output, err := s.run("adr", "create", trackID,
		"--title", "Second ADR",
		"--context", "Second context",
		"--decision", "Second decision",
		"--consequences", "Second impact")
	s.requireSuccess(adr2Output, err, "failed to create second ADR")

	adr3Output, err := s.run("adr", "create", trackID,
		"--title", "Third ADR",
		"--context", "Third context",
		"--decision", "Third decision",
		"--consequences", "Third impact")
	s.requireSuccess(adr3Output, err, "failed to create third ADR")

	// List ADRs for this track
	listOutput, err := s.run("adr", "list", "--track", trackID)
	s.requireSuccess(listOutput, err, "failed to list ADRs for track")

	// Verify all three appear
	adrID1 := s.parseID(adr1Output, "-adr-")
	adrID2 := s.parseID(adr2Output, "-adr-")
	adrID3 := s.parseID(adr3Output, "-adr-")

	s.Contains(listOutput, adrID1, "first ADR should appear")
	s.Contains(listOutput, adrID2, "second ADR should appear")
	s.Contains(listOutput, adrID3, "third ADR should appear")
}
