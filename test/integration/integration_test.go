// +build integration

package integration

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/config"
	grpcService "github.com/vitamin-nn/otus_anti_bruteforce/internal/grpc"
	"google.golang.org/grpc"

	"github.com/stretchr/testify/require"
)

type GrpcAPITestSuite struct {
	suite.Suite
	gClient grpcService.AntiBruteforceServiceClient
	gConn   *grpc.ClientConn
	cfg     *config.Config
}

type modifyFunc func(context.Context, *grpcService.ModifyListRequest, ...grpc.CallOption) (*grpcService.ModifyResponse, error)

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(GrpcAPITestSuite))
}

func (s *GrpcAPITestSuite) TestServerCheck() {
	ctx := context.Background()
	var n int

	for i := 1; i <= 2*s.cfg.RateLimit.Login; i++ {
		if !s.checkRequest(ctx, "test-check", generateRandStr(15), generateRandIP()) {
			n = i
			break
		}
	}
	// здесь проверка не на равенство, т.к. часть токенов может быть регенерировано
	require.GreaterOrEqual(s.T(), n, s.cfg.RateLimit.Login)
}

func (s *GrpcAPITestSuite) TestServerWhiteList() {
	ctx := context.Background()
	checkIP := "192.168.1.1"
	// доходим до лимита
	for i := 1; i <= 2*s.cfg.RateLimit.IP; i++ {
		if !s.checkRequest(ctx, generateRandStr(15), generateRandStr(15), checkIP) {
			break
		}
	}

	// добавляем подсеть в белый список
	s.modifyListRequest(ctx, "192.168.0.0/22", s.gClient.AddToWhiteList)
	time.Sleep(s.cfg.WBSetting.CacheUpdInterval)
	// проверяем, что лимит не ограничивает запросы
	for i := 1; i <= 2*s.cfg.RateLimit.IP; i++ {
		require.True(s.T(), s.checkRequest(ctx, generateRandStr(15), generateRandStr(15), checkIP), "White list doesn't work")
	}

	// удаляем из белого списка
	s.modifyListRequest(ctx, "192.168.0.0/22", s.gClient.RemoveFromWhiteList)
	time.Sleep(s.cfg.WBSetting.CacheUpdInterval)
	var wasFalse bool
	for i := 1; i <= 2*s.cfg.RateLimit.IP; i++ {
		if !s.checkRequest(ctx, generateRandStr(15), generateRandStr(15), checkIP) {
			wasFalse = true
			break
		}
	}
	require.True(s.T(), wasFalse, "Removing from white list doesn't work")
}

func (s *GrpcAPITestSuite) TestServerBlackList() {
	ctx := context.Background()
	checkIP := "192.168.4.4"

	// добавляем подсеть в черный список
	s.modifyListRequest(ctx, "192.168.4.0/24", s.gClient.AddToBlackList)
	time.Sleep(s.cfg.WBSetting.CacheUpdInterval)
	// проверяем, что черный список не пропускает запросы с IP из подсети
	require.False(s.T(), s.checkRequest(ctx, generateRandStr(15), generateRandStr(15), checkIP), "Black list doesn't work")

	// удаляем подсеть из черного списка
	s.modifyListRequest(ctx, "192.168.4.0/24", s.gClient.RemoveFromBlackList)
	time.Sleep(s.cfg.WBSetting.CacheUpdInterval)
	var n int
	for i := 1; i <= 2*s.cfg.RateLimit.IP; i++ {
		if !s.checkRequest(ctx, generateRandStr(15), generateRandStr(15), checkIP) {
			n = i
			break
		}
	}
	require.GreaterOrEqual(s.T(), n, s.cfg.RateLimit.IP)
}

func (s *GrpcAPITestSuite) TestServerClear() {
	ctx := context.Background()
	// доходим до лимита
	for i := 1; i <= 2*s.cfg.RateLimit.Login; i++ {
		if !s.checkRequest(ctx, "test-clear", generateRandStr(15), generateRandIP()) {
			break
		}
	}
	// очищаем счетчик
	s.clearRequest(ctx, "test-clear")
	// проверяем что очистка произошла корректно и мы имеем снова полный лимит
	for i := 1; i <= s.cfg.RateLimit.Login; i++ {
		require.True(s.T(), s.checkRequest(ctx, "test-clear", generateRandStr(15), generateRandIP()))
	}
}

func (s *GrpcAPITestSuite) checkRequest(ctx context.Context, login, passw, ip string) bool {
	req := &grpcService.CheckRequest{
		Login:    login,
		Password: passw,
		Ip:       ip,
	}
	resp, err := s.gClient.Check(ctx, req)
	require.Nil(s.T(), err)
	res, ok := resp.Result.(*grpcService.CheckResponse_Ok)
	require.True(s.T(), ok)
	return res.Ok
}

func (s *GrpcAPITestSuite) clearRequest(ctx context.Context, login string) {
	req := &grpcService.ClearRequest{
		Login: login,
	}
	resp, err := s.gClient.Clear(ctx, req)
	require.Nil(s.T(), err)
	res, ok := resp.Result.(*grpcService.ModifyResponse_Success)
	require.True(s.T(), ok)
	require.True(s.T(), res.Success)
}

func (s *GrpcAPITestSuite) modifyListRequest(ctx context.Context, ip string, f modifyFunc) {
	req := &grpcService.ModifyListRequest{
		Ip: ip,
	}

	resp, err := f(ctx, req)
	require.Nil(s.T(), err)
	res, ok := resp.Result.(*grpcService.ModifyResponse_Success)
	require.True(s.T(), ok)
	require.True(s.T(), res.Success)
}

func (s *GrpcAPITestSuite) SetupSuite() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config file read error: %v", err)
	}
	s.cfg = cfg
}

func (s *GrpcAPITestSuite) SetupTest() {
	cc, err := grpc.Dial(s.cfg.GrpcServer.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	s.gConn = cc
	s.gClient = grpcService.NewAntiBruteforceServiceClient(cc)
}

func (s *GrpcAPITestSuite) TearDownTest() {
	log.Println("Finishing test")
	err := s.gConn.Close()
	require.NoError(s.T(), err)
	log.Println("Connect closed successfuly")
}
