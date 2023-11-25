package ml

import (
	"github.com/dmitryikh/leaves"
	"timeMachine/param"
)

type ML struct {
	modelName  string
	model      *leaves.Ensemble
	coefModel  float64
	coefEngine float64
}

func New(modelName string, coefModel float64, coefEngine float64) ML {
	model, err := leaves.LGEnsembleFromFile(modelName, true)
	if err != nil {
		panic(err)
	}
	return ML{modelName: modelName, model: model, coefModel: coefModel, coefEngine: coefEngine}
}

func (ml ML) GetETAFromML(req *param.ETARequest) *param.ETAResponse {
	etaEngine := float64(req.CurrentETA)
	feature := []float64{
		float64(req.Distance),
		float64(req.Sx),
		float64(req.Sy),
		float64(req.Dx),
		float64(req.Dy),
		float64(req.TrafficLength),
		etaEngine,
		float64(req.Time),
	}
	eta := ml.model.PredictSingle(feature, 0)
	if eta < 500 {
		return nil
	}
	return &param.ETAResponse{ETA: int32(eta*ml.coefModel + etaEngine*ml.coefEngine)}
}
