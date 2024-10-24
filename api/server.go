package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/hex-aragon/go-backend-boilerplate/db/sqlc"
	"github.com/hex-aragon/go-backend-boilerplate/token"
	"github.com/hex-aragon/go-backend-boilerplate/util"
)

//server serves http requests for our banking service.
type Server struct {
	config util.Config
	store db.Store 
	tokenMaker token.Maker 
	router *gin.Engine 
}

//New server creates a new http server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error){
	//tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	
	server := &Server{
		config: config,
		store: store, 
		tokenMaker: tokenMaker,
	}
	

	if v, ok :=	binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server , nil 
}

func (server *Server) setupRouter(){
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.PUT("/accounts/:id", server.putAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	authRoutes.POST("/transfers", server.createTransfer)
	
	server.router = router 
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}