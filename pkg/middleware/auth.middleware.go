package middleware

import (
	"net/http"
	"strings"

	"github.com/Dom-HTG/attendance-management-system/pkg/utils"
	"github.com/gin-gonic/gin"
)

// JWTClaims represents the custom claims in a JWT token
type JWTClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// AuthMiddleware validates the JWT token in the Authorization header and extracts user information
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the Authorization header
		authHeader := ctx.GetHeader("Authorization")

		// Check if the header exists
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header missing",
			})
			ctx.Abort()
			return
		}

		// Extract the token from the "Bearer <token>" format.
		// Use strings.Fields to tolerate extra spaces and EqualFold to accept case-insensitive "Bearer".
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format. expected 'Bearer <token>'",
			})
			ctx.Abort()
			return
		}

		token := strings.TrimSpace(parts[1])

		// Validate the token and extract claims
		claims, err := utils.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid or expired token",
				"details": err.Error(),
			})
			ctx.Abort()
			return
		}

		// Store user information in the context for later use
		ctx.Set("user_id", claims.ID)
		ctx.Set("user_email", claims.Email)
		ctx.Set("user_role", claims.Role)

		ctx.Next()
	}
}

// RoleMiddleware checks if the user has a specific role
// Usage: RoleMiddleware("lecturer") or RoleMiddleware("student")
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the user role from the context (set by AuthMiddleware)
		userRole, exists := ctx.Get("user_role")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "user role not found in context",
			})
			ctx.Abort()
			return
		}

		// Check if the user has the required role
		if userRole.(string) != requiredRole {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "access denied. only " + requiredRole + "s are allowed to access this endpoint",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(ctx *gin.Context) (int, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(int), true
}

// GetUserRoleFromContext retrieves the user role from the context
func GetUserRoleFromContext(ctx *gin.Context) (string, bool) {
	userRole, exists := ctx.Get("user_role")
	if !exists {
		return "", false
	}
	return userRole.(string), true
}

// GetUserEmailFromContext retrieves the user email from the context
func GetUserEmailFromContext(ctx *gin.Context) (string, bool) {
	userEmail, exists := ctx.Get("user_email")
	if !exists {
		return "", false
	}
	return userEmail.(string), true
}
