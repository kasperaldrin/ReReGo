package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"kasperaldrin.com/rerego/pkg/core"
	ml "kasperaldrin.com/rerego/pkg/ml/deep/llm/custom"
	"kasperaldrin.com/rerego/pkg/net"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//true software engineers use a tool for the job that makes an adequate amount of tradeoffs.

	//test.TestRobot(nil)
	//test.TestRouterNInputWaiter(nil)
	// Create the messages

	type Phi3Response struct {
		Text string `json:"generated_text"`
	}

	input := core.NewSubscription[net.WebPostNodeRequest]("query", ctx)
	response := core.NewSubscription[Phi3Response]("response", ctx)

	llm := net.NewWebPostNode[Phi3Response]("llm").
		WithInput(input).
		WithOutput(response)
	go llm.Serve(ctx)

	userInput := input.NewSender()
	userOutput := response.NewReciever()

	go func() {
		messages := []ml.Message{
			{Role: "system", Content: "You are a helpful AI assistant."},
			{Role: "user", Content: "Tell me a joke."},
		}

		// Create the request payload
		requestBody := ml.Phi3Request{
			Messages: messages,
			Phi3GenerationArgs: ml.Phi3GenerationArgs{
				MaxNewTokens:   1500,
				ReturnFullText: false,
				Temperature:    0.7,
				DoSample:       true,
			},
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Printf("Error")
		}

		userInput <- core.NewUserData(net.WebPostNodeRequest{
			URL:         "http://0.0.0.0:8000/generate",
			ContentType: "application/json",
			Body:        bytes.NewBuffer(jsonData),
		})
	}()

	resp := <-userOutput
	if resp.IsError() {
		fmt.Printf("\nError: %s\n", resp.GetError().Error())
	}
	fmt.Printf("\nBot: %s\n", resp.GetData().Text)

	cancel()
	ctx.Done()

}
