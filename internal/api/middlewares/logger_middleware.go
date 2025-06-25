package middlewares

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger(c *fiber.Ctx) {
	// Start time
	startTime := time.Now()

	// Process request
	c.Next()

	// End time
	endTime := time.Now()

	// Execution time
	latency := endTime.Sub(startTime)

	// Request details
	method := string(c.Request().Header.Method())
	statusCode := c.Response().StatusCode()
	path := c.Path()

	// Log request details
	fmt.Printf("[%s] %d | %v | %s | %s\n",
		endTime.Format("2006-01-02 15:04:05"),
		statusCode,
		latency,
		method,
		path,
	)
}
