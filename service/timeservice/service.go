package timeservice

import (
	"github.com/rs/zerolog/log"
	"timeMachine/param"
)

type Repository interface {
	GetTrafficLength() (int32, error)
}

type ML interface {
	GetETAFromML(req *param.ETARequest) *param.ETAResponse
}

type Service struct {
	repo          Repository
	ml            ML
	trafficLength int32
}

func New(repo Repository, ml ML) Service {
	return Service{repo: repo, ml: ml}
}

func (s *Service) SetTrafficLength() error {
	l, e := s.repo.GetTrafficLength()
	if e != nil {
		log.Warn().Msgf("error in getting traffic length: %s", e.Error())
		return e
	}
	s.trafficLength = l
	return nil
}

func (s *Service) GetETA(req *param.ETARequest) *param.ETAResponse {
	req.TrafficLength = s.trafficLength
	eta := s.ml.GetETAFromML(req)
	if eta == nil {
		log.Info().Msgf("%v can't predict eta for this request or eta<500", req)
	} else {
		log.Info().Msgf("%v %d", req, eta)
	}
	return eta
}
