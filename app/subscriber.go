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
	s.log.Debug("%d handler is found", len(handlers))
	s.log.Info("start subscribe")
	queue, err := s.src.GetQueue()
	if err != nil {
		s.log.Error("%d handler is found", len(handlers))
		panic(err)
	}

	for {
		messages, err := queue.GetMessage(1)
		if err != nil {
			s.log.Error(err.Error())
			continue
		}
		if len(messages) == 0 {
			s.log.Info("message queue is empty, re-polling after 10 second")
			time.Sleep(10 * time.Second)
		}
		s.log.Debug("%d message is received", len(messages))
		wg := &sync.WaitGroup{}
		for _, m := range messages {
			wg.Add(1)
			ch := make(chan bool)
			go func(target types.Message, ch chan bool) {
				for {
					select {
					case result := <-ch:
						s.log.Debug("[%s]get all worker done signal", target.GetDeleteID())
						if result {
							s.log.Debug("[%s] message deleting...", target.GetDeleteID())
							if err := queue.DeleteMessage(target); err != nil {
								s.log.Error(err.Error())
							}
							s.log.Debug("[%s] message has been deleted", target.GetDeleteID())
						}
						close(ch)
						return
					default:
						if err := queue.ChangeMessageVisibility(target); err != nil {
							s.log.Error(err.Error())
						}
						time.Sleep(40 * time.Second)
					}
				}
			}(m, ch)

			go func(target types.Message, ch chan bool) {
				defer func(){
					err := recover()
					if err != nil {
						s.log.Error("%s",err)
					}
					wg.Done()
				}()
				s.log.Debug("[%s] worker start", target.GetDeleteID())
				result := true
				start := time.Now()
				for _, handler := range handlers {
					if err := handler.Exec(target, s.src); err != nil {
						s.log.Info("[%s]handler returns some error. stop change visibility for retry", target.GetDeleteID())
						s.log.Error(err.Error())
						result = false
					}
					s.log.Info("[%s]all handler returns no errors. message is processed correctly", target.GetDeleteID())
					s.log.Debug("[%s]worker takes %d msec", target.GetDeleteID(), (time.Now().UnixNano()-start.UnixNano()) / int64(time.Millisecond))
				}
				s.log.Debug("[%s]all worker takes %d msec", target.GetDeleteID(), (time.Now().UnixNano()-start.UnixNano()) / int64(time.Millisecond))
				s.log.Info("[%s] all worker end", target.GetDeleteID())
				ch <- result
			}(m, ch)
		}
		s.log.Info("[subscriber main] wait for processing messages")
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
