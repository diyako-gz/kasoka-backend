package middelwares

import (
	"kasoka/src/common"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

func RoleCheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header missing"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid format"})
			}
			tokenString := parts[1]

			claims, err := common.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
			}

			c.Set("role", claims.Role)
			c.Set("userID", claims.ID)
			c.Set("username", claims.UserName)

			return next(c)
		}
	}
}

func RequireRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userRole, ok := c.Get("role").(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "role not found or invalid"})
			}

			for _, r := range roles {
				if r == userRole {

					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
		}
	}
}
