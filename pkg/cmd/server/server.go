package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	grcp "proto-parking-services/pkg/protocol/grpc"
	"proto-parking-services/pkg/protocol/rest"
	"proto-parking-services/pkg/services/api/v1/controllers"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// initiate a .env variabel
func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env found")
	}
}

// Running gRPC server and http gateway
func RunServer() error {
	ctx := context.Background()

	// Connect to Postgresql
	postUrl := os.Getenv("POSTGRESQL_URL")
	postUser := os.Getenv("POSTGRESQL_USER")
	postDb := os.Getenv("POSTGRESQL_DB")
	postPass := os.Getenv("POSTGRESQL_PASS")
	postPort := os.Getenv("POSTGRESQL_PORT")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=required", postUrl, postPort, postUser, postPass, postDb)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Database Connected")

	v1Api := controllers.NewUserServicesService(db)

	go func() {
		_ = rest.RunServer(ctx, "7980", os.Getenv("PORT"))
	}()

	return grcp.RunServer(ctx, v1Api, "7980")
}
