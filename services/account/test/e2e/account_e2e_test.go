package e2e

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	config "github.com/ashish19912009/zrms/services/account/config/test_config"
	"github.com/ashish19912009/zrms/services/account/internal/handler"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/repository"
	"github.com/ashish19912009/zrms/services/account/internal/service"
	"github.com/ashish19912009/zrms/services/account/pb"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	cfg := config.LoadTestDBConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	// Migrate schema
	err = db.AutoMigrate(&model.Account{})
	if err != nil {
		panic("failed to migrate: " + err.Error())
	}

	// Clean slate before each test run
	db.Exec("TRUNCATE TABLE accounts RESTART IDENTITY CASCADE")

	return db
}

var client pb.AccountServiceClient

func TestMain(m *testing.M) {
	// Load test.env
	err := godotenv.Load("../../config/env/test.env")
	if err != nil {
		log.Fatalf("Error loading test.env file: %v", err)
	}
	// Setup
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db := setupTestDB()   // create a test DB (can be SQLite or in-memory Postgres)
	sqlDB, err := db.DB() // Extract *sql.DB
	if err != nil {
		panic("failed to get SQL DB from GORM: " + err.Error())
	}
	repo := repository.NewRepository(sqlDB)
	svc := service.NewAccountService(repo)
	grpcHandler := handler.NewGRPCHandler(svc)

	server := grpc.NewServer()
	fmt.Printf("GRPCHandler initialized: %+v\n", grpcHandler)
	pb.RegisterAccountServiceServer(server, grpcHandler)

	go server.Serve(lis)
	defer server.Stop()

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client = pb.NewAccountServiceClient(conn)

	code := m.Run()
	os.Exit(code)
}

func TestCreateAndFetchAccount_E2E(t *testing.T) {
	ctx := context.Background()

	req := &pb.CreateAccountRequest{
		Id:       uuid.New().String(),
		MobileNo: "9999999991",
		Name:     "E2E Tests",
		Role:     "admin",
		Status:   "active",
		EmpId:    "EMP-E2E-002",
	}

	//set authorization token
	md := metadata.New(map[string]string{"authorization": "Bearer valid-token"})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	// Create Account
	createResp, err := client.CreateAccount(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Id, createResp.Account.Id)

	// Fetch Account
	getResp, err := client.GetAccountByID(ctx, &pb.GetAccountByIDRequest{Id: req.Id})
	assert.NoError(t, err)
	assert.Equal(t, "E2E Tests", getResp.Account.Name)
}

func TestUpdateAccount_E2E(t *testing.T) {
	ctx := context.Background()

	// Step 1: Create an account first
	createResp, err := client.CreateAccount(ctx, &pb.CreateAccountRequest{
		Id:       uuid.New().String(),
		MobileNo: "9999999995",
		Name:     "E2E Tests",
		Role:     "admin",
		Status:   "active",
		EmpId:    "EMP-E2E-002",
	})
	require.NoError(t, err)
	require.NotNil(t, createResp)

	//set authorization token
	md := metadata.New(map[string]string{"authorization": "Bearer valid-token"})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	// Step 2: Update name and role
	updateResp, err := client.UpdateAccount(ctx, &pb.UpdateAccountRequest{
		Id:   createResp.Account.Id,
		Name: "John Updated",
		Role: "superadmin",
	})
	require.NoError(t, err)
	require.Equal(t, "John Updated", updateResp.Account.Name)
	require.Equal(t, "superadmin", updateResp.Account.Role)
}

func TestGetAccountByID_E2E(t *testing.T) {
	ctx := context.Background()

	//set authorization token
	md := metadata.New(map[string]string{"authorization": "Bearer valid-token"})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	// Create
	createResp, err := client.CreateAccount(ctx, &pb.CreateAccountRequest{
		Id:       uuid.New().String(),
		MobileNo: "7777777777",
		EmpId:    "E103",
		Name:     "Fetch Me",
		Role:     "admin",
		Status:   "active",
	})
	require.NoError(t, err)

	// Fetch by ID
	getResp, err := client.GetAccountByID(ctx, &pb.GetAccountByIDRequest{
		Id: createResp.Account.Id,
	})
	require.NoError(t, err)
	require.Equal(t, "Fetch Me", getResp.Account.Name)
}

func TestListAccounts_E2E(t *testing.T) {
	ctx := context.Background()

	//set authorization token
	md := metadata.New(map[string]string{"authorization": "Bearer valid-token"})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	// Add 3 accounts (feel free to use a loop)
	for i := 1; i <= 3; i++ {
		_, err := client.CreateAccount(ctx, &pb.CreateAccountRequest{
			Id:       uuid.New().String(),
			MobileNo: fmt.Sprintf("900000000%d", i),
			EmpId:    fmt.Sprintf("EMP_%d", i),
			Name:     fmt.Sprintf("User_%d", i),
			Role:     "manager",
			Status:   "active",
		})
		require.NoError(t, err)
	}

	// List with pagination: take 2, skip 0
	listResp, err := client.GetAccounts(ctx, &pb.GetAccountsRequest{
		Skip: 0,
		Take: 2,
	})
	require.NoError(t, err)
	require.Len(t, listResp.Accounts, 2)
}
