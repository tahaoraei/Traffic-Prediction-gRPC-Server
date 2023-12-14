package scheduler

import (
	"github.com/go-co-op/gocron"
	"sync"
	"time"
	_ "time/tzdata"
	"timeMachine/ml"
	"timeMachine/pkg/logger"
)

var log = logger.Get()

type Config struct {
	TrafficLengthInterval    int `koanf:"trafficLenInterval"`
	ModelChangerInterval     int `koanf:"modelChangerInterval"`
	LoadOnlineConfigInterval int `koanf:"loadOnlineConfig"`
}

type Scheduler struct {
	config    Config
	sch       *gocron.Scheduler
	tehranML  *ml.ML
	mashhadML *ml.ML
}

func New(config Config, tehranML *ml.ML, mashhadML *ml.ML) Scheduler {
	return Scheduler{
		config:    config,
		sch:       gocron.NewScheduler(time.UTC),
		tehranML:  tehranML,
		mashhadML: mashhadML,
	}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	wg.Done()

	s.sch.Every(s.config.TrafficLengthInterval).Second().Do(func() {
		s.tehranML.SetTrafficLength(1)
	})
	s.sch.Every(s.config.TrafficLengthInterval).Second().Do(func() {
		s.mashhadML.SetTrafficLength(2)
	})
	s.sch.Every(s.config.ModelChangerInterval).Hour().Do(func() {
		s.tehranML.SetNewModel()
	})
	s.sch.Every(s.config.ModelChangerInterval).Hour().Do(func() {
		s.mashhadML.SetNewModel()
	})
	s.sch.Every(s.config.LoadOnlineConfigInterval).Minute().Do(func() {
		s.tehranML.SetCoefficentModel("tehran")
	})
	s.sch.Every(s.config.LoadOnlineConfigInterval).Minute().Do(func() {
		s.mashhadML.SetCoefficentModel("mashhad")
	})

	s.sch.StartAsync()

	<-done
	s.sch.Stop()
}
