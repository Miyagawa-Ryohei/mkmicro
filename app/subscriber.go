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
	msgChan   chan types.Message
	wg		  *sync.WaitGroup
}

type ProcessManager struct {
	queue types.QueueDriver
	handlers []types.Handler
	log types.Logger
	mu *sync.Mutex
	wg *sync.WaitGroup
	sess types.SessionManager
}

func (p *ProcessManager) changeMessageVisibility(target types.Message, ctx context.Context){
	for {
		select {
		case <-ctx.Done():
			p.log.Info("[%s] worker is done. stop update visibility timeout", target.GetChangeVisibilityID())
			return
		case <-time.After(20 * time.Second)	:
			p.mu.Lock()
			if !(target.IsDeleted()) {
				p.log.Debug("[%s] update visibility", target.GetChangeVisibilityID())
				if err := p.queue.ChangeMessageVisibility(target,60); err != nil {
					p.log.Error(err.Error())
				}
			}
			p.mu.Unlock()
		}
	}
}

func (p *ProcessManager)runWorker(target types.Message,cancel context.CancelFunc) {

	defer func(){
		p.log.Info("[%s] worker done", target.GetDeduplicationID())
		p.wg.Done()
	}()
	p.log.Debug("[%s] worker start", target.GetDeduplicationID())

	result := true
	start := time.Now()

	for _, handler := range p.handlers {
		if err := handler.Exec(target, p.sess); err != nil {
			p.log.Info("[%s] handler returns some error. stop change visibility for retry", target.GetDeduplicationID())
			p.log.Error(err.Error())
			result = false
		} else {
			p.log.Info("[%s] handler returns no errors. message is processed correctly", target.GetDeduplicationID())
		}
		p.log.Debug("worker takes %d msec", (time.Now().UnixNano()-start.UnixNano())/int64(time.Millisecond))
	}

	p.log.Debug("all worker takes %d msec",  (time.Now().UnixNano()-start.UnixNano())/int64(time.Millisecond))
	p.log.Info("[%s] all worker end", target.GetDeduplicationID())

	cancel()
	if result {
		p.log.Debug("[%s] delete message", target.GetDeduplicationID())
		p.mu.Lock()
		defer p.mu.Unlock()
		if err := p.queue.DeleteMessage(target); err != nil {
			p.log.Error("[%s] delete message error : [%s]", target.GetDeduplicationID(), err.Error())
		}
	}
}

func (s *Subscriber) Exec(queue types.QueueDriver,handlers []types.Handler) {
	p := ProcessManager{
		queue :queue,
		handlers: handlers,
		mu: &sync.Mutex{},
		log: s.log,
		wg: s.wg,
		sess: s.src,
	}
	for{
		msg := <- s.msgChan
		s.log.Debug("[%s] start processing message", msg.GetDeduplicationID())
		ctx := context.Background()
		ctxChild, cancel := context.WithCancel(ctx)

		go p.changeMessageVisibility(msg,ctxChild)
		p.runWorker(msg,cancel)
	}
}

func (s *Subscriber) Listen(pollingSize int) {

	defer s.log.Flush()

	h := s.container.Get()
	s.log.Debug("%d handler is found", len(h))
	s.log.Info("start subscribe")
	q, err := s.src.GetQueue()

	if err != nil {
		s.log.Error("cannot get queue")
		panic(err)
	}
	for i := 0; i < 10; i++ {
		go s.Exec(q, h)
	}

	for {
		messages, err := q.GetMessage(pollingSize)
		processID := uuid.New().String()
		if err != nil {
			s.log.Error(err.Error())
			time.Sleep(60 * time.Second)
			continue
		}
		msgs := deduplicationMessages(messages)

		if len(msgs) == 0 {
			s.log.Info("[%s] message queue is empty, re-polling after 10 second", processID)
			time.Sleep(10 * time.Second)
			continue
		}
		s.log.Info("[%s] start message processor", processID)
		s.log.Debug("[%s] %d message is received", processID, len(messages))
		s.log.Debug("[%s] %d message is processed", processID, len(msgs))
		for _, m := range msgs {
			s.wg.Add(1)
			s.msgChan <- m
		}
		s.log.Info("[%s] wait for done",processID)
		s.wg.Wait()
		s.log.Info("[%s] done message processor",processID)
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
		msgChan:   make(chan types.Message),
		wg: 		&sync.WaitGroup{},
		log:       log,
		container: c,
	}
}
