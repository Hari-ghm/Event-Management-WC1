package routes

import (
	"github.com/Hari-ghm/Event-Management-WC1/controllers"
	"github.com/Hari-ghm/Event-Management-WC1/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is running"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	event := r.Group("/event")
	event.Use(middlewares.AuthMiddleware())
	{
		event.POST("/create", controllers.CreateEvent)
		event.GET("/list", controllers.ListEvents)
		event.GET("/:id", controllers.GetEventByID)
		event.PUT("/:id", controllers.UpdateEvent)
		event.DELETE("/:id", controllers.DeleteEvent)
	}
}
