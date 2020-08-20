package grpc

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"
	outErr "github.com/vitamin-nn/otus_anti_bruteforce/internal/error"
)

type modifyFunc func(context.Context, *net.IPNet) error

func (s *AntibruteforceServer) AddToWhiteList(ctx context.Context, req *ModifyListRequest) (*ModifyResponse, error) {
	return s.modify(ctx, req, s.sUseCase.AddToWhiteList)
}

func (s *AntibruteforceServer) AddToBlackList(ctx context.Context, req *ModifyListRequest) (*ModifyResponse, error) {
	return s.modify(ctx, req, s.sUseCase.AddToBlackList)
}

func (s *AntibruteforceServer) RemoveFromWhiteList(ctx context.Context, req *ModifyListRequest) (*ModifyResponse, error) {
	return s.modify(ctx, req, s.sUseCase.DeleteFromWhiteList)
}

func (s *AntibruteforceServer) RemoveFromBlackList(ctx context.Context, req *ModifyListRequest) (*ModifyResponse, error) {
	return s.modify(ctx, req, s.sUseCase.DeleteFromBlackList)
}

func (s *AntibruteforceServer) modify(ctx context.Context, req *ModifyListRequest, f modifyFunc) (*ModifyResponse, error) {
	_, ipNet, err := net.ParseCIDR(req.GetIp())
	if err != nil {
		return nil, outErr.ErrInvalidParams
	}

	err = f(ctx, ipNet)
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
