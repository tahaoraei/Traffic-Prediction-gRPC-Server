package scheduler

import (
	"github.com/go-co-op/gocron"
	"sync"
	"time"
	"timeMachine/service/timeservice"
)

type Scheduler struct {
	sch *gocron.Scheduler
	svc *timeservice.Service
}

func New(svc *timeservice.Service) Scheduler {
	return Scheduler{
		sch: gocron.NewScheduler(time.UTC),
		svc: svc,
	}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	wg.Done()

	s.sch.Every(30).Second().Do(s.svc.SetTrafficLength())
	s.sch.StartAsync()

	<-done
	s.sch.Stop()
}
