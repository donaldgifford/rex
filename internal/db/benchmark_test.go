package db

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// setupBenchmarkDB creates a temporary database for benchmarking
func setupBenchmarkDB(b *testing.B) (*Database, string) {
	b.Helper()

	tmpDir, err := os.MkdirTemp("", "rex-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, ".rex.db")
	db, err := Open(dbPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		b.Fatalf("Failed to open database: %v", err)
	}

	return db, tmpDir
}

// cleanupBenchmarkDB closes database and removes temp directory
func cleanupBenchmarkDB(b *testing.B, db *Database, tmpDir string) {
	b.Helper()
	db.Close()
	os.RemoveAll(tmpDir)
}

// createTestADRs creates N ADRs for benchmarking
func createTestADRs(b *testing.B, db *Database, count int) {
	b.Helper()
	for i := 0; i < count; i++ {
		adr := &ADR{
			Number:   i + 1,
			Title:    fmt.Sprintf("Test ADR %d", i+1),
			Status:   "draft",
			Date:     "2025-11-10",
			FilePath: fmt.Sprintf("docs/adr/%04d-test-adr-%d.md", i+1, i+1),
			Content:  "Test content",
		}
		if err := db.CreateADR(adr); err != nil {
			b.Fatalf("Failed to create ADR: %v", err)
		}
	}
}

// createTestTasks creates N tasks for benchmarking
func createTestTasks(b *testing.B, db *Database, count int) {
	b.Helper()
	for i := 0; i < count; i++ {
		task := &Task{
			ID:             fmt.Sprintf("TASK-%03d", i+1),
			Title:          fmt.Sprintf("Test Task %d", i+1),
			Type:           "core",
			Status:         "planned",
			Priority:       "P2",
			EstimatedHours: 8.0,
			FilePath:       fmt.Sprintf("docs/tasks/core/active/TASK-%03d-test-task-%d.md", i+1, i+1),
			Content:        "Test content",
		}
		if err := db.CreateTask(task); err != nil {
			b.Fatalf("Failed to create task: %v", err)
		}
	}
}

// BenchmarkADRCreate benchmarks ADR creation
// Target: < 200ms
func BenchmarkADRCreate(b *testing.B) {
	db, tmpDir := setupBenchmarkDB(b)
	defer cleanupBenchmarkDB(b, db, tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		adr := &ADR{
			Number:   i + 1,
			Title:    fmt.Sprintf("Benchmark ADR %d", i+1),
			Status:   "draft",
			Date:     "2025-11-10",
			FilePath: fmt.Sprintf("docs/adr/%04d-benchmark-adr-%d.md", i+1, i+1),
			Content:  "Benchmark content",
		}
		if err := db.CreateADR(adr); err != nil {
			b.Fatalf("Failed to create ADR: %v", err)
		}
	}
}

// BenchmarkRFCCreate benchmarks RFC creation
// Target: < 200ms
func BenchmarkRFCCreate(b *testing.B) {
	db, tmpDir := setupBenchmarkDB(b)
	defer cleanupBenchmarkDB(b, db, tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rfc := &RFC{
			Number:      i + 1,
			Title:       fmt.Sprintf("Benchmark RFC %d", i+1),
			Status:      "draft",
			Author:      "Benchmark User",
			CreatedDate: "2025-11-10",
			UpdatedDate: "2025-11-10",
			FilePath:    fmt.Sprintf("docs/rfc/%04d-benchmark-rfc-%d.md", i+1, i+1),
			Content:     "Benchmark content",
		}
		if err := db.CreateRFC(rfc); err != nil {
			b.Fatalf("Failed to create RFC: %v", err)
		}
	}
}

// BenchmarkTaskCreate benchmarks task creation
// Target: < 300ms (includes interactive prompts in real usage)
func BenchmarkTaskCreate(b *testing.B) {
	db, tmpDir := setupBenchmarkDB(b)
	defer cleanupBenchmarkDB(b, db, tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := &Task{
			ID:             fmt.Sprintf("TASK-%03d", i+1),
			Title:          fmt.Sprintf("Benchmark Task %d", i+1),
			Type:           "core",
			Status:         "planned",
			Priority:       "P2",
			EstimatedHours: 8.0,
			FilePath:       fmt.Sprintf("docs/tasks/core/active/TASK-%03d-benchmark-task-%d.md", i+1, i+1),
			Content:        "Benchmark content",
		}
		if err := db.CreateTask(task); err != nil {
			b.Fatalf("Failed to create task: %v", err)
		}
	}
}

// BenchmarkListADRs100 benchmarks listing 100 ADRs
// Target: < 150ms for 100 items
func BenchmarkListADRs100(b *testing.B) {
	db, tmpDir := setupBenchmarkDB(b)
	defer cleanupBenchmarkDB(b, db, tmpDir)

	// Create 100 ADRs
	createTestADRs(b, db, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.ListADRs("")
		if err != nil {
			b.Fatalf("Failed to list ADRs: %v", err)
		}
	}
}

// BenchmarkListTasks100 benchmarks listing 100 tasks
// Target: < 150ms for 100 items
func BenchmarkListTasks100(b *testing.B) {
	db, tmpDir := setupBenchmarkDB(b)
	defer cleanupBenchmarkDB(b, db, tmpDir)

	// Create 100 tasks
	createTestTasks(b, db, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.ListTasks(TaskListOptions{})
		if err != nil {
			b.Fatalf("Failed to list tasks: %v", err)
		}
	}
}

// BenchmarkListTasksFiltered benchmarks filtered task listing
// Target: < 150ms for 100 items with filters
func BenchmarkListTasksFiltered(b *testing.B) {
	db, tmpDir := setupBenchmarkDB(b)
	defer cleanupBenchmarkDB(b, db, tmpDir)

	// Create 100 tasks
	createTestTasks(b, db, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.ListTasks(TaskListOptions{
			Type:     "core",
			Status:   "planned",
			Priority: "P2",
		})
		if err != nil {
			b.Fatalf("Failed to list tasks: %v", err)
		}
	}
}

// BenchmarkTaskStats100 benchmarks task statistics with 100 tasks
// Target: < 300ms for 100 items
func BenchmarkTaskStats100(b *testing.B) {
	db, tmpDir := setupBenchmarkDB(b)
	defer cleanupBenchmarkDB(b, db, tmpDir)

	// Create 100 tasks
	createTestTasks(b, db, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.GetTaskStats()
		if err != nil {
			b.Fatalf("Failed to get task stats: %v", err)
		}
	}
}

// BenchmarkCacheRebuild100 benchmarks cache rebuild with 100 files
// Target: < 3 seconds for 100 files
func BenchmarkCacheRebuild100(b *testing.B) {
	db, tmpDir := setupBenchmarkDB(b)
	defer cleanupBenchmarkDB(b, db, tmpDir)

	// Create test file structure
	docsDir := filepath.Join(tmpDir, "docs")
	adrDir := filepath.Join(docsDir, "adr")
	rfcDir := filepath.Join(docsDir, "rfc")
	tasksDir := filepath.Join(docsDir, "tasks", "core", "active")

	os.MkdirAll(adrDir, 0755)
	os.MkdirAll(rfcDir, 0755)
	os.MkdirAll(tasksDir, 0755)

	// Create 100 test markdown files (mix of ADRs, RFCs, Tasks)
	for i := 0; i < 40; i++ {
		content := fmt.Sprintf(`---
id: %04d
title: Test ADR %d
status: draft
date: 2025-11-10
---

# Test ADR %d

Test content.
`, i+1, i+1, i+1)
		path := filepath.Join(adrDir, fmt.Sprintf("%04d-test-adr-%d.md", i+1, i+1))
		os.WriteFile(path, []byte(content), 0644)
	}

	for i := 0; i < 30; i++ {
		content := fmt.Sprintf(`---
id: %04d
title: Test RFC %d
status: draft
author: Test User
date: 2025-11-10
---

# Test RFC %d

Test content.
`, i+1, i+1, i+1)
		path := filepath.Join(rfcDir, fmt.Sprintf("%04d-test-rfc-%d.md", i+1, i+1))
		os.WriteFile(path, []byte(content), 0644)
	}

	for i := 0; i < 30; i++ {
		content := fmt.Sprintf(`---
id: TASK-%03d
title: Test Task %d
type: core
status: planned
priority: P2
estimated_hours: 8
---

# Test Task %d

Test content.
`, i+1, i+1, i+1)
		path := filepath.Join(tasksDir, fmt.Sprintf("TASK-%03d-test-task-%d.md", i+1, i+1))
		os.WriteFile(path, []byte(content), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := db.Rebuild(docsDir); err != nil {
			b.Fatalf("Failed to rebuild cache: %v", err)
		}
	}
}

// TestPerformanceTargets is a test that measures actual performance and fails if targets are not met
func TestPerformanceTargets(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	t.Run("ADR Creation", func(t *testing.T) {
		db, tmpDir := setupBenchmarkDB(&testing.B{})
		defer cleanupBenchmarkDB(&testing.B{}, db, tmpDir)

		start := time.Now()
		adr := &ADR{
			Number:   1,
			Title:    "Performance Test ADR",
			Status:   "draft",
			Date:     "2025-11-10",
			FilePath: "docs/adr/0001-performance-test.md",
			Content:  "Test content",
		}
		err := db.CreateADR(adr)
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Failed to create ADR: %v", err)
		}

		target := 200 * time.Millisecond
		if duration > target {
			t.Errorf("ADR creation took %v, target is %v", duration, target)
		} else {
			t.Logf("✓ ADR creation: %v (target: %v)", duration, target)
		}
	})

	t.Run("Task Creation", func(t *testing.T) {
		db, tmpDir := setupBenchmarkDB(&testing.B{})
		defer cleanupBenchmarkDB(&testing.B{}, db, tmpDir)

		start := time.Now()
		task := &Task{
			ID:             "TASK-001",
			Title:          "Performance Test Task",
			Type:           "core",
			Status:         "planned",
			Priority:       "P2",
			EstimatedHours: 8.0,
			FilePath:       "docs/tasks/core/active/TASK-001-performance-test.md",
			Content:        "Test content",
		}
		err := db.CreateTask(task)
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}

		target := 300 * time.Millisecond
		if duration > target {
			t.Errorf("Task creation took %v, target is %v", duration, target)
		} else {
			t.Logf("✓ Task creation: %v (target: %v)", duration, target)
		}
	})

	t.Run("List 100 Tasks", func(t *testing.T) {
		db, tmpDir := setupBenchmarkDB(&testing.B{})
		defer cleanupBenchmarkDB(&testing.B{}, db, tmpDir)

		// Create 100 tasks
		for i := 0; i < 100; i++ {
			task := &Task{
				ID:             fmt.Sprintf("TASK-%03d", i+1),
				Title:          fmt.Sprintf("Test Task %d", i+1),
				Type:           "core",
				Status:         "planned",
				Priority:       "P2",
				EstimatedHours: 8.0,
				FilePath:       fmt.Sprintf("docs/tasks/core/active/TASK-%03d-test-%d.md", i+1, i+1),
				Content:        "Test content",
			}
			if err := db.CreateTask(task); err != nil {
				t.Fatalf("Failed to create task: %v", err)
			}
		}

		start := time.Now()
		tasks, err := db.ListTasks(TaskListOptions{})
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}
		if len(tasks) != 100 {
			t.Errorf("Expected 100 tasks, got %d", len(tasks))
		}

		target := 150 * time.Millisecond
		if duration > target {
			t.Errorf("Listing 100 tasks took %v, target is %v", duration, target)
		} else {
			t.Logf("✓ List 100 tasks: %v (target: %v)", duration, target)
		}
	})

	t.Run("Task Stats with 100 Tasks", func(t *testing.T) {
		db, tmpDir := setupBenchmarkDB(&testing.B{})
		defer cleanupBenchmarkDB(&testing.B{}, db, tmpDir)

		// Create 100 tasks
		for i := 0; i < 100; i++ {
			task := &Task{
				ID:             fmt.Sprintf("TASK-%03d", i+1),
				Title:          fmt.Sprintf("Test Task %d", i+1),
				Type:           "core",
				Status:         "planned",
				Priority:       "P2",
				EstimatedHours: 8.0,
				FilePath:       fmt.Sprintf("docs/tasks/core/active/TASK-%03d-test-%d.md", i+1, i+1),
				Content:        "Test content",
			}
			if err := db.CreateTask(task); err != nil {
				t.Fatalf("Failed to create task: %v", err)
			}
		}

		start := time.Now()
		stats, err := db.GetTaskStats()
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Failed to get task stats: %v", err)
		}
		if stats.Total != 100 {
			t.Errorf("Expected 100 total tasks, got %d", stats.Total)
		}

		target := 300 * time.Millisecond
		if duration > target {
			t.Errorf("Task stats took %v, target is %v", duration, target)
		} else {
			t.Logf("✓ Task stats (100 tasks): %v (target: %v)", duration, target)
		}
	})
}
