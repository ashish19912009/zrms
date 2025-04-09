package constants

// List of env variables
var EnvVariable = struct {
	IN_MEMORY_STORE_TYPE string
	ACCESS_TOKEN_TTL     string
	REFRESH_TOKEN_TTL    string
}{
	IN_MEMORY_STORE_TYPE: "IN_MEMORY_STORE_TYPE",
	ACCESS_TOKEN_TTL:     "ACCESS_TOKEN_TTL",
	REFRESH_TOKEN_TTL:    "REFRESH_TOKEN_TTL",
}

// List of Methods
var Methods = struct {
	Login                string
	Logout               string
	AccessToken          string
	RefreshToken         string
	ValidateToken        string
	StoreToken           string
	CheckToken           string
	DeleteToken          string
	GetUserByMobile      string
	VerifyPassword       string
	GenerateAccToken     string
	GenerateRefreshToken string
	LoadConfig           string
	Validate             string
	GetUser              string
}{
	Login:                "Login",
	Logout:               "Logout",
	Validate:             "Validate",
	AccessToken:          "AccessToken",
	RefreshToken:         "RefreshToken",
	ValidateToken:        "ValidateToken",
	StoreToken:           "StoreToken",
	CheckToken:           "CheckToken",
	DeleteToken:          "DeleteToken",
	GetUserByMobile:      "GetUserByMobile",
	VerifyPassword:       "VerifyPassword",
	GenerateAccToken:     "GenerateAccessToken",
	GenerateRefreshToken: "GenerateRefreshToken",
	LoadConfig:           "LoadConfig",
	GetUser:              "GetUser",
}

const (
	// Authentication Messages
	AuthLoginSuccess          = "user successfully logged in"
	AuthLoginFailure          = "invalid credentials provided"
	AuthTokenIssued           = "jwt token issued successfully"
	AuthTokenInvalid          = "invalid or expired token"
	AuthTokenVerified         = "token verified successfully"
	AuthTokenVeriFailed       = "token verification failed"
	AuthAccessDenied          = "access denied: unauthorized user"
	AuthRefreshSuccess        = "access token refreshed successfully"
	AuthRefreshFailure        = "failed to refresh access token"
	AuthRefreshRequired       = "refresh token required"
	AuthAccessRequired        = "access token required"
	FailedToGenerateAct       = "failed to generate access token: %w"
	FailedToGenerateRsh       = "failed to generate refresh token: %w"
	AuthRshTokenInvalid       = "invalid refresh token: %w"
	ErrMissingRefreshToken    = "refresh token is missing in request"
	ErrTokenDeletionFailed    = "failed to delete refresh token from in_memory_DB"
	MsgLogoutSuccess          = "user successfully logged out"
	ErrInvalidRequest         = "invalid request parameters"
	AttemptRefreshToken       = "attempting to refresh token"
	FailedToStoreRshToken     = "failed to store token in in_memory_DB"
	SuccessfulLogin           = "user logged in successfully"
	TokenParamMissing         = "token parameters missing"
	SuccessfulRefreshToken    = "token refreshed successfully"
	SuccessfulTokenValidation = "token validated successfully"
	TokenStoredSuccessfully   = "token stored successfully"
	RedisOperationFailed      = "redis operation failed"
	FailedToDeleteRshToken    = "failed to delete token from in_memory_DB"
	TokenDeleteSuccessfully   = "token deleted successfully"
	GenerateAccessToken       = "generate access token"
	RefreshTokenExistence     = "refresh token existence"
	GenerateRefreshToken      = "generate refresh token"

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
	DBQueryError           = "database query execution failed"
	DBRecordNotFound       = "record not found"
	ErrUnsupportedDatabase = "unsupported database"
	ErrKeyNotFound         = "key not found"
	ErrInvalidConfig       = "invalid config"

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
	MobileNoRequired          = "mobile no required"
	CredentialMissing         = "username or password missing"
	ErrUserNotFound           = "user not found"
	DBQueryFailed             = "db query failed"
	UserFetchedSuccessful     = "user data fetched successfully"
	ErrInvalidPassword        = "invalid password"
	WrongUsernamePassword     = "wrong username and password"
	PasswordMissingFromServer = "password missing from database"

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
