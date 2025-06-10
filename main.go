package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nullndvoid/tempmailgen/server/api"
	"github.com/nullndvoid/tempmailgen/server/db"
	// "github.com/gofiber/utils/v2"  // Used for string copies etc.
)

var (
	postgresUri string = ""
	staticDir   string = ""
	templateDir string = ""
)

func main() {
	err := loadEnvVars()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	templateEngine := html.New(templateDir, ".html")
	app := fiber.New(fiber.Config{
		ServerHeader: "tempmailgen v0.1.0",
		AppName:      "Tempmailgen v0.1.0",
		Views:        templateEngine,
	})
	defer app.Shutdown()

	pool := setupDatabase()
	defer pool.Close()

	queries := db.New(pool)
	queryState := api.NewDbQueryState(context.Background(), queries)

	store := setupStorage(pool)
	setupMiddleware(app, store)
	setupRoutes(app, &queryState)

	log.Println("Starting server on port 3000.")
	app.Listen(":3000")
}

// Creates and returns a PostgreSQL connection pool.
func setupDatabase() *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), postgresUri)
	if err != nil {
		log.Fatalf("Failed to create postgres connection pool: %v", err)
	}

	return pool
}

// Creates and configures PostgreSQL storage for sessions/CSRF.
func setupStorage(pool *pgxpool.Pool) *postgres.Storage {
	return postgres.New(postgres.Config{
		DB:         pool,
		Table:      "fiber_storage",
		Reset:      false,
		GCInterval: 10 * time.Second,
	})
}

// Configures all middleware for the application.
func setupMiddleware(app *fiber.App, store *postgres.Storage) {
	app.Use(favicon.New())
	app.Use(logger.New())
	app.Use(helmet.New())

	setupCSRFAndSessions(store, app)
}

// Configures all application routes.
func setupRoutes(app *fiber.App, queryState *api.DbQueryState) {
	api.RegisterAPIRoutes(app, *queryState)

	app.Get("/", func(c fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Use("/static", static.New(staticDir, static.Config{
		CacheDuration: 5 * time.Minute, // Short cache whilst developing.
		MaxAge:        300,             // 5 minutes.

	}))
}

// Configures session and CSRF middleware using PostgreSQL storage.
// Sessions use secure __Host-session cookies, and CSRF protection validates X-CSRF-Token headers
// with a 30-minute idle timeout.
func setupCSRFAndSessions(store *postgres.Storage, app *fiber.App) {
	sessionMiddleware, sessionStore := session.NewWithStore(session.Config{
		Storage:        store,
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		KeyLookup:      "cookie:__Host-session",
	})

	app.Use(sessionMiddleware)

	// Setup and use CSRF for requests to API endpoints. We will need to pass an X-CSRF-Token on protected endpoints.
	csrfConfig := csrf.Config{
		Session:        sessionStore,
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "__Host-csrf", // Recommended to use the __Host- prefix when serving the app over TLS
		CookieSameSite: "Lax",         // Recommended to set this to Lax or Strict
		CookieSecure:   true,          // Recommended to set to true when serving the app over TLS
		CookieHTTPOnly: false,
		ErrorHandler:   api.CSRFErrorHandler,
		IdleTimeout:    30 * time.Minute,
	}

	app.Use(csrf.New(csrfConfig))
}
