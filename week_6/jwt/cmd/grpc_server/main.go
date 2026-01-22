package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/airats/micro_as_bigtech_course/week_6/jwt/internal/model"
	"github.com/airats/micro_as_bigtech_course/week_6/jwt/internal/utils"
	descAccess "github.com/airats/micro_as_bigtech_course/week_6/jwt/pkg/access_v1"
	descAuth "github.com/airats/micro_as_bigtech_course/week_6/jwt/pkg/auth_v1"
	"github.com/golang/protobuf/ptypes/empty"
	formattedErrors "github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

const (
	gprcPort   = 50061
	authPrefix = "Bearer"

	refreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="

	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 5 * time.Minute
)

type serverAuth struct {
	descAuth.UnimplementedAuthV1Server
}

func (s *serverAuth) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {
	// Проверяем пароль

	refreshToken, err := utils.GenerateToken(
		model.UserInfo{
			Username: req.GetUsername(),
			Role:     "admin",
		},
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)

	if err != nil {
		return nil, formattedErrors.Errorf("failed to generate refresh token: %v", err)
	}

	return &descAuth.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}

func (s *serverAuth) GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	refreshToken, err := utils.GenerateToken(
		model.UserInfo{
			Username: claims.Username,
			Role:     "admin",
		},
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &descAuth.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}

func (s *serverAuth) GetAccessToken(ctx context.Context, req *descAuth.GetAccessTokenRequest) (*descAuth.GetAccessTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	accessToken, err := utils.GenerateToken(
		model.UserInfo{
			Username: claims.Username,
			Role:     "admin",
		},
		[]byte(accessTokenSecretKey),
		accessTokenExpiration,
	)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	return &descAuth.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}

// access server
type serverAccess struct {
	descAccess.UnimplementedAccessV1Server
}

func (s *serverAccess) Check(ctx context.Context, req *descAccess.CheckRequest) (*empty.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no authorization headers")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("no authorization headers")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid auth header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)
	claims, err := utils.VerifyToken(accessToken, []byte(accessTokenSecretKey))
	if err != nil {
		return nil, errors.New("failed verification")
	}

	accessibleRoles := map[string]string{
		"note_": "admin",
	}

	role, ok := accessibleRoles[req.GetEndpointAddress()]
	if !ok {
		return &empty.Empty{}, nil
	}

	if role == claims.Role {
		return &empty.Empty{}, nil
	}

	return nil, errors.New("access denied")
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", gprcPort))
	if err != nil {
		log.Fatalf("failed to listen port: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	descAuth.RegisterAuthV1Server(grpcServer, &serverAuth{})
	descAccess.RegisterAccessV1Server(grpcServer, &serverAccess{})

	log.Println("Running grpc server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to run grpc server: %v", err)
	}
}
