package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jseow5177/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)


type UserModel struct {
	DB *sql.DB
}

// Add a new user record to the users table
func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
		VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// If an error is returned, errors.As() is used to check whether the error has type 
		// *mysql.MySQLError. If it does, the error will be assigned to the mySQLError variable.
		// We can then check if the error relates to our users_uc_email key by checking the contents
		// of the message string. If it does, it returns an ErrDuplicateEmail error.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// Verify whether a user exists with the provided email address
// and password. This will return the relevant user ID if they exist.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashed_password []byte
	// Retrieve the id and hashed password associated with the given email.
	stmt := `SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE`

	// If no matching email exists, or the user is not active, we return ErrInvalidCredentials error.
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashed_password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Otherwise, the password is correct. Return the user ID.
	return 0, nil
}

// Fetch details of a specific user based on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}