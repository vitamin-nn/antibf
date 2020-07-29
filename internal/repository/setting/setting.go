package setting

import (
	"context"
	"net"

	"github.com/vitamin-nn/otus_anti_bruteforce/internal/repository/setting/inmemory"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/repository/setting/psql"
)

type Setting struct {
	inmemRepo inmemory.InMemory
	psqlRepo  *psql.Psql
}

func NewSettingRepo(inmemRepo inmemory.InMemory, psqlRepo *psql.Psql) *Setting {
	return &Setting{
		inmemRepo: inmemRepo,
		psqlRepo:  psqlRepo,
	}
}

func (s *Setting) AddToWhiteList(ctx context.Context, inet *net.IPNet) error {
	err := s.psqlRepo.AddToWhiteList(ctx, inet)
	if err != nil {
		return err
	}
	err = s.inmemRepo.AddToWhiteList(ctx, inet)
	if err != nil {
		return err
	}
	return nil
}

func (s *Setting) AddToBlackList(ctx context.Context, inet *net.IPNet) error {
	err := s.psqlRepo.AddToBlackList(ctx, inet)
	if err != nil {
		return err
	}
	err = s.inmemRepo.AddToBlackList(ctx, inet)
	if err != nil {
		return err
	}
	return nil
}

func (s *Setting) CheckInWhiteList(ctx context.Context, ip net.IP) (bool, error) {
	return s.inmemRepo.CheckInWhiteList(ctx, ip)
}

func (s *Setting) CheckInBlackList(ctx context.Context, ip net.IP) (bool, error) {
	return s.inmemRepo.CheckInBlackList(ctx, ip)
}
