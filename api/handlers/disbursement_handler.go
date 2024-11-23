package handlers

import (
	"net/http"
	"zaid-paper-disbursement/internal/services"

	"github.com/gin-gonic/gin"
)

type DisbursementHandler struct {
	service *services.DisbursementService
}

func NewDisbursementHandler(service *services.DisbursementService) *DisbursementHandler {
	return &DisbursementHandler{service: service}
}

func (h *DisbursementHandler) ProcessDisbursement(c *gin.Context) {
	var req services.DisbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ProcessDisbursement(&req); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "insufficient balance" || err.Error() == "user balance not found" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Disbursement processed successfully",
		"data": gin.H{
			"user_id": req.UserID,
			"amount":  req.Amount,
		},
	})
}