package timeservice

import (
	"github.com/rs/zerolog/log"
	"timeMachine/param"
)

type Repository interface {
	GetTrafficLength(zone int8) (int32, error)
}

type ML interface {
	GetETAFromML(req *param.ETARequest) *param.ETAResponse
}

type Service struct {
	repo          Repository
	ml            ML
	trafficLength int32
	zone          int8
}

func New(repo Repository, ml ML, zone int8) Service {
	return Service{repo: repo, ml: ml, zone: zone}
}

func (s *Service) SetTrafficLength(zone int8) error {
	l, e := s.repo.GetTrafficLength(zone)
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
	}
	return eta
}
