package usecase

import (
	"context"
	"net"

	"github.com/vitamin-nn/antibf/internal/repository/setting"
)

type SettingUseCase struct {
	sRepo *setting.Setting
}

func NewSettingUseCase(sRepo *setting.Setting) *SettingUseCase {
	return &SettingUseCase{
		sRepo: sRepo,
	}
}

func (s SettingUseCase) AddToWhiteList(ctx context.Context, inet *net.IPNet) error {
	return s.sRepo.AddToWhiteList(ctx, inet)
}

func (s SettingUseCase) AddToBlackList(ctx context.Context, inet *net.IPNet) error {
	return s.sRepo.AddToBlackList(ctx, inet)
}

func (s SettingUseCase) DeleteFromWhiteList(ctx context.Context, inet *net.IPNet) error {
	return s.sRepo.DeleteFromWhiteList(ctx, inet)
}

func (s SettingUseCase) DeleteFromBlackList(ctx context.Context, inet *net.IPNet) error {
	return s.sRepo.DeleteFromBlackList(ctx, inet)
}
