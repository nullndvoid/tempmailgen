package main

import (
	"database/sql"
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

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		return
	}

	postgresUri := os.Getenv("POSTGRES_URI")
	if postgresUri == "" {
		log.Printf("POSTGRES_URI not set, please set it in your .env file")
		return
	}

	db, err := sql.Open("postgres", postgresUri)
	if err != nil {
		log.Printf("Error connecting to postgres DB: %v", err)
		return
	}

	store, err := postgres.NewStore(db, []byte("secret"))
	if err != nil {
		log.Printf("Error creating postgres session store: %v", err)
		return
	}

	r := gin.Default()

	r.Use(sessions.Sessions("session", store))

	r.SetTrustedProxies([]string{"127.0.0.0/8", "::1"})
	r.Delims("{{", "}}")

	r.SetFuncMap(template.FuncMap{
		"static": func(file string) string {
			return "/static/" + file
		},
	})

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		log.Println("STATIC_DIR not set, please set it in your .env file")
	}

	r.Use(static.Serve("/static", static.LocalFile(staticDir, true)))

	// Automatically handles rendering of HTML templates if required.
	r.LoadHTMLGlob(staticDir + "/*.html")

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

	// GET /api/email/inbox (Lists emails in inbox)
	// GET /api/email/inbox/search (Fuzzy searches the inbox)
	// DELETE /api/email (Deletes the email address from the database) (CARE SHOULD BE TAKEN SO AS NOT TO DELETE OTHERS)
	// GET /api/email/inbox/<id> (Returns the contents of an email)

	r.Run(":8080")
}

// GET /api/email (Creates a new temporary email)
func emailCreateEndpoint(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
