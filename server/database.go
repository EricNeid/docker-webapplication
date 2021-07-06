package server

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/encoding/wkt"
)

const tablePositions = "positions"
const tableUsers = "users"

// start of position schema operations

func CreateTablePositions(logger *log.Logger, db *pgxpool.Pool) error {
	logger.Println("Creating table positions")
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s
			(
				id bigserial,
				position GEOGRAPHY(POINT, 4326) NOT NULL
			)`,
			tablePositions,
		),
	)
	return err
}

func AddPosition(logger *log.Logger, db *pgxpool.Pool, position Position) (int64, error) {
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			"INSERT INTO %s (position) VALUES (ST_GeomFromText('%s'))",
			tablePositions,
			wkt.MarshalString(position.Position),
		),
	)

	return 1, err
}

func DeletePosition(logger *log.Logger, db *pgxpool.Pool, id int64) error {
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			`DELETE FROM %s WHERE id=%d`,
			tablePositions,
			id,
		),
	)
	return err
}

// GetPosition returns the position that is ascoiated with the given id.
// If no position exists, pgx.ErrNoRows is returned.
func GetPosition(logger *log.Logger, db *pgxpool.Pool, id int64) (Position, error) {
	var position orb.Point
	err := db.QueryRow(
		context.Background(),
		fmt.Sprintf(
			`SELECT ST_AsBinary(position) FROM %s WHERE id=%d`,
			tablePositions,
			id,
		),
	).Scan(wkb.Scanner(&position))
	return Position{Position: position}, err
}

func GetPositions(logger *log.Logger, db *pgxpool.Pool) ([]Position, error) {
	var positions []Position
	// query all rows
	rows, err := db.Query(
		context.Background(),
		fmt.Sprintf(
			`SELECT ST_AsBinary(position) FROM %s`,
			tablePositions,
		),
	)
	if err != nil {
		return positions, err
	}
	defer rows.Close()

	// collect result
	for rows.Next() {
		var position orb.Point
		err = rows.Scan(wkb.Scanner(&position))
		if err != nil {
			return positions, err
		}
		positions = append(positions, Position{Position: position})
	}

	return positions, err
}

// start of user schema operations

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

func AddUser(logger *log.Logger, db *pgxpool.Pool, user User) (int64, error) {
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
func GetUser(logger *log.Logger, db *pgxpool.Pool, id int64) (User, error) {
	var name string
	err := db.QueryRow(
		context.Background(),
		fmt.Sprintf(
			`SELECT username FROM %s WHERE id=%d`,
			tableUsers,
			id,
		),
	).Scan(&name)
	return User{Name: name}, err
}

func GetUsers(logger *log.Logger, db *pgxpool.Pool) ([]User, error) {
	var users []User
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
		users = append(users, User{Name: name})
	}

	return users, err
}
