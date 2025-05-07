package integration

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"testing"

// 	"github.com/joho/godotenv"
// 	"github.com/stretchr/testify/assert"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"

// 	pb "github.com/ashish19912009/zrms/services/account/pb"
// )

// var grpcClient pb.AccountServiceClient

// func TestMain(m *testing.M) {
// 	if err := godotenv.Load("../../../config/env/test.env"); err != nil {
// 		log.Fatalf("Failed to load .env.test file: %v", err)
// 	}

// 	conn, err := grpc.NewClient(os.Getenv("GRPC_SERVER_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("failed to connect to gRPC server: %v", err)
// 	}
// 	defer conn.Close()

// 	grpcClient = pb.NewAccountServiceClient(conn)

// 	os.Exit(m.Run())
// }

// func TestCreateAccount_Integration(t *testing.T) {
// 	ctx := context.Background()
// 	resp, err := grpcClient.CreateAccount(ctx, &pb.CreateAccountRequest{
// 		Id:       "acc-integ-01",
// 		MobileNo: "9000000001",
// 		Name:     "Integration User",
// 		Role:     "admin",
// 		Status:   "active",
// 		EmpId:    "EMP9001",
// 	})
// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.Equal(t, "acc-integ-01", resp.Account.Id)
// }

// func TestGetAccountByID_Integration(t *testing.T) {
// 	ctx := context.Background()
// 	resp, err := grpcClient.GetAccountByID(ctx, &pb.GetAccountByIDRequest{Id: "acc-integ-01"})
// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.Equal(t, "acc-integ-01", resp.Account.Id)
// }

// func TestUpdateAccount_Integration(t *testing.T) {
// 	ctx := context.Background()
// 	resp, err := grpcClient.UpdateAccount(ctx, &pb.UpdateAccountRequest{
// 		Id:     "acc-integ-01",
// 		Name:   "Updated Integration User",
// 		Status: "inactive",
// 	})
// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.Equal(t, "Updated Integration User", resp.Account.Name)
// 	assert.Equal(t, "inactive", resp.Account.Status)
// }

// func TestListAccounts_Integration(t *testing.T) {
// 	ctx := context.Background()
// 	resp, err := grpcClient.GetAccounts(ctx, &pb.GetAccountsRequest{Skip: 0, Take: 5})
// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.GreaterOrEqual(t, len(resp.Accounts), 1)
// }

// func TestCreateAccount_InvalidInput(t *testing.T) {
// 	ctx := context.Background()
// 	_, err := grpcClient.CreateAccount(ctx, &pb.CreateAccountRequest{
// 		Id:       "",
// 		MobileNo: "",
// 	})
// 	assert.Error(t, err)
// }

// func TestUpdateAccount_InvalidInput(t *testing.T) {
// 	ctx := context.Background()
// 	_, err := grpcClient.UpdateAccount(ctx, &pb.UpdateAccountRequest{
// 		Id: "",
// 	})
// 	assert.Error(t, err)
// }

// func TestListAccounts_ZeroTake(t *testing.T) {
// 	ctx := context.Background()
// 	_, err := grpcClient.GetAccounts(ctx, &pb.GetAccountsRequest{Skip: 0, Take: 0})
// 	assert.Error(t, err)
// }

// var grpcClient pb.AccountServiceClient

// func TestMain(m *testing.M) {
// 	// _ = godotenv.Load("config/env/.env.test")

// 	// if err := db.SeedDB(); err != nil {
// 	// 	log.Fatalf("DB seeding failed: %v", err)
// 	// }

// 	// go startTestGRPCServer()
// 	// time.Sleep(time.Second)

// 	// conn, err := grpc.NewClient("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	// if err != nil {
// 	// 	log.Fatalf("failed to connect to test grpc server: %v", err)
// 	// }
// 	// defer conn.Close()

// 	// grpcClient = pb.NewAccountServiceClient(conn)

// 	// code := m.Run()
// 	// os.Exit(code)
// 	shutdown, err := startTestGRPCServerUnsafe()
// 	if err != nil {
// 		log.Fatalf("Failed to start test gRPC server: %v", err)
// 	}
// 	defer shutdown()

// 	code := m.Run()
// 	os.Exit(code)
// }

// var grpcAddr = "localhost:50055"

// func startTestGRPCServer(t *testing.T) func() {
// 	t.Helper()

// 	// Load env file
// 	err := godotenv.Load("../../config/env/.env.local")
// 	if err != nil {
// 		t.Fatalf("failed to load .env.local: %v", err)
// 	}

// 	dbURL := os.Getenv("DATABASE_URL")
// 	if dbURL == "" {
// 		t.Fatal("DATABASE_URL not set")
// 	}

// 	// Connect to the database
// 	db, err := sql.Open("postgres", dbURL)
// 	if err != nil {
// 		t.Fatalf("failed to connect to DB: %v", err)
// 	}

