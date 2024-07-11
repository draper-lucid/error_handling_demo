package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/lucidhq/go-logger"
)

// callDatasource represents upstream.Get(), eg. the generic http layer.
func callDatasource(ctx context.Context) (err error) {

	// Imagine this is an error returned from the http package in the go standard library
	err = errors.New("document not found")

	msg := fmt.Sprintf("%s while calling sprocket GET http://sprocket/widgets", "error")

	var mykey StandardField = logger.LPAccountID
	// This simple logging call may include any or all canonical fields for logging automatically when my PR is merged.
	//log.Errorw(ctx, msg)

	// When go-logger is updated, this will no longer be necessary since LPAccountID is a standard field.
	// For now I want to show what a standard field looks like in the logs as opposed to metadata.
	log.Errorw(ctx, msg, mykey, 23)

	// If we need data to be present in structured logging, but not searchable, any keys that are not canonical
	// will be stored in the metadata field.
	// Eg.
	log.Errorw(ctx, msg, mykey, 23, "service_name", "sprocket")

	return HttpUpstreamError{Status: 404, Message: err.Error(), UrlFragment: "/23/ABCDE/fancywidgets", err: err}
}

// callPackage represents the business-logic layer that bridges the handler and the upstream data sources.
// this function returns a wrapped error that was returned from the http upstream
func callPackage(ctx context.Context) (err error) {

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

		// At the presentation (endpoint layer) we don't want to expose the actual guts of the system
		// For this demo I wanted to show that we can cast down to the underlying error type and check for it.
		// Optionally in callPackage() we could wrap the upstream error as a PackageError{} to further define/abstract it,
		//   then check for PackageError here and decide what to render.
		ue := &HttpUpstreamError{}
		if errors.As(err, ue) {
			if ue.Status > 399 && ue.Status < 500 {
				// Note we have access to any property that was properly assigned at origin.
				// we present the error the way we need to, relating the information we want the way we want.
				return ClientError{Message: "sorry, couldn't get widgets", err: err}
			}

		}

		// Note that when we are *not* matching the error we **wrap** it, not swallow it
		return fmt.Errorf("failed to get widgets: %w", err)
	}
	return nil
}

type StandardField string
