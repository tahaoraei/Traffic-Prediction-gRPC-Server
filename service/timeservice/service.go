package timeservice

import (
	"fmt"
	"github.com/dmitryikh/leaves"
	"github.com/rs/zerolog/log"
	"timeMachine/param"
)

type Repository interface {
	GetTrafficLength() (int32, error)
}

type Service struct {
	repo          Repository
	trafficLength int32
	modelName     string
	model         *leaves.Ensemble
}

func New(modelName string, repo Repository) Service {
	model, err := leaves.LGEnsembleFromFile(modelName, true)
	if err != nil {
		panic(err)
	}
	if model == nil {
		fmt.Println("model is nil")
	} else {
		fmt.Println(model.Name())
	}
	return Service{repo: repo, modelName: modelName, model: model}
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

func (s *Service) GetETA(req param.ETARequest) param.ETAResponse {
	feature := []float64{float64(req.Distance), float64(req.CurrentETA), float64(req.Sx), float64(req.Sy), float64(req.Dx), float64(req.Dy), float64(req.Time), float64(s.trafficLength)}
	eta := s.model.PredictSingle(feature, 0)
	log.Info().Msgf("%+v traffic length %d ETA is: %d", req, s.trafficLength, eta)
	return param.ETAResponse{ETA: int32(eta)}
}
