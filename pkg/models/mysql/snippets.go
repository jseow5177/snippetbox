package mysql

import (
	"database/sql"
	"errors"

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

	// INSERT SQL statement.
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
	// SELECT SQL statement.
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Use QueryRow() on the connection pool to execute the SQL statement, 
	// passing in the id variables as the value of placeholder parameter. This
	// returns a pointer to a sql.Row object which holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new, empty Snippet struct
	s := new(models.Snippet)

	// Use rows.Scan() to copy the values from each field in sql.Row to the
  // corresponding field in the Snippet struct. Notice that the arguments
  // to row.Scan are *pointers* to the place you want to copy the data into,
  // and the number of arguments must be exactly the same as the number of
  // columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// Return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([] *models.Snippet, error) {
	// SELECT SQL statement.
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	// Use the Query() method on the connection pool to execute the SQL
	// statement. This returns a sql.Rows resultset containing the query result.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Defer rows.Close() to ensure the sql.Rows resultset is always
	// properly closed before the Latest() method returns.
	// If a resultset remains open, the underlying data connection will remain open.
	defer rows.Close()

	// Initialize an empty slice to hold the models.Snippets objects.
	snippets := []*models.Snippet{}
	
	// Use rows.Next to iterate through the rows in the resultset. This 
	// prepares the *first* (and then each subsequent) row to be acted on by the
	// rows.Scan() method.
	// If iteration over all the rows complete, the resultset automatically closes
	// and frees-up the underlying database connection.
	for rows.Next() {
		s := new(models.Snippet)

		// Use rows.Scan() to copy the values from each field in the row to the new Snippet object
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// Append to the list of snippets
		snippets = append(snippets, s)
	}

	// rows.Next() either returns true on success, or false if there is no next result row
	// or an error happened while preparing it.
	// To distinguish between two cases of false, we need to explicitly check for errors by calling
	// rows.Err(). This retrieves any errors that was encountered during the iteration.
	// It is important to call this, do not assume that a successful iteration was completed
	// over the whole resultset.
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return snippets, nil
}