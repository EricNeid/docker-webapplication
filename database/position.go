package database

import (
	"context"
	"fmt"
	"log"

	"github.com/EricNeid/go-webserver/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/encoding/wkt"
)

const tablePositions = "positions"

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

func AddPosition(logger *log.Logger, db *pgxpool.Pool, position model.Position) (int64, error) {
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
func GetPosition(logger *log.Logger, db *pgxpool.Pool, id int64) (model.Position, error) {
	var position orb.Point
	err := db.QueryRow(
		context.Background(),
		fmt.Sprintf(
			`SELECT ST_AsBinary(position) FROM %s WHERE id=%d`,
			tablePositions,
			id,
		),
	).Scan(wkb.Scanner(&position))
	return model.Position{Position: position}, err
}

func GetPositions(logger *log.Logger, db *pgxpool.Pool) ([]model.Position, error) {
	var positions []model.Position
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
		positions = append(positions, model.Position{Position: position})
	}

	return positions, err
}
