package grpc

import (
	context "context"
	"net"

	log "github.com/sirupsen/logrus"
	outErr "github.com/vitamin-nn/otus_anti_bruteforce/internal/error"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/usecase"
	"google.golang.org/grpc"
)

type AntibruteforceServer struct {
	rUseCase *usecase.RateLimitUseCase
}

func NewAntibruteforceServer(rUseCase *usecase.RateLimitUseCase) *AntibruteforceServer {
	s := new(AntibruteforceServer)
	s.rUseCase = rUseCase
	return s
}

func (s *AntibruteforceServer) CheckRequest(ctx context.Context, req *CheckAuthRequest) (*CheckAuthResponse, error) {
	login := req.GetLogin()
	if login == "" {
		return nil, outErr.ErrEmptyLogin
	}

	passwd := req.GetPassword()
	if passwd == "" {
		return nil, outErr.ErrEmptyPassword
	}

	ip := net.ParseIP(req.GetIp())
	if ip == nil {
		return nil, outErr.ErrEmptyIP
	}

	ok, err := s.rUseCase.CheckRequest(ctx, login, passwd, ip)
	if err != nil {
		oErr, ok := err.(outErr.OutError)
		if !ok {
			oErr = outErr.ErrInternal
			log.Errorf("unknown error: %v", err)
		}
		resp := &CheckAuthResponse{
			Result: &CheckAuthResponse_Error{
				Error: oErr.Error(),
			},
		}
		return resp, nil
	}

	resp := &CheckAuthResponse{
		Result: &CheckAuthResponse_Ok{
			Ok: ok,
		},
	}
	return resp, nil
}

func (s *AntibruteforceServer) Run(addr string) error {
	gs := grpc.NewServer()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	RegisterAntiBruteforceServiceServer(gs, s)
	return gs.Serve(l)
}
