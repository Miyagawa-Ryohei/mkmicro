package app

import (
	"context"
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

		msgs := deduplicationMessages(messages)

		if len(msgs) == 0 {
			s.log.Info("message queue is empty, re-polling after 10 second")
			time.Sleep(10 * time.Second)
		}

		s.log.Debug("%d message is received", len(messages))
		s.log.Debug("%d message is processed", len(msgs))
		wg := &sync.WaitGroup{}
		for _, m := range msgs {
			wg.Add(1)
			processor := func() {
				mu := &sync.Mutex{}
				s.log.Debug("start processing message %s", m.GetDeleteID())
				ctx := context.Background()
				ctxChild, cancel := context.WithCancel(ctx)

				go func(queue types.QueueDriver, target types.Message, mu *sync.Mutex){
					for {
						select {
						case <-ctxChild.Done():
							return
						case <-time.After(30 * time.Second)	:
							mu.Lock()
							if !(target.IsDeleted()) {
								if err := queue.ChangeMessageVisibility(target,60); err != nil {
									s.log.Error(err.Error())
								}
							}
							mu.Unlock()
						}
					}
				}(queue, m, mu)

				go func(target types.Message, mu *sync.Mutex) {
					defer func(){
						s.log.Info("worker done [%s]", target.GetID())
						wg.Done()
					}()
					s.log.Debug("worker start [%s]", target.GetID())

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

					if result {
						s.log.Debug("delete message %s", target.GetDeleteID())
						cancel()
						mu.Lock()
						defer mu.Unlock()
						if err := queue.DeleteMessage(target); err != nil {
							s.log.Error("delete message error : [%s]", err.Error())
						} else {
							target.SetDeleted(true)
						}
					}
				}(m, mu)
			}
			processor()
		}
		s.log.Info("wait for done [%s]",processID)
		wg.Wait()
		s.log.Info("done [%s]",processID)
		time.Sleep(1 * time.Second)
	}
}

func deduplicationMessages(message []types.Message) []types.Message{
	if len(message) < 1 {
		return []types.Message{}
	}
	msg := []types.Message{}
	for _, m1 := range message {
		m1Key := m1.GetDeduplicationID()
		find := false
		for _, m2 := range msg {
			m2Key := m2.GetDeduplicationID()
			if m2Key == m1Key{
				find = true
			}
		}
		if !find {
			msg = append(msg,m1)
		}
	}
	return msg
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
