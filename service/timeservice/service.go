package timeservice

import (
	"github.com/rs/zerolog/log"
	"sync"
	"timeMachine/ml"
	"timeMachine/param"
)

const (
	minXTehran  = 5683551
	minYTehran  = 4216113
	maxXTehran  = 5774664
	maxYTehran  = 4277874
	minXMashhad = 6583521
	minYMashhad = 4309213
	maxXMashhad = 6656388
	maxYMashhad = 4383510
)

type Repository interface {
	GetTrafficLength(zone int8) (int32, error)
}

//type ML interface {
//	GetETAFromML(req *param.ETARequest) *param.ETAResponse
//}

type Service struct {
	tehranML  *ml.ML
	mashhadML *ml.ML
}

func New(tehranML *ml.ML, mashhadML *ml.ML) Service {
	return Service{tehranML: tehranML, mashhadML: mashhadML}
}

func (s *Service) GetETA(req *param.ETARequest) *param.ETAResponse {
	var eta *param.ETAResponse
	var wg sync.WaitGroup

	processRequest := func() {
		defer wg.Done()

		if req.Sx > minXTehran && req.Sy > minYTehran && req.Sx < maxXTehran && req.Sy < maxYTehran &&
			req.Dy > minXTehran && req.Dy > minYTehran && req.Dx < maxXTehran && req.Dy < maxYTehran {
			if req.Time < 360 || req.CurrentETA < 500 {
				eta = &param.ETAResponse{ETA: req.CurrentETA}
			} else {
				eta = s.tehranML.GetETAFromML(req)
			}
		} else if req.Sx > minXMashhad && req.Sy > minYMashhad && req.Sx < maxXMashhad && req.Sy < maxYMashhad &&
			req.Dx > minXMashhad && req.Dy > minYMashhad && req.Dx < maxXMashhad && req.Dy < maxYMashhad {
			if req.Time < 390 || req.CurrentETA < 500 {
				eta = &param.ETAResponse{ETA: req.CurrentETA}
			} else {
				eta = s.mashhadML.GetETAFromML(req)
			}
		} else {
			eta = &param.ETAResponse{ETA: req.CurrentETA}
		}
	}

	wg.Add(1)
	go processRequest()

	wg.Wait()

	if eta == nil {
		log.Info().Msgf("%v can't predict eta for this request or eta<500", req)
	}

	return eta
}
