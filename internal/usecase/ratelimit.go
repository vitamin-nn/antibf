package usecase

import (
	"context"
	"net"

	"github.com/vitamin-nn/otus_anti_bruteforce/internal/ratelimit"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/repository/setting"
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
	// TODO: наверное, стоит все 3 проверки запускать в отдельных горутинах и при получении хотябы одного false сразу отдавать ответ
	if !uc.rlLogin.IsAllow(login) || !uc.rlPasswd.IsAllow(passwd) || !uc.rlIP.IsAllow(ip.String()) {
		return false, nil
	}
	return true, nil
}
