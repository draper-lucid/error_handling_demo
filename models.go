package main

import "fmt"

type ClientError struct {
	ServiceName string `json:"service"`
	Message     string `json:"message"`
	err         error
}

func (g ClientError) Error() string {
	return "bad request: " + g.Message
}

type HttpUpstreamError struct {
	ServiceName string
	Status      int
	Message     string
	UrlFragment string
	err         error
}

func (h HttpUpstreamError) Error() string {
	return fmt.Sprintf("%s in %s [%d]: %s", h.err.Error(), h.ServiceName, h.Status, h.UrlFragment)
}
