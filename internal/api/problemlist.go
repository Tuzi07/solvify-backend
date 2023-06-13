package api

import (
	"net/http"

	"github.com/Tuzi07/solvify-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func (server *Server) setupProblemListRoutes() {
	problemGroup := server.router.Group("/api/problemlists")
	{
		problemGroup.POST("", server.createProblemList)
		problemGroup.GET("/:id", server.getProblemList)
		problemGroup.GET("", server.listProblemLists)
		problemGroup.POST("/:id", server.updateProblemList)
		problemGroup.DELETE("/:id", server.deleteProblemList)
	}
}

func (server *Server) createProblemList(ctx *gin.Context) {
	var arg db.CreateProblemListParams

	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	problemList, err := server.db.CreateProblemList(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, problemList)
}

func (server *Server) getProblemList(ctx *gin.Context) {
	id := ctx.Param("id")

	problemList, err := server.db.GetProblemList(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, problemList)
}

type listProblemListsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=30"`
}

func (server *Server) listProblemLists(ctx *gin.Context) {
	var req listProblemListsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListProblemListsParams{
		Limit: req.PageSize,
		Skip:  (req.PageID - 1) * req.PageSize,
	}

	problemLists, err := server.db.ListProblemLists(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, problemLists)
}

func (server *Server) updateProblemList(ctx *gin.Context) {
	var arg db.UpdateProblemListParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id := ctx.Param("id")
	if err := server.db.UpdateProblemList(arg, id); err != nil {
		if err.Error() == "problem list not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, arg)
}

func (server *Server) deleteProblemList(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := server.db.DeleteProblemList(id); err != nil {
		if err.Error() == "problem list not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
