package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lucidhq/go-logger"
)

var log *logger.Logger

// callHttp simulates an underlying error. This is returned for all scenarios.
func callHttp() error {
	return fmt.Errorf("document not found")
}

// callUpstreamStruct receives an underlying raw error and then sugars it with contextualized properties.
func callUpstreamStruct() error {
	if err := callHttp(); err != nil {
		return ClientError{"cannot connect", err}
	}
	return nil
}

// callUpstreamPtError acts similarly to struct, but instead returns a pointer to the error
func callUpstreamPtError() error {
	if err := callHttp(); err != nil {
		return &ClientError{"cannot connect", err}
	}
	return nil
}

// webServiceStructError simulates a handler calling a provider package getting an error back
func webServiceStructError() error {
	if err := callUpstreamStruct(); err != nil {
		return err
	}
	return nil
}

// webServiceReturningPtrError simulates a handler calling a provider package getting an error pointer back
func webServiceReturningPtrError() error {
	if err := callUpstreamPtError(); err != nil {

		return err
	}
	return nil
}
func main() {
	fmt.Println("\n\n")
	demoStruct()
	fmt.Println("\n-----")
	demoPtr()
	fmt.Println("======")
	demoRecommendation()
}

func demoStruct() {
	if err := webServiceStructError(); err != nil {
		var e ClientError
		if e, ok := err.(ClientError); ok {
			fmt.Printf("[struct] (type-assertion) an upstream service call failed: %s\n", e)
		} else if err != nil {
			fmt.Printf("err was not nil, but failed type assertion\n")
		}
		fmt.Printf("e: %v\n", e)
		var upstreamError ClientError
		if errors.As(err, &upstreamError) {
			fmt.Printf(
				"[struct] (errors.As) an upstream service call failed: %s\n", upstreamError.Message,
			)
		}
		fmt.Printf("upstreamError: %v", upstreamError)

	}
}
func demoPtr() {
	if err := webServiceReturningPtrError(); err != nil {
		e := &ClientError{}
		if e, ok := err.(*ClientError); ok {
			fmt.Printf("[ptr] (type-assertion) an upstream service call failed: %s\n", e.Message)
		} else if err != nil {
			fmt.Printf("[ptr] err was not nil, but failed type assertion\n")
		}
		fmt.Printf("Notice all of e's props are gone: e.Message: %s\n", e.Message)

		upstreamError := &ClientError{}
		if errors.As(err, upstreamError) {
			fmt.Printf(
				"([ptr] errors.As) an upstream service call failed: %s\n",
				upstreamError.Message,
			)
		} else {
			fmt.Printf("[ptr] error is not *ClientError\n")
		}
	}
}

// demoRecommendation demonstrates the behavior of errorCatch() middleware at the application/service level.
func demoRecommendation() {
	var mykey StandardField = logger.LPAccountID
	// Note that these values are set in context when their value is available, which in our case
	// is when this function is run.
	ctx := context.TODO()
	ctx = context.WithValue(ctx, mykey, 23)
	ctx = context.WithValue(ctx, mykey, "hello")
	ctx = context.WithValue(ctx, logger.LPProjectID, "ABCDE")

	if err := getWidgets(ctx); err != nil {

		ce := ClientError{}
		if errors.As(err, &ce) {
			b, _ := json.Marshal(ce)
			fmt.Println(string(b))
		}
	}
}
func init() {
	l, _ := logger.Default()
	log = l
}
