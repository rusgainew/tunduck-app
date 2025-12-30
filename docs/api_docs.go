package docs

// SwaggerInfo информация о Swagger API
var SwaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

func init() {
	SwaggerInfo.Version = "1.0.0"
	SwaggerInfo.Host = "localhost:8080"
	SwaggerInfo.BasePath = ""
	SwaggerInfo.Schemes = []string{"http", "https"}
	SwaggerInfo.Title = "Tunduc API System"
	SwaggerInfo.Description = "Enterprise API for managing ESF documents, organizations, and users with caching, rate limiting, and health monitoring"
}

// @title Tunduc API System
// @version 1.0.0
// @description Enterprise API for managing ESF documents, organizations, and users with caching, rate limiting, and health monitoring
// @contact.name API Support
// @contact.url https://github.com/rusgainew/tunduck-app
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @basePath /
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT token in format "Bearer <token>"

// RegisterRequest represents user registration request
// @x-swagger-router-model io.swagger.tunduc.models.RegisterRequest
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"SecurePassword123!"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePassword123!"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

// UserResponse represents user data in response
type UserResponse struct {
	ID       int    `json:"id" example:"1"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status     string            `json:"status" example:"UP"`
	Timestamp  string            `json:"timestamp" example:"2025-12-28T10:30:00Z"`
	Uptime     string            `json:"uptime" example:"2h15m30s"`
	Components []ComponentHealth `json:"components"`
}

// ComponentHealth represents individual component health
type ComponentHealth struct {
	Name         string `json:"name" example:"PostgreSQL"`
	Status       string `json:"status" example:"UP"`
	ResponseTime string `json:"response_time" example:"3ms"`
	Message      string `json:"message" example:"Database connected successfully"`
	LastChecked  string `json:"last_checked" example:"2025-12-28T10:30:00Z"`
}

// ErrorResponse represents standard error response
type ErrorResponse struct {
	Error     string `json:"error" example:"Unauthorized"`
	Message   string `json:"message" example:"Invalid credentials"`
	Timestamp string `json:"timestamp" example:"2025-12-28T10:30:00Z"`
}

// RateLimitErrorResponse represents rate limit error
type RateLimitErrorResponse struct {
	Error   string `json:"error" example:"Too Many Requests"`
	Message string `json:"message" example:"Rate limit exceeded. Please try again later."`
	Reset   int64  `json:"reset" example:"1735382400"`
}
