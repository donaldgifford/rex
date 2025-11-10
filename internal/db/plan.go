package db

import (
	"database/sql"
	"fmt"
)

// Plan represents a plan or roadmap document
type Plan struct {
	ID          int
	Number      int
	Title       string
	Status      string
	CreatedDate string
	FilePath    string
	Content     string
}

// CreatePlan inserts a new plan into the database
func (db *Database) CreatePlan(plan *Plan) error {
	result, err := db.conn.Exec(`
		INSERT INTO plans (number, title, status, created_date, file_path, content)
		VALUES (?, ?, ?, ?, ?, ?)
	`, plan.Number, plan.Title, plan.Status, plan.CreatedDate, plan.FilePath, plan.Content)

	if err != nil {
		return fmt.Errorf("insert plan: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert ID: %w", err)
	}

	plan.ID = int(id)
	return nil
}

// GetPlan retrieves a plan by number
func (db *Database) GetPlan(number int) (*Plan, error) {
	var plan Plan
	err := db.conn.QueryRow(`
		SELECT id, number, title, status, created_date, file_path, content
		FROM plans
		WHERE number = ?
	`, number).Scan(&plan.ID, &plan.Number, &plan.Title, &plan.Status, &plan.CreatedDate, &plan.FilePath, &plan.Content)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("plan %d not found", number)
	}

	if err != nil {
		return nil, fmt.Errorf("query plan: %w", err)
	}

	return &plan, nil
}

// ListPlans returns all plans, optionally filtered by status
func (db *Database) ListPlans(status string) ([]*Plan, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT id, number, title, status, created_date, file_path, content
			FROM plans
			WHERE status = ?
			ORDER BY number DESC
		`
		args = append(args, status)
	} else {
		query = `
			SELECT id, number, title, status, created_date, file_path, content
			FROM plans
			ORDER BY number DESC
		`
	}

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query plans: %w", err)
	}
	defer rows.Close()

	var plans []*Plan
	for rows.Next() {
		var plan Plan
		err := rows.Scan(&plan.ID, &plan.Number, &plan.Title, &plan.Status, &plan.CreatedDate, &plan.FilePath, &plan.Content)
		if err != nil {
			return nil, fmt.Errorf("scan plan: %w", err)
		}
		plans = append(plans, &plan)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate plans: %w", err)
	}

	return plans, nil
}

// UpdatePlan updates an existing plan
func (db *Database) UpdatePlan(plan *Plan) error {
	result, err := db.conn.Exec(`
		UPDATE plans
		SET title = ?, status = ?, created_date = ?, file_path = ?, content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE number = ?
	`, plan.Title, plan.Status, plan.CreatedDate, plan.FilePath, plan.Content, plan.Number)

	if err != nil {
		return fmt.Errorf("update plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("plan %d not found", plan.Number)
	}

	return nil
}

// DeletePlan deletes a plan by number
func (db *Database) DeletePlan(number int) error {
	result, err := db.conn.Exec("DELETE FROM plans WHERE number = ?", number)
	if err != nil {
		return fmt.Errorf("delete plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("plan %d not found", number)
	}

	return nil
}

// GetPlanStats returns statistics about plans
func (db *Database) GetPlanStats() (map[string]int, error) {
	rows, err := db.conn.Query(`
		SELECT status, COUNT(*) as count
		FROM plans
		GROUP BY status
	`)
	if err != nil {
		return nil, fmt.Errorf("query plan stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("scan stats: %w", err)
		}
		stats[status] = count
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stats: %w", err)
	}

	return stats, nil
}
