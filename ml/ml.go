package ml

import (
	"fmt"
	"github.com/dmitryikh/leaves"
	"timeMachine/param"
	"timeMachine/pkg/logger"
	"timeMachine/pkg/util"
)

const (
	TehranWorkday   = "ml/model/tehran-workdays-20231213.txt"
	TehranThursday  = "ml/model/tehran-thursday-20231213.txt"
	TehranFriday    = "ml/model/tehran-friday-20231213.txt"
	MashhadWorkday  = "ml/model/mashhad-workdays-20231213.txt"
	MashhadThursday = "ml/model/mashhad-thursday-20231213.txt"
	MashhadFriday   = "ml/model/mashhad-friday-20231213.txt"
)

var log = logger.Get()

type Repository interface {
	GetTrafficLength(zone int8) (int32, error)
	GetOnlineConfig(city string) (float64, error)
}

type ML struct {
	repo          Repository
	modelCity     string
	model         *leaves.Ensemble
	coefModel     float64
	trafficLength int32
	zone          int8
}

func New(repo Repository, modelCity string, zone int8, coefModel float64) (*ML, error) {
	m := ML{repo: repo, modelCity: modelCity, zone: zone, model: nil, coefModel: coefModel}
	err := m.SetNewModel()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (ml *ML) SetTrafficLength(zone int8) error {
	l, e := ml.repo.GetTrafficLength(zone)
	if e != nil {
		log.Warn().Msgf("error in getting traffic length: %s", e.Error())
		return e
	}
	ml.trafficLength = l
	log.Info().Msgf("SetTrafficLength is: %d", l)
	return nil
}

func (ml *ML) SetNewModel() error {
	weekday := util.Weekday("Asia/Tehran")
	var (
		m   *leaves.Ensemble
		err error
	)

	// TODO: fix this shitty code
	if ml.modelCity == "tehran" {
		switch weekday {
		case 4:
			m, err = leaves.LGEnsembleFromFile(TehranThursday, true)
		case 5:
			m, err = leaves.LGEnsembleFromFile(TehranFriday, true)
		default:
			m, err = leaves.LGEnsembleFromFile(TehranWorkday, true)
		}
	} else if ml.modelCity == "mashhad" {
		switch weekday {
		case 4:
			m, err = leaves.LGEnsembleFromFile(MashhadThursday, true)
		case 5:
			m, err = leaves.LGEnsembleFromFile(MashhadFriday, true)
		default:
			m, err = leaves.LGEnsembleFromFile(MashhadWorkday, true)
		}
	} else {
		return fmt.Errorf("wrong city name for loading model")
	}

	ml.model = m
	if err != nil {
		return fmt.Errorf("faild to load model for weekday %d: %v", weekday, err)
	}
	log.Info().Msgf("Change model for city %s to weekday number %d", ml.modelCity, weekday)
	return nil
}

func (ml *ML) SetCoefficentModel(city string) error {
	coef, err := ml.repo.GetOnlineConfig(city)
	if err != nil {
		log.Warn().Msgf("can't get online config")
		return err
	}
	ml.coefModel = coef
	return nil
}

func (ml *ML) GetETAFromML(req *param.ETARequest) *param.ETAResponse {
	req.TrafficLength = ml.trafficLength

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

	return &param.ETAResponse{ETA: int32(eta*ml.coefModel + etaEngine*(1.-ml.coefModel))}
}
