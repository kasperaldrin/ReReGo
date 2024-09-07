package net

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"kasperaldrin.com/rerego/pkg/core"
)

type APICallerNode struct {
	core.NodeBase[http.Request, http.Response]
	// Recieves APICall definition and returns WebData
}

// T is the return type of the GET call.
type WebGetNode[T any] struct {
	core.NodeBase[string, T]
}

type WebResponseObject[T any] struct {
	core.DataObjectInterface[T]
	ResponseCode int
	Response     http.Response
}

func NewWebGetNode[T any](label string) WebGetNode[T] {
	return WebGetNode[T]{
		NodeBase: core.NewNodeBase[string, T](label).
			UseFunctionWhenAny(func(
				input core.DataObjectInterface[string],
			) core.DataObjectInterface[T] {
				response, err := http.Get(input.GetData())
				if err != nil {
					return core.NewErrorData[T](err)
				}
				defer response.Body.Close()
				var responseObject T
				body, err := io.ReadAll(response.Body)
				if err != nil {
					return core.NewErrorData[T](err)
				}
				err = json.Unmarshal(body, &responseObject)
				if err != nil {
					return core.NewErrorData[T](err)
				}
				nd := core.NewData("GetOutput", responseObject)
				wr := WebResponseObject[T]{
					DataObjectInterface: nd,
					ResponseCode:        response.StatusCode,
					Response:            *response,
				}
				return wr
			},
			),
	}
}
func (n WebGetNode[T]) UseFunction() {} // Cannot replace default function.

type WebPostNodeRequest struct {
	URL         string
	ContentType string
	Body        io.Reader
}

type WebPostNode[T any] struct {
	core.NodeBase[WebPostNodeRequest, T]
}

func NewWebPostNode[T any](label string) WebPostNode[T] {
	return WebPostNode[T]{
		NodeBase: core.NewNodeBase[WebPostNodeRequest, T](label).
			UseFunctionWhenAny(func(
				input core.DataObjectInterface[WebPostNodeRequest],
			) core.DataObjectInterface[T] {
				response, err := http.Post(input.GetData().URL, input.GetData().ContentType, input.GetData().Body)
				if err != nil {
					fmt.Printf("78: %s\n", err)
					return core.NewErrorData[T](err)
				}
				defer response.Body.Close()
				var responseObject T
				body, err := io.ReadAll(response.Body)
				if err != nil {
					fmt.Printf("66: %s\n", err)
					return core.NewErrorData[T](err)
				}
				err = json.Unmarshal(body, &responseObject)
				if err != nil {
					fmt.Printf("91: %s\n", err)
					return core.NewErrorData[T](err)
				}
				nd := core.NewData("GetOutput", responseObject)
				wr := WebResponseObject[T]{
					DataObjectInterface: nd,
					ResponseCode:        response.StatusCode,
					Response:            *response,
				}
				return wr
			},
			),
	}
}

/*
func NewAPICallerNode(label string, request http.Request) APICallerNode {
	return APICallerNode{
		NodeBase: NodeBase[http.Request, http.Response]{
			// Initialize the embedded NodeBase struct directly
			_ID:           uuid.New().String(),
			Label:         label,
			Subscriptions: make([]<-chan DataObject[http.Request], 0),
			Produces:      make([]chan<- DataObject[http.Response], 0),
			Function: func(input DataObject[http.Request]) DataObject[http.Response] {
				// Use a type switch to handle the conversion properly
				return nil
			},
		},
	}
}

func (n APICallerNode[T]) UseFunction() {} // Cannot replace default function.

type WebData interface {
	JSON | HTML
}

type JSON struct {
	Text string
}

type HTML struct {
	DataObject[string]
}
*/
