package server

import (
	"log"
	"os"
	"testing"

	"github.com/EricNeid/go-webserver/internal/integrationtest"
	"github.com/EricNeid/go-webserver/internal/verify"
	"github.com/jackc/pgx/v4"
)

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
		err = createTableUsers(logger, db)
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
		result, err := getUser(logger, db, id)
		// verify
		verify.Ok(t, err)
		verify.Equals(t, "testuser", result.Name)
	})

	t.Run("Getting all users", func(t *testing.T) {
		// action
		result, err := getUsers(logger, db)
		// verify
		verify.Ok(t, err)
		verify.Equals(t, 1, len(result))
		verify.Equals(t, "testuser", result[0].Name)
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
