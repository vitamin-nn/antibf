package setting

import (
	"context"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/repository/setting/cache"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/repository/setting/psql"
)

type Setting struct {
	cache    cache.Cache
	psqlRepo *psql.Psql
}

func NewSettingRepo(ctx context.Context, psqlRepo *psql.Psql, cacheUpdInterval time.Duration) *Setting {
	s := &Setting{
		psqlRepo: psqlRepo,
	}
	err := s.UpdateCache(ctx)
	if err != nil {
		log.Errorf("update cache error: %v", err)
	}

	go func() {
		ticker := time.NewTicker(cacheUpdInterval)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := s.UpdateCache(ctx)
				if err != nil {
					log.Errorf("update cache error: %v", err)
				}
			}
		}
	}()

	return s
}

func (s *Setting) AddToWhiteList(ctx context.Context, inet *net.IPNet) error {
	err := s.psqlRepo.AddToWhiteList(ctx, inet)
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

	return nil
}

func (s *Setting) DeleteFromWhiteList(ctx context.Context, inet *net.IPNet) error {
	return s.psqlRepo.DeleteFromWhiteList(ctx, inet)
}

func (s *Setting) DeleteFromBlackList(ctx context.Context, inet *net.IPNet) error {
	return s.psqlRepo.DeleteFromBlackList(ctx, inet)
}

func (s *Setting) CheckInWhiteList(ctx context.Context, ip net.IP) (bool, error) {
	return s.cache.CheckInWhiteList(ctx, ip)
}

func (s *Setting) CheckInBlackList(ctx context.Context, ip net.IP) (bool, error) {
	return s.cache.CheckInBlackList(ctx, ip)
}

func (s *Setting) UpdateCache(ctx context.Context) error {
	c := cache.NewSettingCache()
	whiteList, err := s.psqlRepo.GetWhiteList(ctx)
	if err != nil {
		return err
	}
	blackList, err := s.psqlRepo.GetBlackList(ctx)
	if err != nil {
		return err
	}
	c.SetWhiteList(whiteList)
	c.SetBlackList(blackList)
	s.cache = c

	return nil
}
