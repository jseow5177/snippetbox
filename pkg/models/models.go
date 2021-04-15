package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

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