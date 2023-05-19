package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (uint64, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 day' * $3) RETURNING id`
	// _, err := m.DB.Exec(stmt, title, content, expires)
	// if err != nil {
	// 	return 0, err
	// }
	var id uint64
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (m *SnippetModel) Get(id uint64) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > CURRENT_TIMESTAMP AND id = $1`

	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > CURRENT_TIMESTAMP ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

// type ExampleModel struct {
// 	DB *sql.DB
// 	InsertStmt *sql.Stmt
// }

// example of transactions
// func (m *ExampleModel) ExampleTransaction() error {
// 	tx, err := m.DB.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()
// 	_, err = tx.Exec("INSERT INTO ...")
// 	if err != nil {
// 		return err
// 	}
// 	_, err = tx.Exec("UPDATE ...")
// 	if err != nil {
// 		return err
// 	}
// 	err = tx.Commit()
// 	return err
// }


// example of pre-preparing statements
// func NewExampleModel(db *sql.DB) (*ExampleModel, error) {
// 	insertStmt, err := db.Prepare("INSERT INTO ...")
// 	if err != nil {
// 	return nil, err
// 	}
// 	return &ExampleModel{db, insertStmt}, nil
// 	}
// 	func (m *ExampleModel) Insert(args...) error {
// 	_, err := m.InsertStmt.Exec(args...)
// 	return err
// 	}
// 	func main() {
// 	db, err := sql.Open(...)
// 	if err != nil {
// 	errorLog.Fatal(err)
// 	}
// 	defer db.Close()
// 	exampleModel, err := NewExampleModel(db)
// 	if err != nil {
// 	errorLog.Fatal(err)
// 	}
// 	defer exampleModel.InsertStmt.Close()
// 	}