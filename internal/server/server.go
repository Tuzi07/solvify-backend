package server

import (
	"log"
	"net/http"

	"github.com/Tuzi07/solvify-backend/internal/database"
	"github.com/Tuzi07/solvify-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Run(db database.Database) error {
	router := buildAPIRouter(db)
	err := router.Run(":8080")
	return err
}

func buildAPIRouter(db database.Database) *gin.Engine {
	router := gin.Default()

	router.POST("/api/problem", createProblemHandler(db))
	router.GET("/api/problem/:id", getProblemHandler(db))
	router.POST("/api/problem/:id", updateProblemHandler(db))

	return router
}

func createProblemHandler(db database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var problem models.Problem
		err := c.BindJSON(&problem)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		problem.ID = uuid.New().String()

		err = db.AddProblem(problem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, problem)
	}
}

func getProblemHandler(db database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")

		problem, err := db.GetProblem(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, problem)
	}
}

func updateProblemHandler(db database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")

		var problem models.Problem
		err := c.BindJSON(&problem)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		problem.ID = id

		err = db.UpdateProblem(problem)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, problem)
	}
}
