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
	Connection *pgx.Conn
}

type Request struct {
	Username string `json:"username"`
	RoleARN  string `json:"roleARN"`
}

func (handler *Handler) handleRequest(ctx context.Context, request *Request) error {
	log.Println("Creating database role...")
	_, err := handler.Connection.Exec(ctx, `CREATE ROLE `+request.Username+` WITH LOGIN`)
	if err != nil {
		return err
	}
	log.Println("Created.")

	log.Println("Grant IAM role access to database role...")
	_, err = handler.Connection.Exec(ctx, `AWS IAM GRANT `+request.Username+` TO '`+request.RoleARN+`'`)
	if err != nil {
		return err
	}
	log.Println("Granted.")

	log.Println("Granting schema usage to database role...")
	_, err = handler.Connection.Exec(ctx, `GRANT USAGE ON SCHEMA people TO `+request.Username)
	if err != nil {
		return err
	}
	log.Println("Granted.")

	log.Println("Granting table usage to database role...")
	_, err = handler.Connection.Exec(ctx, `GRANT SELECT, INSERT ON ALL TABLES IN SCHEMA people TO `+request.Username)
	if err != nil {
		return err
	}
	log.Println("Granted.")

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
		Connection: conn,
	}
	lambda.Start(handler.handleRequest)
}
