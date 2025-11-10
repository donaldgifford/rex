package db

import (
	"database/sql"
	"fmt"
)

// Task represents a task with all metadata
type Task struct {
	ID             string
	Title          string
	Type           string
	Status         string
	Priority       string
	EstimatedHours float64
	ActualHours    *float64
	StartedDate    *string
	CompletedDate  *string
	Assignee       *string
	Phase          *int
	FilePath       string
	Content        string
	BlockedBy      []string
	Blocks         []string
	RelatedTo      []string
	Tags           []string
}

// CreateTask inserts a new task into the database
func (db *Database) CreateTask(task *Task) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert task
	_, err = tx.Exec(`
		INSERT INTO tasks (
			id, title, type, status, priority, estimated_hours, actual_hours,
			started_date, completed_date, assignee, phase, file_path, content
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.ID, task.Title, task.Type, task.Status, task.Priority,
		task.EstimatedHours, task.ActualHours, task.StartedDate, task.CompletedDate,
		task.Assignee, task.Phase, task.FilePath, task.Content,
	)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	// Insert relationships
	if err := db.insertTaskRelationships(tx, task); err != nil {
		return err
	}

	// Insert tags
	if err := db.insertTaskTags(tx, task.ID, task.Tags); err != nil {
		return err
	}

	return tx.Commit()
}

// insertTaskRelationships inserts task relationships
func (db *Database) insertTaskRelationships(tx *sql.Tx, task *Task) error {
	for _, blockedByID := range task.BlockedBy {
		_, err := tx.Exec(`
			INSERT INTO task_relationships (task_id, related_task_id, relationship_type)
			VALUES (?, ?, ?)
		`, task.ID, blockedByID, "blocked_by")
		if err != nil {
			return fmt.Errorf("insert blocked_by relationship: %w", err)
		}
	}

	for _, blocksID := range task.Blocks {
		_, err := tx.Exec(`
			INSERT INTO task_relationships (task_id, related_task_id, relationship_type)
			VALUES (?, ?, ?)
		`, task.ID, blocksID, "blocks")
		if err != nil {
			return fmt.Errorf("insert blocks relationship: %w", err)
		}
	}

	for _, relatedID := range task.RelatedTo {
		_, err := tx.Exec(`
			INSERT INTO task_relationships (task_id, related_task_id, relationship_type)
			VALUES (?, ?, ?)
		`, task.ID, relatedID, "related_to")
		if err != nil {
			return fmt.Errorf("insert related_to relationship: %w", err)
		}
	}

	return nil
}

// insertTaskTags inserts task tags
func (db *Database) insertTaskTags(tx *sql.Tx, taskID string, tags []string) error {
	for _, tag := range tags {
		_, err := tx.Exec(`
			INSERT INTO task_tags (task_id, tag)
			VALUES (?, ?)
		`, taskID, tag)
		if err != nil {
			return fmt.Errorf("insert tag: %w", err)
		}
	}
	return nil
}

// GetTask retrieves a task by ID
func (db *Database) GetTask(id string) (*Task, error) {
	var task Task
	err := db.conn.QueryRow(`
		SELECT id, title, type, status, priority, estimated_hours, actual_hours,
		       started_date, completed_date, assignee, phase, file_path, content
		FROM tasks
		WHERE id = ?
	`, id).Scan(
		&task.ID, &task.Title, &task.Type, &task.Status, &task.Priority,
		&task.EstimatedHours, &task.ActualHours, &task.StartedDate, &task.CompletedDate,
		&task.Assignee, &task.Phase, &task.FilePath, &task.Content,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task %s not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("query task: %w", err)
	}

	// Load relationships
	if err := db.loadTaskRelationships(&task); err != nil {
		return nil, err
	}

	// Load tags
	if err := db.loadTaskTags(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

// loadTaskRelationships loads all relationships for a task
func (db *Database) loadTaskRelationships(task *Task) error {
	rows, err := db.conn.Query(`
		SELECT related_task_id, relationship_type
		FROM task_relationships
		WHERE task_id = ?
	`, task.ID)
	if err != nil {
		return fmt.Errorf("query relationships: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var relatedID, relType string
		if err := rows.Scan(&relatedID, &relType); err != nil {
			return fmt.Errorf("scan relationship: %w", err)
		}

		switch relType {
		case "blocked_by":
			task.BlockedBy = append(task.BlockedBy, relatedID)
		case "blocks":
			task.Blocks = append(task.Blocks, relatedID)
		case "related_to":
			task.RelatedTo = append(task.RelatedTo, relatedID)
		}
	}

	return rows.Err()
}

// loadTaskTags loads all tags for a task
func (db *Database) loadTaskTags(task *Task) error {
	rows, err := db.conn.Query(`
		SELECT tag
		FROM task_tags
		WHERE task_id = ?
	`, task.ID)
	if err != nil {
		return fmt.Errorf("query tags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return fmt.Errorf("scan tag: %w", err)
		}
		task.Tags = append(task.Tags, tag)
	}

	return rows.Err()
}

// TaskListOptions holds filtering options for listing tasks
type TaskListOptions struct {
	Type     string
	Status   string
	Priority string
	Tags     []string
}

// ListTasks returns all tasks, optionally filtered
func (db *Database) ListTasks(opts TaskListOptions) ([]*Task, error) {
	query := `
		SELECT DISTINCT t.id, t.title, t.type, t.status, t.priority, t.estimated_hours,
		       t.actual_hours, t.started_date, t.completed_date, t.assignee, t.phase,
		       t.file_path, t.content
		FROM tasks t
	`

	var whereClauses []string
	var args []interface{}

	if opts.Type != "" {
		whereClauses = append(whereClauses, "t.type = ?")
		args = append(args, opts.Type)
	}

	if opts.Status != "" {
		whereClauses = append(whereClauses, "t.status = ?")
		args = append(args, opts.Status)
	}

	if opts.Priority != "" {
		whereClauses = append(whereClauses, "t.priority = ?")
		args = append(args, opts.Priority)
	}

	if len(opts.Tags) > 0 {
		query += " INNER JOIN task_tags tt ON t.id = tt.task_id"
		placeholders := "?"
		for i := 1; i < len(opts.Tags); i++ {
			placeholders += ", ?"
		}
		whereClauses = append(whereClauses, fmt.Sprintf("tt.tag IN (%s)", placeholders))
		for _, tag := range opts.Tags {
			args = append(args, tag)
		}
	}

	if len(whereClauses) > 0 {
		query += " WHERE "
		for i, clause := range whereClauses {
			if i > 0 {
				query += " AND "
			}
			query += clause
		}
	}

	query += " ORDER BY t.id DESC"

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID, &task.Title, &task.Type, &task.Status, &task.Priority,
			&task.EstimatedHours, &task.ActualHours, &task.StartedDate, &task.CompletedDate,
			&task.Assignee, &task.Phase, &task.FilePath, &task.Content,
		)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		// Load relationships
		if err := db.loadTaskRelationships(&task); err != nil {
			return nil, err
		}

		// Load tags
		if err := db.loadTaskTags(&task); err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}

// UpdateTask updates an existing task
func (db *Database) UpdateTask(task *Task) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update task
	result, err := tx.Exec(`
		UPDATE tasks
		SET title = ?, type = ?, status = ?, priority = ?, estimated_hours = ?,
		    actual_hours = ?, started_date = ?, completed_date = ?, assignee = ?,
		    phase = ?, file_path = ?, content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		task.Title, task.Type, task.Status, task.Priority, task.EstimatedHours,
		task.ActualHours, task.StartedDate, task.CompletedDate, task.Assignee,
		task.Phase, task.FilePath, task.Content, task.ID,
	)

	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task %s not found", task.ID)
	}

	// Delete and re-insert relationships
	if _, err := tx.Exec("DELETE FROM task_relationships WHERE task_id = ?", task.ID); err != nil {
		return fmt.Errorf("delete relationships: %w", err)
	}

	if err := db.insertTaskRelationships(tx, task); err != nil {
		return err
	}

	// Delete and re-insert tags
	if _, err := tx.Exec("DELETE FROM task_tags WHERE task_id = ?", task.ID); err != nil {
		return fmt.Errorf("delete tags: %w", err)
	}

	if err := db.insertTaskTags(tx, task.ID, task.Tags); err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteTask deletes a task by ID
