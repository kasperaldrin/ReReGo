package test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"kasperaldrin.com/rerego/pkg/core"
)

// Testing the Router functionality.

func TestRouterNInputWaiter(t *testing.T) {
	// Waiting for N inputs to arrive.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Divider
	//dividend / divisor
	divident := core.NewSubscription[float64]("divident", ctx)
	divisor := core.NewSubscription[float64]("divisor", ctx)

	type Quotient struct {
		Message string  `json:"message"`
		Value   float64 `json:"value"`
	}
	quotient := core.NewSubscription[Quotient]("quotient", ctx)

	dividentIn := divident.NewSender()
	divisorIn := divisor.NewSender()

	divider := core.
		NewNodeBase[float64, Quotient]("divider").
		WithInput(divident, divisor).
		WithOutput(quotient).
		UseFunctionWhenAll(
			func(input ...core.DataObjectInterface[float64]) core.DataObjectInterface[Quotient] {
				if len(input) != 2 {
					return core.NewErrorData[Quotient](
						errors.New(fmt.Sprintf("Expected 2 inputs, got %d", len(input))))
				}

				var divident float64
				var divisor float64
				for _, in := range input {
					switch in.GetFrom() {
					case "divident":
						divident = in.GetData()
					case "divisor":
						divisor = in.GetData()
					}
				}
				if divisor == 0 {
					return core.NewErrorData[Quotient](
						errors.New(fmt.Sprintf("Division by Zero (%f/%f)", divident, divisor)))
				}
				result := divident / divisor
				return core.NewData("quotient", Quotient{
					Message: fmt.Sprintf("%f divided by %f is %f", divident, divisor, result),
					Value:   result,
				})
			},
		)
	go divider.Serve(ctx)

	go func() {
		for _, i := range [][]float64{
			{2, 4}, {5, 2}, {1, 1}, {1, 0}, {1000, 10},
		} {
			dividentIn <- core.NewUserData(i[0])
			divisorIn <- core.NewUserData(i[1])
			time.Sleep(time.Second / 2)

		}
		//time.Sleep(time.Second)
		//cancel()
	}()

	sumDivided := 0.0
	go func() {
		i := 0
		for q := range quotient.NewReciever() {

			if q.IsError() {
				fmt.Println(q.GetError().Error())
			} else {
				t.Logf("%s\n", q.GetData().Message)
				sumDivided += q.GetData().Value
			}
			if i == 4 {
				cancel()
				return
			}
			i++
		}
	}()

	//toBeDivided := core.Router[float64]{
	//}
	<-ctx.Done()

	// 2/4 =
	desired := 2.0/4.0 + 5.0/2.0 + 1.0 + 1000.0/10.0
	if sumDivided != desired {
		t.Errorf("Expected %f, got %f", desired, sumDivided)
	}

}

/*

func TestMultitypeSource(t *testing.T){
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()



}
*/

/*
func TestRouterVariousTypesWaiter(t *testing.T) {
	// Waiting for various types of inputs to arrive.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

}*/
