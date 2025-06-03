package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.0/8", "::1"})

	r.GET("/temp-mail", func(c *gin.Context) {
		email := "<temporary_email@example.com>"
		c.JSON(http.StatusOK, gin.H{
			"email": email,
		})
	})

	r.Run(":8080")
}
