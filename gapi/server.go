package gapi

import (
	"fmt"

	db "github.com/shouta0715/simple-bank/db/sqlc"
	"github.com/shouta0715/simple-bank/pb"
	"github.com/shouta0715/simple-bank/token"
	"github.com/shouta0715/simple-bank/util"
	"github.com/shouta0715/simple-bank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store
	maker           token.Maker
	config          util.Config
	taskDistributor worker.TaskDistributor
}

// setup gRPC server
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannnot create token maker: %w", err)
	}

	server := &Server{
		store:           store,
		maker:           tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
