package ml

import (
	"github.com/dmitryikh/leaves"
	"timeMachine/param"
)

type ML struct {
	modelName string
	model     *leaves.Ensemble
}

func New(modelName string) ML {
	model, err := leaves.LGEnsembleFromFile(modelName, true)
	if err != nil {
		panic(err)
	}
	return ML{modelName: modelName, model: model}
}

func (ml ML) GetETAFromML(req *param.ETARequest) *param.ETAResponse {
	feature := []float64{float64(req.Distance), float64(req.Sx), float64(req.Sy), float64(req.Dx), float64(req.Dy), float64(req.TrafficLength), float64(req.CurrentETA), float64(req.Time)}
	eta := ml.model.PredictSingle(feature, 0)
	if eta < 500 {
		return nil
	}
	return &param.ETAResponse{ETA: int32(eta)}
}
