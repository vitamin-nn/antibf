package usecase

import (
	"context"
	"net"

	"github.com/vitamin-nn/antibf/internal/ratelimit"
	"github.com/vitamin-nn/antibf/internal/repository/setting"
)

type RateLimitUseCase struct {
	rlLogin  *ratelimit.RateLimit
	rlPasswd *ratelimit.RateLimit
	rlIP     *ratelimit.RateLimit
	s        *setting.Setting
}

func NewRateLimitUseCase(rlLogin *ratelimit.RateLimit, rlPasswd *ratelimit.RateLimit, rlIP *ratelimit.RateLimit, s *setting.Setting) *RateLimitUseCase {
	return &RateLimitUseCase{
		rlLogin:  rlLogin,
		rlPasswd: rlPasswd,
		rlIP:     rlIP,
		s:        s,
	}
}

func (uc RateLimitUseCase) CheckRequest(ctx context.Context, login, passwd string, ip net.IP) (bool, error) {
	// check black
	isBlack, err := uc.s.CheckInBlackList(ctx, ip)
	if err != nil {
		return false, err
	}
	if isBlack {
		return false, nil
	}
	// check white
	isWhite, err := uc.s.CheckInWhiteList(ctx, ip)
	if err != nil {
		return false, err
	}
	if isWhite {
		return true, nil
	}

	// check login rate
	if !uc.rlLogin.Allow(login) || !uc.rlPasswd.Allow(passwd) || !uc.rlIP.Allow(ip.String()) {
		return false, nil
	}

	return true, nil
}

func (uc RateLimitUseCase) ClearRequest(ctx context.Context, login string, ip net.IP) error {
	if login != "" {
		uc.rlLogin.Clear(login)
	}

	if ip != nil {
		uc.rlIP.Clear(ip.String())
	}

	return nil
}
