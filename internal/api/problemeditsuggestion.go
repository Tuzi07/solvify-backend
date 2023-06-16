package api

import (
	"net/http"

	"github.com/Tuzi07/solvify-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func (server *Server) setupProblemEditSuggestionsRoutes() {
	problemEditSuggestionGroup := server.router.Group("/api/problem-edit-suggestions")
	{
		problemEditSuggestionGroup.POST("", server.createProblemEditSuggestion)
		problemEditSuggestionGroup.GET("/:id", server.getProblemEditSuggestion)
		problemEditSuggestionGroup.DELETE("/:id", server.deleteProblemEditSuggestion)
		problemEditSuggestionGroup.GET("/:id/accept", server.acceptProblemEditSuggestion)
	}
	problemGroup := server.router.Group("/api/problems")
	{
		problemGroup.GET("/:id/suggestions", server.listSuggestionsOfProblem)
	}
}

func (server *Server) createProblemEditSuggestion(ctx *gin.Context) {
	var arg db.CreateProblemEditSuggestionParams

	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	suggestion, err := server.db.CreateProblemEditSuggestion(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, suggestion)
}

func (server *Server) getProblemEditSuggestion(ctx *gin.Context) {
	id := ctx.Param("id")

	suggestion, err := server.db.GetProblemEditSuggestion(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, suggestion)
}

func (server *Server) deleteProblemEditSuggestion(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := server.db.DeleteProblemEditSuggestion(id); err != nil {
		if err.Error() == "problem edit suggestion not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (server *Server) acceptProblemEditSuggestion(ctx *gin.Context) {
	id := ctx.Param("id")

	problem, err := server.db.AcceptProblemEditSuggestion(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, problem)
}

func (server *Server) listSuggestionsOfProblem(ctx *gin.Context) {
	id := ctx.Param("id")

	suggestions, err := server.db.ListSuggestionsOfProblem(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, suggestions)
}
