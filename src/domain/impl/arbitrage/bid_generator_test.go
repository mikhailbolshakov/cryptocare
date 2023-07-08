package arbitrage

import (
	kitTestSuite "github.com/mikhailbolshakov/cryptocare/src/kit/test/suite"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type bidGenTestSuite struct {
	kitTestSuite.Suite
	svc *bidGeneratorImpl
}

func (s *bidGenTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
}

func TestGenTestSuite(t *testing.T) {
	suite.Run(t, new(bidGenTestSuite))
}

func (s *bidGenTestSuite) SetupTest() {
	s.svc = NewBidGenerator(nil).(*bidGeneratorImpl)
}

func (s *bidGenTestSuite) Test_GenBin() {
	bin := s.svc.getBid()
	s.NotEmpty(bin)
	s.NotEmpty(bin.Id)
	s.NotEmpty(bin.TrgAsset)
	s.NotEmpty(bin.SrcAsset)
	s.NotEmpty(bin.Rate)
	s.NotEmpty(bin.MinLimit)
	s.NotEmpty(bin.MaxLimit)
	s.NotEmpty(bin.Available)
	s.Greater(bin.MaxLimit, bin.MinLimit)
	s.Greater(bin.MaxLimit*bin.Rate, bin.Available)
	s.Greater(bin.Available, bin.MinLimit*bin.Rate)
}
