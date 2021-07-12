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

const tableVehicleState = "vehicle_state"

func CreateTablePositions(logger *log.Logger, db *pgxpool.Pool) error {
	logger.Printf("Creating table %s\n", tableVehicleState)
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s
			(
				id bigserial,
				position GEOGRAPHY(POINT, 4326) NOT NULL
			)`,
			tableVehicleState,
		),
	)
	return err
}

func addPosition(logger *log.Logger, db *pgxpool.Pool, state vehicleState) (int64, error) {
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			"INSERT INTO %s (position) VALUES (ST_GeomFromText('%s'))",
			tableVehicleState,
			wkt.MarshalString(state.Position),
		),
	)

	return 1, err
}

func deletePosition(logger *log.Logger, db *pgxpool.Pool, id int64) error {
	_, err := db.Exec(
		context.Background(),
		fmt.Sprintf(
			`DELETE FROM %s WHERE id=%d`,
			tableVehicleState,
			id,
		),
	)
	return err
}

// getPosition returns the position that is ascoiated with the given id.
// If no position exists, pgx.ErrNoRows is returned.
func getPosition(logger *log.Logger, db *pgxpool.Pool, id int64) (vehicleState, error) {
	var position orb.Point
	err := db.QueryRow(
		context.Background(),
		fmt.Sprintf(
			`SELECT ST_AsBinary(position) FROM %s WHERE id=%d`,
			tableVehicleState,
			id,
		),
	).Scan(wkb.Scanner(&position))
	return vehicleState{Position: position}, err
}

func getPositions(logger *log.Logger, db *pgxpool.Pool) ([]vehicleState, error) {
	var states []vehicleState
	// query all rows
	rows, err := db.Query(
		context.Background(),
		fmt.Sprintf(
			`SELECT ST_AsBinary(position) FROM %s`,
			tableVehicleState,
		),
	)
	if err != nil {
		return states, err
	}
	defer rows.Close()

	// collect result
	for rows.Next() {
		var position orb.Point
		err = rows.Scan(wkb.Scanner(&position))
		if err != nil {
			return states, err
		}
		states = append(states, vehicleState{Position: position})
	}

	return states, err
}
