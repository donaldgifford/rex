package db

import (
	"database/sql"
	"fmt"
)

// ADR represents an Architecture Decision Record
type ADR struct {
	ID       int
	Number   int
	Title    string
	Status   string
	Date     string
	FilePath string
	Content  string
}

// CreateADR inserts a new ADR into the database
func (db *Database) CreateADR(adr *ADR) error {
	result, err := db.conn.Exec(`
		INSERT INTO adrs (number, title, status, date, file_path, content)
		VALUES (?, ?, ?, ?, ?, ?)
	`, adr.Number, adr.Title, adr.Status, adr.Date, adr.FilePath, adr.Content)

	if err != nil {
		return fmt.Errorf("insert ADR: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert ID: %w", err)
	}

	adr.ID = int(id)
	return nil
}

// GetADR retrieves an ADR by number
func (db *Database) GetADR(number int) (*ADR, error) {
	var adr ADR
	err := db.conn.QueryRow(`
		SELECT id, number, title, status, date, file_path, content
		FROM adrs
		WHERE number = ?
	`, number).Scan(&adr.ID, &adr.Number, &adr.Title, &adr.Status, &adr.Date, &adr.FilePath, &adr.Content)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ADR %d not found", number)
	}

	if err != nil {
		return nil, fmt.Errorf("query ADR: %w", err)
	}

	return &adr, nil
}

// ListADRs returns all ADRs, optionally filtered by status
func (db *Database) ListADRs(status string) ([]*ADR, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT id, number, title, status, date, file_path, content
			FROM adrs
			WHERE status = ?
			ORDER BY number DESC
		`
		args = append(args, status)
	} else {
		query = `
			SELECT id, number, title, status, date, file_path, content
			FROM adrs
			ORDER BY number DESC
		`
	}

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query ADRs: %w", err)
	}
	defer rows.Close()

	var adrs []*ADR
	for rows.Next() {
		var adr ADR
		err := rows.Scan(&adr.ID, &adr.Number, &adr.Title, &adr.Status, &adr.Date, &adr.FilePath, &adr.Content)
		if err != nil {
			return nil, fmt.Errorf("scan ADR: %w", err)
		}
		adrs = append(adrs, &adr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate ADRs: %w", err)
	}

	return adrs, nil
}

// UpdateADR updates an existing ADR
func (db *Database) UpdateADR(adr *ADR) error {
	result, err := db.conn.Exec(`
		UPDATE adrs
		SET title = ?, status = ?, date = ?, file_path = ?, content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE number = ?
	`, adr.Title, adr.Status, adr.Date, adr.FilePath, adr.Content, adr.Number)

	if err != nil {
		return fmt.Errorf("update ADR: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ADR %d not found", adr.Number)
	}

	return nil
}

// DeleteADR deletes an ADR by number
func (db *Database) DeleteADR(number int) error {
	result, err := db.conn.Exec("DELETE FROM adrs WHERE number = ?", number)
	if err != nil {
		return fmt.Errorf("delete ADR: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ADR %d not found", number)
	}

	return nil
}

// GetADRStats returns statistics about ADRs
func (db *Database) GetADRStats() (map[string]int, error) {
	rows, err := db.conn.Query(`
		SELECT status, COUNT(*) as count
		FROM adrs
		GROUP BY status
	`)
	if err != nil {
		return nil, fmt.Errorf("query ADR stats: %w", err)
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
