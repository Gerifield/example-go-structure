package example2

import (
	"context"
	"fmt"
)

type AppInterface interface {
	Index(_ context.Context, _ indexInput) (indexOutput, error)
	Greeting(_ context.Context, in greetingInput) (greetingOutput, error)
}

type application struct {
}

func NewApplication() *application {
	return &application{}
}

type indexInput struct{}
type indexOutput struct {
	message string
}

func (a *application) Index(_ context.Context, _ indexInput) (indexOutput, error) {
	return indexOutput{message: "hello index"}, nil
}

type greetingInput struct {
	name string
}
type greetingOutput struct {
	message string
}

func (a *application) Greeting(_ context.Context, in greetingInput) (greetingOutput, error) {
	return greetingOutput{message: fmt.Sprintf("Hello %s", in.name)}, nil
}
