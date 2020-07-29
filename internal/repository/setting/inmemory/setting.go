package inmemory

import (
	"context"
	"net"
	"sync"
)

type InMemory struct {
	ipWhiteList []*net.IPNet
	ipBlackList []*net.IPNet
	mutex       *sync.RWMutex
}

func NewSettingRepo() InMemory {
	return InMemory{
		mutex: new(sync.RWMutex),
	}
}

func (s *InMemory) AddToWhiteList(ctx context.Context, inet *net.IPNet) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ipWhiteList = append(s.ipWhiteList, inet)
	return nil
}

func (s *InMemory) AddToBlackList(ctx context.Context, inet *net.IPNet) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ipWhiteList = append(s.ipWhiteList, inet)
	return nil
}

func (s *InMemory) CheckInWhiteList(ctx context.Context, ip net.IP) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, inet := range s.ipWhiteList {
		if inet.Contains(ip) {
			return true, nil
		}
	}
	return false, nil
}

func (s *InMemory) CheckInBlackList(ctx context.Context, ip net.IP) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, inet := range s.ipBlackList {
		if inet.Contains(ip) {
			return true, nil
		}
	}
	return false, nil
}
