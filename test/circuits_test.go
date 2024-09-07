package test

import (
	"context"
	"testing"

	"kasperaldrin.com/rerego/pkg/core"
)

func TestCircuitBasic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := core.NewCircuit("testCircuit", ctx)

}
