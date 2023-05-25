package server

import (
	"log"
	"net/http"

	"github.com/Tuzi07/solvify-backend/internal/database"
	"github.com/Tuzi07/solvify-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func Run(db database.Database) error {
	router := buildAPIRouter(db)
	err := router.Run(":8080")
	return err
}

func buildAPIRouter(db database.Database) *gin.Engine {
	router := gin.Default()

	router.POST("/api/problem", requestHandlerOfCreateProblem(db))
	router.GET("/api/problem/:id", requestHandlerOfGetProblem(db))
	router.POST("/api/problem/:id", requestHandlerOfUpdateProblem(db))
	router.DELETE("/api/problem/:id", requestHandlerOfDeleteProblem(db))

	return router
}

func requestHandlerOfCreateProblem(db database.Database) gin.HandlerFunc {
	return func(requestContext *gin.Context) {
		var problem models.Problem
		err := requestContext.BindJSON(&problem)
		if err != nil {
			emptyResponse := gin.H{}
			requestContext.JSON(http.StatusBadRequest, emptyResponse)
			return
		}

		err = db.CreateProblem(&problem)
		if err != nil {
			emptyResponse := gin.H{}
			requestContext.JSON(http.StatusInternalServerError, emptyResponse)
			return
		}

		requestContext.JSON(http.StatusOK, problem)
	}
}

func requestHandlerOfGetProblem(db database.Database) gin.HandlerFunc {
	return func(requestContext *gin.Context) {
		id := requestContext.Params.ByName("id")

		var problem models.Problem

		err := db.GetProblem(&problem, id)
		if err != nil {
			log.Println(err)
			emptyResponse := gin.H{}
			requestContext.JSON(http.StatusInternalServerError, emptyResponse)
			return
		}

		requestContext.JSON(http.StatusOK, problem)
	}
}

func requestHandlerOfUpdateProblem(db database.Database) gin.HandlerFunc {
	return func(requestContext *gin.Context) {
		id := requestContext.Params.ByName("id")

		var problem models.Problem
		err := requestContext.BindJSON(&problem)
		if err != nil {
			emptyResponse := gin.H{}
			requestContext.JSON(http.StatusBadRequest, emptyResponse)
			return
		}

		err = db.UpdateProblem(&problem, id)
		if err != nil {
			log.Println(err)
			emptyResponse := gin.H{}
			requestContext.JSON(http.StatusInternalServerError, emptyResponse)
			return
		}

		requestContext.JSON(http.StatusOK, problem)
	}
}

func requestHandlerOfDeleteProblem(db database.Database) gin.HandlerFunc {
	return func(requestContext *gin.Context) {
		id := requestContext.Params.ByName("id")

		err := db.DeleteProblem(id)
		if err != nil {
			log.Println(err)
			emptyResponse := gin.H{}
			requestContext.JSON(http.StatusInternalServerError, emptyResponse)
			return
		}

		requestContext.JSON(http.StatusNoContent, nil)
	}
}
