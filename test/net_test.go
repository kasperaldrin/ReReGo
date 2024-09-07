package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"kasperaldrin.com/rerego/pkg/core"
	"kasperaldrin.com/rerego/pkg/net"
)

// Mock types and core functions/interfaces for testing purposes
type MockDataObject[T any] struct {
	data T
}

func (m MockDataObject[T]) GetData() T {
	return m.data
}

func (m MockDataObject[T]) GetError() error {
	return nil
}

func (m MockDataObject[T]) IsError() bool {
	return false
}

func (m MockDataObject[T]) GetLabel() string {
	return "MockDataObject"
}

func TestNewWebGetNode_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a mock HTTP server
	mockResponse := map[string]string{"key": "value"}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	urlChannel := core.NewSubscription[string]("url", ctx)
	resultOutput := core.NewSubscription[map[string]string]("result", ctx)

	userUrlInput := urlChannel.NewSender()
	responseOutput := resultOutput.NewReciever()

	// Instantiate WebGetNode with the mock server
	node := net.NewWebGetNode[map[string]string]("TestNode").
		WithInput(urlChannel).
		WithOutput(resultOutput)

	go node.Serve(ctx)

	userUrlInput <- core.NewUserData(mockServer.URL)

	result := <-responseOutput

	// Assertions
	assert.False(t, result.IsError())
	assert.Equal(t, http.StatusOK, result.(net.WebResponseObject[map[string]string]).ResponseCode)
	assert.Equal(t, mockResponse, result.GetData())
}

// TODO: Lägg till tester för WebGetNode
/*

func TestNewWebGetNode_HttpError(t *testing.T) {
	// Create a mock HTTP server that returns an error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// Create a MockDataObject with the mock server's URL
	input := MockDataObject[string]{data: mockServer.URL}

	// Instantiate WebGetNode with the mock server
	node := NewWebGetNode[map[string]string]("TestNode")

	// Execute the node's function
	result := node.NodeBase.Execute(input)

	// Assertions
	assert.True(t, result.IsError())
	assert.Equal(t, "500 Internal Server Error", result.GetError().Error())
}

func TestNewWebGetNode_InvalidJson(t *testing.T) {
	// Create a mock HTTP server that returns invalid JSON
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer mockServer.Close()

	// Create a MockDataObject with the mock server's URL
	input := MockDataObject[string]{data: mockServer.URL}

	// Instantiate WebGetNode with the mock server
	node := NewWebGetNode[map[string]string]("TestNode")

	// Execute the node's function
	result := node.NodeBase.Execute(input)

	// Assertions
	assert.True(t, result.IsError())
	assert.Contains(t, result.GetError().Error(), "invalid character")
}

func TestNewWebGetNode_RequestError(t *testing.T) {
	// Create a MockDataObject with an invalid URL
	input := MockDataObject[string]{data: "http://invalid-url"}

	// Instantiate WebGetNode with the invalid URL
	node := NewWebGetNode[map[string]string]("TestNode")

	// Execute the node's function
	result := node.NodeBase.Execute(input)

	// Assertions
	assert.True(t, result.IsError())
	assert.Contains(t, result.GetError().Error(), "no such host")
}

*/
