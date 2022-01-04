package app

import (
	"github.com/Miyagawa-Ryohei/mkmicro/container"
	"github.com/Miyagawa-Ryohei/mkmicro/infra"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"sync"
	"time"
)

type Subscriber struct {
	log types.Logger
	src types.SessionManager
}

func (s *Subscriber) Listen() {

	defer s.log.Flush()

	c := container.GetHandlerContainer()
	handlers := c.Get()
	s.log.Debugf("%d handler is found", len(handlers))
	s.log.Info("start subscribe")
	queue, err := s.src.GetQueue()
	if err != nil {
		s.log.Errorf("%d handler is found", len(handlers))
		panic(err)
	}

	for {
		messages, err := queue.GetMessage(1)
		if err != nil {
			s.log.Error(err.Error())
			continue
		}
		if len(messages) == 0 {
			s.log.Info("message queue is empty")
			time.Sleep(10 * time.Second)
		}
		s.log.Debugf("%d message is received", len(messages))
		wg := &sync.WaitGroup{}
		for _, m := range messages {
			wg.Add(1)
			ch := make(chan bool)
			go func(target types.Message, ch chan bool) {
				for {
					select {
					case result := <-ch:
						if result {
							if err := queue.DeleteMessage(target); err != nil {
								s.log.Error(err.Error())
							}
						}
						close(ch)
						return
					default:
						if err := queue.ChangeMessageVisibility(target); err != nil {
							s.log.Error(err.Error())
						}
					}
				}
			}(m, ch)

			go func(target types.Message, ch chan bool) {
				s.log.Debugf("[%s] worker start", target.GetDeleteID())
				result := true
				for _, handler := range handlers {
					if err := handler.Exec(target, s.src); err != nil {
						s.log.Error(err.Error())
						result = false
					}
				}
				ch <- result
				s.log.Debugf("[%s] all worker end", target.GetDeleteID())
				wg.Done()
			}(m, ch)
		}
		wg.Wait()
	}
}

func NewSubscriber(src types.SessionManager, logger types.Logger) *Subscriber {
	log := logger
	if log == nil {
		log = infra.DefaultLogger
	}
	return &Subscriber{
		src: src,
		log: log,
	}
}
