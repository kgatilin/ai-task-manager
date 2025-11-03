#!/bin/bash

echo "=== Testing Iteration Reordering ==="
echo ""

echo "Initial order (all should have rank 500 by default):"
./dw task-manager iteration list | head -15
echo ""

echo "Setting custom ranks for testing:"
./dw task-manager iteration update 10 --name "Phase 1: LLM Agent Integration (RANK 100)"
./dw task-manager iteration update 11 --name "Phase 1: AC Enforcement (RANK 200)"
./dw task-manager iteration update 12 --name "Backlog Management (RANK 300)"

echo ""
echo "After setting names (order should still be by number since ranks are same):"
./dw task-manager iteration list | grep -E "^(10|11|12) "
echo ""

echo "Now test in TUI:"
echo "1. Run: ./dw task-manager tui"
echo "2. Press Tab until you're in Iterations section"
echo "3. Select iteration #10 with j/k"
echo "4. Press Shift+J (capital J) to move down"
echo "5. Iteration #10 should swap positions with #11"
echo "6. Exit TUI and run: ./dw task-manager iteration list | grep -E '^(10|11|12) '"
echo "7. Order should be: 11, 10, 12 (because we swapped ranks)"
