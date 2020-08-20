package grpc

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"
	outErr "github.com/vitamin-nn/otus_anti_bruteforce/internal/error"
)

func (s *AntibruteforceServer) Check(ctx context.Context, req *CheckRequest) (*CheckResponse, error) {
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
		resp := &CheckResponse{
			Result: &CheckResponse_Error{
				Error: oErr.Error(),
			},
		}
		return resp, nil
	}

	resp := &CheckResponse{
		Result: &CheckResponse_Ok{
			Ok: ok,
		},
	}
	return resp, nil
}

func (s *AntibruteforceServer) Clear(ctx context.Context, req *ClearRequest) (*ModifyResponse, error) {
	login := req.GetLogin()
	ip := net.ParseIP(req.GetIp())

	if login == "" && ip == nil {
		return nil, outErr.ErrEmptyClearParams
	}

	err := s.rUseCase.ClearRequest(ctx, login, ip)
	if err != nil {
		oErr, ok := err.(outErr.OutError)
		if !ok {
			oErr = outErr.ErrInternal
			log.Errorf("unknown error: %v", err)
		}
		resp := &ModifyResponse{
			Result: &ModifyResponse_Error{
				Error: oErr.Error(),
			},
		}
		return resp, nil
	}

	resp := &ModifyResponse{
		Result: &ModifyResponse_Success{
			Success: true,
		},
	}
	return resp, nil
}
