package api

import (
	"github.com/Tuzi07/solvify-backend/internal/db"
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
	}

	server.setupUserRoutes()
	server.setupProblemRoutes()
	// setupProblemListRoutes()
}

func (server *Server) Start() error {
	return server.router.Run(":8080")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
