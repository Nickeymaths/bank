package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	ClientIp  string
	UserAgent string
}

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardForHeader          = "x-forwarded-for"
)

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtd := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgent := md.Get(grpcGatewayUserAgentHeader); len(userAgent) > 0 {
			mtd.UserAgent = userAgent[0]
		}
		if userAgent := md.Get(userAgentHeader); len(userAgent) > 0 {
			mtd.UserAgent = userAgent[0]
		}
		if clientIP := md.Get(xForwardForHeader); len(clientIP) > 0 {
			mtd.ClientIp = clientIP[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtd.ClientIp = p.Addr.String()
	}
	return mtd
}
