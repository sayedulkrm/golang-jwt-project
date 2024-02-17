package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	routes "github.com/sayedulkrm/golang-jwt-project/routes"
)

func main() {

	port := os.Getenv("PORT")

	fmt.Println(port)

	if port == "" {
		port = "8000"
	}

	router := gin.New()

	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/", func(c *gin.Context) {
		// Set the content type to HTML
		c.Header("Content-Type", "text/html")

		// Send the HTML response with an h1 tag
		html := `<h1>Server is working. To See Frontend <a href="http://localhost:3000"> Click Here </a></h1>`
		c.String(200, html)
	})

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})

	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})

	})

	router.Run(":" + port)

}
