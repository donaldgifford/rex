package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/donaldgifford/rex/internal/parser"
)

// Rebuild scans all markdown files and rebuilds the SQLite cache
func (db *Database) Rebuild(docsDir string) error {
	// Start a transaction for performance
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Clear existing data
	if err := db.clearTables(tx); err != nil {
		return fmt.Errorf("clear tables: %w", err)
	}

	// Rebuild ADRs
	adrDir := filepath.Join(docsDir, "adr")
	if err := db.rebuildADRs(tx, adrDir); err != nil {
		return fmt.Errorf("rebuild ADRs: %w", err)
	}

	// Rebuild RFCs
	rfcDir := filepath.Join(docsDir, "rfc")
	if err := db.rebuildRFCs(tx, rfcDir); err != nil {
		return fmt.Errorf("rebuild RFCs: %w", err)
	}

	// Rebuild Tasks
	tasksDir := filepath.Join(docsDir, "tasks")
	if err := db.rebuildTasks(tx, tasksDir); err != nil {
		return fmt.Errorf("rebuild tasks: %w", err)
	}

	// Rebuild Plans
	plansDir := filepath.Join(docsDir, "plans")
	if err := db.rebuildPlans(tx, plansDir); err != nil {
		return fmt.Errorf("rebuild plans: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	// Update last rebuild timestamp (after transaction commits)
	if err := db.SetLastRebuild(); err != nil {
		return fmt.Errorf("set last rebuild: %w", err)
	}

	return nil
}

// clearTables removes all existing data from cache tables
func (db *Database) clearTables(tx *sql.Tx) error {
	tables := []string{"task_tags", "task_relationships", "tasks", "adrs", "rfcs", "plans"}
	for _, table := range tables {
		if _, err := tx.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return fmt.Errorf("clear %s: %w", table, err)
		}
	}
	return nil
}

// rebuildADRs scans and imports all ADR markdown files
func (db *Database) rebuildADRs(tx *sql.Tx, adrDir string) error {
	if _, err := os.Stat(adrDir); os.IsNotExist(err) {
		return nil // Directory doesn't exist, skip
	}

	return filepath.Walk(adrDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		adr, body, err := parser.ParseADR(content)
		if err != nil {
			// Skip files that don't have valid frontmatter
			return nil
		}

		_, err = tx.Exec(`
			INSERT INTO adrs (number, title, status, date, file_path, content)
			VALUES (?, ?, ?, ?, ?, ?)
		`, adr.Number, adr.Title, adr.Status, adr.Date, path, string(body))
		if err != nil {
			return fmt.Errorf("insert ADR %s: %w", path, err)
		}

		return nil
	})
}

// rebuildRFCs scans and imports all RFC markdown files
func (db *Database) rebuildRFCs(tx *sql.Tx, rfcDir string) error {
	if _, err := os.Stat(rfcDir); os.IsNotExist(err) {
		return nil // Directory doesn't exist, skip
	}

	return filepath.Walk(rfcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		rfc, body, err := parser.ParseRFC(content)
		if err != nil {
			// Skip files that don't have valid frontmatter
			return nil
		}

		_, err = tx.Exec(`
			INSERT INTO rfcs (number, title, status, author, created_date, updated_date, file_path, content)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, rfc.Number, rfc.Title, rfc.Status, rfc.Author, rfc.CreatedDate, rfc.UpdatedDate, path, string(body))
		if err != nil {
			return fmt.Errorf("insert RFC %s: %w", path, err)
		}

		return nil
	})
}

// rebuildTasks scans and imports all Task markdown files
func (db *Database) rebuildTasks(tx *sql.Tx, tasksDir string) error {
	if _, err := os.Stat(tasksDir); os.IsNotExist(err) {
		return nil // Directory doesn't exist, skip
	}

	return filepath.Walk(tasksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		task, body, err := parser.ParseTask(content)
		if err != nil {
			// Skip files that don't have valid frontmatter
			return nil
		}

		// Insert task
		_, err = tx.Exec(`
			INSERT INTO tasks (
				id, title, type, status, priority, estimated_hours, actual_hours,
				started_date, completed_date, assignee, phase, file_path, content
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			task.ID, task.Title, task.Type, task.Status, task.Priority,
			task.EstimatedHours, task.ActualHours, task.StartedDate, task.CompletedDate,
			task.Assignee, task.Phase, path, string(body),
		)
		if err != nil {
			return fmt.Errorf("insert task %s: %w", path, err)
		}

		// Insert task relationships
		for _, blockedByID := range task.BlockedBy {
			_, err = tx.Exec(`
				INSERT INTO task_relationships (task_id, related_task_id, relationship_type)
				VALUES (?, ?, ?)
			`, task.ID, blockedByID, "blocked_by")
			if err != nil {
				return fmt.Errorf("insert blocked_by relationship: %w", err)
			}
		}

		for _, blocksID := range task.Blocks {
			_, err = tx.Exec(`
				INSERT INTO task_relationships (task_id, related_task_id, relationship_type)
				VALUES (?, ?, ?)
			`, task.ID, blocksID, "blocks")
			if err != nil {
				return fmt.Errorf("insert blocks relationship: %w", err)
			}
		}

		for _, relatedID := range task.RelatedTo {
			_, err = tx.Exec(`
				INSERT INTO task_relationships (task_id, related_task_id, relationship_type)
				VALUES (?, ?, ?)
			`, task.ID, relatedID, "related_to")
			if err != nil {
				return fmt.Errorf("insert related_to relationship: %w", err)
			}
		}

		// Insert task tags
		for _, tag := range task.Tags {
			_, err = tx.Exec(`
				INSERT INTO task_tags (task_id, tag)
				VALUES (?, ?)
			`, task.ID, tag)
			if err != nil {
				return fmt.Errorf("insert tag: %w", err)
			}
		}

		return nil
	})
}

// rebuildPlans scans and imports all Plan markdown files
func (db *Database) rebuildPlans(tx *sql.Tx, plansDir string) error {
	if _, err := os.Stat(plansDir); os.IsNotExist(err) {
		return nil // Directory doesn't exist, skip
	}

	return filepath.Walk(plansDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		plan, body, err := parser.ParsePlan(content)
		if err != nil {
			// Skip files that don't have valid frontmatter
			return nil
		}

		_, err = tx.Exec(`
			INSERT INTO plans (number, title, status, created_date, file_path, content)
			VALUES (?, ?, ?, ?, ?, ?)
		`, plan.Number, plan.Title, plan.Status, plan.CreatedDate, path, string(body))
		if err != nil {
			return fmt.Errorf("insert plan %s: %w", path, err)
		}

		return nil
	})
}

// NeedsRebuild checks if the cache needs to be rebuilt
// Returns true if:
// - Cache doesn't exist (.rex.db missing)
// - Any markdown file is newer than last rebuild timestamp
func (db *Database) NeedsRebuild(docsDir string) (bool, error) {
	lastRebuild, err := db.GetLastRebuild()
	if err != nil {
		return false, fmt.Errorf("get last rebuild: %w", err)
	}

	// If never rebuilt, rebuild
	if lastRebuild.IsZero() {
		return true, nil
	}

	// Check if any markdown file is newer than last rebuild
	needsRebuild := false
	err = filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			if info.ModTime().After(lastRebuild) {
				needsRebuild = true
				return filepath.SkipAll // Stop walking, we found a newer file
			}
		}

		return nil
	})

	if err != nil && err != filepath.SkipAll {
		return false, fmt.Errorf("walk docs: %w", err)
	}

	return needsRebuild, nil
}
