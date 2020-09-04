package grpc

import (
	context "context"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/antibf/internal/usecase"
	"google.golang.org/grpc"
)

type AntibruteforceServer struct {
	rUseCase *usecase.RateLimitUseCase
	sUseCase *usecase.SettingUseCase
}

func NewAntibruteforceServer(rUseCase *usecase.RateLimitUseCase, sUseCase *usecase.SettingUseCase) *AntibruteforceServer {
	s := new(AntibruteforceServer)
	s.rUseCase = rUseCase
	s.sUseCase = sUseCase

	return s
}

func (s *AntibruteforceServer) Run(addr string) error {
	gs := grpc.NewServer(unaryInterceptor())
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	RegisterAntiBruteforceServiceServer(gs, s)

	return gs.Serve(l)
}

func unaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		defer log.Infof(
			"%s %v",
			info.FullMethod,
			time.Since(start),
		)

		resp, err := handler(ctx, req)
		if err != nil {
			log.Errorf("method %q throws error: %v", info.FullMethod, err)
		}

		return resp, err
	})
}
