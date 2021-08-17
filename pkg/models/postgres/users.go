package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"snippetbox.hachiko.app/pkg/models"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf("INSERT INTO users (name, email, hashed_password, created) VALUES ('%s', '%s', '%s', CURRENT_TIMESTAMP) RETURNING id;", name, email, string(hashedPassword))

	_, err = m.DB.Exec(stmt)
	if err != nil {
		var pqSQLError *pq.Error
		if errors.As(err, &pqSQLError) {
			if pqSQLError.Code == "23505" && strings.Contains(pqSQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := fmt.Sprintf("SELECT id, hashed_password FROM users WHERE email = '%s' AND active = TRUE", email)
	row := m.DB.QueryRow(stmt)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
