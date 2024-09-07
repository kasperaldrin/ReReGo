package core

import (
	"context"
	"time"

	"github.com/google/uuid"
)

/*

router
route

share
slot
space
channel
Nexus
Relay
Conduit
Hub
Stream
Pivot
Vector
Node
Portal
Splice
*/

// A router consists of many Subscriptions, and manages the flow of data between them. It acts as a subscription.
// It can take multiple values and output them to multiple values. One value can become many or many can become one. Three can become two, and two can become three, etc.
type Router[T any, O any] struct {
	_ID     string
	Label   string
	Inputs  map[string]<-chan DataObjectInterface[T]
	Outputs map[string]chan<- DataObjectInterface[O]
}

/*
func (r *Router[T]) AddShare(label string) {

}*/

func NewRouter[T any, O any](label string) Router[T, O] {
	return Router[T, O]{
		_ID:     uuid.New().String(),
		Label:   label,
		Inputs:  make(map[string]<-chan DataObjectInterface[T]),
		Outputs: make(map[string]chan<- DataObjectInterface[O]),
	}
}

type Subscription[T any] interface {
	NewReciever() <-chan DataObjectInterface[T]
	NewSender() chan<- DataObjectInterface[T]
	Serve(context.Context)
}

type SubscriptionObject[T any] struct {
	_ID         string
	Label       string
	Sources     []<-chan DataObjectInterface[T]
	Listeners   []chan DataObjectInterface[T]
	AddListener chan chan DataObjectInterface[T]
	AddSource   chan (<-chan DataObjectInterface[T])
}

func (s *SubscriptionObject[T]) NewReciever() <-chan DataObjectInterface[T] {
	c := make(chan DataObjectInterface[T])
	s.AddListener <- c
	return c
}

func (s *SubscriptionObject[T]) NewSender() chan<- DataObjectInterface[T] {
	c := make(chan DataObjectInterface[T])
	s.AddSource <- c
	return c
}

func (s *SubscriptionObject[T]) Serve(ctx context.Context) {
	defer func() {
		for _, listener := range s.Listeners {
			if listener != nil {
				close(listener)
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case newListener := <-s.AddListener:
			//fmt.Printf("\nAdding listener\n")
			s.Listeners = append(s.Listeners, newListener)
		case newSource := <-s.AddSource:
			//fmt.Printf("\nAdding source, recieving from new source\n")
			s.Sources = append(s.Sources, newSource)
		default:
			for _, source := range s.Sources {
				select {
				case value, ok := <-source:
					value.SetId(s._ID)
					value.SetFrom(s.Label)
					//fmt.Printf("\n<sub:%s> recieved value %v", s.Label, value)
					if !ok {
						return
					}
					for _, listener := range s.Listeners {
						if listener != nil {
							//fmt.Printf("\n<sub:%s> sending value %v\n", s.Label, value)
							select {
							case listener <- value:
							case <-ctx.Done():
								return
							}
						}
					}
				default:
					time.Sleep(time.Millisecond * 10) // TODO: remove when stable
				}
			}
		}
	}
}

func NewSubscription[T any](label string, ctx context.Context) Subscription[T] {
	subscription := &SubscriptionObject[T]{
		_ID:         uuid.New().String(),
		Label:       label,
		Sources:     make([]<-chan DataObjectInterface[T], 0),
		Listeners:   make([]chan DataObjectInterface[T], 0),
		AddListener: make(chan chan DataObjectInterface[T]),
		AddSource:   make(chan (<-chan DataObjectInterface[T])),
	}
	go subscription.Serve(ctx)
	return subscription
}
