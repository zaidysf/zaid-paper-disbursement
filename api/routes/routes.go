package routes

import (
	"database/sql"
	"zaid-paper-disbursement/api/handlers"
	"zaid-paper-disbursement/api/middlewares"
	"zaid-paper-disbursement/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// Initialize services
	disbursementService := services.NewDisbursementService(db)

	// Initialize handlers
	disbursementHandler := handlers.NewDisbursementHandler(disbursementService)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Add middleware
		v1.Use(middlewares.Logger())

		// Disbursement endpoint
		v1.POST("/disbursement", disbursementHandler.ProcessDisbursement)
	}
}