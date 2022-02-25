package app

import (
	"github.com/Miyagawa-Ryohei/mkmicro/infra"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Subscriber struct {
	log       types.Logger
	src       types.SessionManager
	container types.HandlerContainer
}

func (s *Subscriber) Listen(pollingSize int) {

	defer s.log.Flush()

	handlers := s.container.Get()
	s.log.Debug("%d handler is found", len(handlers))
	s.log.Info("start subscribe")
	queue, err := s.src.GetQueue()

	if err != nil {
		s.log.Error("%d handler is found", len(handlers))
		panic(err)
	}

	for {
		messages, err := queue.GetMessage(pollingSize)
		processID := uuid.New().String()
		if err != nil {
			s.log.Error(err.Error())
			time.Sleep(60 * time.Second)
			continue
		}
		s.log.Info("start message processor [%s]",processID)
		if len(messages) == 0 {
			s.log.Info("message queue is empty, re-polling after 10 second")
			time.Sleep(10 * time.Second)
		}
		msgs := []types.Message{}

		for _, m := range messages {
			exist := false
			for _, m2 := range msgs {
				if m.GetDeduplicationID() == m2.GetDeduplicationID() {
					exist = true
				}
			}
			if !exist {
				msgs = append(msgs, m)
			}
		}

		s.log.Debug("%d message is received", len(messages))
		wg := &sync.WaitGroup{}
		mu := &sync.Mutex{}
		for _, m := range msgs {
			wg.Add(1)
			go func(target types.Message) {
				s.log.Debug("start processing message %s", target.GetDeleteID())
				done := new(bool)
				*done = false
				go ChangeMessageVisibility(queue, target, mu, done, s.log)
				go func(target types.Message, mu *sync.Mutex, done *bool) {
					defer wg.Done()
					s.log.Debug("worker start [%s]", target.GetDeleteID())
					result := true
					start := time.Now()
					for _, handler := range handlers {
						if err := handler.Exec(target, s.src); err != nil {
							s.log.Info("[%s]handler returns some error. stop change visibility for retry")
							s.log.Error(err.Error())
							result = false
						} else {
							s.log.Info("all handler returns no errors. message is processed correctly")
						}
						s.log.Debug("worker takes %d msec", (time.Now().UnixNano()-start.UnixNano())/int64(time.Millisecond))
					}
					s.log.Debug("all worker takes %d msec",  (time.Now().UnixNano()-start.UnixNano())/int64(time.Millisecond))
					s.log.Info("all worker end" )
					mu.Lock()
					defer mu.Unlock()
					*done = true
					if result {
						if err := queue.ChangeMessageVisibility(target, 10); err != nil {
							s.log.Error(err.Error())
						}
						s.log.Debug("delete message %s", target.GetDeleteID())
						if err := queue.DeleteMessage(target); err != nil {
							s.log.Error(err.Error())
						} else {
							target.SetDeleted(true)
						}
					}
				}(target, mu, done)
			}(m)
		}
		s.log.Info("wait for done [%s]",processID)
		wg.Wait()
		s.log.Info("done [%s]",processID)
		time.Sleep(1 * time.Second)
	}
}

func NewSubscriber(src types.SessionManager, logger types.Logger, c types.HandlerContainer) *Subscriber {
	log := logger
	if log == nil {
		log = infra.DefaultLogger
	}
	return &Subscriber{
		src:       src,
		log:       log,
		container: c,
	}
}

func ChangeMessageVisibility(queue types.QueueDriver, target types.Message, mu *sync.Mutex, done *bool, log types.Logger) {
	defer mu.Unlock()
	for {
		time.Sleep(40 * time.Second)
		mu.Lock()
		if *done {
			return
		}
		if !(target.IsDeleted()) {
			log.Debug("change message visibility %s", target.GetID())
			if err := queue.ChangeMessageVisibility(target,60); err != nil {
				log.Error(err.Error())
				return
			}
		} else {
			log.Debug("message is deleted %d", target.GetID())
			return
		}
		mu.Unlock()
	}
}
