package services

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"ogugu/database"
)

func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "ogugutest",
			"POSTGRES_PASSWORD": "ogugutest",
			"POSTGRES_USER":     "ogugutest",
		},
	}

	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	require.NoError(t, err)

	port, err := dbContainer.MappedPort(context.Background(), "5432")
	require.NoError(t, err)

	dbstr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", "ogugutest", "ogugutest",
		fmt.Sprintf("localhost:%s", port.Port()), "ogugutest")

	db, err := database.New("postgres", dbstr)
	require.NoError(t, err)

	migrateDB(t, dbstr)

	return db, func() {
		err := dbContainer.Terminate(context.Background())
		require.NoError(t, err)
	}
}

func migrateDB(t *testing.T, dbconnstr string) {
	// magic file path, not good at all. will update
	filepath := "file:///home/victor/dev/projects/ogugu/migrations"
	m, err := migrate.New(
		filepath,
		dbconnstr)

	require.NoError(t, err)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}
}
