package core

import (
	"context"
	"time"

	"github.com/google/uuid"
)

/*
	Base.go contain all the base structure for circuits.

	A circuit is really just a list of nodes.

	Circuit, DataObject and Node are the main structures, and all other structures are based on these.
*/

type ConditionalFunction struct {
}

type NodeBase[S, P any] struct {
	_ID           string                          `json:"id"` // used internally
	Label         string                          `json:"label"`
	Subscriptions []<-chan DataObjectInterface[S] `json:"subscriptions"`
	Produces      []chan<- DataObjectInterface[P] `json:"produces"`
	Condition     RUNCONDITION
	Function      func(input ...DataObjectInterface[S]) DataObjectInterface[P]
}

func NewNodeBase[S any, P any](label string) NodeBase[S, P] {
	id := uuid.New().String()
	return NodeBase[S, P]{
		_ID:           id,
		Label:         label,
		Subscriptions: make([]<-chan DataObjectInterface[S], 0),
		Produces:      make([]chan<- DataObjectInterface[P], 0),
		Condition:     ALL,
	}
}

type Node interface {
	GetID() string // Get the ID of the node
}

func (n NodeBase[S, P]) GetID() string {
	return n._ID
}

func (n NodeBase[S, P]) WithInput(
	subscriptions ...Subscription[S]) NodeBase[S, P] {
	for _, sub := range subscriptions {
		n.Subscriptions = append(n.Subscriptions, sub.NewReciever())
	}
	return n
}

func (n NodeBase[S, P]) WithRouter() NodeBase[S, P] {
	return n
}

func (n NodeBase[S, P]) WithOutput(
	subscriptions ...Subscription[P]) NodeBase[S, P] {
	for _, sub := range subscriptions {
		n.Produces = append(n.Produces, sub.NewSender())
	}
	return n
}

/*
func (n NodeBase[S, P]) UseFunction(
	f func(input ...DataObjectInterface[S]) DataObjectInterface[P]) NodeBase[S, P] {
	n.Function = f
	return n
}*/

type RUNCONDITION int

const (
	NONE RUNCONDITION = iota
	ALL
	ANY
)

func (n NodeBase[S, P]) UseFunctionWhenAll(
	f func(input ...DataObjectInterface[S]) DataObjectInterface[P]) NodeBase[S, P] {

	n.Condition = ALL
	n.Function = f

	return n
}

func (n NodeBase[S, P]) UseFunctionWhenAny(
	f func(input DataObjectInterface[S]) DataObjectInterface[P]) NodeBase[S, P] {

	n.Condition = ANY
	n.Function = func(input ...DataObjectInterface[S]) DataObjectInterface[P] {
		return f(input[0])
	}

	return n
}

func (n NodeBase[S, P]) Serve(ctx context.Context) {
	if n.Function == nil {
		panic("Function is not set in NodeBase. Use UseFunction, UseFunctionWhenAll, or UseFunctionWhenAny to set it before calling Serve.")
	}

	ready := make([]DataObjectInterface[S], 0)
	readyCounter := 0
	var processedData DataObjectInterface[P]
	for {
		for _, subscription := range n.Subscriptions {
			select {
			case <-ctx.Done():
				return
			case data := <-subscription:
				ready = append(ready, data)
				readyCounter++

				if n.Condition == ALL {
					if readyCounter == len(n.Subscriptions) {
						processedData = n.Function(ready...)
						for _, produce := range n.Produces {
							//fmt.Printf("\nNode '%s' is sending %v to a channel\n",
							//	n.Label, processedData)
							produce <- processedData
						}
						ready = make([]DataObjectInterface[S], 0)
						readyCounter = 0
					}
				} else if n.Condition == ANY {
					processedData = n.Function(data)
					for _, produce := range n.Produces {
						//fmt.Printf("\nNode '%s' is sending %v to a channel\n",
						//	n.Label, processedData)
						produce <- processedData
					}
				}

				//mt.Printf("\nNode '%s' is working with data %v\n", n.Label, data)
				//fmt.Printf("Node '%s' has created %v",
				//	n.Label, processedData)

			default:
			}
		}
		time.Sleep(time.Millisecond * 10)
	}

}

//func (n NodeBase) Subscribe() {
//	n.Subscriptions = append(n.Subscriptions, <-Subscription[any].Subscribe())
//}
