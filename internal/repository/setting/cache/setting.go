package cache

import (
	"context"
	"net"
	"sync"
)

type Cache struct {
	ipWhiteList []*net.IPNet
	ipBlackList []*net.IPNet
	mutex       *sync.RWMutex
}

func NewSettingCache() Cache {
	return Cache{
		mutex: new(sync.RWMutex),
	}
}

func (s *Cache) CheckInWhiteList(ctx context.Context, ip net.IP) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, inet := range s.ipWhiteList {
		if inet.Contains(ip) {
			return true, nil
		}
	}

	return false, nil
}

func (s *Cache) CheckInBlackList(ctx context.Context, ip net.IP) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, inet := range s.ipBlackList {
		if inet.Contains(ip) {
			return true, nil
		}
	}

	return false, nil
}

func (s *Cache) SetWhiteList(list []*net.IPNet) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ipWhiteList = list
}

func (s *Cache) SetBlackList(list []*net.IPNet) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ipBlackList = list
}
