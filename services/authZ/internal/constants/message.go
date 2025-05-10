package constants

// List of env variables
var EnvVariable = struct {
	IN_MEMORY_STORE_TYPE string
}{
	IN_MEMORY_STORE_TYPE: "type",
}

// List of Methods
var Methods = struct {
	NewAuthZService      string
	IsAuthorized         string
	IsAuthorizedBatch    string
	StoreWithTTL         string
	Check                string
	DeleteToken          string
	Store                string
	Get                  string
	GetDirectPermissions string
	GetAccountRole       string
	GetRolePermissions   string
}{
	NewAuthZService:      "NewAuthZService",
	IsAuthorized:         "IsAuthorized",
	IsAuthorizedBatch:    "IsAuthorizedBatch",
	StoreWithTTL:         "StoreWithTTL",
	Store:                "Store",
	Check:                "Check",
	DeleteToken:          "DeleteToken",
	Get:                  "Get",
	GetDirectPermissions: "GetDirectPermissions",
	GetAccountRole:       "GetAccountRole",
	GetRolePermissions:   "GetRolePermissions",
}

const (
	FailedToLoadConfig        = "failed to load config: %v"
	FailedToStartCache        = "Failed to start cache db"
	FailedToStartService      = "failed to initialize AuthZ service: %v"
	GRPCServerRunning         = "✅ AuthZ gRPC server running on %s in %s enviroment"
	FailedToStartServer       = "failed to serve gRPC server: %v"
	FailedToListen            = "failed to listen: %v"
	FailedIniStrManager       = "Failed to initialize store manager: %v"
	StoreConfigNil            = "Store config is nil after loading, cannot proceed"
	FailedPreparePolicy       = "failed to prepare policy query: %w"
	FailedFetchAccount        = "failed to fetch account info: %w"
	FailedFetchRolePermission = "failed to fetch role permissions: %w"
	FailedFetchDPermission    = "failed to fetch direct permissions: %w"
	EvaluationErr             = "policy evaluation error: %w"
	FailedOPAEval             = "OPA evaluation error for resource %s action %s: %w"
	RegoEvalFailed            = "rego evaluation failed: %w"
	DecompressionFailed       = "zstd decompression failed: %w"
	ZSTDEncodingFailed        = "failed to create zstd encoder: %w"
	ZSTDDecodingFailed        = "failed to create zstd decoder: %w"
	AccoutNotAssociated       = "account not associated with any franchise"
	PolicyDenied              = "Policy returned no result, denying access"
	PolicyDenied_2            = "No 'allow' binding found in OPA result, denying access"
	RegoEvaluationErr         = "Invalid type for 'allow' binding in policy result"
	NonBoolean                = "policy returned non-boolean 'allow'"
	IdMismatch                = "franchise ID or account ID mismatch"
	InvalidAssociation        = "invalid franchise or account association"
	NoResources               = "no resources provided for batch authorization"
	ErrUnsupportedDatabase    = "unsupported database"
	ErrKeyNotFound            = "key not found"
	ErrInvalidConfig          = "invalid config"
	// Config error handling messages
	ConfigOverride          = "overriding config type with environment variable: %s"
	FailedToParse           = "failed to parse YAML config"
	UnsupportedDatabaseType = "unsupported database type"
	YAMLFileIssue           = "store config section is missing or nil in YAML file"
	// Store error
	TypeSpecify                   = "%w: type must be specified"
	InvalidConfig                 = "invalid store configuration"
	FallbackLightning             = "falling back to LightningDB due to missing config"
	FallbackLightningDueToFailure = "falling back to LightningDB due to store initialization failure"
	FailedToStoreDecision         = "failed to decision in in_memory_DB"
	FailedToCheckDecision         = "failed to check decision in in_memory_DB"
	DecisionNotExists             = "decision not exists in in_memory_DB"
	InvalidCacheValue             = "invalid cache value type"
	FailedToMarshal               = "failed to marshal proto: %w"
	RedisOperationFailed          = "redis operation failed"
	FailedToDeleteRshToken        = "failed to delete token from in_memory_DB"
	TokenDeleteSuccessfully       = "token deleted successfully"
	// in memory DB Type
	TypeKey       = ""
	RedisType     = "redis"
	MemcachedType = "memcached"
	DragonflyType = "dragonfly"
	BadgerType    = "badger"
	LightningType = "lightning"

	InvalidRedisConfig     = "invalid redis config"
	InvalidMemcachedConfig = "invalid memcached config"
	InvalidDragonflyConfig = "invalid dragonfly config"
	InvalidBadgerConfig    = "invalid badger config"
	CapacityReached        = "cache capacity reached"
	FailedToStoreCache     = "failed to store in cache: %w"
	VerficationFailed      = "cache verification failed: %w"
	FailedToRetrieve       = "failed to retrieve cached value: %w"
	InvalidCacheValueType  = "invalid cache value type: %T"
	DoesnotMatch           = "stored data doesn't match original"
	FailedToUnmarshal      = "failed to unmarshal %w"
	FailedToRead           = "failed to read config file: %w"

	// DButils
	Repository                  = "Repository"
	DBConnectionSuccess         = "connected to the database successfully"
	DBConnectionFailure         = "failed to connect to the database"
	DBConnectionNil             = "database connection is nil"
	DBQueryError                = "database query execution failed"
	DBRecordNotFound            = "record not found"
	BuildInsertQuery            = "something went wrong inside query insert builder function"
	BuildUpdateQuery            = "something went wrong inside query update builder function"
	BuildSelectQuery            = "something went wrong inside query select builder function"
	BuildDeleteQuery            = "something went wrong inside query delete builder function"
	WrongFetchingData           = "something went wrong while fetching data from cache memory"
	DBQueryFailed               = "db query failed"
	FailedToRetrv               = "failed to retrive data from last query performed"
	NoColumProvided             = "no columns provided"
	UnauthorizedSchema          = "unauthorized schema: %s"
	UnauthorizedTable           = "unauthorized table: %s"
	UnauthorizedCloumn          = "unauthorized column: %s"
	UnauthorizedConditionColumn = "unauthorized condition column: %s"
	UnauthorizedReturningColumn = "unauthorized returning column: %s"
	UnauthorizedJoinTable       = "unauthorized join table: %s"
)
