package main

import (
	"example.com/event-booking/db"
	"example.com/event-booking/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8000")
}
