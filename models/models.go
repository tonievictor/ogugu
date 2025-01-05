package models

import (
	"database/sql"
	"time"

	"go.opentelemetry.io/otel"
)

var db *sql.DB

var tracer = otel.Tracer("models")

const dbtimeout = time.Second * 3

func New(d *sql.DB) {
	db = d
}
