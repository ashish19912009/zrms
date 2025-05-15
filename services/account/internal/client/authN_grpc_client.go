package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/logger"
	"github.com/ashish19912009/zrms/services/account/internal/mapper"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	pb_authn "github.com/ashish19912009/zrms/services/authN/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthNClient interface {
	VerifyToken(ctx context.Context, access_token model.Token) (*model.AuthClaims, error)
	Close() error
}

type authNClient struct {
	client pb_authn.AuthServiceClient
	conn   *grpc.ClientConn
	err    error
}

func NewAuthNServiceClient(host string, port string) (AuthNClient, error) {
	env := os.Getenv("ENV")
	if host == "" || port == "" {
		logger.Fatal("host and port must be set", nil, nil)
	}
	address := fmt.Sprintf("%s:%s", host, port)

	// Choose transport credentials based on env
	var opts []grpc.DialOption
	if env == "prod" {
		creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
		if err != nil {
			logger.Fatal("Failed to load TLS credentials: %v", err, nil)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		logger.Fatal("Failed to connect to AuthN service: %v", err, nil)
	}

	client := pb_authn.NewAuthServiceClient(conn)
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return &authNClient{
		client: client,
		conn:   conn,
	}, nil
}

func (authnClient *authNClient) VerifyToken(ctx context.Context, access_token model.Token) (*model.AuthClaims, error) {
	token := mapper.VerifyTokenFromModelToPb(access_token)
	authClaims, err := authnClient.client.VerifyToken(ctx, token)
	if err != nil {
		logger.Error("something went wrong while verifying token on authN services", err, nil)
		return nil, err
	}
	aClaims, err := mapper.VerifyTokenFromPbToModel(authClaims)
	if err != nil {
		logger.Error("something went wrong while converting from PB to Model", err, nil)
		return nil, err
	}
	return aClaims, nil
}

func (authnClient *authNClient) Close() error {
	return authnClient.conn.Close()
}
