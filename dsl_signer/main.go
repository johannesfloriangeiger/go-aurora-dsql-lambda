package dsl_signer

import (
	"context"
	"fmt"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jackc/pgx/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetConnection(ctx context.Context, username string, region string, clusterEndpoint string) (*pgx.Conn, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}

	var action string
	if username == "admin" {
		action = "DbConnectAdmin"
	} else {
		action = "DbConnect"
	}

	endpoint := "https://" + clusterEndpoint
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	values := req.URL.Query()
	values.Set("Action", action)
	req.URL.RawQuery = values.Encode()

	uri, _, err := v4.NewSigner().PresignHTTP(ctx, creds, req, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", "dsql", region, time.Now())
	if err != nil {
		return nil, err
	}

	var sb strings.Builder
	sb.WriteString("postgres://")
	sb.WriteString(clusterEndpoint)
	sb.WriteString(":5432/postgres?user=" + username + "&sslmode=verify-full")
	url := sb.String()
	connConfig, err := pgx.ParseConfig(url)
	connConfig.Password = uri[len("https://"):]
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		if err != nil {
			return nil, err
		}
	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)

	return conn, err
}
