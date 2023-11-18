package timeservice

import (
	"github.com/dmitryikh/leaves"
	"timeMachine/param"
)

type Service struct {
	modelName string
	model     *leaves.Ensemble
}

func New(modelName string) Service {
	loadTransformation := true
	model, err := leaves.LGEnsembleFromFile(modelName, loadTransformation)
	if err != nil {
		panic(err)
	}
	return Service{modelName: modelName, model: model}
}

func (s Service) GetETA(req param.ETARequest) param.ETAResponse {
	feature := []float64{float64(req.Distance), float64(req.CurrentETA), float64(req.StraightETA), float64(req.Sx), float64(req.Sy), float64(req.Dx), float64(req.Dy), float64(req.Time)}
	eta := s.model.PredictSingle(feature, 0)
	return param.ETAResponse{ETA: int(eta)}
}
