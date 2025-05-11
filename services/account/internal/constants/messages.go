package constants

// List of env variables
var EnvVariable = struct {
	IN_MEMORY_STORE_TYPE string
	ACCESS_TOKEN_TTL     string
	REFRESH_TOKEN_TTL    string
}{
	IN_MEMORY_STORE_TYPE: "type",
	ACCESS_TOKEN_TTL:     "ACCESS_TOKEN_TTL",
	REFRESH_TOKEN_TTL:    "REFRESH_TOKEN_TTL",
}

// List of Methods
var Methods = struct {
	CheckDBConn                  string
	CreateFranchise              string
	CreateOwner                  string
	UpdateFranchise              string
	UpdateFranchiseStatus        string
	DeleteFranchise              string
	GetFranchiseByID             string
	GetFranchiseByBusinessName   string
	GetAllFranchises             string
	GetFranchiseOwner            string
	AddFranchiseDocument         string
	UpdateFranchiseDocument      string
	GetAllFranchiseDocuments     string
	GetAllFranchiseAccounts      string
	CreateFranchiseAccount       string
	UpdateFranchiseAccount       string
	GetAccountByID               string
	GetFranchiseAddressByID      string
	AddFranchiseAddress          string
	UpdateFranchiseAddress       string
	AddFranchiseRole             string
	UpdateFranchiseRole          string
	GetAllFranchiseRoles         string
	AddPermissionsToRole         string
	UpdatePermissionsToRole      string
	GetAllPermissionsToRole      string
	GetFranchiseOwnerByID        string
	GetFranchiseAccountByID      string
	CheckIfOwnerExistsByAadharID string
	CreateNewOwner               string
	UpdateOwner                  string
}{
	CheckDBConn:                  "CheckDBConn",
	CreateFranchise:              "CreateFranchise",
	CreateOwner:                  "CreateOwner",
	UpdateFranchise:              "UpdateFranchise",
	UpdateFranchiseStatus:        "UpdateFranchiseStatus",
	DeleteFranchise:              "DeleteFranchise",
	GetFranchiseByID:             "GetFranchiseByID",
	GetFranchiseByBusinessName:   "GetFranchiseByBusinessName",
	GetAllFranchises:             "GetAllFranchises",
	GetFranchiseOwner:            "GetFranchiseOwner",
	AddFranchiseDocument:         "AddFranchiseDocument",
	UpdateFranchiseDocument:      "UpdateFranchiseDocument",
	GetAllFranchiseDocuments:     "GetAllFranchiseDocuments",
	GetAllFranchiseAccounts:      "GetAllFranchiseAccounts",
	CreateFranchiseAccount:       "CreateFranchiseAccount",
	UpdateFranchiseAccount:       "UpdateFranchiseAccount",
	GetAccountByID:               "GetAccountByID",
	GetFranchiseAddressByID:      "GetFranchiseAddressByID",
	AddFranchiseAddress:          "AddFranchiseAddress",
	UpdateFranchiseAddress:       "UpdateFranchiseAddress",
	AddFranchiseRole:             "AddFranchiseRole",
	UpdateFranchiseRole:          "UpdateFranchiseRole",
	GetAllFranchiseRoles:         "GetAllFranchiseRoles",
	AddPermissionsToRole:         "AddPermissionsToRole",
	UpdatePermissionsToRole:      "UpdatePermissionsToRole",
	GetAllPermissionsToRole:      "GetAllPermissionsToRole",
	GetFranchiseOwnerByID:        "GetFranchiseOwnerByID",
	GetFranchiseAccountByID:      "GetFranchiseAccountByID",
	CheckIfOwnerExistsByAadharID: "CheckIfOwnerExistsByAadharID",
	CreateNewOwner:               "CreateNewOwner",
	UpdateOwner:                  "UpdateOwner",
}

const (
	Layer  = "layer"
	Method = "method"
)

const (
	Handler    = "handler"
	Repository = "repository"
)

