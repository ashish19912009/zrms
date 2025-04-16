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
	CheckDBConn           string
	CreateFranchise       string
	CreateOwner           string
	UpdateFranchise       string
	UpdateFranchiseStatus string
	DeleteFranchise       string
	GetFranchiseByID      string
	GetAllFranchises      string
	GetFranchiseOwner     string
	GetFranchiseDocuments string
	GetFranchiseAccounts  string
	UpdateAccount         string
	GetAccountByID        string
}{
	CheckDBConn:           "CheckDBConn",
	CreateFranchise:       "CreateFranchise",
	CreateOwner:           "CreateOwner",
	UpdateFranchise:       "UpdateFranchise",
	UpdateFranchiseStatus: "UpdateFranchiseStatus",
	DeleteFranchise:       "DeleteFranchise",
	GetFranchiseByID:      "GetFranchiseByID",
	GetAllFranchises:      "GetAllFranchises",
	GetFranchiseOwner:     "GetFranchiseOwner",
	GetFranchiseDocuments: "GetFranchiseDocuments",
	GetFranchiseAccounts:  "GetFranchiseAccounts",
	UpdateAccount:         "UpdateAccount",
	GetAccountByID:        "GetAccountByID",
}

const (
	Layer  = "layer"
	Method = "method"
)

const (
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
