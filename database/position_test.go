package database

import (
	"log"
	"os"
	"testing"

	"github.com/EricNeid/go-webserver/internal/integrationtest"
	"github.com/EricNeid/go-webserver/internal/verify"
	"github.com/EricNeid/go-webserver/model"
	"github.com/jackc/pgx/v4"
)

func TestPositionSchemaIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test")
	}
	// arrange
	integrationtest.Setup()
	defer integrationtest.Cleanup()
	logger := log.New(os.Stdout, "test: ", log.LstdFlags)
	db, _ := integrationtest.GetDbConnectionPool()

	var err error
	t.Run("creating table", func(t *testing.T) {
		// action
		err = CreateTablePositions(logger, db)
		// verify
		verify.Ok(t, err)
	})

	var id int64
	t.Run("adding position", func(t *testing.T) {
		// arrange
		position := model.Position{Position: [2]float64{20, 30}}
		// action
		id, err = AddPosition(logger, db, position)
		// verify
		verify.Ok(t, err)
		verify.Assert(t, id > 0, "no id returned")
	})

	t.Run("getting position by id", func(t *testing.T) {
		// action
		position, err := GetPosition(logger, db, id)
		// verify
		verify.Ok(t, err)
		verify.Assert(t, position.Position.X()-20.0 < 0.1, "Unexpected value returned")
		verify.Assert(t, position.Position.Y()-30.0 < 0.1, "Unexpected value returned")
	})

	t.Run("Getting all positions", func(t *testing.T) {
		// action
		positions, err := GetPositions(logger, db)
		// verify
		verify.Ok(t, err)
		verify.Equals(t, 1, len(positions))
	})

	t.Run("delete position by id", func(t *testing.T) {
		// action
		err := DeletePosition(logger, db, id)
		// verify
		verify.Ok(t, err)
	})

	t.Run("getting position by id, should return nil", func(t *testing.T) {
		// action
		_, err := GetPosition(logger, db, id)
		// verify
		verify.Equals(t, pgx.ErrNoRows, err)
	})
}
