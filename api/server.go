package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/shouta0715/simple-bank/db/sqlc"
	"github.com/shouta0715/simple-bank/token"
	"github.com/shouta0715/simple-bank/util"
)

type Server struct {
	store  db.Store
	maker  token.Maker
	router *gin.Engine
	config util.Config
}

// setup api server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannnot create token maker: %w", err)
	}

	server := &Server{store: store, maker: tokenMaker, config: config}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// accounts api
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	// transfer api
	router.POST("/transfers", server.createTransfer)

	// users api
	router.POST("/users", server.createUser)

	// auth api
	router.POST("/auth/login", server.login)

	server.router = router

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
