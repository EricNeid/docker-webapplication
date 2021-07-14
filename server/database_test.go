package server

import (
	"log"
	"os"
	"testing"

	"github.com/EricNeid/go-webserver/internal/integrationtest"
	"github.com/EricNeid/go-webserver/internal/verify"
	"github.com/jackc/pgx/v4"
)

func TestVehicleStateSchemaIntegration(t *testing.T) {
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
	t.Run("add", func(t *testing.T) {
		// arrange
		position := vehicleState{Position: [2]float64{20, 30}, Timestamp: "2008-02-01T09:00:22+05"}
		// action
		id, err = addPosition(logger, db, position)
		// verify
		verify.Ok(t, err)
		verify.Assert(t, id > 0, "no id returned")
	})

	t.Run("get by id", func(t *testing.T) {
		// action
		position, err := getPosition(logger, db, id)
		// verify
		verify.Ok(t, err)
		verify.Assert(t, position.Position.X()-20.0 < 0.1, "Unexpected value returned")
		verify.Assert(t, position.Position.Y()-30.0 < 0.1, "Unexpected value returned")
	})

	t.Run("get all", func(t *testing.T) {
		// action
		positions, err := getPositions(logger, db)
		// verify
		verify.Ok(t, err)
		verify.Equals(t, 1, len(positions))
	})

	t.Run("delete by id", func(t *testing.T) {
		// action
		err := deletePosition(logger, db, id)
		// verify
		verify.Ok(t, err)
	})

	t.Run("get by id, should return nil", func(t *testing.T) {
		// action
		_, err := getPosition(logger, db, id)
		// verify
		verify.Equals(t, pgx.ErrNoRows, err)
	})
}

func TestUserSchemaIntegration(t *testing.T) {
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
		err = CreateTableUsers(logger, db)
		// verify
		verify.Ok(t, err)
	})

	var id int64
	t.Run("adding user", func(t *testing.T) {
		// arrange
		user := user{Name: "testuser"}
		// action
		id, err = addUser(logger, db, user)
		// verify
		verify.Ok(t, err)
		verify.Assert(t, id > 0, "no id returned")
	})

	t.Run("getting user by id", func(t *testing.T) {
		// action
		user, err := getUser(logger, db, id)
		// verify
		verify.Ok(t, err)
		verify.Equals(t, "testuser", user.Name)
	})

	t.Run("Getting all users", func(t *testing.T) {
		// action
		users, err := getUsers(logger, db)
		// verify
		verify.Ok(t, err)
		verify.Equals(t, 1, len(users))
		verify.Equals(t, "testuser", users[0].Name)
	})

	t.Run("delete user by id", func(t *testing.T) {
		// action
		err := deleteUser(logger, db, id)
		// verify
		verify.Ok(t, err)
	})

	t.Run("getting user by id, should return nil", func(t *testing.T) {
		// action
		_, err := getUser(logger, db, id)
		// verify
		verify.Equals(t, pgx.ErrNoRows, err)
	})
}
