package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// @BasePath /booking-service/v1

// Health godoc
// @Summary Liveness
// @Schemes
// @Description Liveness
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {array} api.HealthResponse
// @Router /health/live [get]
func (server *Server) Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}

// @BasePath /booking-service/v1

// Health godoc
// @Summary Readiness
// @Schemes
// @Description Readiness
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {string} api.HealthResponse
// @Router /health/ready [get]
func (server *Server) Ready(ctx *gin.Context) {

	// Check connection with database.
	err := server.store.PingDB()
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN"})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}
