package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"waterordersystem/GOBACKEND/utils"

	_ "github.com/jackc/pgx/v5/stdlib"

	_ "waterordersystem/GOBACKEND/modules/logistic/application"
	_ "waterordersystem/GOBACKEND/modules/order/application"
	_ "waterordersystem/GOBACKEND/modules/organization/application"
	_ "waterordersystem/GOBACKEND/modules/payment/application"
	_ "waterordersystem/GOBACKEND/modules/user/application"
)

func init() {
	dbSslEnabled := "disable"
	if utils.Env.DATABASE_SSL_ENABLED {
		dbSslEnabled = "enable"
	}

	conn, err := sql.Open("pgx", fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		utils.Env.DATABASE_USER,
		utils.Env.DATABASE_PASSWORD,
		utils.Env.DATABASE_HOST,
		utils.Env.DATABASE_PORT,
		utils.Env.DATABASE_NAME,
		dbSslEnabled,
	))
	if err != nil {
		log.Fatalln("Failed to open db:", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		panic("database unreachable: " + err.Error())
	}

	fmt.Println("Database connection OK")

}

func main() {
	fmt.Printf("Server started on port %d\n", utils.Env.PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", utils.Env.PORT), utils.Mux)
	log.Fatalf("Server failed: %v", err)
}
