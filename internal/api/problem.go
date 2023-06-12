package api

import (
	"net/http"

	"github.com/Tuzi07/solvify-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func (server *Server) setupProblemRoutes() {
	problemGroup := server.router.Group("/api/problems")
	{
		problemGroup.POST("/tf", server.createTFProblem)
		problemGroup.POST("/mtf", server.createMTFProblem)
		problemGroup.POST("/mc", server.createMCProblem)
		problemGroup.POST("/ms", server.createMSProblem)

		problemGroup.POST("/solve-tf", server.solveTFProblem)
		problemGroup.POST("/solve-mtf", server.solveMTFProblem)
		problemGroup.POST("/solve-mc", server.solveMCProblem)
		problemGroup.POST("/solve-ms", server.solveMSProblem)

		problemGroup.GET("/:id", server.getProblem)
		problemGroup.GET("", server.listProblems)
		problemGroup.POST("/:id", server.updateProblem)
		problemGroup.DELETE("/:id", server.deleteProblem)

		problemGroup.POST("/vote/", server.voteProblem)
	}
}

func (server *Server) createTFProblem(ctx *gin.Context) {
	var arg db.CreateTFProblemParams

	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	problem, err := server.db.CreateTFProblem(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, problem)
}

func (server *Server) createMTFProblem(ctx *gin.Context) {
	var arg db.CreateMTFProblemParams

	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	problem, err := server.db.CreateMTFProblem(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, problem)
}

func (server *Server) createMCProblem(ctx *gin.Context) {
	var arg db.CreateMCProblemParams

	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	problem, err := server.db.CreateMCProblem(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, problem)
}

func (server *Server) createMSProblem(ctx *gin.Context) {
	var arg db.CreateMSProblemParams

	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	problem, err := server.db.CreateMSProblem(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, problem)
}

func (server *Server) getProblem(ctx *gin.Context) {
	id := ctx.Param("id")

	problem, err := server.db.GetProblem(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, problem)
}

type listProblemsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listProblems(ctx *gin.Context) {
	var req listProblemsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListProblemsParams{
		Limit: req.PageSize,
		Skip:  (req.PageID - 1) * req.PageSize,
	}

	problems, err := server.db.ListProblems(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, problems)
}

func (server *Server) updateProblem(ctx *gin.Context) {
	var arg db.UpdateProblemParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id := ctx.Param("id")
	if err := server.db.UpdateProblem(arg, id); err != nil {
		if err.Error() == "problem not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, arg)
}

func (server *Server) deleteProblem(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := server.db.DeleteProblem(id); err != nil {
		if err.Error() == "problem not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (server *Server) solveTFProblem(ctx *gin.Context) {
	var arg db.SolveTFProblemParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.db.SolveTFProblem(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) solveMTFProblem(ctx *gin.Context) {
	var arg db.SolveMTFProblemParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.db.SolveMTFProblem(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) solveMCProblem(ctx *gin.Context) {
	var arg db.SolveMCProblemParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.db.SolveMCProblem(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) solveMSProblem(ctx *gin.Context) {
	var arg db.SolveMSProblemParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.db.SolveMSProblem(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) voteProblem(ctx *gin.Context) {
	var arg db.VoteProblemParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.db.VoteProblem(arg); err != nil {
		if err.Error() == "problem not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, arg)
}
