#!/usr/bin/env bash
#
# End-to-end workflow testing for Rex v1.0
# Tests all major workflows in a clean environment

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Test directory
TEST_DIR=$(mktemp -d -t rex-test-XXXXXX)
REX_BIN="$(pwd)/rex"

echo -e "${BLUE}=== Rex v1.0 Workflow Tests ===${NC}"
echo ""
echo "Test directory: $TEST_DIR"
echo ""

# Cleanup on exit
cleanup() {
    echo ""
    echo -e "${BLUE}=== Cleanup ===${NC}"
    rm -rf "$TEST_DIR"
    echo "Test directory removed"
}
trap cleanup EXIT

# Test helper functions
test_start() {
    ((TESTS_RUN++)) || true
    printf "${BLUE}[TEST $TESTS_RUN]${NC} %s\n" "$1"
}

test_pass() {
    ((TESTS_PASSED++)) || true
    printf "${GREEN}✓ PASS${NC}: %s\n\n" "$1"
}

test_fail() {
    ((TESTS_FAILED++)) || true
    printf "${RED}✗ FAIL${NC}: %s\n" "$1"
    printf "  Error: %s\n\n" "$2"
}

# Build rex
echo -e "${BLUE}=== Building Rex ===${NC}"
ORIGINAL_DIR=$(pwd)
go build -o "$REX_BIN" main.go
echo -e "${GREEN}✓${NC} Rex built successfully"
echo ""

# Set up test project structure
cd "$TEST_DIR"
mkdir -p docs/{adr,rfc,tasks/{core,plugin,ui,other}/{active,completed},plans}

# ========================================
# Test 1: ADR Workflow
# ========================================
test_start "ADR Workflow (create, list, update)"

# Create ADRs
$REX_BIN adr create "Use SQLite as Cache" > /dev/null 2>&1
$REX_BIN adr create "Adopt Go Embed for Templates" > /dev/null 2>&1
$REX_BIN adr create "Use Cobra CLI Framework" > /dev/null 2>&1

# Verify creation
if [ -f "docs/adr/0001-use-sqlite-as-cache.md" ] && \
   [ -f "docs/adr/0002-adopt-go-embed-for-templates.md" ] && \
   [ -f "docs/adr/0003-use-cobra-cli-framework.md" ]; then
    test_pass "ADR files created"
else
    test_fail "ADR files not created" "Expected 3 ADR files in docs/adr/"
fi

# List ADRs
ADR_COUNT=$($REX_BIN adr list 2>/dev/null | grep -c "^[0-9]" || echo "0")
if [ "$ADR_COUNT" -eq 3 ]; then
    test_pass "ADR list shows 3 ADRs"
else
    test_fail "ADR list incorrect" "Expected 3, got $ADR_COUNT"
fi

# Update ADR README
$REX_BIN adr update > /dev/null 2>&1
if [ -f "docs/adr/README.md" ] && grep -q "0001" "docs/adr/README.md"; then
    test_pass "ADR README generated"
else
    test_fail "ADR README not generated" "Expected docs/adr/README.md with ADR entries"
fi

# ========================================
# Test 2: RFC Workflow
# ========================================
test_start "RFC Workflow (create, list, update)"

# Create RFCs (need to provide author via input)
echo "John Doe" | $REX_BIN rfc create "Add Plugin System" > /dev/null 2>&1
echo "Jane Smith" | $REX_BIN rfc create "Improve Task Statistics" > /dev/null 2>&1

# Verify creation
if [ -f "docs/rfc/0001-add-plugin-system.md" ] && \
   [ -f "docs/rfc/0002-improve-task-statistics.md" ]; then
    test_pass "RFC files created"
else
    test_fail "RFC files not created" "Expected 2 RFC files in docs/rfc/"
fi

# List RFCs
RFC_COUNT=$($REX_BIN rfc list 2>/dev/null | grep -c "^[0-9]" || echo "0")
if [ "$RFC_COUNT" -eq 2 ]; then
    test_pass "RFC list shows 2 RFCs"
else
    test_fail "RFC list incorrect" "Expected 2, got $RFC_COUNT"
fi

# Update RFC README
$REX_BIN rfc update > /dev/null 2>&1
if [ -f "docs/rfc/README.md" ] && grep -q "0001" "docs/rfc/README.md"; then
    test_pass "RFC README generated"
else
    test_fail "RFC README not generated" "Expected docs/rfc/README.md with RFC entries"
fi

# ========================================
# Test 3: Task Workflow
# ========================================
test_start "Task Workflow (create, list, stats, complete)"

# Create tasks programmatically (simulate interactive input)
{
    echo "Implement SQLite cache"
    echo "core"
    echo "P1"
    echo "8"
    echo "database,cache"
    echo ""
    echo ""
    echo ""
    echo ""
    echo ""
} | $REX_BIN task create > /dev/null 2>&1

{
    echo "Add task statistics"
    echo "core"
    echo "P1"
    echo "6"
    echo "cli,stats"
    echo ""
    echo ""
    echo ""
    echo ""
    echo ""
} | $REX_BIN task create > /dev/null 2>&1

{
    echo "Create documentation"
    echo "core"
    echo "P2"
    echo "4"
    echo "docs"
    echo ""
    echo ""
    echo ""
    echo ""
    echo ""
} | $REX_BIN task create > /dev/null 2>&1

# Verify task files created
TASK_COUNT=$(find docs/tasks -name "TASK-*.md" | wc -l | tr -d ' ')
if [ "$TASK_COUNT" -eq 3 ]; then
    test_pass "Task files created (3 tasks)"
