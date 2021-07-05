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
		user := model.User{Name: "testuser"}
		// action
		id, err = AddUser(logger, db, user)
		// verify
		verify.Ok(t, err)
		verify.Assert(t, id > 0, "no id returned")
	})

	t.Run("getting user by id", func(t *testing.T) {
		// action
		user, err := GetUser(logger, db, id)
		// verify
		verify.Ok(t, err)
		verify.Equals(t, "testuser", user.Name)
	})

	t.Run("Getting all users", func(t *testing.T) {
		// action
		users, err := GetUsers(logger, db)
		// verify
		verify.Ok(t, err)
		verify.Equals(t, 1, len(users))
		verify.Equals(t, "testuser", users[0].Name)
	})

	t.Run("delete user by id", func(t *testing.T) {
		// action
		err := DeleteUser(logger, db, id)
		// verify
		verify.Ok(t, err)
	})

	t.Run("getting user by id, should return nil", func(t *testing.T) {
		// action
		_, err := GetUser(logger, db, id)
		// verify
		verify.Equals(t, pgx.ErrNoRows, err)
	})
}
