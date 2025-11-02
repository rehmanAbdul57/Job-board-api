package middlewares

import (
	"example.com/job-board/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or Invalid Token",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		userId, role, err := utils.VerifyToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		fmt.Println("Middleware OK - userId:", userId, "role:", role)

		c.Locals("user_id", int(userId))
		c.Locals("role", role)
		return c.Next()
	}
}
