package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jackc/pgx/v5"
	"github.com/johannesfloriangeiger/go-aurora-dsql-lambda/dsql_signer"
	"log"
	"os"
)

type Handler struct {
	conn *pgx.Conn
}

func (handler *Handler) handleRequest(ctx context.Context) error {
	log.Println("Create schema...")
	_, err := handler.conn.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS people`)
	if err != nil {
		return err
	}
	log.Println("Created.")

	log.Println("Create table...")
	_, err = handler.conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS people.people (
			id  VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255)
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Created.")

	return nil
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	conn, err := dsql_signer.GetConnection(context.TODO(), "admin", cfg.Region, os.Getenv("CLUSTER_ENDPOINT"))
	if err != nil {
		panic(err)
	}

	handler := Handler{
		conn,
	}
	lambda.Start(handler.handleRequest)
}
