package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/shouta0715/simple-bank/db/sqlc"
	"github.com/shouta0715/simple-bank/pb"
	"github.com/shouta0715/simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	// getterが使用できる。
	hashedPassword, err := util.HashPassword(req.GetPassword())

	if err != nil {

		return nil, status.Errorf(codes.Internal, "failed to hashed password %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {

		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":

				return nil, status.Errorf(codes.InvalidArgument, "username %s already exists", req.GetUsername())
			}
		}

		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}
