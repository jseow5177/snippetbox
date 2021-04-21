package mysql

import (
	"database/sql"

	"github.com/jseow5177/snippetbox/pkg/models"
)


type UserModel struct {
	DB *sql.DB
}

// Add a new user record to the users table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Verify whether a user exists with the provided email address
// and password. This will return the relevant user ID if they exist.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Fetch details of a specific user based on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}