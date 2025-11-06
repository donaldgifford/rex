package generator

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	// Create tables
	schema := `
	CREATE TABLE IF NOT EXISTS adrs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		number INTEGER NOT NULL UNIQUE,
		title TEXT NOT NULL,
		status TEXT NOT NULL,
		date TEXT NOT NULL,
		file_path TEXT NOT NULL UNIQUE,
		content TEXT
	);

	CREATE TABLE IF NOT EXISTS rfcs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		number INTEGER NOT NULL UNIQUE,
		title TEXT NOT NULL,
		status TEXT NOT NULL,
		author TEXT,
		created_date TEXT NOT NULL,
		updated_date TEXT NOT NULL,
		file_path TEXT NOT NULL UNIQUE,
		content TEXT
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		type TEXT NOT NULL,
		status TEXT NOT NULL,
		priority TEXT NOT NULL,
		estimated_hours REAL,
		actual_hours REAL,
		started_date TEXT,
		completed_date TEXT,
		assignee TEXT,
		phase INTEGER,
		file_path TEXT NOT NULL UNIQUE,
		content TEXT
	);

	CREATE TABLE IF NOT EXISTS plans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		number INTEGER NOT NULL UNIQUE,
		title TEXT NOT NULL,
		status TEXT NOT NULL,
		created_date TEXT NOT NULL,
		file_path TEXT NOT NULL UNIQUE,
		content TEXT
	);
	`

	_, err = db.Exec(schema)
	require.NoError(t, err)

	return db
}

func TestGetNextADRNumber(t *testing.T) {
	tests := []struct {
		name        string
		setupData   []int
		wantNumber  int
	}{
		{
			name:        "empty database",
			setupData:   []int{},
			wantNumber:  1,
		},
		{
			name:        "one existing ADR",
			setupData:   []int{1},
			wantNumber:  2,
		},
		{
			name:        "multiple ADRs",
			setupData:   []int{1, 2, 3, 5},
			wantNumber:  6,
		},
		{
			name:        "non-sequential ADRs",
			setupData:   []int{1, 5, 10},
			wantNumber:  11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			// Insert setup data
			for _, num := range tt.setupData {
				_, err := db.Exec(
					"INSERT INTO adrs (number, title, status, date, file_path) VALUES (?, ?, ?, ?, ?)",
					num, "Test", "draft", "2025-11-05", "test-"+string(rune('0'+num))+".md",
				)
				require.NoError(t, err)
			}

			got, err := GetNextADRNumber(db)
			require.NoError(t, err)
			assert.Equal(t, tt.wantNumber, got)
		})
	}
}

func TestGetNextRFCNumber(t *testing.T) {
	tests := []struct {
		name        string
		setupData   []int
		wantNumber  int
	}{
		{
			name:        "empty database",
			setupData:   []int{},
			wantNumber:  1,
		},
		{
			name:        "one existing RFC",
			setupData:   []int{1},
			wantNumber:  2,
		},
		{
			name:        "multiple RFCs",
			setupData:   []int{1, 2, 3},
			wantNumber:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			// Insert setup data
			for _, num := range tt.setupData {
				_, err := db.Exec(
					"INSERT INTO rfcs (number, title, status, author, created_date, updated_date, file_path) VALUES (?, ?, ?, ?, ?, ?, ?)",
					num, "Test", "draft", "Test Author", "2025-11-01", "2025-11-05", "test-"+string(rune('0'+num))+".md",
				)
				require.NoError(t, err)
			}

			got, err := GetNextRFCNumber(db)
			require.NoError(t, err)
			assert.Equal(t, tt.wantNumber, got)
		})
	}
}

func TestGetNextTaskNumber(t *testing.T) {
	tests := []struct {
		name        string
		setupData   []string
		wantNumber  int
	}{
		{
			name:        "empty database",
			setupData:   []string{},
			wantNumber:  1,
		},
		{
			name:        "one existing task",
			setupData:   []string{"TASK-001"},
			wantNumber:  2,
		},
		{
			name:        "multiple tasks",
			setupData:   []string{"TASK-001", "TASK-002", "TASK-003"},
			wantNumber:  4,
		},
		{
			name:        "non-sequential tasks",
			setupData:   []string{"TASK-001", "TASK-005", "TASK-010"},
			wantNumber:  11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			// Insert setup data
			for i, id := range tt.setupData {
				_, err := db.Exec(
					"INSERT INTO tasks (id, title, type, status, priority, estimated_hours, file_path) VALUES (?, ?, ?, ?, ?, ?, ?)",
					id, "Test", "core", "planned", "P2", 4, "test-"+string(rune('0'+i))+".md",
				)
				require.NoError(t, err)
			}

			got, err := GetNextTaskNumber(db)
			require.NoError(t, err)
			assert.Equal(t, tt.wantNumber, got)
		})
	}
}

func TestGenerateTaskID(t *testing.T) {
	tests := []struct {
		name   string
		number int
		want   string
	}{
		{
			name:   "single digit",
			number: 1,
			want:   "TASK-001",
		},
		{
			name:   "double digit",
			number: 42,
			want:   "TASK-042",
		},
		{
			name:   "triple digit",
			number: 123,
			want:   "TASK-123",
		},
		{
			name:   "four digit",
			number: 1234,
			want:   "TASK-1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateTaskID(tt.number)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetNextPlanNumber(t *testing.T) {
	tests := []struct {
		name        string
		setupData   []int
		wantNumber  int
	}{
		{
			name:        "empty database",
			setupData:   []int{},
			wantNumber:  1,
		},
		{
			name:        "one existing plan",
			setupData:   []int{1},
			wantNumber:  2,
		},
		{
			name:        "multiple plans",
			setupData:   []int{1, 2, 3},
			wantNumber:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			// Insert setup data
			for _, num := range tt.setupData {
				_, err := db.Exec(
					"INSERT INTO plans (number, title, status, created_date, file_path) VALUES (?, ?, ?, ?, ?)",
					num, "Test", "draft", "2025-11-01", "test-"+string(rune('0'+num))+".md",
				)
				require.NoError(t, err)
			}

			got, err := GetNextPlanNumber(db)
			require.NoError(t, err)
			assert.Equal(t, tt.wantNumber, got)
		})
	}
}
