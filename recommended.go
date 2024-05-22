package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func callDatasource(ctx context.Context) (err error) {
	// Log should be performed once at the point of failure with context containing all relevant information
	// the go-logger will handle most of this for you
	sn := ctx.Value("service-name").(string)
	a := ctx.Value("account_id").(int)
	p := ctx.Value("project_id").(string)

	err = errors.New("document not found")
	fmt.Printf(
		"%s while executing GET http://sprocket/widgets [service: %s account_id: %d project_id: %s]\n",
		err.Error(), sn, a, p,
	)
	return HttpUpstreamError{Status: 404, Message: err.Error(), ServiceName: sn, UrlFragment: "/23/ABCDE/fancywidgets", err: err}
}

func callPackage(ctx context.Context) (err error) {
	ctx = context.WithValue(ctx, "service-name", "sprocket")
	if err = callDatasource(ctx); err != nil {
		// When receiving an error on the callstack, wrap the error for messaging
		// but keep the original error intact. %w wraps this error "around" the other
		// error.
		// Keep wrapping the error down the stack until you either want to render it
		// or re-shape it.
		return fmt.Errorf("sprocket failed getting widgets: %w", err)
	}
	return nil
}

// getWidgets simulates a handler calling a provider package getting an error back
func getWidgets(ctx context.Context) (err error) {

	if err = callPackage(ctx); err != nil {

		// Note to get at the underlying error we can unwrap via the private `err` property
		// on structs matching the error interface
		e := errors.Unwrap(err)
		log.Print(e)

		// At the presentation (endpoint layer) we don't want to expose the actual guts of the system
		ue := &HttpUpstreamError{}
		if errors.As(err, ue) {
			if ue.Status > 399 && ue.Status < 500 {
				// Note we have access to any property that was properly assigned at origin.
				// we present the error the way we need to, relating the information we want the way we want.
				return ClientError{ServiceName: ue.ServiceName, Message: "sorry, couldn't get widgets", err: err}
			}

		}

		// Note that when we are *not* matching the error we **wrap** it, not swallow it
		return fmt.Errorf("failed to get widgets: %w", err)
	}
	return nil
}

// demoRecommendation demonstrates the behavior of errorCatch() middleware at the application/service level.
func demoRecommendation() {

	ctx := context.TODO()
	ctx = context.WithValue(ctx, "account_id", 23)
	ctx = context.WithValue(ctx, "project_id", "ABCDE")

	if err := getWidgets(ctx); err != nil {

		ce := &ClientError{}
		if errors.As(err, ce) {
			b, _ := json.Marshal(ce)
			fmt.Println(string(b))
		}
	}
}
