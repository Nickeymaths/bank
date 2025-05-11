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
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMarker, err := token.NewPasetoMarker(config.SymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		tokenMaker: tokenMarker,
		store:      store,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", currencyValidator)
	}

	server.SetupRouter()

	return server, nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()

	router.POST("/login", server.loginUser)
	router.POST("/users", server.createUser)

	authorizedRouter := router.Group("/").Use(authMiddleWare(server.tokenMaker))

	authorizedRouter.POST("/accounts", server.createAccount)
	authorizedRouter.GET("/accounts/:id", server.getAccount)
	authorizedRouter.PUT("/accounts/:id", server.updateAccount)
	authorizedRouter.DELETE("/accounts/:id", server.deleteAccount)
	authorizedRouter.GET("/accounts", server.listAccounts)
	authorizedRouter.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(serverAddress string) error {
	return server.router.Run(serverAddress)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
