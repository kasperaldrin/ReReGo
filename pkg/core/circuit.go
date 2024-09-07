package core

import (
	"context"
	"reflect"
)

/*
	Circuits are programs consisting of BaseNodes and Subscriptions(Shares/Routers), including DataObjects.

	This is the second Level abstraction of ReReGo, and is used to define the structure of a program or an agent.

	The third level of abstraction are fluids of order N, where N is the depth of circuits within the fluid.
*/

type CircuitSubscription struct {
	Label string
	Type  any
}

type CircuitNode struct {
	Label       string
	InputLabel  []string
	OutputLabel []string
	Function    func(input ...DataObjectInterface[any]) DataObjectInterface[any]
}

type Circuit struct {
	Label         string                `json:"label"`
	Subscriptions []CircuitSubscription `json:"subscriptions"`
	Nodes         []CircuitNode         `json:"nodes"`
}

func NewCircuit(label string, ctx context.Context) *Circuit {
	newCircuit := Circuit{
		Label:         label,
		Subscriptions: make([]CircuitSubscription, 0),
		Nodes:         make([]CircuitNode, 0),
	}
	newCircuit.Serve(ctx)
	return &newCircuit
}

func (c *Circuit) Serve(ctx context.Context) {
	cirquitState := CircuitInner{
		Label:         c.Label,
		Subscriptions: make(map[string]Subscription[any]),
		Nodes:         make(map[string]NodeBase[any, any]),
	}
	for _, sub := range c.Subscriptions {
		cirquitState.Subscriptions[sub.Label] = NewSubscription[any](sub.Label, ctx)
	}
	for _, node := range c.Nodes {
		cirquitState.Nodes[node.Label] = NewNodeBase[any, any](node.Label)
		for _, inLabel := range node.InputLabel {
			cirquitState.Nodes[node.Label].WithInput(cirquitState.Subscriptions[inLabel])
		}
		for _, outLabel := range node.OutputLabel {
			cirquitState.Nodes[node.Label].WithOutput(cirquitState.Subscriptions[outLabel])
		}

	}

}

func (c Circuit) WithSubscription(label string, t reflect.Type) Circuit {
	c.Subscriptions = append(c.Subscriptions, CircuitSubscription{
		Label: label,
		Type:  t,
	})
	return c
}

func (c Circuit) WithNode(
	label string,
	inType reflect.Type,
	outType reflect.Type,
	f func(...DataObjectInterface[reflect.Type]) DataObjectInterface[reflect.Type]) Circuit {
	c.Nodes = append(c.Nodes, CircuitNode{
		Label:      label,
		InputType:  inType,
		OutputType: outType,
		Function:   f,
	})
	return c
}

func (c *Circuit) AddSubscription(label string, t reflect.Type) {
	c.Subscriptions = append(c.Subscriptions, CircuitSubscription{
		Label: label,
		Type:  t,
	})
}

func (c *Circuit) SetNode(
	label string,
	inType reflect.Type,
	outType reflect.Type,
	f func(...DataObjectInterface[reflect.Type]) DataObjectInterface[reflect.Type]) {
	c.Nodes = append(c.Nodes, CircuitNode{
		Label:      label,
		InputType:  inType,
		OutputType: outType,
		Function:   f,
	})
}

type CircuitInner struct {
	Label         string                        `json:"label"`
	Subscriptions map[string]Subscription[any]  `json:"subscriptions"`
	Nodes         map[string]NodeBase[any, any] `json:"nodes"`
}

type RepeaterNode struct {
}
