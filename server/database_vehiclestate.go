package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/geojson"
)

const tableVehicleState = "vehicle_state"

func createTableVehicleState(logger *log.Logger, db *pgxpool.Pool) error {
	logger.Printf("Creating table %s\n", tableVehicleState)
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s
			(
				id              bigserial,
				position        GEOGRAPHY(POINT, 4326) NOT NULL,
				state_timestamp TIMESTAMP
			)`,
			tableVehicleState,
		),
	)
	return err
}

func addVehicleState(logger *log.Logger, db *pgxpool.Pool, position orb.Point, timestamp time.Time) (int64, error) {
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			"INSERT INTO %s (position, state_timestamp) VALUES (ST_GeomFromWKB($1), $2)",
			tableVehicleState,
		),
		wkb.Value(position),
		timestamp,
	)

	return 1, err
}

func deleteVehicleState(logger *log.Logger, db *pgxpool.Pool, id int64) error {
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			`DELETE FROM %s WHERE id=$1`,
			tableVehicleState,
		),
		id,
	)
	return err
}

// getVehicleState returns the position that is ascoiated with the given id.
// If no position exists, pgx.ErrNoRows is returned.
func getVehicleState(logger *log.Logger, db *pgxpool.Pool, id int64) (vehicleState, error) {
	var position orb.Point
	var timestamp time.Time
	var err error

	err = db.QueryRow(
		context.Background(),
		fmt.Sprintf(
			`SELECT ST_AsBinary(position), state_timestamp FROM %s WHERE id=$1`,
			tableVehicleState,
		),
		id,
	).Scan(wkb.Scanner(&position), &timestamp)
	if err == pgx.ErrNoRows {
		err = ErrorNotFound // return custom error
	}
	return vehicleState{Position: *geojson.NewGeometry(position), Timestamp: timestamp}, err
}

func getVehicleStates(logger *log.Logger, db *pgxpool.Pool) ([]vehicleState, error) {
	var states []vehicleState

	var position orb.Point
	var timestamp time.Time
	var err error
	// query all rows
	rows, err := db.Query(
		context.Background(),
		fmt.Sprintf(
			`SELECT ST_AsBinary(position), state_timestamp FROM %s`,
			tableVehicleState,
		),
	)
	if err != nil {
		return states, err
	}
	defer rows.Close()

	// collect result
	for rows.Next() {
		err = rows.Scan(wkb.Scanner(&position), &timestamp)
		if err != nil {
			return states, err
		}
		states = append(states, vehicleState{Position: *geojson.NewGeometry(position), Timestamp: timestamp})
	}

	return states, err
}
