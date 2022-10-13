package client

import (
	"context"
	"grpc-go/middleware/server"
)

type Authentication struct {
	AccessKey string
	SecretKey string
}

func NewAuthentication(ak, sk string) *Authentication {
	return &Authentication{
		AccessKey: ak,
		SecretKey: sk,
	}
}

func (c *Authentication) buildCredential() map[string]string {
	return map[string]string{
		server.ClientHeaderAccessKey: c.AccessKey,
		server.ClientHeaderSecretKey: c.SecretKey,
	}
}

func (c *Authentication) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return c.buildCredential(), nil
}

func (c *Authentication) RequireTransportSecurity() bool { return false }
