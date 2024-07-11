package main

import "fmt"

type ClientError struct {
	Message string `json:"message"`
	err     error
}

func (g ClientError) Error() string {
	return "bad request: " + g.Message
}

type HttpUpstreamError struct {
	Status      int
	Message     string
	UrlFragment string
	err         error
}

func (h HttpUpstreamError) Error() string {
	return fmt.Sprintf("%s [status: %d]in : %s", h.err.Error(), h.Status, h.UrlFragment)
}
