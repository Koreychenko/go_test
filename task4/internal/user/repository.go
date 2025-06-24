package user

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) Get(ctx context.Context, userId int) (*User, error) {
	var user User

	rows, err := r.db.QueryContext(ctx, "SELECT id, first_name, last_name FROM user WHERE id = ?", userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &user, nil
}

func (r *Repository) Delete(ctx context.Context, userId int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM user WHERE id = ?", userId)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	var err error

	result, err := r.db.ExecContext(ctx, "INSERT INTO user (first_name, last_name) VALUES (?, ?)", user.FirstName, user.LastName)

	if err != nil {
		return err
	}

	lastInsertId, err := result.LastInsertId()

	if err != nil {
		return err
	}

	user.ID = &lastInsertId

	return nil
}

func (r *Repository) Update(ctx context.Context, user *User) error {
	var err error

	_, err = r.db.ExecContext(ctx, "UPDATE user SET first_name = ?, last_name = ?) WHERE id = ?", user.FirstName, user.LastName, user.ID)

	if err != nil {
		return err
	}

	return nil
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}
