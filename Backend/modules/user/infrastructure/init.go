package repository

import (
	"context"
	"fmt"
	"log"
	"waterordersystem/GOBACKEND/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func init() {

	dbSslEnabled := "disable"
	if utils.Env.DATABASE_SSL_ENABLED {
		dbSslEnabled = "enable"
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		utils.Env.DATABASE_USER,
		utils.Env.DATABASE_PASSWORD,
		utils.Env.DATABASE_HOST,
		utils.Env.DATABASE_PORT,
		utils.Env.DATABASE_NAME,
		dbSslEnabled,
	)

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to create database pool:", err)
	}

	Query = New(pool)
}

var Query *Queries
