package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stvp/tempredis"

	"github.com/declantraynor/go-events-service/interfaces/web"
)

func startRedis(port string) *tempredis.Server {
	server, err := tempredis.Start(
		tempredis.Config{
			"port": port,
		},
	)

	if err != nil {
		log.Fatal("Unable to start tempredis for test")
	}

	return server
}

func stopRedis(server *tempredis.Server) {
	err := server.Kill()
	if err != nil {
		log.Fatal("Problem killing tempredis server during test")
	}
}

type TestServer struct {
	server *httptest.Server
}

func (t *TestServer) serveCreate(webservice *web.WebService) {
	t.server = httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			webservice.Create(res, req)
		}))
}

func TestUnableToConnectToRedis(t *testing.T) {
	os.Setenv("REDIS_PORT_6379_TCP_ADDR", "127.0.0.1")
	os.Setenv("REDIS_PORT_6379_TCP_PORT", "12313")

	testserver := TestServer{}
	if err := run(testserver.serveCreate); err == nil {
		t.Errorf("expected error due to redis connection error")
	}
}

func TestCreateEndToEnd(t *testing.T) {
	redis := startRedis("12313")
	defer stopRedis(redis)

	os.Setenv("REDIS_PORT_6379_TCP_ADDR", "127.0.0.1")
	os.Setenv("REDIS_PORT_6379_TCP_PORT", "12313")

	testserver := TestServer{}
	run(testserver.serveCreate)

	cases := []struct {
		name           string
		timestamp      string
		expectedStatus int
	}{
		{"test", "2015-02-18T13:26:00+00:00", http.StatusCreated},
		{"test", "2015/02/18", http.StatusBadRequest},
		{"test", "2015-02-18T13:26:00-08:00", http.StatusBadRequest},
	}

	for _, c := range cases {
		req := fmt.Sprintf(`{"name": "%s", "timestamp": "%s"}`, c.name, c.timestamp)
		res, err := http.Post(
			testserver.server.URL, "application/json", strings.NewReader(req))

		if err != nil {
			t.Errorf("unexpected error")
		}

		if res.StatusCode != c.expectedStatus {
			t.Errorf("expected status code %d, got %d", c.expectedStatus, res.StatusCode)
		}
	}
}