func (db *Database) DeleteTask(id string) error {
	result, err := db.conn.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task %s not found", id)
	}

	return nil
}

// TaskStats holds statistics about tasks
type TaskStats struct {
	ByType     map[string]int
	ByStatus   map[string]int
	ByPriority map[string]int
	Total      int
	TotalHours struct {
		Estimated float64
		Actual    float64
	}
}

// GetTaskStats returns comprehensive statistics about tasks
func (db *Database) GetTaskStats() (*TaskStats, error) {
	stats := &TaskStats{
		ByType:     make(map[string]int),
		ByStatus:   make(map[string]int),
		ByPriority: make(map[string]int),
	}

	// Count by type
	rows, err := db.conn.Query(`
		SELECT type, COUNT(*) as count
		FROM tasks
		GROUP BY type
	`)
	if err != nil {
		return nil, fmt.Errorf("query type stats: %w", err)
	}
	for rows.Next() {
		var taskType string
		var count int
		if err := rows.Scan(&taskType, &count); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan type stats: %w", err)
		}
		stats.ByType[taskType] = count
		stats.Total += count
	}
	rows.Close()

	// Count by status
	rows, err = db.conn.Query(`
		SELECT status, COUNT(*) as count
		FROM tasks
		GROUP BY status
	`)
	if err != nil {
		return nil, fmt.Errorf("query status stats: %w", err)
	}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan status stats: %w", err)
		}
		stats.ByStatus[status] = count
	}
	rows.Close()

	// Count by priority
	rows, err = db.conn.Query(`
		SELECT priority, COUNT(*) as count
		FROM tasks
		GROUP BY priority
	`)
	if err != nil {
		return nil, fmt.Errorf("query priority stats: %w", err)
	}
	for rows.Next() {
		var priority string
		var count int
		if err := rows.Scan(&priority, &count); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan priority stats: %w", err)
		}
		stats.ByPriority[priority] = count
	}
	rows.Close()

	// Sum hours
	var estimatedHours, actualHours sql.NullFloat64
	err = db.conn.QueryRow(`
		SELECT
			COALESCE(SUM(estimated_hours), 0) as total_estimated,
			COALESCE(SUM(actual_hours), 0) as total_actual
		FROM tasks
	`).Scan(&estimatedHours, &actualHours)
	if err != nil {
		return nil, fmt.Errorf("query hours: %w", err)
	}

	if estimatedHours.Valid {
		stats.TotalHours.Estimated = estimatedHours.Float64
	}
	if actualHours.Valid {
		stats.TotalHours.Actual = actualHours.Float64
	}

	return stats, nil
}
