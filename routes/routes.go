package routes

import (
	"example.com/event-booking/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine)  {
	// Create server gruop to add authentication in middleware
	authenticated := server.Group("/")
	authenticated.Use(middleware.Authenticate)

	// Get all events
	authenticated.GET("/events", getEvents)

	// Get a specific event by id
	authenticated.GET("/events/:id", getEvent)

	// Create a new event
	authenticated.POST("/events", createEvent)

	// Update an existing event by id
	authenticated.PUT("/events/:id", updateEvent)

	// Delete an existing event by id
	authenticated.DELETE("/events/:id", deleteEvent)

	// Register an event
	authenticated.POST("/events/:id/register", registerEvent)

	// Get registered event by id
	authenticated.GET("/events/register", getRegistrationByID)

	// Cancel an event
	authenticated.DELETE("/events/:id/register", cancelEvent)

	// Sign up users
	server.POST("/signup", signup)

	// Login user
	server.POST("/login", login)
}