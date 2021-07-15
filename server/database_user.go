package server

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const tableUser = "application_user"

func createTableUsers(logger *log.Logger, db *pgxpool.Pool) error {
	logger.Printf("creating table %s\n", tableUser)
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s
			(
				id bigserial,
				username varchar NOT NULL
			)`,
			tableUser,
		),
	)
	return err
}

func addUser(logger *log.Logger, db *pgxpool.Pool, user user) (int64, error) {
	var id int64
	err := db.QueryRow(
		context.Background(),
		fmt.Sprintf(
			`INSERT INTO %s (username) VALUES ($1) RETURNING id`,
			tableUser,
		),
		user.Name,
	).Scan(&id)
	return id, err
}

func deleteUser(logger *log.Logger, db *pgxpool.Pool, id int64) error {
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			`DELETE FROM %s WHERE id=$1`,
			tableUser,
		),
		id,
	)
	return err
}

// getUser returns the user that is ascoiated with the given id.
// If no users exists, ErrorNotFound is returned.
func getUser(logger *log.Logger, db *pgxpool.Pool, id int64) (user, error) {
	var user user
	err := db.QueryRow(
		context.Background(),
		fmt.Sprintf(
			`SELECT username FROM %s WHERE id=$1`,
			tableUser,
		),
		id,
	).Scan(&user.Name)
	if err == pgx.ErrNoRows {
		err = ErrorNotFound // return custom error
	}
	return user, err
}

func getUsers(logger *log.Logger, db *pgxpool.Pool) ([]user, error) {
	var users []user
	// query all rows
	rows, err := db.Query(
		context.Background(),
		fmt.Sprintf(
			`SELECT username FROM %s`,
			tableUser,
		),
	)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	// collect result
	for rows.Next() {
		var user user
		err = rows.Scan(&user.Name)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, err
}
