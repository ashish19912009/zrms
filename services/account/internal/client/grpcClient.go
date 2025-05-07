package client

type GRPC_AuthZ struct {
	client pb.AuthZServiceClient
}

func NewGRPC_AuthZ(cPath pb.AuthZServiceClient) (*GRPC_AuthZ,error) {
	authZHost := os.Getenv("AUTHZ_SERVICE_HOST")
	authZPort := os.Getenv("AUTHZ_SERVICE_PORT")
	
	env := os.Getenv("ENV")

	if authZHost == "" || authZPort == "" {
		logger.Fatal("AUTHZ_SERVICE_HOST and AUTHZ_SERVICE_PORT must be set",nil)
	}
	address := fmt.Sprintf("%s:%s",authZHost,authZPort)
	
	// Choose transport credentials based on env
	var opts []grpc.DialOption
    if env == "prod" {
        creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
        if err != nil {
            log.Fatalf("Failed to load TLS credentials: %v", err)
        }
        opts = append(opts, grpc.WithTransportCredentials(creds))
    } else {
        opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
    }
	// Dial the authz service
    conn, err := grpc.Dial(address, opts...)
    if err != nil {
        log.Fatalf("Failed to connect to AuthZ service: %v", err)
    }
    defer conn.Close()

	client := pb.NewAuthZServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

	return &GRPC_AuthZ{
		client: client
	},nil
}

func(authzClient *GRPC_AuthZ) CheckAccess(AccountID,franchiseID, resource, action string){
	res, err := authzClient.client.CheckAccess(ctx, &pb.CheckAccessRequest{
        AccountId:   "acc-123",
        FranchiseId: "fr-001",
        Resource:    "menu",
        Action:      "edit",
    })
    if err != nil {
        log.Fatalf("CheckAccess failed: %v", err)
    }
	return res
}