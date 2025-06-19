package gapi

import (
	db "github.com/Nickeymaths/bank/db/sqlc"
	"github.com/Nickeymaths/bank/pb"
	"github.com/Nickeymaths/bank/token"
	"github.com/Nickeymaths/bank/util"
	"github.com/Nickeymaths/bank/worker"
)

type Server struct {
	pb.UnimplementedBankServer
	config          util.Config
	tokenMaker      token.Maker
	store           db.Store
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMarker, err := token.NewPasetoMarker(config.SymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:          config,
		tokenMaker:      tokenMarker,
		store:           store,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
