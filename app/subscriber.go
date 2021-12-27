package app

import (
	"fmt"
	"github.com/Miyagawa-Ryohei/mkmicro/entity"
	"github.com/Miyagawa-Ryohei/mkmicro/repository"
	"log"
	"sync"
)

type Subscriber struct {
	src entity.SessionManager
	dist []entity.SessionManager
	factory entity.SessionManagerFactory
}

func (s *Subscriber) Listen() {

	repo := repository.GetHandlerRepository()
	handlers := repo.Get()
	log.Print(fmt.Sprintf("%d handler is found", len(handlers)))
	log.Print("start subscribe")
	queue, err := s.src.GetQueue()
	if err != nil {
		log.Fatal(err)
	}

	for _, handler := range handlers {
		queueConfig := handler.GetResultQueueConfig()
		sessionConfig := handler.GetResultStorageConfig()
		if queueConfig != nil && sessionConfig != nil{
			f, err := s.factory.CreateWithConfig(*queueConfig,*sessionConfig)
			if err != nil {
				log.Fatal(err)
			}
			s.dist = append(s.dist, f)
		} else {
			s.dist = append(s.dist, s.src)
		}
	}

	for {
		messages, err := queue.GetMessage(1)
		if err != nil {
			log.Print(err)
			continue
		}
		if len(messages) == 0 {
			log.Print("empty message is received")
		}
		wg := &sync.WaitGroup{}
		for _, m := range messages {
			wg.Add(1)

			go func(target entity.Message) {

				result := true

				for index, handler := range handlers {
					if err := handler.Exec(target, s.dist[index]); err != nil {
						result = false
					}
				}
				if result {
					_ = queue.DeleteMessage(target)
				}

				wg.Done()
			}(m)
		}
		wg.Wait()
	}
}

func (s *Subscriber) addSession(){

}

func (s *Subscriber) updateSession(){
	for _, sess := range s.dist{
		sess.UpdateSession()
	}
}

func (s *Subscriber) PushResultMessage(result []byte){

}
func (s *Subscriber) PutResultFile(name string, root string, data []byte){

}


func NewSubscriber(src entity.SessionManager, factory entity.SessionManagerFactory) *Subscriber {
	return &Subscriber{
		src: src,
		dist: []entity.SessionManager{},
		factory: factory,
	}
}
