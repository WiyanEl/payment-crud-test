package middleware

import (
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			strings.TrimSpace(parts[1]) == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Missing or invalid authorization header",
			})
		}

		secret := os.Getenv("JWT_SECRET")
		if len(secret) < 32 {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Authentication configuration error",
			})
		}

		claims := new(jwt.RegisteredClaims)

		token, err := jwt.ParseWithClaims(
			strings.TrimSpace(parts[1]),
			claims,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(secret), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
		)

		if err != nil || token == nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired token",
			})
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil || userID <= 0 {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Invalid token subject",
			})
		}

		c.Locals("userID", userID)

		return c.Next()
	}
}
