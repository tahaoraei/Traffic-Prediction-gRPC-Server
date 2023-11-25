package ml

import (
	"fmt"
	"github.com/dmitryikh/leaves"
	"timeMachine/param"
)

type ML struct {
	modelName  string
	model      *leaves.Ensemble
	coefModel  float64
	coefEngine float64
}

func New(modelName string, coefModel float64, coefEngine float64) (*ML, error) {
	model, err := leaves.LGEnsembleFromFile(modelName, true)
	if err != nil {
		return nil, fmt.Errorf("faild to load model: %v", err)
	}
	return &ML{modelName: modelName, model: model, coefModel: coefModel, coefEngine: coefEngine}, nil
}

func (ml *ML) GetETAFromML(req *param.ETARequest) *param.ETAResponse {
	etaEngine := float64(req.CurrentETA)
	features := []float64{
		float64(req.Distance),
		float64(req.Sx),
		float64(req.Sy),
		float64(req.Dx),
		float64(req.Dy),
		float64(req.TrafficLength),
		etaEngine,
		float64(req.Time),
	}
	eta := ml.model.PredictSingle(features, 0)
	if eta < 500 {
		return nil
	}
	return &param.ETAResponse{ETA: int32(eta*ml.coefModel + etaEngine*ml.coefEngine)}
}
