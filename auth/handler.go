package auth

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func Register(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RegisterRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
			})
		}

		req.Name = strings.TrimSpace(req.Name)
		req.Email = strings.TrimSpace(req.Email)

		if req.Name == "" {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Name is required",
			})
		}

		if req.Email == "" {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Email is required",
			})
		}
		if req.Password == "" {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Password is required",
			})
		}

		if len(req.Password) < 8 {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Password must be at least 8 characters",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to process password",
			})
		}

		var userID int64

		err = db.QueryRow(
			`INSERT INTO users (name, email, password)
			VALUES ($1, $2, $3)
			RETURNING id`,
			req.Name,
			req.Email,
			string(hashedPassword),
		).Scan(&userID)

		if err != nil {
			log.Printf("Register user error: %v", err)

			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to register user",
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"success": true,
			"message": "User registered successfully",
			"data": fiber.Map{
				"id":    userID,
				"name":  req.Name,
				"email": req.Email,
			},
		})
	}
}

func Login(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
			})
		}

		req.Email = strings.TrimSpace(req.Email)

		if req.Email == "" {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Email is required",
			})
		}

		if req.Password == "" {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Password is required",
			})
		}

		var userID int64
		var name, passwordHash string

		err := db.QueryRow(
			`SELECT id, name, password
			FROM users
			WHERE email = $1`,
			req.Email,
		).Scan(&userID, &name, &passwordHash)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(401).JSON(fiber.Map{
					"success": false,
					"message": "Invalid email or password",
				})
			}

			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Internal server error",
			})
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(passwordHash),
			[]byte(req.Password),
		); err != nil {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Invalid email or password",
			})
		}
		
		secret := os.Getenv("JWT_SECRET")
		if len(secret) < 32 {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Authentication configuration error",
			})
		}

		expiresIn, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
		if err != nil || expiresIn <= 0 {
			expiresIn = 24
		}

		now := time.Now()

		claims := jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiresIn) * time.Hour)),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to generate token",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Login successful",
			"data": fiber.Map{
					"token": tokenString,
					"user": fiber.Map{
							"id":    userID,
							"name":  name,
							"email": req.Email,
					},
			},
	})
	}
}