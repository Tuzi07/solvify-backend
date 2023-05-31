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
		problemGroup.POST("", server.createTFProblem)
		problemGroup.GET("/:id", server.getTFProblem)
		problemGroup.POST("/:id", server.updateTFProblem)
		problemGroup.DELETE("/:id", server.deleteTFProblem)
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

func (server *Server) getTFProblem(ctx *gin.Context) {
	id := ctx.Params.ByName("id")

	problem, err := server.db.GetTFProblem(id)
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

func (server *Server) updateTFProblem(ctx *gin.Context) {
	var arg db.UpdateTFProblemParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id := ctx.Params.ByName("id")

	err := server.db.UpdateTFProblem(arg, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, arg)
}

func (server *Server) deleteTFProblem(ctx *gin.Context) {
	id := ctx.Params.ByName("id")

	err := server.db.DeleteTFProblem(id)
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
