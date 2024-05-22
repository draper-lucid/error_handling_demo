package main

import (
	"errors"
	"fmt"
)

// callHttp simulates an underlying error. This is returned for all scenarios.
func callHttp() error {
	return fmt.Errorf("document not found")
}

// callUpstreamStruct receives an underlying raw error and then sugars it with contextualized properties.
func callUpstreamStruct() error {
	if err := callHttp(); err != nil {
		return ClientError{"sprocket", "cannot connect", err}
	}
	return nil
}

// callUpstreamPtrError acts similarly to struct, but instead returns a pointer to the error
func callUpstreamPtError() error {
	if err := callHttp(); err != nil {
		return &ClientError{"sprocket", "cannot connect", err}
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

// webServiceStructError simulates a handler calling a provider package getting an error pointer back
func webServiceReturningPtrError() error {
	if err := callUpstreamPtError(); err != nil {

		return err
	}
	return nil
}
func main() {

	if err := webServiceStructError(); err != nil {
		if e, ok := err.(ClientError); ok {
			fmt.Printf("[struct] (type-assertion) an upstream service call failed: %s", e)
		} else if err != nil {
			fmt.Printf("err was not nil, but failed type assertion\n")
		}

		var upstreamError ClientError
		if errors.As(err, &upstreamError) {
			fmt.Printf(
				"[struct] (errors.As) an upstream service call failed: %s", upstreamError.Message,
			)
		}

	}
	fmt.Println("-----")
	if err := webServiceReturningPtrError(); err != nil {
		if e, ok := err.(*ClientError); ok {
			fmt.Printf("[ptr] (type-assertion) an upstream service call failed: %s:  %s\n", e.ServiceName, e.Message)
		} else if err != nil {
			fmt.Printf("[ptr] err was not nil, but failed type assertion\n")
		}

		upstreamError := &ClientError{}
		if errors.As(err, upstreamError) {
			fmt.Printf(
				"([ptr] errors.As) an upstream service call failed: %s: %s", upstreamError.ServiceName,
				upstreamError.Message,
			)
		} else {
			fmt.Printf("[ptr] error is not *ClientError\n")
		}

	}

	fmt.Println("======")
	demoRecommendation()
}
