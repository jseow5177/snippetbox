package mysql

import (
	"database/sql"

	"github.com/jseow5177/snippetbox/pkg/models"
)

// Go methods for executing database queries:
// DB.Query() is used for SELECT queries which return multiple rows.
// DB.QueryRow() is used for SELECT queries which return a single row.
// DB.Exec() is used for statements which do not return rows (like INSERT and DELETE).

// Define a SnippetModel type which wraps a sql.DB connection pool.
// It implements methods to access and manipulate data related to Snippet in the database.
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	// Write the SQL statement we want to execute.
	// The ? character is used to indicate placeholder parameters.
	// Using parameterized queries instead of string interpolation can help
	// to prevent SQL injection.
	// Behind the scenes, DB.Exec() creates a prepared statement before passing in the parameters.
	// See https://en.wikipedia.org/wiki/Prepared_statement for more on prepare statements.
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use the Exec() method on the connection pool to execute the statement.
	// The first parameter is the SQL statement, followed by the title, content and expiry values for
	// the placeholder parameters.
	// This method returns a sql.Result object, which contains basic information about what happened 
	// when the query is executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result object to get the ID of our
	// newly inserted record in the snippets table.
	// Warning: Not all databases support this feature. E.g. PostgresSQL
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has a type of int64.
	// So we need to convert it back to int type.
	return int(id), nil
}

// Return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([] *models.Snippet, error) {
	return nil, nil
}