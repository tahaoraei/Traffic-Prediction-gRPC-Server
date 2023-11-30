package timeservice

import (
	"errors"
	"testing"
	"timeMachine/param"
)

type mockRepo struct{}

func (repo *mockRepo) GetTrafficLength(zone int8) (int32, error) {
	if zone == 1 {
		return 10, nil
	}
	return 0, errors.New("zone not found")
}

type mockML struct{}

func (ml *mockML) GetETAFromML(req *param.ETARequest) *param.ETAResponse {
	if req.Distance == 10 {
		return &param.ETAResponse{ETA: 20}
	}
	return nil
}

func TestNew(t *testing.T) {
	repo := &mockRepo{}
	ml := &mockML{}
	zone := int8(1)

	service := New(repo, ml, zone)

	if service.zone != zone {
		t.Errorf("New() zone = %d, want %d", service.zone, zone)
	}
}

func TestService_SetTrafficLength(t *testing.T) {
	repo := &mockRepo{}
	ml := &mockML{}
	zone := int8(1)

	service := New(repo, ml, zone)

	err := service.SetTrafficLength(zone)
	if err != nil {
		t.Errorf("SetTrafficLength() error = %v, want nil", err)
	}

	if service.trafficLength != 10 {
		t.Errorf("SetTrafficLength() trafficLength = %d, want 10", service.trafficLength)
	}

	// Test with a zone not found in the repo
	zone = int8(2)
	service = New(repo, ml, zone)

	err = service.SetTrafficLength(zone)
	if err == nil {
		t.Error("SetTrafficLength() returned nil error, want non-nil error")
	}
}

func TestService_GetETA(t *testing.T) {
	repo := &mockRepo{}
	ml := &mockML{}
	zone := int8(1)

	service := New(repo, ml, zone)

	req := &param.ETARequest{
		Distance: 5,
	}

	eta := service.GetETA(req)
	if eta != nil {
		t.Errorf("GetETA() eta = %v, want nil", eta)
	}

	req.Distance = 10

	eta = service.GetETA(req)

	if eta == nil {
		t.Error("GetETA() returned nil eta, want non-nil eta")
	}

	if eta.ETA != 20 {
		t.Errorf("GetETA() ETA = %d, want 20", eta.ETA)
	}
}
