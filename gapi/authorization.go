package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/shouta0715/simple-bank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("metadata is not provided")
	}

	values := md.Get(authorizationHeader)

	if len(values) == 0 {
		return nil, fmt.Errorf("authorization token is not provided")
	}

	authHeader := values[0]

	// <authorization-type> <token>
	fields := strings.Fields(authHeader)

	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])

	if authType != authorizationBearer {
		return nil, fmt.Errorf("authorization type must be bearer")
	}

	accessToken := fields[1]
	payload, err := server.maker.VerifyToken(accessToken)

	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	if !hasPermission(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("access denied")
	}

	return payload, nil
}

func hasPermission(payloadRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if payloadRole == role {
			return true
		}
	}

	return false
}
