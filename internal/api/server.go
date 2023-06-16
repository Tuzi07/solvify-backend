package api

import (
	"github.com/Tuzi07/solvify-backend/internal/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	db     db.Database
	router *gin.Engine
}

func NewServer(db db.Database) *Server {
	server := &Server{db: db}
	server.buildAPIRouter()

	return server
}

func (server *Server) buildAPIRouter() {
	server.router = gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("language", validLanguage)
		v.RegisterValidation("languages", validLanguages)
		v.RegisterValidation("field_to_order_problems", validFieldToOrderProblems)
		v.RegisterValidation("level_of_education", validLevelOfEducation)
	}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	server.router.Use(cors.New(config))

	server.setupProblemRoutes()
	server.setupUserRoutes()
	server.setupProblemListRoutes()
	server.setupLabelsRoutes()
	server.setupProblemEditSuggestionsRoutes()
	server.setupProblemAttemptRoutes()
}

func (server *Server) Start() error {
	return server.router.Run(":8080")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
