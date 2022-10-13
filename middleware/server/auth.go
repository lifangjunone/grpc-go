package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	ClientHeaderAccessKey = "client-id"
	ClientHeaderSecretKey = "client-secret"
)

type AuthInterceptor struct{}

func NewAuthUnaryInterceptor() grpc.UnaryServerInterceptor {
	return (&AuthInterceptor{}).AuthUnaryIntercept
}

func NewAuthStreamInterceptor() grpc.StreamServerInterceptor {
	return (&AuthInterceptor{}).AuthStreamIntercept
}

func NewClientAuth(ak, sk string) (md metadata.MD) {
	return metadata.MD{
		ClientHeaderAccessKey: []string{ak},
		ClientHeaderSecretKey: []string{sk},
	}
}

// AuthUnaryIntercept 认证拦截器
func (a *AuthInterceptor) AuthUnaryIntercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// 初始化Grpc metadata上下文
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("ctx is not a grpc incoming context")
	}
	// 获取凭证信息
	clientId, clientSecret := a.getAuthInfo(md)
	// 验证凭证
	if err = a.verifyAuth(clientId, clientSecret); err != nil {
		return nil, fmt.Errorf("auth verify is failed")
	}
	return handler(ctx, req)
}

// AuthStreamIntercept 认证拦截器
func (a *AuthInterceptor) AuthStreamIntercept(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// 获取上下文
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return fmt.Errorf("ctx is not a grpc incoming context")
	}
	// 获取认证信息
	ak, sk := a.getAuthInfo(md)
	// 验证Auth
	err := a.verifyAuth(ak, sk)
	if err != nil {
		return fmt.Errorf("verify auth info is failed")
	}
	return handler(srv, ss)
}

// getAuthInfo 获取凭证
func (a *AuthInterceptor) getAuthInfo(md metadata.MD) (clientId, clientSecret string) {
	clientIdList := md[ClientHeaderAccessKey]
	clientSecretList := md[ClientHeaderSecretKey]
	if len(clientIdList) > 0 && len(clientSecretList) > 0 {
		clientId = clientIdList[0]
		clientSecret = clientSecretList[0]
		return
	}
	return
}

// verifyAuth 验证凭证
func (a *AuthInterceptor) verifyAuth(clientId, clientSecret string) error {
	if !(clientId == "admin" && clientSecret == "123456") {
		return status.Errorf(codes.Unauthenticated, "Auth verify is failed")
	}
	return nil
}
