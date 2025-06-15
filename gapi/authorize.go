package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nickeymaths/bank/pb"
	"github.com/Nickeymaths/bank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	bearerHeader        = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context, req *pb.UpdateUserRequest) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorize header format")
	}

	if bearerHeader != strings.ToLower(fields[0]) {
		return nil, fmt.Errorf("unsupported authentication method")
	}

	payload, err := server.tokenMaker.VerifyToken(fields[1])
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err.Error())
	}

	return payload, nil
}
