package api

import (
	db "github.com/Nickeymaths/bank/db/sqlc"
	"github.com/Nickeymaths/bank/token"
	"github.com/Nickeymaths/bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config      util.Config
	tokenMarker token.Maker
	store       db.Store
	router      *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMarker, err := token.NewPasetoMarker(config.SymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:      config,
		tokenMarker: tokenMarker,
		store:       store,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", currencyValidator)
	}

	server.SetupRouter()

	return server, nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/login", server.loginUser)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.PUT("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.GET("/accounts", server.listAccounts)
	router.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(serverAddress string) error {
	return server.router.Run(serverAddress)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
