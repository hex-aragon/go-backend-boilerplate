package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/hex-aragon/go-backend-boilerplate/db/sqlc"
)

//server serves http requests for our banking service.
type Server struct {
	store db.Store 
	router *gin.Engine 
}

//New server creates a new http server and setup routing
func NewServer(store db.Store) *Server{
	server := &Server{store: store}
	router := gin.Default()

	if v, ok :=	binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)
	router.POST("/accounts", server.createAccount)
	router.PUT("/accounts/:id", server.putAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router 
	return server 
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}