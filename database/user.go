package database

import (
	"context"
	"fmt"
	"log"

	"github.com/EricNeid/go-webserver/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

const tableUsers = "users"

func CreateTableUsers(logger *log.Logger, db *pgxpool.Pool) error {
	logger.Println(fmt.Sprintf("CreateTableUsers %s", tableUsers))
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s
			(
				id bigserial,
				username varchar NOT NULL
			)`,
			tableUsers,
		),
	)
	return err
}

func AddUser(logger *log.Logger, db *pgxpool.Pool, user model.User) (int64, error) {
	var id int64
	err := db.QueryRow(
		context.Background(),
		fmt.Sprintf(
			`INSERT INTO %s (username) VALUES ('%s') RETURNING id`,
			tableUsers,
			user.Name,
		),
	).Scan(&id)
	return id, err
}

func DeleteUser(logger *log.Logger, db *pgxpool.Pool, id int64) error {
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			`DELETE FROM %s WHERE id=%d`,
			tableUsers,
			id,
		),
	)
	return err
}

// GetUser returns the user that is ascoiated with the given id.
// If no users exists, pgx.ErrNoRows is returned.
func GetUser(logger *log.Logger, db *pgxpool.Pool, id int64) (model.User, error) {
	var name string
	err := db.QueryRow(
		context.Background(),
		fmt.Sprintf(
			`SELECT username FROM %s WHERE id=%d`,
			tableUsers,
			id,
		),
	).Scan(&name)
	return model.User{Name: name}, err
}

func GetUsers(logger *log.Logger, db *pgxpool.Pool) ([]model.User, error) {
	var users []model.User
	// query all rows
	rows, err := db.Query(
		context.Background(),
		fmt.Sprintf(
			`SELECT username FROM %s`,
			tableUsers,
		),
	)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	// collect result
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return users, err
		}
		users = append(users, model.User{Name: name})
	}

	return users, err
}
