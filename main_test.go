package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/stvp/tempredis"

	"github.com/declantraynor/go-events-service/interfaces/web"
)

type TestServer struct {
	server *httptest.Server
}

func (t *TestServer) serveCreate(webservice *web.WebService) {
	t.server = httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			webservice.Create(res, req)
		}))
}

func (t *TestServer) serveCount(webservice *web.WebService) {
	t.server = httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			webservice.Count(res, req)
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

func TestCreateEventEndToEnd(t *testing.T) {
	redis := startRedis("12313")
	defer stopRedis(redis)

	os.Setenv("REDIS_PORT_6379_TCP_ADDR", "127.0.0.1")
	os.Setenv("REDIS_PORT_6379_TCP_PORT", "12313")

	testserver := TestServer{}
	run(testserver.serveCreate)

	cases := []struct {
		name            string
		timestamp       string
		expectedStatus  int
		expectedContent string
	}{
		{
			"test",
			"2015-02-18T13:26:00+00:00",
			http.StatusCreated,
			`{}`,
		},
		{
			"test",
			"2015/02/18",
			http.StatusBadRequest,
			`{"error": "2015/02/18 does not conform to ISO8601"}`,
		},
		{
			"test",
			"2015-02-18T13:26:00-08:00",
			http.StatusBadRequest,
			`{"error": "2015-02-18T13:26:00-08:00 is not UTC"}`,
		},
	}

	for _, c := range cases {
		req := fmt.Sprintf(`{"name": "%s", "timestamp": "%s"}`, c.name, c.timestamp)
		res, err := http.Post(
			testserver.server.URL, "application/json", strings.NewReader(req))

		if err != nil {
			t.Errorf("unexpected error")
		}

		assertJSONResponse(t, res, c.expectedStatus, c.expectedContent)
	}
}

func TestCountEventsEndToEnd(t *testing.T) {
	redis := startRedis("12313")
	defer stopRedis(redis)

	populateRedis("127.0.0.1:12313")

	os.Setenv("REDIS_PORT_6379_TCP_ADDR", "127.0.0.1")
	os.Setenv("REDIS_PORT_6379_TCP_PORT", "12313")

	testserver := TestServer{}
	run(testserver.serveCount)

	cases := []struct {
		from            string
		to              string
		expectedStatus  int
		expectedContent string
	}{
		{
			"2015-01-01T00:00:00+00:00",
			"2015-01-01T01:22:59+00:00",
			http.StatusOK,
			`{"a": 2}`,
		},
		{
			"2015-01-01T01:23:00+00:00",
			"2015-01-02T19:23:44+00:00",
			http.StatusOK,
			`{"a": 2, "b": 3}`,
		},
		{
			"2015-01-02T19:23:45+00:00",
			"2015-01-03T13:16:12+00:00",
			http.StatusOK,
			`{"a": 2, "c": 3}`,
		},
		{
			"2015-01-03T13:16:13+00:00",
			"2015-01-03T23:59:00+00:00",
			http.StatusOK,
			`{"c": 1, "b": 1}`,
		},
		{
			"2015-02-01",
			"2015-01-03T23:59:00+00:00",
			http.StatusBadRequest,
			`{"error": "2015-02-01 does not conform to ISO8601"}`,
		},
		{
			"2015-02-01T13:16:13+00:00",
			"02/01/2015",
			http.StatusBadRequest,
			`{"error": "02/01/2015 does not conform to ISO8601"}`,
		},
		{
			"2015-01-03T13:16:13-05:00",
			"2015-01-03T23:59:00+00:00",
			http.StatusBadRequest,
			`{"error": "2015-01-03T13:16:13-05:00 is not UTC"}`,
		},
		{
			"2015-01-03T13:16:13+00:00",
			"2015-01-03T23:59:00+05:00",
			http.StatusBadRequest,
			`{"error": "2015-01-03T23:59:00+05:00 is not UTC"}`,
		},
		{
			"2015-01-04T13:16:13+00:00",
			"2015-01-03T23:59:00+00:00",
			http.StatusBadRequest,
			`{"error": "2015-01-04T13:16:13+00:00 is later than 2015-01-03T23:59:00+00:00"}`,
		},
	}

	for _, c := range cases {
		params := url.Values{}
		params.Set("from", c.from)
		params.Set("to", c.to)

		reqURL := fmt.Sprintf("%s?%s", testserver.server.URL, params.Encode())
		res, err := http.Get(reqURL)

		if err != nil {
			t.Errorf("unexpected error")
		}

		assertJSONResponse(t, res, c.expectedStatus, c.expectedContent)
	}
}

func assertJSONResponse(t *testing.T, res *http.Response, status int, content string) {
	if res.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		t.Error("expected JSON response")
	}

	if res.StatusCode != status {
		t.Errorf("expected status code %d, got %d", status, res.StatusCode)
	}

	contentExpected := make(map[string]interface{})
	json.Unmarshal([]byte(content), &contentExpected)

	defer res.Body.Close()
	responseBytes, _ := ioutil.ReadAll(res.Body)
	contentReceived := make(map[string]interface{})
	if err := json.Unmarshal(responseBytes, &contentReceived); err != nil {
		t.Error("response content is not valid JSON")
	}

	if !reflect.DeepEqual(contentExpected, contentReceived) {
		t.Errorf("expected content %q, got %q\n", content, string(responseBytes))
	}
}

func populateRedis(address string) {
	fixtures := []struct {
		name      string
		timestamp string
	}{
		{"a", "2015-01-01T00:00:00+00:00"},
		{"a", "2015-01-01T00:01:05+00:00"},

		{"a", "2015-01-01T01:23:00+00:00"},
		{"a", "2015-01-01T01:23:00+00:00"},
		{"b", "2015-01-01T02:56:10+00:00"},
		{"b", "2015-01-01T02:56:11+00:00"},
		{"b", "2015-01-01T02:57:34+00:00"},

		{"a", "2015-01-02T19:23:45+00:00"},
		{"a", "2015-01-02T20:18:34+00:00"},
		{"c", "2015-01-02T20:20:37+00:00"},
		{"c", "2015-01-02T22:10:14+00:00"},
		{"c", "2015-01-02T22:23:45+00:00"},

		{"c", "2015-01-03T13:16:13+00:00"},
		{"b", "2015-01-03T14:42:12+00:00"},
	}

	conn, _ := redis.Dial("tcp", address)
	conn.Send("MULTI")

	for i, f := range fixtures {
		eventKey := fmt.Sprintf("events:%d", i+1)
		eventTimestamp, _ := time.Parse(time.RFC3339, f.timestamp)
		indexKey := fmt.Sprintf("events:%s:by-timestamp", f.name)
		conn.Send("SADD", "event_names", f.name)
		conn.Send("HMSET", eventKey, "name", f.name, "timestamp", eventTimestamp.Unix())
		conn.Send("ZADD", indexKey, eventTimestamp.Unix(), eventKey)
	}

	conn.Do("EXEC")
}

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
