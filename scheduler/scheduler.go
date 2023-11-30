package scheduler

import (
	"github.com/go-co-op/gocron"
	"sync"
	"time"
	"traffic-prediction-grpc-server/service/timeservice"
)

type Scheduler struct {
	sch        *gocron.Scheduler
	tehranSvc  *timeservice.Service
	mashhadSvc *timeservice.Service
}

func New(tehranSvc *timeservice.Service, mashhadSvc *timeservice.Service) Scheduler {
	return Scheduler{
		sch:        gocron.NewScheduler(time.UTC),
		tehranSvc:  tehranSvc,
		mashhadSvc: mashhadSvc,
	}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	wg.Done()

	s.sch.Every(30).Second().Do(s.tehranSvc.SetTrafficLength(1))
	s.sch.Every(30).Second().Do(s.mashhadSvc.SetTrafficLength(2))
	s.sch.StartAsync()

	<-done
	s.sch.Stop()
}
