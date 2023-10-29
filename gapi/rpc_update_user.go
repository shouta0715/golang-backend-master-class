package gapi

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/shouta0715/simple-bank/db/sqlc"
	"github.com/shouta0715/simple-bank/pb"
	"github.com/shouta0715/simple-bank/util"
	"github.com/shouta0715/simple-bank/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violation := validateUpdateUserRequest(req)

	if violation != nil {
		return nil, invalidArgumentError(violation)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's data")
	}

	arg := db.UpdateUserParams{
		Username: authPayload.Username,
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error hashing password: %v", err)
		}

		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}
		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}

	}

	user, err := server.store.UpdateUser(ctx, arg)

	if err != nil {
		if errors.Is(err, db.ErrorRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "error updating user: %v", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if validator.ValidateUsername(req.GetUsername()) != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "username",
			Description: "invalid username",
		})
	}

	if req.Password != nil {
		if validator.ValidatePassword(req.GetPassword()) != nil {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       "password",
				Description: "invalid password",
			})
		}
	}

	if req.Email != nil {
		if validator.ValidateEmail(req.GetEmail()) != nil {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       "email",
				Description: "invalid email",
			})
		}
	}

	if req.FullName != nil {
		if validator.ValidateFullName(req.GetFullName()) != nil {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       "full_name",
				Description: "invalid full_name",
			})
		}
	}

	return violations

}
