package models

import (
	"errors"
	"time"
)

// Go handles the conversion of data types from SQL to native Go types.
// Source code: https://golang.org/src/database/sql/driver/types.go

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// Return this error if a user tries to login with an incorrect email or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// Return this error if a user tries to signup with an email address that is already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

// Database model of Snippet.
// Go doesn't do very well in managing NULL values in database records.
// Go cannot convert NULL into a string.
// To solve this, we can use sql.NullString instead (https://golang.org/pkg/database/sql/#NullString)
// But, as a rule, try to avoid NULL values altogether. Use either NOT NULL contraints along with DEFAULT
// values.
type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}


// Database model of User
type User struct {
	ID int
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
	Active bool
}