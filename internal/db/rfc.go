package db

import (
	"database/sql"
	"fmt"
)

// RFC represents an RFC (Request for Comments)
type RFC struct {
	ID          int
	Number      int
	Title       string
	Status      string
	Author      string
	CreatedDate string
	UpdatedDate string
	FilePath    string
	Content     string
}

// CreateRFC inserts a new RFC into the database
func (db *Database) CreateRFC(rfc *RFC) error {
	result, err := db.conn.Exec(`
		INSERT INTO rfcs (number, title, status, author, created_date, updated_date, file_path, content)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, rfc.Number, rfc.Title, rfc.Status, rfc.Author, rfc.CreatedDate, rfc.UpdatedDate, rfc.FilePath, rfc.Content)

	if err != nil {
		return fmt.Errorf("insert RFC: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert ID: %w", err)
	}

	rfc.ID = int(id)
	return nil
}

// GetRFC retrieves an RFC by number
func (db *Database) GetRFC(number int) (*RFC, error) {
	var rfc RFC
	err := db.conn.QueryRow(`
		SELECT id, number, title, status, author, created_date, updated_date, file_path, content
		FROM rfcs
		WHERE number = ?
	`, number).Scan(&rfc.ID, &rfc.Number, &rfc.Title, &rfc.Status, &rfc.Author, &rfc.CreatedDate, &rfc.UpdatedDate, &rfc.FilePath, &rfc.Content)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("RFC %d not found", number)
	}

	if err != nil {
		return nil, fmt.Errorf("query RFC: %w", err)
	}

	return &rfc, nil
}

// ListRFCs returns all RFCs, optionally filtered by status
func (db *Database) ListRFCs(status string) ([]*RFC, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT id, number, title, status, author, created_date, updated_date, file_path, content
			FROM rfcs
			WHERE status = ?
			ORDER BY number DESC
		`
		args = append(args, status)
	} else {
		query = `
			SELECT id, number, title, status, author, created_date, updated_date, file_path, content
			FROM rfcs
			ORDER BY number DESC
		`
	}

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query RFCs: %w", err)
	}
	defer rows.Close()

	var rfcs []*RFC
	for rows.Next() {
		var rfc RFC
		err := rows.Scan(&rfc.ID, &rfc.Number, &rfc.Title, &rfc.Status, &rfc.Author, &rfc.CreatedDate, &rfc.UpdatedDate, &rfc.FilePath, &rfc.Content)
		if err != nil {
			return nil, fmt.Errorf("scan RFC: %w", err)
		}
		rfcs = append(rfcs, &rfc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate RFCs: %w", err)
	}

	return rfcs, nil
}

// UpdateRFC updates an existing RFC
func (db *Database) UpdateRFC(rfc *RFC) error {
	result, err := db.conn.Exec(`
		UPDATE rfcs
		SET title = ?, status = ?, author = ?, created_date = ?, updated_date = ?, file_path = ?, content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE number = ?
	`, rfc.Title, rfc.Status, rfc.Author, rfc.CreatedDate, rfc.UpdatedDate, rfc.FilePath, rfc.Content, rfc.Number)

	if err != nil {
		return fmt.Errorf("update RFC: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("RFC %d not found", rfc.Number)
	}

	return nil
}

// DeleteRFC deletes an RFC by number
func (db *Database) DeleteRFC(number int) error {
	result, err := db.conn.Exec("DELETE FROM rfcs WHERE number = ?", number)
	if err != nil {
		return fmt.Errorf("delete RFC: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("RFC %d not found", number)
	}

	return nil
}

// GetRFCStats returns statistics about RFCs
func (db *Database) GetRFCStats() (map[string]int, error) {
	rows, err := db.conn.Query(`
		SELECT status, COUNT(*) as count
		FROM rfcs
		GROUP BY status
	`)
	if err != nil {
		return nil, fmt.Errorf("query RFC stats: %w", err)
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
