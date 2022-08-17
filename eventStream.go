package goose

import (
	"fmt"
	"io"
	"net/http"
)

// EventStream a stream to write outgoing SSE data to
type EventStream struct {
	writer http.ResponseWriter
}

// Begin start this SSE event stream by sending the HTTP headers
func (me *EventStream) Begin(stream chan string) error {

	if err := me.writeHeaders(); err != nil {
		return err
	}

	for s := range stream {

		if _, err := me.writer.Write(processData(s)); err != nil {
			return err
		}

		if f, ok := me.writer.(http.Flusher); ok {
			f.Flush()
		}
	}

	return me.Close()
}

func (me *EventStream) writeHeaders() error {

	me.writer.Header().Set("Content-Type", "text/event-stream")
	me.writer.Header().Set("Cache-Control", "no-cache")
	me.writer.Header().Set("Connection", "keep-alive")
	me.writer.WriteHeader(http.StatusOK)
	return nil
}

// Close close this SSE event stream and the underlying http.ResponseWriter
func (me *EventStream) Close() error {

	if f, ok := me.writer.(http.Flusher); ok {
		f.Flush()
	}

	if c, ok := me.writer.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func processData(data string) []byte {

	return ([]byte)(fmt.Sprintf("data: %s\n\n", data))
}

// NewEventStream create and initialize a new SSE event stream
func NewEventStream(w http.ResponseWriter) *EventStream {

	return &EventStream{writer: w}
}
