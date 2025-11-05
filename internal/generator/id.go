package generator

import (
	"database/sql"
	"fmt"
)

// GetNextADRNumber returns the next available ADR number
func GetNextADRNumber(db *sql.DB) (int, error) {
	var maxNumber sql.NullInt64
	err := db.QueryRow("SELECT MAX(number) FROM adrs").Scan(&maxNumber)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("query max ADR number: %w", err)
	}

	if !maxNumber.Valid {
		return 1, nil
	}

	return int(maxNumber.Int64) + 1, nil
}

// GetNextRFCNumber returns the next available RFC number
func GetNextRFCNumber(db *sql.DB) (int, error) {
	var maxNumber sql.NullInt64
	err := db.QueryRow("SELECT MAX(number) FROM rfcs").Scan(&maxNumber)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("query max RFC number: %w", err)
	}

	if !maxNumber.Valid {
		return 1, nil
	}

	return int(maxNumber.Int64) + 1, nil
}

// GetNextTaskNumber returns the next available task number
func GetNextTaskNumber(db *sql.DB) (int, error) {
	var maxID sql.NullString
	err := db.QueryRow("SELECT MAX(id) FROM tasks WHERE id LIKE 'TASK-%'").Scan(&maxID)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("query max task ID: %w", err)
	}

	if !maxID.Valid {
		return 1, nil
	}

	// Parse TASK-NNN format
	var number int
	_, err = fmt.Sscanf(maxID.String, "TASK-%d", &number)
	if err != nil {
		return 1, nil // If parsing fails, start from 1
	}

	return number + 1, nil
}

// GenerateTaskID generates a task ID in the format TASK-NNN
func GenerateTaskID(number int) string {
	return fmt.Sprintf("TASK-%03d", number)
}

// GetNextPlanNumber returns the next available plan number
func GetNextPlanNumber(db *sql.DB) (int, error) {
	var maxNumber sql.NullInt64
	err := db.QueryRow("SELECT MAX(number) FROM plans").Scan(&maxNumber)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("query max plan number: %w", err)
	}

	if !maxNumber.Valid {
		return 1, nil
	}

	return int(maxNumber.Int64) + 1, nil
}
