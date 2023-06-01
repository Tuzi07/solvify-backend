package api

import (
	"net/http"

	"github.com/Tuzi07/solvify-backend/internal/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func (server *Server) setupProblemRoutes() {
	problemGroup := server.router.Group("/api/problems")
	{
		problemGroup.POST("/tf", server.createTFProblem)
		problemGroup.POST("/mtf", server.createMTFProblem)
		problemGroup.POST("/mc", server.createMCProblem)
		problemGroup.POST("/ms", server.createMSProblem)

		problemGroup.GET("/:id", server.getProblem)
		problemGroup.GET("", server.listProblems)
		problemGroup.POST("/:id", server.updateProblem)
		problemGroup.DELETE("/:id", server.deleteProblem)
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
	id := ctx.Params.ByName("id")

	problem, err := server.db.GetProblem(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

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

	id := ctx.Params.ByName("id")

	err := server.db.UpdateProblem(arg, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, arg)
}

func (server *Server) deleteProblem(ctx *gin.Context) {
	id := ctx.Params.ByName("id")

	err := server.db.DeleteProblem(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