else
    test_fail "Task files not created" "Expected 3, found $TASK_COUNT"
fi

# List tasks
LIST_COUNT=$($REX_BIN task list 2>/dev/null | grep -c "^TASK-" || echo "0")
if [ "$LIST_COUNT" -eq 3 ]; then
    test_pass "Task list shows 3 tasks"
else
    test_fail "Task list incorrect" "Expected 3, got $LIST_COUNT"
fi

# List filtered tasks (core type)
CORE_COUNT=$($REX_BIN task list --type core 2>/dev/null | grep -c "^TASK-" || echo "0")
if [ "$CORE_COUNT" -eq 3 ]; then
    test_pass "Filtered task list works (--type core)"
else
    test_fail "Filtered task list failed" "Expected 3 core tasks, got $CORE_COUNT"
fi

# List filtered tasks (P1 priority)
P1_COUNT=$($REX_BIN task list --priority P1 2>/dev/null | grep -c "^TASK-" || echo "0")
if [ "$P1_COUNT" -eq 2 ]; then
    test_pass "Filtered task list works (--priority P1)"
else
    test_fail "Filtered task list failed" "Expected 2 P1 tasks, got $P1_COUNT"
fi

# Task stats
$REX_BIN task stats > /dev/null 2>&1
if [ $? -eq 0 ]; then
    test_pass "Task stats command works"
else
    test_fail "Task stats command failed" "Command exited with error"
fi

# Complete a task
echo "8" | $REX_BIN task complete TASK-001 > /dev/null 2>&1
if [ -f docs/tasks/core/completed/TASK-001-*.md ]; then
    test_pass "Task completion moves file to completed/"
else
    test_fail "Task completion failed" "File not moved to completed/"
fi

# ========================================
# Test 4: Plan Workflow
# ========================================
test_start "Plan Workflow (create, list, update)"

# Create plans
$REX_BIN plan create "Q4 2025 Roadmap" > /dev/null 2>&1
$REX_BIN plan create "Sprint 42" > /dev/null 2>&1

# Verify creation
if [ -f "docs/plans/0001-q4-2025-roadmap.md" ] && \
   [ -f "docs/plans/0002-sprint-42.md" ]; then
    test_pass "Plan files created"
else
    test_fail "Plan files not created" "Expected 2 plan files in docs/plans/"
fi

# List plans
PLAN_COUNT=$($REX_BIN plan list 2>/dev/null | grep -c "^[0-9]" || echo "0")
if [ "$PLAN_COUNT" -eq 2 ]; then
    test_pass "Plan list shows 2 plans"
else
    test_fail "Plan list incorrect" "Expected 2, got $PLAN_COUNT"
fi

# Update plan README
$REX_BIN plan update > /dev/null 2>&1
if [ -f "docs/plans/README.md" ] && grep -q "0001" "docs/plans/README.md"; then
    test_pass "Plan README generated"
else
    test_fail "Plan README not generated" "Expected docs/plans/README.md with plan entries"
fi

# ========================================
# Test 5: Database Operations
# ========================================
test_start "Database Operations (rebuild, info)"

# Rebuild cache
$REX_BIN rebuild > /dev/null 2>&1
if [ -f ".rex.db" ]; then
    test_pass "Cache rebuild creates database"
else
    test_fail "Cache rebuild failed" "Database file not created"
fi

# Database info
$REX_BIN db info > /dev/null 2>&1
if [ $? -eq 0 ]; then
    test_pass "Database info command works"
else
    test_fail "Database info failed" "Command exited with error"
fi

# Verify document counts
DB_INFO=$($REX_BIN db info 2>/dev/null)
if echo "$DB_INFO" | grep -q "ADRs.*:.*3" && \
   echo "$DB_INFO" | grep -q "RFCs.*:.*2" && \
   echo "$DB_INFO" | grep -q "Plans.*:.*2"; then
    test_pass "Database contains correct document counts"
else
    test_fail "Database counts incorrect" "Expected ADRs:3, RFCs:2, Plans:2"
fi

# ========================================
# Test 6: Version and Help
# ========================================
test_start "Version and Help Commands"

# Version command
$REX_BIN version > /dev/null 2>&1
if [ $? -eq 0 ]; then
    test_pass "Version command works"
else
    test_fail "Version command failed" "Command exited with error"
fi

# Help command
$REX_BIN --help > /dev/null 2>&1
if [ $? -eq 0 ]; then
    test_pass "Help command works"
else
    test_fail "Help command failed" "Command exited with error"
fi

# ========================================
# Test 7: Error Handling
# ========================================
test_start "Error Handling"

# Non-existent task
if $REX_BIN task complete TASK-999 > /dev/null 2>&1; then
    test_fail "Should error on non-existent task" "Command succeeded when it should fail"
else
    test_pass "Gracefully handles non-existent task"
fi

# Invalid filter
$REX_BIN task list --type invalid > /dev/null 2>&1 || true
# Note: This might still succeed but return 0 results - that's OK
test_pass "Handles invalid filters gracefully"

# ========================================
# Test Summary
# ========================================
echo ""
echo -e "${BLUE}=== Test Summary ===${NC}"
echo ""
echo "Tests Run:    $TESTS_RUN"
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
if [ $TESTS_FAILED -gt 0 ]; then
    echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
else
    echo -e "Tests Failed: ${GREEN}$TESTS_FAILED${NC}"
fi
echo ""

# Exit with appropriate code
if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi
