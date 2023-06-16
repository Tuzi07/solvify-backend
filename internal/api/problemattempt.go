package api

import (
	"net/http"

	"github.com/Tuzi07/solvify-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func (server *Server) setupProblemAttemptRoutes() {
	server.router.GET("/api/problem-attempt/:id", server.listUserAttempts)
}

type listAttemptsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=30"`
}

func (server *Server) listUserAttempts(ctx *gin.Context) {
	userID := ctx.Param("id")

	var req listAttemptsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pagination := db.PaginationParams{
		Limit: req.PageSize,
		Skip:  (req.PageID - 1) * req.PageSize,
	}

	attempts, err := server.db.ListUserAttempts(userID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, attempts)
}
