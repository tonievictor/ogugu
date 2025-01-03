package models

import (
	"database/sql"
	"time"
)

var db *sql.DB

const dbtimeout = time.Second * 3

func New(d *sql.DB) {
	db = d
}
