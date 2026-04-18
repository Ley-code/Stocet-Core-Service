package config

// Role constants
const (
	ADMIN_ROLE      = "admin"
	USER_ROLE       = "user"
	RESEARCHER_ROLE = "researcher"
)

// Subscription status constants
const (
	SUBSCRIPTION_FREE    = "free"
	SUBSCRIPTION_PAID    = "paid"
	SUBSCRIPTION_PREMIUM = "premium"
)

// Permission constants
const (
	PERMISSION_READ_NEWS                        = "read:news"
	PERMISSION_DELETE_NEWS                      = "delete:news"
	PERMISSION_MANAGE_CAPITAL_SERVICE_PROVIDERS = "manage:capital_service_providers"
	PERMISSION_MANAGE_EDUCATION                 = "manage:education"
	PERMISSION_MANAGE_USERS                     = "manage:users"
	PERMISSION_VIEW_ANALYTICS                   = "view:analytics"
	PERMISSION_MANAGE_ALERTS                    = "manage:alerts"
)

// API version constants
const (
	API_VERSION = "v1"
	API_PREFIX  = "/api/v1"
)

// Feature flag constants (for subscription model)
const (
	FEATURE_ADVANCED_ANALYTICS = "advanced_analytics"
	FEATURE_PREMIUM_ALERTS     = "premium_alerts"
	FEATURE_EXPORT_DATA        = "export_data"
	FEATURE_API_ACCESS         = "api_access"
)

// JWT constants
const (
	JWT_ISSUER           = "stocet-core-service"
	JWT_EXPIRATION_HOURS = 24
)

// Pagination constants
const (
	DEFAULT_PAGE_SIZE = 10
	MAX_PAGE_SIZE     = 100
)

// Database constants
const (
	DB_TIMEOUT_SECONDS = 30
)

// Chat constants
const (
	CHAT_MESSAGE_MAX_LENGTH = 1000
	CHAT_CHANNEL_GENERAL    = "general"
)

// Alert constants
const (
	ALERT_CONDITION_ABOVE  = "above"
	ALERT_CONDITION_BELOW  = "below"
	ALERT_STATUS_ACTIVE    = "active"
	ALERT_STATUS_TRIGGERED = "triggered"
	ALERT_STATUS_DISABLED  = "disabled"
)

// News category constants
const (
	NEWS_CATEGORY_TREASURY_BILLS  = "Treasury Bills"
	NEWS_CATEGORY_CAPITAL_MARKETS = "Capital Markets"
	NEWS_CATEGORY_MACROECONOMICS  = "Macroeconomics"
	NEWS_CATEGORY_POLITICS        = "Politics"
)

// Capital service provider status (ECMA-style)
const (
	CSP_STATUS_ACTIVE    = "Active"
	CSP_STATUS_SUSPENDED = "Suspended"
	CSP_STATUS_REVOKED   = "Revoked"
	CSP_STATUS_PENDING   = "Pending"
	CSP_STATUS_INACTIVE  = "Inactive"
)
