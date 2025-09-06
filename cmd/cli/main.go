/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"ogugu/cmd/cli/cmd"
	"ogugu/internal/database"

	"github.com/tonievictor/dotenv"
)

func main() {
	dotenv.Config()

	db, err := database.New("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("unable to initialize database", err.Error())
		os.Exit(1)
	}

	if err := cmd.Execute(db); err != nil {
		os.Exit(1)
	}
}
