package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/nullndvoid/tempmailgen/server/db"
)

// Holds a context and an instance of sqlc's Queries to talk to the database.
// Passed in locals to all our API routes.
type DbQueryState struct {
	Ctx context.Context
	// Provides database query operations and is used to interact with the database layer
	// through the generated SQLC query methods.
	Queries *db.Queries
}

// Creates a new DbQueryState instance.
func NewDbQueryState(ctx context.Context, queries *db.Queries) DbQueryState {
	return DbQueryState{Ctx: ctx, Queries: queries}
}

// Registers various API routes for us.
func RegisterAPIRoutes(app *fiber.App, queryState DbQueryState) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Add ratelimiting to API routes only.
	api.Use(
		limiter.New(limiter.Config{
			Max:               20,
			Expiration:        30 * time.Second,
			LimiterMiddleware: limiter.SlidingWindow{},
		}))

	api.Use(func(c fiber.Ctx) error {
		// Passes our query state to handlers. This is useful because we can then just
		// use query.Create<Resource> and it is all type checked and handled for us.
		c.Locals("db", queryState)
		c.Next()

		return nil
	})

	v1.Post("/email", CreateNewEmailHandler)
}

func CSRFErrorHandler(c fiber.Ctx, err error) error {
	// Log the error so we can track who is trying to perform CSRF attacks.
	// TODO: Make some sort of admin console project for managing my server.

	fmt.Printf("CSRF Error: %v Request: %v From: %v\n", err, c.OriginalURL(), c.IP())

	// Check accepted content types
	switch c.Accepts("html", "json") {

	case "json":
		// Return a 403 Forbidden response for JSON requests.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "403 Forbidden",
		})

	case "html":
		// Return a 403 Forbidden response for HTML requests.
		return c.Status(fiber.StatusForbidden).Render("error", fiber.Map{
			"Title":     "Error",
			"Error":     "403 Forbidden (CSRF Token not given!)",
			"ErrorCode": "403",
		})

	default:
		// Return a 403 Forbidden response for all other requests.
		return c.Status(fiber.StatusForbidden).SendString("403 Forbidden")
	}
}
