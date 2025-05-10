package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ashish19912009/zrms/services/authN/internal/logger"
	"github.com/ashish19912009/zrms/services/authN/internal/mapper"
	"github.com/ashish19912009/zrms/services/authN/internal/model"
	"github.com/ashish19912009/zrms/services/authZ/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthZClient interface {
	CheckAccess(ctx context.Context, accountID, franchiseID, resource, action string) (*model.CheckAccessResponse, error)
	BatchCheckAccess(ctx context.Context, accountID, franchiseID string, resources []*model.ResourceAction) (*model.BatchCheckAccessResponse, error)
	Close() error
}

type authZClient struct {
	client pb.AuthZServiceClient
	conn   *grpc.ClientConn
	err    error
}

func NewAuthZServiceClient(host string, port string) (AuthZClient, error) {
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
		logger.Fatal("Failed to connect to AuthZ service: %v", err, nil)
	}

	client := pb.NewAuthZServiceClient(conn)
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return &authZClient{
		client: client,
		conn:   conn,
	}, nil
}

func (authzClient *authZClient) CheckAccess(ctx context.Context, accountID, franchiseID, resource, action string) (*model.CheckAccessResponse, error) {

	aM := mapper.CheckAccessFromModelToPb(accountID, franchiseID, resource, action)
	res, err := authzClient.client.CheckAccess(ctx, aM)
	if err != nil {
		logger.Error("CheckAccess failed: %v", err, map[string]interface{}{
			"layer":  "client",
			"method": "CheckAccess",
		})
		return nil, err
	}
	accessPb := mapper.CheckAccessFromPbToModel(res)
	return accessPb, nil
}

func (authzClient *authZClient) BatchCheckAccess(ctx context.Context, accountID, franchiseID string, resources []*model.ResourceAction) (*model.BatchCheckAccessResponse, error) {
	req := &model.BatchCheckAccess{
		AccountID:   accountID,
		FranchiseID: franchiseID,
		Resources:   resources,
	}
	aM, err := mapper.BatchCheckAccessFromModelToPb(req)
	if err != nil {
		logger.Error("Something went wrong while converting from model to PB: %v", err, map[string]interface{}{
			"layer":  "client",
			"method": "BatchCheckAccess",
		})
	}
	res, err := authzClient.client.BatchCheckAccess(ctx, aM)
	if err != nil {
		logger.Error("Batch CheckAccess failed: %v", err, map[string]interface{}{
			"layer":  "client",
			"method": "BatchCheckAccess",
		})
	}
	response, err := mapper.BatchCheckAccessFromPbToModel(res)
	return response, nil
}

func (authzClient *authZClient) Close() error {
	return authzClient.conn.Close()
}