// 	repo := repository.NewRepository(db)
// 	svc := service.NewAccountService(repo)
// 	h := handler.NewGRPCHandler(svc)

// 	// Start gRPC server
// 	listener, err := net.Listen("tcp", grpcAddr)
// 	if err != nil {
// 		t.Fatalf("failed to listen: %v", err)
// 	}

// 	grpcServer := grpc.NewServer()
// 	pb.RegisterAccountServiceServer(grpcServer, h)

// 	go func() {
// 		if err := grpcServer.Serve(listener); err != nil {
// 			log.Fatalf("failed to serve: %v", err)
// 		}
// 	}()

// 	// Return shutdown function
// 	return func() {
// 		grpcServer.Stop()
// 		db.Close()
// 	}
// }

// // This is used in TestMain (no *testing.T)
// func startTestGRPCServerUnsafe() (shutdown func(), err error) {
// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to listen: %w", err)
// 	}

// 	grpcServer := grpc.NewServer()
// 	// register your services here...

// 	go func() {
// 		if err := grpcServer.Serve(lis); err != nil {
// 			log.Fatalf("gRPC server error: %v", err)
// 		}
// 	}()

// 	shutdown = func() {
// 		grpcServer.Stop()
// 	}
// 	return shutdown, nil
// }

// // This is used in actual tests
// func startTestGRPCServer(t *testing.T) func() {
// 	shutdown, err := startTestGRPCServerUnsafe()
// 	if err != nil {
// 		t.Fatalf("failed to start gRPC server: %v", err)
// 	}
// 	return shutdown
// }

// func TestCreateAccount_Success(t *testing.T) {
// 	t.Cleanup(func() {
// 		err := db.SeedDB()
// 		if err != nil {
// 			t.Fatalf("cleanup failed: %v", err)
// 		}
// 	})

// 	ctx := context.Background()
// 	req := &pb.CreateAccountRequest{
// 		Id:       "test-int-001",
// 		Name:     "Integration User",
// 		MobileNo: "9000000000",
// 		Role:     "admin",
// 		Status:   "active",
// 		EmpId:    "EMPINT001",
// 	}

// 	resp, err := grpcClient.CreateAccount(ctx, req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, req.Id, resp.Account.Id)
// 	assert.Equal(t, req.MobileNo, resp.Account.MobileNo)
// }

// func TestGetAccountByID_Success(t *testing.T) {
// 	t.Cleanup(func() {
// 		_ = db.SeedDB()
// 	})

// 	ctx := context.Background()

// 	createReq := &pb.CreateAccountRequest{
// 		Id:       "test-int-002",
// 		Name:     "Get User",
// 		MobileNo: "9111111111",
// 		Role:     "delivery",
// 		Status:   "active",
// 		EmpId:    "EMPINT002",
// 	}
// 	_, _ = grpcClient.CreateAccount(ctx, createReq)

// 	resp, err := grpcClient.GetAccount(ctx, &pb.GetAccountRequest{Id: "test-int-002"})
// 	assert.NoError(t, err)
// 	assert.Equal(t, "test-int-002", resp.Account.Id)
// 	assert.Equal(t, "Get User", resp.Account.Name)
// }

// func TestListAccounts_Pagination(t *testing.T) {
// 	t.Cleanup(func() {
// 		_ = db.SeedDB()
// 	})

// 	ctx := context.Background()
// 	resp, err := grpcClient.GetAccounts(ctx, &pb.GetAccountsRequest{Skip: 0, Take: 10})
// 	assert.NoError(t, err)
// 	assert.LessOrEqual(t, len(resp.Accounts), 10)
// }

// integration/account/account_test.go
// package account_test

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"testing"
// 	"time"

// 	pb "github.com/ashish19912009/zrms/services/account/pb"
// 	"github.com/stretchr/testify/assert"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// var client pb.AccountServiceClient

// func TestMain(m *testing.M) {
// 	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// conn, err := grpc.NewClient(ctx, "localhost:50051",
// 	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	// 	grpc.WithBlock(),
// 	// )
// 	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("failed to connect: %v", err)
// 	}

// 	client = pb.NewAccountServiceClient(conn)

// 	code := m.Run()

// 	conn.Close() // Clean up
// 	os.Exit(code)
// }

// func TestCreateAndGetAccount(t *testing.T) {
// 	ctx := context.Background()

// 	created, err := client.CreateAccount(ctx, &pb.CreateAccountRequest{
// 		MobileNo: "9999999999",
// 		Name:     "Integration User",
// 		Role:     "admin",
// 		Status:   "active",
// 		EmpId:    "EMP999",
// 	})
// 	assert.NoError(t, err)
// 	assert.NotNil(t, created.Account)

// 	got, err := client.GetAccount(ctx, &pb.GetAccountRequest{Id: created.Account.Id})
// 	assert.NoError(t, err)
// 	assert.Equal(t, created.Account.Id, got.Account.Id)
// }
