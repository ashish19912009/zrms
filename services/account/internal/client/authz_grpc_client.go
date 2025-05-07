package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/logger"
	"github.com/ashish19912009/zrms/services/account/internal/mapper"
	"github.com/ashish19912009/zrms/services/authZ/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthZClient struct {
	client pb.AuthZServiceClient
	conn   *grpc.ClientConn
	err    error
}

func NewAuthZServiceClient(host string, port int) *AuthZClient {
	env := os.Getenv("ENV")
	if host == "" || port == 0 {
		logger.Fatal("host and port must be set", nil, nil)
	}
	address := fmt.Sprintf("%s:%d", host, port)

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
	defer conn.Close()

	client := pb.NewAuthZServiceClient(conn)
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return &AuthZClient{
		client: client,
		conn:   conn,
		err:    err,
	}
}

func (authzClient *AuthZClient) CheckAccess(ctx context.Context, accountID, franchiseID, resource, action string) (*pb.CheckAccessResponse, error) {

	aM := mapper.CheckAccessFromModelToPb(accountID, franchiseID, resource, action)
	res, err := authzClient.client.CheckAccess(ctx, aM)
	if err != nil {
		logger.Fatal("CheckAccess failed: %v", err, nil)return

	}
	// &model.CheckAccessResponse{
	// 	Allowed:       allowed,
	// 	Reason:        reason,
	// 	IssuedAt:      issued_at,
	// 	ExpiresAt:     expires_at,
	// 	PolicyVersion: policy_version,
	// }
	accessPb, err := mapper.CheckAccessFromPbToModel(res)
	return accessPb, nil
	return res, nil
}

func (authzClient *AuthZClient) BatchCheckAccess(ctx context.Context, accountID, franchiseID string, resources []*pb.ResourceAction) (*pb.BatchCheckAccessResponse, error) {
	res, err := authzClient.client.BatchCheckAccess(ctx, &pb.BatchCheckAccessRequest{
		AccountId:   accountID,
		FranchiseId: franchiseID,
		Resources:   resources,
	})
	if err != nil {
		logger.Fatal("CheckAccess failed: %v", err, nil)
	}
	return res, nil
}
