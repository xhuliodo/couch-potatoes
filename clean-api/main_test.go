package main

import (
	"net"
	"net/http"
	"testing"
	"time"
)

type recordingCloser struct {
	closed bool
}

func (c *recordingCloser) Close() error {
	c.closed = true
	return nil
}

func TestShutdownWaitsForInFlightRequestsAndClosesDriver(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %s", err)
	}
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusOK)
	})}
	go func() { _ = server.Serve(listener) }()
	defer server.Close()

	response := make(chan error, 1)
	go func() {
		resp, err := http.Get("http://" + listener.Addr().String() + "/")
		if err == nil {
			resp.Body.Close()
		}
		response <- err
	}()
	<-started

	driver := &recordingCloser{}
	done := make(chan struct{})
	go func() {
		shutdownServer(server, driver)
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("shutdown returned while a request was still in flight")
	case <-time.After(100 * time.Millisecond):
	}

	close(release)
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("shutdown did not return after the in-flight request completed")
	}

	if err := <-response; err != nil {
		t.Errorf("in-flight request failed: %s", err)
	}
	if !driver.closed {
		t.Error("database driver was not closed")
	}
}
