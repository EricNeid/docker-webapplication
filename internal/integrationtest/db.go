package integrationtest

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Setup() {
	stopAndRemoveTestDb()
	startTestDb()
}

func Cleanup() {
	stopAndRemoveTestDb()
}

func GetDbConnectionPool() (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/localdb")
}

func startTestDb() {
	fmt.Println(
		"docker run -d -p 5432:5432 --name postgis-test-db -e POSTGRES_USER=postgres -e POSTGRES_PASS=postgres -e POSTGRES_DBNAME=localdb kartoza/postgis:9.6-2.4",
	)
	cmd := exec.Command(
		"docker", "run",
		"-d",
		"-p", "5432:5432",
		"--name", "postgis-test-db",
		"-e", "POSTGRES_USER=postgres",
		"-e", "POSTGRES_PASS=postgres",
		"-e", "POSTGRES_DBNAME=localdb",
		"kartoza/postgis:9.6-2.4",
	)
	cmd.Run()
	time.Sleep(10 * time.Second)
}

func stopAndRemoveTestDb() {
	fmt.Println("docker stop postgis-test-db")
	cmd := exec.Command("docker", "stop", "postgis-test-db")
	cmd.Run()
	time.Sleep(5 * time.Second)

	fmt.Println("docker rm postgis-test-db")
	cmd = exec.Command("docker", "rm", "postgis-test-db")
	cmd.Run()
	time.Sleep(5 * time.Second)
}
