package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/fastscripts/testing_in_go/data"
	"golang.org/x/crypto/bcrypt"
)

type SQLiteConn struct {
	DB *sql.DB
}

// AllUsers returns all users as a slice of *data.User
func (m *SQLiteConn) AllUsers() ([]*data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, is_admin, created_at, updated_at
	from users order by last_name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*data.User

	for rows.Next() {
		var user data.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, rows.Err()
}

// GetUser returns one user by id
func (m *SQLiteConn) GetUser(id int) (*data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, is_admin, created_at, updated_at
		from users where id = ?`

	var user data.User
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail returns one user by email address
func (m *SQLiteConn) GetUserByEmail(email string) (*data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, is_admin, created_at, updated_at
		from users where email = ?`

	var user data.User
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates one user in the database
func (m *SQLiteConn) UpdateUser(u data.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
		email = ?,
		first_name = ?,
		last_name = ?,
		is_admin = ?,
		updated_at = ?
		where id = ?`

	_, err := m.DB.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		u.IsAdmin,
		time.Now(),
		u.ID,
	)
	return err
}

// DeleteUser deletes one user from the database, by id
func (m *SQLiteConn) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, `delete from users where id = ?`, id)
	return err
}

// InsertUser inserts a new user into the database, and returns the ID of the newly inserted row
func (m *SQLiteConn) InsertUser(user data.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	result, err := m.DB.ExecContext(ctx, `insert into users
		(email, first_name, last_name, password, is_admin, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?, ?)`,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.IsAdmin,
		now,
		now,
	)
	if err != nil {
		return 0, err
	}

	newID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(newID), nil
}

// ResetPassword is the method we will use to change a user's password.
func (m *SQLiteConn) ResetPassword(id int, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = m.DB.ExecContext(ctx, `update users set password = ? where id = ?`, hashedPassword, id)
	return err
}

// InsertUserImage inserts a user profile image into the database.
func (m *SQLiteConn) InsertUserImage(i data.UserImage) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, `insert into user_images
		(user_id, file_name, created_at, updated_at) values (?, ?, ?, ?)`,
		i.UserID,
		i.FileName,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	newID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(newID), nil
}
