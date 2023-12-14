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

type Scheduler struct {
	sch       *gocron.Scheduler
	tehranML  *ml.ML
	mashhadML *ml.ML
}

func New(tehranML *ml.ML, mashhadML *ml.ML) Scheduler {
	return Scheduler{
		sch:       gocron.NewScheduler(time.UTC),
		tehranML:  tehranML,
		mashhadML: mashhadML,
	}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	wg.Done()

	s.sch.Every(30).Second().Do(func() {
		s.tehranML.SetTrafficLength(1)
	})
	s.sch.Every(30).Second().Do(func() {
		s.mashhadML.SetTrafficLength(2)
	})
	s.sch.Every(3).Hour().Do(func() {
		s.tehranML.SetNewModel()
	})
	s.sch.Every(3).Hour().Do(func() {
		s.mashhadML.SetNewModel()
	})

	s.sch.StartAsync()

	<-done
	s.sch.Stop()
}
