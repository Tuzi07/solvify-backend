package api

import (
	"net/http"

	"github.com/Tuzi07/solvify-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func (server *Server) setupUserRoutes() {
	userGroup := server.router.Group("/api/users")
	{
		userGroup.POST("", server.createUser)
		userGroup.GET("/:id", server.getUser)
		userGroup.POST("/:id", server.updateUser)
		userGroup.DELETE("/:id", server.deleteUser)

		userGroup.POST("/email", server.getUserByEmail)
		userGroup.POST("/check-email-unique", server.checkEmailUnique)
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var arg db.CreateUserParams

	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.db.CreateUser(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) getUser(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := server.db.GetUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) updateUser(ctx *gin.Context) {
	var arg db.UpdateUserParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id := ctx.Param("id")
	if err := server.db.UpdateUser(arg, id); err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, arg)
}

func (server *Server) deleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := server.db.DeleteUser(id); err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (server *Server) checkEmailUnique(ctx *gin.Context) {
	var arg db.CheckEmailUniqueParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	isUnique, err := server.db.CheckEmailUnique(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"is_unique": isUnique})
}

func (server *Server) getUserByEmail(ctx *gin.Context) {
	var arg db.GetUserByEmailParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.db.GetUserByEmail(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
