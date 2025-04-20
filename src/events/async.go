package events

import (
	"log"
)

type AsyncEventManager struct {
	Subscribers map[string][]func(interface{}) error
	Events      chan *AsyncEvent
}

type AsyncEvent struct {
	Key   string
	Value interface{}
}

func (manager *AsyncEventManager) Emit(key string, value interface{}) error {
	manager.Events <- &AsyncEvent{
		Key:   key,
		Value: value,
	}
	return nil
}

func (manager *AsyncEventManager) On(key string, foo func(interface{}) error) {
	subscribers, ok := manager.Subscribers[key]
	if !ok {
		subscribers = make([]func(interface{}) error, 0)
	}
	subscribers = append(subscribers, foo)
	manager.Subscribers[key] = subscribers
}

func (manager *AsyncEventManager) Run() error {
	log.Println("[INFO] Event manager has started")
	for event := range manager.Events {
		subscribers, ok := manager.Subscribers[event.Key]
		if !ok {
			return nil
		}
		for _, subscriber := range subscribers {
			subscriber(event.Value)
		}
	}
	log.Println("[STOPPED] Event manager has stopped")
	return nil
}

func NewAsyncEventManager() *AsyncEventManager {
	return &AsyncEventManager{
		Subscribers: make(map[string][]func(interface{}) error),
		Events:      make(chan *AsyncEvent, 1024),
	}
}
