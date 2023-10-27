package gapi

import (
	"context"
	"database/sql"

	db "github.com/shouta0715/simple-bank/db/sqlc"
	"github.com/shouta0715/simple-bank/pb"
	"github.com/shouta0715/simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())

	if err != nil {
		if err == sql.ErrNoRows {

			return nil, status.Errorf(codes.NotFound, "cannot find user: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)

	if err != nil {

		return nil, status.Errorf(codes.Unauthenticated, "invalid password: %v", err)
	}

	accessToken, accessPayload, err := server.maker.CreateToken(user.Username, server.config.AccessTokenDuration)

	if err != nil {

		return nil, status.Errorf(codes.Internal, "cannot create access token: %v", err)
	}

	refreshToken, refreshPayload, err := server.maker.CreateToken(user.Username, server.config.RefreshTokenDuration)

	if err != nil {

		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %v", err)
	}

	mtdt := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %v", err)
	}

	rsp := &pb.LoginResponse{
		User:                  convertUser(user),
		SessionId:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt.Time),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiresAt.Time),
	}
	return rsp, nil
}
