package ml

import (
	"testing"
	"timeMachine/param"
)

func TestNew(t *testing.T) {
	modelName := "../mashhad-20231125.txt"
	coefModel := 0.
	coefEngine := 1.

	ml, err := New(modelName, coefModel, coefEngine)
	if err != nil {
		t.Errorf("New() error=%v, want nil", err)
	}

	if ml.modelName != modelName {
		t.Errorf("New() modelName=%s, want %s", ml.modelName, modelName)
	}

	if ml.coefModel != coefModel {
		t.Errorf("New() coefModel = %f, want %f", ml.coefModel, coefModel)
	}

	if ml.coefEngine != coefEngine {
		t.Errorf("New() coefEngine = %f, want %f", ml.coefEngine, coefEngine)
	}
}

func TestML_GetETAFromML(t *testing.T) {
	modelName := "../mashhad-20231125.txt"
	coefModel := 0.
	coefEngine := 1.

	ml, err := New(modelName, coefModel, coefEngine)
	if err != nil {
		t.Fatalf("New() error=%v, want nil", err)
	}

	req := &param.ETARequest{
		Distance:      10,
		Sx:            0,
		Sy:            0,
		Dx:            20,
		Dy:            20,
		TrafficLength: 5,
		CurrentETA:    100,
		Time:          10,
	}

	etaResponse := ml.GetETAFromML(req)
	if etaResponse == nil {
		t.Fatalf("GetETAFromML() is nil, want %d", req.CurrentETA)
	}
}