const (

	// Config error handling messages
	ConfigOverride          = "overriding config type with environment variable: %s"
	FailedToParse           = "failed to parse YAML config"
	UnsupportedDatabaseType = "unsupported database type"

	// Store error
	TypeSpecify                   = "%w: type must be specified"
	InvalidConfig                 = "invalid store configuration"
	FallbackLightning             = "falling back to LightningDB due to missing config"
	FallbackLightningDueToFailure = "falling back to LightningDB due to store initialization failure"

	// Logger Info
	LoginAttempt  = "login attempt"
	AccessToken   = "accessToken"
	RefreshToken  = "refreshToken"
	ValidateToken = "validate token"
	StoreToken    = "store token"
	CheckToken    = "check token"
	DeleteToken   = "delete token"
	LoadConfig    = "load config"

	// Validation Messages
	ValidationMissingCredentials = "username or password cannot be empty"
	ValidationInvalidEmail       = "invalid email format"
	ValidationMissingToken       = "authorization token is required"
	ValidationInvalidRole        = "user role is invalid"

	// Database Messages
	DBConnectionSuccess    = "connected to the database successfully"
	DBConnectionFailure    = "failed to connect to the database"
	DBConnectionNil        = "database connection is nil"
	DBQueryError           = "database query execution failed"
	DBRecordNotFound       = "record not found"
	ErrUnsupportedDatabase = "unsupported database"
	ErrKeyNotFound         = "key not found"
	ErrInvalidConfig       = "invalid config"
	BuildInsertQuery       = "something went wrong inside query insert builder function"
	BuildUpdateQuery       = "something went wrong inside query update builder function"
	BuildSelectQuery       = "something went wrong inside query select builder function"
	BuildDeleteQuery       = "something went wrong inside query delete builder function"

	// System & Server Messages
	SystemStartup  = "auth service is starting up..."
	SystemShutdown = "auth service is shutting down..."
	SystemError    = "unexpected system error occurred"

	// gRPC Messages
	GRPCRequestReceived = "received gRPC request: %s"
	GRPCResponseSent    = "sent gRPC response: %s"
	GRPCInternalError   = "internal gRPC error occurred"
	GRPCUnauthorized    = "gRPC request unauthorized"
	GRPCInvalidArgument = "invalid gRPC argument provided"

	// service Layer
	UserDataMissing = "user data missing"

	// Repository Layer
	MobileNoRequired            = "mobile no required"
	CredentialMissing           = "username or password missing"
	ErrUserNotFound             = "user not found"
	DBQueryFailed               = "db query failed"
	FailedToRetrv               = "failed to retrive data from last query performed"
	UserFetchedSuccessful       = "user data fetched successfully"
	ErrInvalidPassword          = "invalid password"
	WrongUsernamePassword       = "wrong username and password"
	PasswordMissingFromServer   = "password missing from database"
	NoColumProvided             = "no columns provided"
	UnauthorizedSchema          = "unauthorized schema: %s"
	UnauthorizedTable           = "unauthorized table: %s"
	UnauthorizedCloumn          = "unauthorized column: %s"
	UnauthorizedConditionColumn = "unauthorized condition column: %s"
	UnauthorizedReturningColumn = "unauthorized returning column: %s"
	UnauthorizedJoinTable       = "unauthorized join table"
	FailedToBeginTransaction    = "failed to begin transaction: %w"
	BusinessAlreadyExist        = "business registered with same name for the same franchise owner"
	FranchiseOwnerExist         = "a person is already registered with the same aadhar no"
	InvalidTimestamp            = "invalid timestamp format: %v"
	FailedToCreateOwner         = "failed to create franchise owner: %v"
	FailedToCreateFranchsie     = "failed to create franchise : %v"
	MappingFromProtoToModel     = "something went wrong while mapping proto to model"
	MappingFromModelToProto     = "something went wrong while mapping model to proto"
	SomethinWentWrongOnNew      = "somthing went wrong while creating: %s"
	SomethinWentWrongOnUpdate   = "somthing went wrong while updating: %s"
	MissingPermission           = "missing permissions: %v"

	// Token layer
	ErrInvalidToken            = "invalid token"
	ErrUnexpectedSigningMethod = "unexpected signing method"

	// in memory DB Type
	TypeKey       = ""
	RedisType     = "redis"
	MemcachedType = "memcached"
	DragonflyType = "dragonfly"
	BadgerType    = "badger"
	LightningType = "lightning"
	Access_token  = "access_token"
	Refresh_token = "refresh_token"

	InvalidRedisConfig     = "invalid redis config"
	InvalidMemcachedConfig = "invalid memcached config"
	InvalidDragonflyConfig = "invalid dragonfly config"
	InvalidBadgerConfig    = "invalid badger config"
)
