package payments

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CreatePaymentRequest struct {
	CustomerNumber string `json:"customer_number"`
	CustomerName string `json:"customer_name"`
	Amount int64  `json:"amount"`
	Method string `json:"method"`
}

type Payment struct {
	ID int64  `json:"id"`
	CustomerNumber string `json:"customer_number"`
	CustomerName string `json:"customer_name"`
	Amount int64  `json:"amount"`
	Method string `json:"method"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type UpdatePaymentRequest struct {
	CustomerNumber string `json:"customer_number"`
	CustomerName string `json:"customer_name"`
	Amount int64  `json:"amount"`
	Method string `json:"method"`
	Status string `json:"status"`
}

func Create(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(int64)
		if !ok || userID <= 0 {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}

		var req CreatePaymentRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
			})
		}

		req.CustomerNumber = strings.TrimSpace(req.CustomerNumber)
		req.CustomerName = strings.TrimSpace(req.CustomerName)
		req.Method = strings.ToLower(strings.TrimSpace(req.Method))

		if req.CustomerNumber == "" {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Customer number is required",
			})
		}

		if req.CustomerName == "" {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Customer name is required",
			})
		}

		if req.Amount <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Positive amount is required",
			})
		}

		switch req.Method {
		case "transfer", "ewallet", "va":
		default:
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Method must be transfer, ewallet, or va",
			})
		}

		var payment struct {
			ID int64  `json:"id"`
			CustomerNumber string `json:"customer_number"`
			CustomerName string `json:"customer_name"`
			Amount int64  `json:"amount"`
			Method string `json:"method"`
			Status string `json:"status"`
		}

		err := db.QueryRow(`
			INSERT INTO payments (
				user_id,
				customer_number,
				customer_name,
				amount,
				method
			)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING
				id,
				customer_number,
				customer_name,
				amount,
				method,
				status
		`,
			userID,
			req.CustomerNumber,
			req.CustomerName,
			req.Amount,
			req.Method,
		).Scan(
			&payment.ID,
			&payment.CustomerNumber,
			&payment.CustomerName,
			&payment.Amount,
			&payment.Method,
			&payment.Status,
		)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to create payment",
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"success": true,
			"message": "Payment created successfully",
			"data":    payment,
		})
	}
}

func List(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(int64)
		if !ok || userID <= 0 {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}

		page, err := strconv.Atoi(c.Query("page", "1"))
		if err != nil || page < 1 {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Page must be a positive integer",
			})
		}

		limit, err := strconv.Atoi(c.Query("limit", "10"))
		if err != nil || limit < 1 || limit > 100 {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Limit must be between 1 and 100",
			})
		}

		offset := (page - 1) * limit

		var total int64
		err = db.QueryRow(
			`SELECT COUNT(*) FROM payments WHERE user_id = $1`,
			userID,
		).Scan(&total)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to count payments",
			})
		}

		rows, err := db.Query(`
			SELECT
				id,
				customer_number,
				customer_name,
				amount,
				method,
				status,
				created_at
			FROM payments
			WHERE user_id = $1
			ORDER BY id DESC
			LIMIT $2 OFFSET $3
		`, userID, limit, offset)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to retrieve payments",
			})
		}
		defer rows.Close()

		payments := make([]Payment, 0)

		for rows.Next() {
			var payment Payment
			var createdAt interface{}

			err := rows.Scan(
				&payment.ID,
				&payment.CustomerNumber,
				&payment.CustomerName,
				&payment.Amount,
				&payment.Method,
				&payment.Status,
				&createdAt,
			)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"success": false,
					"message": "Failed to read payment data",
				})
			}

			if t, ok := createdAt.(interface{ Format(string) string }); ok {
				payment.CreatedAt = t.Format("2006-01-02T15:04:05Z07:00")
			}

			payments = append(payments, payment)
		}

		if err := rows.Err(); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to iterate payments",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    payments,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		})
	}
}


func GetByID(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(int64)
		if !ok || userID <= 0 {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}

		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil || id <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid payment ID",
			})
		}

		var payment Payment

		err = db.QueryRow(`
			SELECT
				id,
				customer_number,
				customer_name,
				amount,
				method,
				status,
				created_at
			FROM payments
			WHERE id = $1 AND user_id = $2
		`, id, userID).Scan(
			&payment.ID,
			&payment.CustomerNumber,
			&payment.CustomerName,
			&payment.Amount,
			&payment.Method,
			&payment.Status,
			&payment.CreatedAt,
		)

		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Payment not found",
			})
		}

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to retrieve payment",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    payment,
		})
	}
}

func Update(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(int64)
		if !ok || userID <= 0 {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}

		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil || id <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid payment ID",
			})
		}

		var req UpdatePaymentRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
			})
		}

		req.CustomerNumber = strings.TrimSpace(req.CustomerNumber)
		req.CustomerName = strings.TrimSpace(req.CustomerName)
		req.Method = strings.ToLower(strings.TrimSpace(req.Method))
		req.Status = strings.ToLower(strings.TrimSpace(req.Status))

		if req.CustomerNumber == "" ||
			req.CustomerName == "" ||
			req.Amount <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Customer number, customer name, and positive amount are required",
			})
		}

		switch req.Method {
		case "transfer", "ewallet", "va":
		default:
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Method must be transfer, ewallet, or va",
			})
		}

		switch req.Status {
		case "pending", "success", "failed":
		default:
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Status must be pending, success, or failed",
			})
		}

		var payment Payment

		err = db.QueryRow(`
			UPDATE payments
			SET
				customer_number = $1,
				customer_name = $2,
				amount = $3,
				method = $4,
				status = $5,
				updated_at = NOW()
			WHERE id = $6 AND user_id = $7
			RETURNING
				id,
				customer_number,
				customer_name,
				amount,
				method,
				status,
				created_at
		`,
			req.CustomerNumber,
			req.CustomerName,
			req.Amount,
			req.Method,
			req.Status,
			id,
			userID,
		).Scan(
			&payment.ID,
			&payment.CustomerNumber,
			&payment.CustomerName,
			&payment.Amount,
			&payment.Method,
			&payment.Status,
			&payment.CreatedAt,
		)

		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Payment not found",
			})
		}

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to update payment",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Payment updated successfully",
			"data":    payment,
		})
	}
}


func Delete(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(int64)
		if !ok || userID <= 0 {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}

		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil || id <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid payment ID",
			})
		}

		result, err := db.Exec(`
			DELETE FROM payments
			WHERE id = $1 AND user_id = $2
		`, id, userID)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to delete payment",
			})
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to verify deleted payment",
			})
		}

		if rowsAffected == 0 {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Payment not found",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Payment deleted successfully",
		})
	}
}
