package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jackc/pgx/v5"
	"github.com/johannesfloriangeiger/go-aurora-dsql-lambda/dsql_signer"
	"os"
)

type Handler struct {
	conn *pgx.Conn
}

type Request struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (handler *Handler) handleRequest(ctx context.Context, request *Request) error {
	_, err := handler.conn.Exec(ctx, `INSERT INTO people.people (id, name) VALUES ($1, $2)`, request.ID, request.Name)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	conn, err := dsql_signer.GetConnection(context.TODO(), os.Getenv("USERNAME"), cfg.Region, os.Getenv("CLUSTER_ENDPOINT"))
	if err != nil {
		panic(err)
	}

	handler := Handler{
		conn,
	}
	lambda.Start(handler.handleRequest)
}
