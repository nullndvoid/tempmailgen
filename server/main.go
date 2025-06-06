package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	postgresUri string = ""
	staticDir   string = ""
)

type AppState struct {
	// Shared connection to the database, used for adding temporary emails as well as handling sessions.
	db *sql.DB
}

// Loads the environment variables from the `.env` file.
func loadEnvVars() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("loading .env file failed: %v", err)
	}

	postgresUri = os.Getenv("POSTGRES_URI")
	if postgresUri == "" {
		return errors.New("POSTGRES_URI not set, please set it in your .env file")
	}

	staticDir = os.Getenv("STATIC_DIR")
	if staticDir == "" {
		return errors.New("STATIC_DIR not set, please set it in your .env file")
	}

	return nil
}

// Connects to our postgres db. Doesn't do much yet.
func setupPostgresConnection() (*sql.DB, error) {
	if postgresUri == "" {
		return nil, fmt.Errorf("seems the postgres connection URI is unset, check you have called loadEnvVars")
	}

	return sql.Open("postgres", postgresUri)
}

func setupGinRouter(state *AppState) (*gin.Engine, error) {
	r := gin.Default()

	store, err := postgres.NewStore(state.db, []byte("secret"))
	if err != nil {
		return nil, err
	}

	r.Use(sessions.Sessions("session", store))

	r.SetTrustedProxies([]string{"127.0.0.0/8", "::1"})
	r.Delims("{{", "}}")

	r.SetFuncMap(template.FuncMap{
		"static": func(file string) string {
			return "/static/" + file
		},
	})

	r.Use(static.Serve("/static", static.LocalFile(staticDir, true)))

	// Automatically handles rendering of HTML templates if required.
	r.LoadHTMLGlob(staticDir + "/*.html")

	return r, nil
}

// GET /api/email/inbox (Lists emails in inbox)
// GET /api/email/inbox/search (Fuzzy searches the inbox)
// DELETE /api/email (Deletes the email address from the database) (CARE SHOULD BE TAKEN SO AS NOT TO DELETE OTHERS)
// GET /api/email/inbox/<id> (Returns the contents of an email)
func addRoutes(r *gin.Engine) {
	r.GET("/temp-mail", func(c *gin.Context) {
		email := "<temporary_email@example.com>"
		c.JSON(http.StatusOK, gin.H{
			"email": email,
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Temporary Email Service",
		})
	})

	{
		api := r.Group("/api")
		api.GET("/email", emailCreateEndpoint)
	}
}

func main() {
	err := loadEnvVars()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	db, err := setupPostgresConnection()
	if err != nil {
		log.Printf("Error connecting to postgres DB: %v", err)
		return
	}

	state := AppState{
		db,
	}

	// This just configures the router etc.
	r, err := setupGinRouter(&state)
	if err != nil {
		log.Fatalf("Something went wrong configuring the Gin router: %v", err)
	}

	// Add our routes to Gin.
	addRoutes(r)

	r.Run(":8080")
}

// GET /api/email (Creates a new temporary email)
func emailCreateEndpoint(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
