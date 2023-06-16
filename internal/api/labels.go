package api

import (
	"net/http"

	"github.com/Tuzi07/solvify-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func (server *Server) setupLabelsRoutes() {
	subjectGroup := server.router.Group("/api/subjects")
	{
		subjectGroup.POST("", server.createSubject)
		subjectGroup.GET("/:id", server.getSubject)
		subjectGroup.GET("", server.listSubjects)
	}
	topicGroup := server.router.Group("/api/topics")
	{
		topicGroup.POST("/create", server.createTopic)
		topicGroup.GET("/:id", server.getTopic)
		topicGroup.POST("/list", server.listTopics)
	}
	subtopicGroup := server.router.Group("/api/subtopics")
	{
		subtopicGroup.POST("/create", server.createSubtopic)
		subtopicGroup.GET("/:id", server.getSubtopic)
		subtopicGroup.POST("/list", server.listSubtopics)
	}
}

func (server *Server) createSubject(ctx *gin.Context) {
	var arg db.CreateSubjectParams

	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subject, err := server.db.CreateSubject(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subject)
}

func (server *Server) getSubject(ctx *gin.Context) {
	id := ctx.Param("id")

	subject, err := server.db.GetSubject(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subject)
}

func (server *Server) listSubjects(ctx *gin.Context) {
	subjects, err := server.db.ListSubjects()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subjects)
}

func (server *Server) createTopic(ctx *gin.Context) {
	var arg db.CreateTopicParams

	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	topic, err := server.db.CreateTopic(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, topic)
}

func (server *Server) getTopic(ctx *gin.Context) {
	id := ctx.Param("id")

	topic, err := server.db.GetTopic(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, topic)
}

func (server *Server) listTopics(ctx *gin.Context) {
	var arg db.ListTopicsParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	topics, err := server.db.ListTopics(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, topics)
}

func (server *Server) createSubtopic(ctx *gin.Context) {
	var arg db.CreateSubtopicParams

	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subtopic, err := server.db.CreateSubtopic(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subtopic)
}

func (server *Server) getSubtopic(ctx *gin.Context) {
	id := ctx.Param("id")

	subtopic, err := server.db.GetSubtopic(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subtopic)
}

func (server *Server) listSubtopics(ctx *gin.Context) {
	var arg db.ListSubtopicsParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	subtopics, err := server.db.ListSubtopics(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, subtopics)
}
