package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/yuuban007/simplebank/db/sqlc"
	"github.com/yuuban007/simplebank/token"
	"github.com/yuuban007/simplebank/util"
)

// Server server HTTP requests for banking service.
type Server struct {
	config     util.Config
	store      db.Store
	TokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		TokenMaker: tokenMaker,
	}

	// binding validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
}

// Start run the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}
