package main

import (
	"context"
	"fmt"
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
	ID string `json:"id"`
}

type Response struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (handler *Handler) handleRequest(ctx context.Context, request *Request) (*Response, error) {
	query, err := handler.conn.Query(ctx, `SELECT id, name FROM people.people WHERE id = $1`, request.ID)
	if err != nil {
		return nil, err
	}

	if query.Next() {
		values, err := query.Values()
		if err != nil {
			return nil, err
		}

		return &Response{
			ID:   values[0].(string),
			Name: values[1].(string),
		}, nil
	} else {
		return nil, fmt.Errorf("no entry found for ID: %s", request.ID)
	}
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
