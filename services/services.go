package services

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"ogugu/database"
)

func SetupTestDB(t *testing.T, mdir string) (*sql.DB, func()) {
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

	db, err := database.New("pgx", dbstr)
	require.NoError(t, err)

	migrateDB(t, dbstr, mdir)

	return db, func() {
		err := dbContainer.Terminate(context.Background())
		require.NoError(t, err)
	}
}

func migrateDB(t *testing.T, dbconnstr, mdir string) {
	// magic file path, not good at all. will update
	// mdir is the migration directory
	m, err := migrate.New(mdir, dbconnstr)
	require.NoError(t, err)
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}
}
