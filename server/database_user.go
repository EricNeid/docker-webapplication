package server

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

const tableUser = "application_user"

func CreateTableUsers(logger *log.Logger, db *pgxpool.Pool) error {
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
			`INSERT INTO %s (username) VALUES ('%s') RETURNING id`,
			tableUser,
			user.Name,
		),
	).Scan(&id)
	return id, err
}

func deleteUser(logger *log.Logger, db *pgxpool.Pool, id int64) error {
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			`DELETE FROM %s WHERE id=%d`,
			tableUser,
			id,
		),
	)
	return err
}

// getUser returns the user that is ascoiated with the given id.
// If no users exists, pgx.ErrNoRows is returned.
func getUser(logger *log.Logger, db *pgxpool.Pool, id int64) (user, error) {
	var name string
	err := db.QueryRow(
		context.Background(),
		fmt.Sprintf(
			`SELECT username FROM %s WHERE id=%d`,
			tableUser,
			id,
		),
	).Scan(&name)
	return user{Name: name}, err
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
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return users, err
		}
		users = append(users, user{Name: name})
	}

	return users, err
}
